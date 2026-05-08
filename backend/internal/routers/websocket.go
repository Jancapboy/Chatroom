package routers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/Jancapboy/Chatroom/backend/global"
	"github.com/Jancapboy/Chatroom/backend/internal/chat"
	"github.com/Jancapboy/Chatroom/backend/internal/dao"
	"github.com/Jancapboy/Chatroom/backend/internal/model"
	"github.com/Jancapboy/Chatroom/backend/internal/service"
	"github.com/Jancapboy/Chatroom/backend/internal/simulation"
	"github.com/Jancapboy/Chatroom/backend/pkg/ws_protocol"
	"github.com/gin-gonic/gin"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

// ========== 原有广播式WebSocket（兼容） ==========

func WebsocketHandler(c *gin.Context) {
	conn, err := websocket.Accept(c.Writer, c.Request, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		log.Println("websocket accept error:", err)
		return
	}

	ID, exists := c.Get("UserID")
	if !exists {
		log.Println("jwt middleware error")
	}
	svc := service.New(c.Request.Context())
	userGot, err := svc.UserGet(ID.(uint64))
	if err != nil {
		log.Printf("get user error: %v", err)
		conn.Close(websocket.StatusInternalError, "user not found")
		return
	}
	user := chat.NewUser(userGot, conn)

	go user.SendMessage(c.Request.Context())

	chat.Broadcaster.UserEntering(user)
	msg := chat.NewUserEnterMessage(user)
	chat.Broadcaster.Broadcast(msg)
	log.Println("user:", user.Nickname, "joins chat")

	err = user.ReceiveMessage(c.Request.Context())

	chat.Broadcaster.UserLeaving(user)
	msg = chat.NewUserLeaveMessage(user)
	chat.Broadcaster.Broadcast(msg)
	log.Println("user:", user.Nickname, "leaves chat")

	if err == nil {
		conn.Close(websocket.StatusNormalClosure, "")
	} else {
		log.Println("read from client error:", err)
		conn.Close(websocket.StatusInternalError, "Read from client error")
	}
}

// ========== 新增房间隔离WebSocket ==========

// RoomHub 单个房间的WebSocket Hub
type RoomHub struct {
	roomID     string
	clients    map[*websocket.Conn]*WSClient
	broadcast  chan *ws_protocol.ServerMessage
	register   chan *WSClient
	unregister chan *WSClient
	mu         sync.RWMutex
}

type WSClient struct {
	conn     *websocket.Conn
	userID   uint64
	nickname string
	roomID   string
}

// 全局房间Hub管理器
var roomHubs = make(map[string]*RoomHub)
var roomHubsMu sync.RWMutex

func getOrCreateRoomHub(roomID string) *RoomHub {
	roomHubsMu.Lock()
	defer roomHubsMu.Unlock()

	if hub, ok := roomHubs[roomID]; ok {
		return hub
	}

	hub := &RoomHub{
		roomID:     roomID,
		clients:    make(map[*websocket.Conn]*WSClient),
		broadcast:  make(chan *ws_protocol.ServerMessage, 256),
		register:   make(chan *WSClient),
		unregister: make(chan *WSClient),
	}
	roomHubs[roomID] = hub
	go hub.run()
	return hub
}

func (h *RoomHub) run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.conn] = client
			h.mu.Unlock()
			log.Printf("[WS] 用户 %s 加入房间 %s", client.nickname, h.roomID)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.conn]; ok {
				delete(h.clients, client.conn)
			}
			h.mu.Unlock()
			client.conn.Close(websocket.StatusNormalClosure, "")
			log.Printf("[WS] 用户 %s 离开房间 %s", client.nickname, h.roomID)

		case msg := <-h.broadcast:
			h.mu.RLock()
			clients := make(map[*websocket.Conn]*WSClient, len(h.clients))
			for k, v := range h.clients {
				clients[k] = v
			}
			h.mu.RUnlock()

			for conn := range clients {
				if err := wsjson.Write(context.Background(), conn, msg); err != nil {
					log.Printf("[WS] 广播消息失败: %v", err)
				}
			}
		}
	}
}

// WebsocketRoomHandler 房间隔离的WebSocket处理器
func WebsocketRoomHandler(c *gin.Context) {
	roomID := c.Param("id")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "room id required"})
		return
	}

	// 验证房间存在
	svc := service.New(c.Request.Context())
	room, err := svc.RoomGet(roomID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
		return
	}

	// 升级WebSocket
	conn, wsErr := websocket.Accept(c.Writer, c.Request, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if wsErr != nil {
		log.Println("websocket accept error:", wsErr)
		return
	}

	// 获取用户信息（JWT已验证）
	var userID uint64 = 0
	var nickname string = "游客"
	if id, exists := c.Get("UserID"); exists {
		userID, _ = id.(uint64)
		userGot, err := svc.UserGet(userID)
		if err == nil && userGot != nil {
			nickname = userGot.Nickname
		}
	}

	// 创建房间Hub
	hub := getOrCreateRoomHub(roomID)

	client := &WSClient{
		conn:     conn,
		userID:   userID,
		nickname: nickname,
		roomID:   roomID,
	}

	hub.register <- client

	// 发送欢迎消息
	welcome := ws_protocol.NewSystemMessage(roomID, "user_joined", nickname+" 加入房间")
	hub.broadcast <- welcome

	// 如果房间在运行中，启动引擎
	if room.Status == "running" {
		engine := simulation.GlobalEngineManager.Get(roomID)
		if engine == nil {
			d := dao.New(global.DBEngine)
			agents, _ := d.AgentListByRoom(roomID)
			roomModel := &model.Room{
				ID:             room.ID,
				Name:           room.Name,
				Topic:          room.Topic,
				Status:         room.Status,
				CurrentPhase:   room.CurrentPhase,
				CurrentRound:   room.CurrentRound,
				MaxRounds:      room.MaxRounds,
				ConsensusScore: room.ConsensusScore,
			}
			engine = simulation.GlobalEngineManager.GetOrCreate(roomModel, agents, hub.broadcast)
			go engine.Run()
		}
	}

	// 监听客户端消息
	ctx := c.Request.Context()
	for {
		var msg ws_protocol.ClientMessage
		err := wsjson.Read(ctx, conn, &msg)
		if err != nil {
			log.Printf("[WS] 读取客户端消息失败: %v", err)
			break
		}

		switch msg.Type {
		case "user_message":
			var payload ws_protocol.UserMessagePayload
			if err := json.Unmarshal(msg.Payload, &payload); err == nil {
				// 保存并广播用户消息
				engine := simulation.GlobalEngineManager.Get(roomID)
				if engine != nil {
					engine.HandleUserMessage(userID, nickname, payload.Content)
				} else {
					// 房间未启动引擎，仅广播
					hub.broadcast <- ws_protocol.NewUserMessage(roomID, userID, nickname, payload.Content, "", 0)
				}
			}

		case "command":
			var payload ws_protocol.CommandPayload
			if err := json.Unmarshal(msg.Payload, &payload); err == nil {
				handleCommand(room, hub, payload.Command, roomID)
			}
		}
	}

	hub.unregister <- client
	leaveMsg := ws_protocol.NewSystemMessage(roomID, "user_left", nickname+" 离开房间")
	hub.broadcast <- leaveMsg
}

// handleCommand 处理WS命令
func handleCommand(room *service.RoomResponse, hub *RoomHub, command, roomID string) {
	switch command {
	case "pause":
		engine := simulation.GlobalEngineManager.Get(roomID)
		if engine != nil {
			engine.Pause()
		}
		hub.broadcast <- ws_protocol.NewSystemMessage(roomID, "paused", "推演已暂停")

	case "resume":
		engine := simulation.GlobalEngineManager.Get(roomID)
		if engine != nil {
			engine.Resume()
		} else {
			// 重新创建引擎并运行
			d := dao.New(global.DBEngine)
			agents, _ := d.AgentListByRoom(roomID)
			roomModel := &model.Room{
				ID:             room.ID,
				Name:           room.Name,
				Topic:          room.Topic,
				Status:         "running",
				CurrentPhase:   room.CurrentPhase,
				CurrentRound:   room.CurrentRound,
				MaxRounds:      room.MaxRounds,
				ConsensusScore: room.ConsensusScore,
			}
			engine = simulation.GlobalEngineManager.GetOrCreate(roomModel, agents, hub.broadcast)
			go engine.Run()
		}
		hub.broadcast <- ws_protocol.NewSystemMessage(roomID, "resumed", "推演已恢复")
	}
}

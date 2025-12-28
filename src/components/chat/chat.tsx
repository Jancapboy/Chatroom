import "./chat.scss";
import Input from "../input/input";
import Message from "../message/message";
import { useEffect, useLayoutEffect, useRef, useState } from "react";
import wsHostname from "../../api/wshost";
import IUserInfo from "../../types/IUserInfo";
import { message, Radio } from "antd";

interface IChatParams {
  token: string;
  userInfo: IUserInfo;
}

interface IMessageInfo {
  user: {
    user_id: number;
    nickname: string;
  };
  message_content: string;
  send_time: number; // 秒
}

// 新增：聊天模式类型
type ChatMode = "public" | "ai";

function Chat(params: IChatParams) {
  const [connectReady, setConnectReady] = useState<boolean>(false);
  const [msgList, setMsgList] = useState<IMessageInfo[]>([]);
  const [chatMode, setChatMode] = useState<ChatMode>("public"); // 新增：聊天模式状态
  const ws = useRef<WebSocket | null>(null);
  const msgWindow = useRef<HTMLDivElement>(null);

  // 新增：AI API地址（根据你的后端地址调整）
  const AI_API_URL = "http://192.168.40.132:4001/api/ai/chat";

  useLayoutEffect(() => {
    // 只有在公开聊天模式才连接WebSocket
    if (chatMode === "public" && params.token !== "") {
      ws.current = new WebSocket(wsHostname + "/?token=" + params.token);
      ws.current.onopen = () => {
        setConnectReady(true);
      };
      ws.current.onclose = () => {
        setConnectReady(false);
      };
      ws.current.onerror = (e) => {
        console.error("WebSocket错误:", e);
        message.error("连接聊天室失败", 2);
      };
      ws.current.onmessage = handleReceiveMsg;
    } else {
      // AI模式下关闭WebSocket连接
      if (ws.current) {
        ws.current.close();
        ws.current = null;
        setConnectReady(false);
      }
    }

    return () => {
      ws.current?.close();
    };
  }, [params.token, chatMode]);

  useEffect(() => {
    if (msgWindow.current) {
      msgWindow.current.scrollTop = msgWindow.current.scrollHeight;
    }
  }, [msgList]);

  const nowTime = new Date();

  // 新增：调用AI API的函数
  const callAIApi = async (content: string): Promise<string> => {
  try {
    const response = await fetch(AI_API_URL, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        messages: [
          {
            role: "user",
            content: content,
          },
        ],
      }),
    });

    if (!response.ok) {
      throw new Error(`HTTP错误: ${response.status}`);
    }

    const data = await response.json();
    
    // 现在只处理统一格式 { success: true, content: "..." }
    if (data.success === true && data.content) {
      return data.content;
    } else {
      throw new Error(data.error || "AI服务返回失败");
    }
  } catch (error) {
    console.error("调用AI失败:", error);
    throw error;
  }
};
  
  const sendMsg = async (Content: string) => {
    if (chatMode === "public") {
      // 公开聊天模式：通过WebSocket发送
      ws.current?.send(
        JSON.stringify({
          message_content: Content,
        })
      );
      let msgData: IMessageInfo = {
        user: {
          user_id: params.userInfo.user_id,
          nickname: params.userInfo.nickname,
        },
        message_content: Content,
        send_time: new Date().getTime(),
      };
      setMsgList((list) => list.concat([msgData]));
    } else {
      // AI聊天模式
      // 1. 添加用户消息
      const userMsg: IMessageInfo = {
        user: {
          user_id: params.userInfo.user_id,
          nickname: params.userInfo.nickname,
        },
        message_content: Content,
        send_time: new Date().getTime(),
      };
      setMsgList((list) => list.concat([userMsg]));

      // 2. 添加AI"思考中"消息
      const thinkingMsg: IMessageInfo = {
        user: {
          user_id: -1, // 使用-1表示AI
          nickname: "AI助手",
        },
        message_content: "正在思考...",
        send_time: new Date().getTime(),
      };
      const newMsgList = [...msgList, userMsg, thinkingMsg];
      setMsgList(newMsgList);

      // 3. 调用AI API
      try {
        const aiResponse = await callAIApi(Content);
        
        // 4. 替换"思考中"消息为AI回复
        setMsgList((list) => {
          const updatedList = [...list];
          // 找到最后一个消息（应该是思考中消息）
          const lastIndex = updatedList.length - 1;
          if (lastIndex >= 0 && updatedList[lastIndex].user.user_id === -1) {
            updatedList[lastIndex] = {
              user: {
                user_id: -1,
                nickname: "AI助手",
              },
              message_content: aiResponse,
              send_time: new Date().getTime(),
            };
          }
          return updatedList;
        });
      } catch (error) {
        // 出错时替换为错误消息
        setMsgList((list) => {
          const updatedList = [...list];
          const lastIndex = updatedList.length - 1;
          if (lastIndex >= 0 && updatedList[lastIndex].user.user_id === -1) {
            updatedList[lastIndex] = {
              user: {
                user_id: -1,
                nickname: "AI助手",
              },
              message_content: "抱歉，AI助手暂时无法响应，请稍后再试。",
              send_time: new Date().getTime(),
            };
          }
          return updatedList;
        });
        message.error("AI服务暂时不可用", 2);
      }
    }
  };

  const handleReceiveMsg = (event: MessageEvent) => {
    // 只在公开聊天模式下处理WebSocket消息
    if (chatMode !== "public") return;
    
    let data = JSON.parse(event.data) as IMessageInfo;
    data.send_time *= 1000;
    if (data.user.user_id === params.userInfo.user_id) {
      return;
    }
    setMsgList((list) => list.concat([data]));
  };

  // 新增：处理模式切换
  const handleModeChange = (mode: ChatMode) => {
    setChatMode(mode);
    setMsgList([]); // 切换时清空消息
    
    // 可选：AI模式下添加欢迎消息
    if (mode === "ai") {
      setTimeout(() => {
        const welcomeMsg: IMessageInfo = {
          user: {
            user_id: -1,
            nickname: "AI助手",
          },
          message_content: "你好！我是AI助手，可以帮你解答各种问题。请问有什么可以帮你的？",
          send_time: new Date().getTime(),
        };
        setMsgList([welcomeMsg]);
      }, 100);
    }
  };

  return (
  <div className="chat">
    {/* ========== 新增：AI模式切换按钮 ========== */}
    <div style={{
      padding: "10px 15px", 
      background: "white", 
      borderBottom: "1px solid #eee",
      display: "flex",
      justifyContent: "center",
      gap: "10px"
    }}>
      <button 
        onClick={() => handleModeChange("public")}
        style={{
          padding: "8px 20px",
          background: chatMode === "public" ? "#1890ff" : "#f5f5f5",
          color: chatMode === "public" ? "white" : "#333",
          border: "1px solid #d9d9d9",
          borderRadius: "4px",
          cursor: "pointer"
        }}
      >
        公开聊天
      </button>
      <button 
        onClick={() => handleModeChange("ai")}
        style={{
          padding: "8px 20px",
          background: chatMode === "ai" ? "#52c41a" : "#f5f5f5",
          color: chatMode === "ai" ? "white" : "#333",
          border: "1px solid #d9d9d9",
          borderRadius: "4px",
          cursor: "pointer"
        }}
      >
        🤖 AI助手
      </button>
      
      <div style={{
        marginLeft: "20px",
        fontSize: "14px",
        color: "#666",
        display: "flex",
        alignItems: "center"
      }}>
        <span style={{
          display: "inline-block",
          width: "8px",
          height: "8px",
          borderRadius: "50%",
          background: chatMode === "public" 
            ? (connectReady ? "#52c41a" : "#faad14") 
            : "#1890ff",
          marginRight: "8px"
        }}></span>
        {chatMode === "public" 
          ? (connectReady ? "已连接" : "连接中...") 
          : "AI助手模式"}
      </div>
    </div>
    {/* ========== 新增结束 ========== */}

    <div className="messages" ref={msgWindow}>
      {msgList.map((msg, index) => {
        // 状态判断逻辑
        let status: "receive" | "sending" | "sent" | "fail" | "ai";
        
        if (msg.user.user_id === params.userInfo.user_id) {
          status = "sent";
        } else if (msg.user.user_id === -1) {
          status = "ai";
        } else {
          status = "receive";
        }
        
        return (
          <Message
            key={index}
            content={msg.message_content}
            nickname={msg.user.nickname}
            nowTime={nowTime}
            sentTime={msg.send_time}
            status={status}
          />
        );
      })}
    </div>
    <Input sendMessage={sendMsg} />
  </div>
);
}

export default Chat;
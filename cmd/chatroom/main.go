package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Jancapboy/Chatroom/global"
	"github.com/Jancapboy/Chatroom/internal/chat"
	"github.com/Jancapboy/Chatroom/internal/model"
	"github.com/Jancapboy/Chatroom/internal/routers"
	"github.com/Jancapboy/Chatroom/internal/setting"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func init() {
	err := setupSettings()
	if err != nil {
		log.Fatalf("init.setupSettings err: %v", err)
	}
	err = setupDBEngine()
	if err != nil {
		log.Fatalf("init.setupDBEngine err: %v", err)
	}
	err = SetupTableModel(global.DBEngine, &model.User{})
	if err != nil {
		log.Fatalf("init.setupTableModel err: %v", err)
	}

	go chat.Broadcaster.Start()
}

func main() {
	gin.SetMode(global.ServerSettings.RunMode)
	router := routers.NewRouter()
	s := &http.Server{
		Addr:           ":" + global.ServerSettings.HttpPort,
		Handler:        router,
		ReadTimeout:    global.ServerSettings.ReadTimeout,
		WriteTimeout:   global.ServerSettings.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	if err := s.ListenAndServe(); err != nil {
		log.Fatalf("ListenAndServe err: %v", err)
	}
}

func setupSettings() error {
	settings, err := setting.NewSetting()
	if err != nil {
		return err
	}

	// 读取各个配置节
	err = settings.ReadSection("Server", &global.ServerSettings)
	if err != nil {
		return err
	}
	err = settings.ReadSection("Database", &global.DatabaseSettings)
	if err != nil {
		return err
	}
	err = settings.ReadSection("JWT", &global.JWTSettings)
	if err != nil {
		return err
	}
	err = settings.ReadSection("Chatroom", &global.ChatroomSettings)
	if err != nil {
		return err
	}

	// 新增：读取 AI 配置
	aiConfig, err := settings.GetAIConfig()
	if err != nil {
		log.Printf("警告：读取AI配置失败，将使用默认值: %v", err)
		// 设置默认 AI 配置
		global.AISettings = &setting.AISettingS{
			Model:   "deepseek-chat",
			BaseURL: "https://api.deepseek.com",
		}
	} else {
		global.AISettings = aiConfig
		log.Printf("AI配置加载成功: model=%s", global.AISettings.Model)
	}

	// 转换时间单位
	global.ServerSettings.ReadTimeout *= time.Second
	global.ServerSettings.WriteTimeout *= time.Second
	global.JWTSettings.Expire *= time.Second
	return nil
}
func setupDBEngine() error {
	var err error
	global.DBEngine, err = model.NewDBEngine(global.DatabaseSettings)
	if err != nil {
		return err
	}
	return nil
}

func SetupTableModel(db *gorm.DB, models interface{}) error {
	err := db.AutoMigrate(models)
	if err != nil {
		return err
	}
	return nil
}

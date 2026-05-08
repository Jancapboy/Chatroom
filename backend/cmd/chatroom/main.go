package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Jancapboy/Chatroom/backend/global"
	"github.com/Jancapboy/Chatroom/backend/internal/chat"
	"github.com/Jancapboy/Chatroom/backend/internal/model"
	"github.com/Jancapboy/Chatroom/backend/internal/routers"
	"github.com/Jancapboy/Chatroom/backend/internal/setting"
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

	// 注册所有数据表
	err = SetupTableModel(global.DBEngine,
		&model.User{},
		&model.Room{},
		&model.RoomAgent{},
		&model.Message{},
		&model.AgentTemplate{},
		&model.RoomSnapshot{},
		&model.AgentMemory{},
	)
	if err != nil {
		log.Fatalf("init.setupTableModel err: %v", err)
	}

	// 初始化Agent模板数据
	err = initAgentTemplates()
	if err != nil {
		log.Printf("init.initAgentTemplates warning: %v", err)
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

	aiConfig, err := settings.GetAIConfig()
	if err != nil {
		log.Printf("警告：读取AI配置失败，将使用默认值: %v", err)
		global.AISettings = &setting.AISettingS{
			Model:   "deepseek-chat",
			BaseURL: "https://api.deepseek.com",
		}
	} else {
		global.AISettings = aiConfig
		log.Printf("AI配置加载成功: model=%s", global.AISettings.Model)
	}

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

func SetupTableModel(db *gorm.DB, models ...interface{}) error {
	err := db.AutoMigrate(models...)
	if err != nil {
		return err
	}
	return nil
}

// initAgentTemplates 初始化5个预设Agent模板
func initAgentTemplates() error {
	templates := []model.AgentTemplate{
		{
			ID:                   "tpl-arch",
			Name:                 "系统架构师",
			Role:                 "architect",
			Personality:          "理性、严谨，关注技术可行性和资源约束",
			Expertise:              `["系统架构","技术选型","资源规划"]`,
			DefaultModel:         "deepseek-chat",
			SystemPromptTemplate: "你是一位资深的系统架构师。你在讨论中的关注点是：技术可行性、系统架构合理性、资源消耗、扩展性。你会用工程思维分析问题，给出具体的架构建议。当前讨论主题：%s。",
			SortOrder:            1,
		},
		{
			ID:                   "tpl-risk",
			Name:                 "风险官",
			Role:                 "risk_officer",
			Personality:          "谨慎、质疑，关注安全、伦理和边界条件",
			Expertise:              `["风险评估","安全","合规","伦理"]`,
			DefaultModel:         "deepseek-chat",
			SystemPromptTemplate: "你是一位严格的风险官。你的职责是识别所有潜在风险：安全风险、合规风险、伦理风险、财务风险。你对\"一票否决\"权非常谨慎，只在真正不可接受的风险时使用。当前讨论主题：%s。",
			SortOrder:            2,
		},
		{
			ID:                   "tpl-strat",
			Name:                 "策略家",
			Role:                 "strategist",
			Personality:          "果断、全局视野，关注目标达成和效率",
			Expertise:              `["战略规划","竞争分析","资源优化"]`,
			DefaultModel:         "deepseek-chat",
			SystemPromptTemplate: "你是一位高瞻远瞩的策略家。你的关注点是：目标达成路径、竞争优势、资源最优配置、时机把握。你善于从全局角度思考，给出方向性建议。当前讨论主题：%s。",
			SortOrder:            3,
		},
		{
			ID:                   "tpl-analyst",
			Name:                 "数据分析师",
			Role:                 "analyst",
			Personality:          "客观、数据驱动，用数字说话",
			Expertise:              `["数据分析","量化评估","概率计算"]`,
			DefaultModel:         "deepseek-chat",
			SystemPromptTemplate: "你是一位冷静的数据分析师。你要求每个观点都要有数据支撑。你会主动计算概率、估算数值、寻找量化依据。如果缺乏数据，你会指出这一点。当前讨论主题：%s。",
			SortOrder:            4,
		},
		{
			ID:                   "tpl-exec",
			Name:                 "执行者",
			Role:                 "executor",
			Personality:          "务实、注重细节，关注落地和时间线",
			Expertise:              `["项目管理","执行落地","时间管理","细节把控"]`,
			DefaultModel:         "deepseek-chat",
			SystemPromptTemplate: "你是一位务实的执行者。你的关注点：具体落地步骤、时间节点、执行细节、依赖关系。你会把抽象方案转化为可执行的计划。当前讨论主题：%s。",
			SortOrder:            5,
		},
	}

	for _, tpl := range templates {
		var existing model.AgentTemplate
		result := global.DBEngine.Where("id = ?", tpl.ID).First(&existing)
		if result.Error != nil {
			// 不存在则创建
			if err := global.DBEngine.Create(&tpl).Error; err != nil {
				log.Printf("创建Agent模板 %s 失败: %v", tpl.Name, err)
			} else {
				log.Printf("初始化Agent模板: %s", tpl.Name)
			}
		}
	}
	return nil
}

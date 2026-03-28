package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Jancapboy/Chatroom/internal/service"
	"github.com/gin-gonic/gin"
)

// 独立测试服务器 - 只验证3D接口，不需要MySQL
func main() {
	// 从环境变量或凭证文件读取
	secretID := os.Getenv("TENCENT_SECRET_ID")
	secretKey := os.Getenv("TENCENT_SECRET_KEY")

	if secretID == "" || secretKey == "" {
		// 尝试从凭证文件读取
		data, err := os.ReadFile(os.ExpandEnv("$HOME/.openclaw/.tencent_credentials"))
		if err == nil {
			for _, line := range splitLines(string(data)) {
				if len(line) > 20 && line[:20] == "TENCENT_SECRET_ID=" {
					secretID = line[19:]
				} else if len(line) > 19 && line[:19] == "TENCENT_SECRET_KEY=" {
					secretKey = line[19:]
				}
			}
		}
	}

	// 简单解析凭证文件
	if secretID == "" {
		data, _ := os.ReadFile(os.ExpandEnv("$HOME/.openclaw/.tencent_credentials"))
		lines := splitLines(string(data))
		for _, l := range lines {
			if len(l) > 0 {
				for i := range l {
					if l[i] == '=' {
						key := l[:i]
						val := l[i+1:]
						switch key {
						case "TENCENT_SECRET_ID":
							secretID = val
						case "TENCENT_SECRET_KEY":
							secretKey = val
						}
						break
					}
				}
			}
		}
	}

	if secretID == "" || secretKey == "" {
		log.Fatal("请设置 TENCENT_SECRET_ID 和 TENCENT_SECRET_KEY")
	}

	config := &service.Hunyuan3DConfig{
		SecretID:  secretID,
		SecretKey: secretKey,
	}

	r := gin.Default()

	// CORS
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "3d-test-server"})
	})

	// 3D生成
	r.POST("/api/3d/generate", func(c *gin.Context) {
		var req struct {
			Prompt       string `json:"prompt" binding:"required"`
			ResultFormat string `json:"result_format"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": err.Error()})
			return
		}
		format := req.ResultFormat
		if format == "" {
			format = "GLB"
		}
		jobId, err := service.SubmitHunyuan3DJob(config, &service.Hunyuan3DSubmitRequest{
			Prompt:       req.Prompt,
			ResultFormat: format,
		})
		if err != nil {
			c.JSON(500, gin.H{"code": 1, "msg": err.Error()})
			return
		}
		log.Printf("✅ 任务提交成功: %s", jobId)
		c.JSON(200, gin.H{"code": 0, "msg": "任务已提交", "job_id": jobId})
	})

	// 3D查询
	r.POST("/api/3d/query", func(c *gin.Context) {
		var req struct {
			JobId string `json:"job_id" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": err.Error()})
			return
		}
		result, err := service.QueryHunyuan3DJob(config, req.JobId)
		if err != nil {
			c.JSON(500, gin.H{"code": 1, "msg": err.Error()})
			return
		}
		resp := gin.H{"code": 0, "status": result.Response.Status}
		if result.Response.Status == "DONE" && len(result.Response.ResultFile3Ds) > 0 {
			files := make([]gin.H, 0)
			for _, f := range result.Response.ResultFile3Ds {
				files = append(files, gin.H{"type": f.Type, "url": f.Url})
			}
			resp["files"] = files
		}
		c.JSON(200, resp)
	})

	fmt.Println("🚀 3D测试服务器启动: http://localhost:4002")
	r.Run(":4002")
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := range s {
		if s[i] == '\n' {
			line := s[start:i]
			if len(line) > 0 && line[len(line)-1] == '\r' {
				line = line[:len(line)-1]
			}
			lines = append(lines, line)
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

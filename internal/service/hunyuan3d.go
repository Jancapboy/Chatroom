package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	Hunyuan3DHost    = "ai3d.tencentcloudapi.com"
	Hunyuan3DService = "ai3d"
	Hunyuan3DVersion = "2025-05-13"
	Hunyuan3DRegion  = "ap-guangzhou"
)

// Hunyuan3DConfig 混元3D配置
type Hunyuan3DConfig struct {
	SecretID  string
	SecretKey string
}

// Hunyuan3DSubmitRequest 提交任务请求
type Hunyuan3DSubmitRequest struct {
	Prompt       string `json:"Prompt,omitempty"`
	ImageUrl     string `json:"ImageUrl,omitempty"`
	ResultFormat string `json:"ResultFormat,omitempty"`
}

// Hunyuan3DResponse 腾讯云API通用响应
type Hunyuan3DResponse struct {
	Response struct {
		JobId     string `json:"JobId"`
		Status    string `json:"Status"`
		RequestId string `json:"RequestId"`
		Error     *struct {
			Code    string `json:"Code"`
			Message string `json:"Message"`
		} `json:"Error,omitempty"`
		ResultFile3Ds []struct {
			Type string `json:"Type"`
			Url  string `json:"Url"`
		} `json:"ResultFile3Ds,omitempty"`
	} `json:"Response"`
}

// tc3Sign 生成TC3-HMAC-SHA256签名
func tc3Sign(secretKey, date, service, toSign string) string {
	hmacSHA256 := func(key []byte, data string) []byte {
		h := hmac.New(sha256.New, key)
		h.Write([]byte(data))
		return h.Sum(nil)
	}

	secretDate := hmacSHA256([]byte("TC3"+secretKey), date)
	secretService := hmacSHA256(secretDate, service)
	secretSigning := hmacSHA256(secretService, "tc3_request")
	h := hmac.New(sha256.New, secretSigning)
	h.Write([]byte(toSign))
	return hex.EncodeToString(h.Sum(nil))
}

// tc3Request 发送腾讯云API请求
func tc3Request(config *Hunyuan3DConfig, action string, params interface{}) (*Hunyuan3DResponse, error) {
	payload, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("marshal params: %w", err)
	}

	timestamp := time.Now().Unix()
	date := time.Unix(timestamp, 0).UTC().Format("2006-01-02")

	// 1. 规范请求
	hashedPayload := sha256Hex(string(payload))
	canonicalRequest := fmt.Sprintf("POST\n/\n\ncontent-type:application/json\nhost:%s\n\ncontent-type;host\n%s",
		Hunyuan3DHost, hashedPayload)

	// 2. 待签名字符串
	credentialScope := fmt.Sprintf("%s/%s/tc3_request", date, Hunyuan3DService)
	hashedCanonical := sha256Hex(canonicalRequest)
	stringToSign := fmt.Sprintf("TC3-HMAC-SHA256\n%d\n%s\n%s",
		timestamp, credentialScope, hashedCanonical)

	// 3. 签名
	signature := tc3Sign(config.SecretKey, date, Hunyuan3DService, stringToSign)

	// 4. Authorization
	authorization := fmt.Sprintf(
		"TC3-HMAC-SHA256 Credential=%s/%s, SignedHeaders=content-type;host, Signature=%s",
		config.SecretID, credentialScope, signature)

	// 5. 发请求
	req, err := http.NewRequest("POST", "https://"+Hunyuan3DHost+"/",
		strings.NewReader(string(payload)))
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	req.Header.Set("Authorization", authorization)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Host", Hunyuan3DHost)
	req.Header.Set("X-TC-Action", action)
	req.Header.Set("X-TC-Version", Hunyuan3DVersion)
	req.Header.Set("X-TC-Timestamp", fmt.Sprintf("%d", timestamp))
	req.Header.Set("X-TC-Region", Hunyuan3DRegion)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http do: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	var result Hunyuan3DResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	if result.Response.Error != nil {
		return &result, fmt.Errorf("api error: %s - %s",
			result.Response.Error.Code, result.Response.Error.Message)
	}

	return &result, nil
}

func sha256Hex(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// SubmitHunyuan3DJob 提交3D生成任务
func SubmitHunyuan3DJob(config *Hunyuan3DConfig, req *Hunyuan3DSubmitRequest) (string, error) {
	if req.ResultFormat == "" {
		req.ResultFormat = "GLB"
	}
	result, err := tc3Request(config, "SubmitHunyuanTo3DJob", req)
	if err != nil {
		return "", err
	}
	return result.Response.JobId, nil
}

// QueryHunyuan3DJob 查询任务状态
func QueryHunyuan3DJob(config *Hunyuan3DConfig, jobId string) (*Hunyuan3DResponse, error) {
	params := map[string]string{"JobId": jobId}
	return tc3Request(config, "QueryHunyuanTo3DJob", params)
}

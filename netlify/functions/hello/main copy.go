package main

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Rule struct {
	Prefix   string
	Protocol string
}

func getUrl(str string) string {

	ma := []Rule{
		{Prefix: "/0/", Protocol: ""},
		{Prefix: "/1/", Protocol: "http://"},
		{Prefix: "/2/", Protocol: "https://"},
	}

	for _, v := range ma {
		if strings.HasPrefix(str, v.Prefix) {
			return strings.Replace(str, v.Prefix, v.Protocol, 1)
		}
	}

	return ""
}

func proxy(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// 目标 URL
	targetURL := getUrl(request.Path) // 替换为你要转发的目标 URL

	// 创建一个新的 HTTP 请求
	httpReq, err := http.NewRequest(
		request.HTTPMethod, // 使用原始请求的方法
		targetURL,
		bytes.NewBuffer([]byte(request.Body)), // 使用原始请求的 Body
	)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Failed to create HTTP request",
		}, err
	}

	// 设置请求头
	for key, value := range request.Headers {
		httpReq.Header.Set(key, value)
	}

	// 发起 HTTP 请求
	client := &http.Client{}
	httpResp, err := client.Do(httpReq)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Failed to send HTTP request",
		}, err
	}
	defer httpResp.Body.Close()

	// 读取响应体
	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Failed to read response body",
		}, err
	}

	// 将 http.Header 转换为 map[string]string
	headers := make(map[string]string)
	for key, values := range httpResp.Header {
		// 只取第一个值，如果有多个值，可以根据需要处理
		if len(values) > 0 {
			headers[key] = values[0]
		}
	}

	// 构建 API Gateway 响应
	return &events.APIGatewayProxyResponse{
		StatusCode: httpResp.StatusCode,
		Headers:    headers,
		Body:       string(respBody),
	}, nil
}

func main() {
	lambda.Start(proxy)
}

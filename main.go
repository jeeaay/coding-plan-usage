package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type usageDetail struct {
	Usage     string `json:"usage"`
	ResetTime string `json:"resetTime"`
}

type usageResponse struct {
	Hours5 usageDetail `json:"hours5"`
	Week   usageDetail `json:"week"`
	Month  usageDetail `json:"month"`
}

func shouldPrintAgentBrowserOutput(client *AliyunCodingLoginClient) bool {
	_ = client
	return false
}

func main() {
	client, output, err := NewAliyunCodingLoginClient()
	if err != nil {
		if strings.TrimSpace(output) != "" {
			fmt.Print(output)
		}
		fmt.Printf("Error executing command: %v\n", err)
		os.Exit(1)
	}
	if shouldPrintAgentBrowserOutput(client) && strings.TrimSpace(output) != "" {
		fmt.Print(output)
	}

	result, runOutput, err := client.Run()
	if shouldPrintAgentBrowserOutput(client) && strings.TrimSpace(runOutput) != "" {
		fmt.Print(runOutput)
	}
	if err != nil {
		fmt.Printf("Error running aliyun login flow: %v\n", err)
		os.Exit(1)
	}

	hasUsage := strings.TrimSpace(result.Hours5Usage) != "" || strings.TrimSpace(result.WeekUsage) != "" || strings.TrimSpace(result.MonthUsage) != ""
	if hasUsage {
		response := usageResponse{
			Hours5: usageDetail{Usage: result.Hours5Usage, ResetTime: result.Hours5ResetTime},
			Week:   usageDetail{Usage: result.WeekUsage, ResetTime: result.WeekResetTime},
			Month:  usageDetail{Usage: result.MonthUsage, ResetTime: result.MonthResetTime},
		}
		payload, marshalErr := json.MarshalIndent(response, "", "  ")
		if marshalErr != nil {
			fmt.Printf("Error marshaling usage response: %v\n", marshalErr)
			os.Exit(1)
		}
		fmt.Println(string(payload))

		closeOutput, closeErr := client.Close()
		if shouldPrintAgentBrowserOutput(client) && strings.TrimSpace(closeOutput) != "" {
			fmt.Print(closeOutput)
		}
		if closeErr != nil {
			fmt.Printf("Error closing session: %v\n", closeErr)
			os.Exit(1)
		}
	} else {
		fmt.Printf("Already logged in: %t\n", result.AlreadyLoggedIn)
		fmt.Printf("Entered login page: %t\n", result.EnteredLogin)
		if strings.TrimSpace(result.ScreenshotPath) != "" {
			fmt.Println("请使用阿里云 App 扫码完成登录后，再次执行此程序以查询用量。")
			fmt.Printf("Login screenshot: %s\n", result.ScreenshotPath)
		}
		fmt.Printf("Scan completed: %t\n", result.ScanCompleted)
	}
}

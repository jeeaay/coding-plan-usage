package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	aliyunHomeURL        = "https://www.aliyun.com/"
	aliyunDirectLoginURL = "https://account.aliyun.com/login/login.htm"
	codingPlanDetailURL  = "https://bailian.console.aliyun.com/cn-beijing/?tab=coding-plan#/efm/detail"
)

var loginRefRegex = regexp.MustCompile(`\[ref=(e\d+)\]`)

type AliyunCodingLoginResult struct {
	AlreadyLoggedIn bool
	EnteredLogin    bool
	ScreenshotPath  string
	ScanCompleted   bool
	Hours5ResetTime string
	Hours5Usage     string
	WeekResetTime   string
	WeekUsage       string
	MonthResetTime  string
	MonthUsage      string
}

type AliyunCodingLoginClient struct {
	session *AgentBrowserSession
}

func NewAliyunCodingLoginClient() (*AliyunCodingLoginClient, string, error) {
	session, output, err := NewAgentBrowserSession(aliyunHomeURL)
	if err != nil {
		return nil, output, err
	}

	return &AliyunCodingLoginClient{
		session: session,
	}, output, nil
}

func (c *AliyunCodingLoginClient) Close() (string, error) {
	if c.session == nil {
		return "", nil
	}
	return c.session.Close()
}

func (c *AliyunCodingLoginClient) Run() (*AliyunCodingLoginResult, string, error) {
	result := &AliyunCodingLoginResult{}
	logs := make([]string, 0, 8)

	waitOutput, err := c.WaitWindowReady()
	if strings.TrimSpace(waitOutput) != "" {
		logs = append(logs, waitOutput)
	}
	if err != nil {
		return nil, strings.Join(logs, "\n"), err
	}

	snapshotOutput, err := c.SnapshotInteractive()
	if c.shouldPrintSnapshotOutput() && strings.TrimSpace(snapshotOutput) != "" {
		logs = append(logs, snapshotOutput)
	}
	if err != nil {
		return nil, strings.Join(logs, "\n"), err
	}

	result.AlreadyLoggedIn = c.IsLoggedInFromSnapshot(snapshotOutput)
	if result.AlreadyLoggedIn {
		openOutput, openErr := c.OpenCodingPlanDetail()
		if strings.TrimSpace(openOutput) != "" {
			logs = append(logs, openOutput)
		}
		if openErr != nil {
			return nil, strings.Join(logs, "\n"), openErr
		}
		detailSnapshot, detailErr := c.SnapshotAll()
		if c.shouldPrintSnapshotOutput() && strings.TrimSpace(detailSnapshot) != "" {
			logs = append(logs, detailSnapshot)
		}
		if detailErr != nil {
			return nil, strings.Join(logs, "\n"), detailErr
		}
		c.FillUsageFromSnapshot(result, detailSnapshot)
		if result.Hours5Usage == "" && result.WeekUsage == "" && result.MonthUsage == "" {
			fallbackSnapshot, fallbackErr := c.SnapshotInteractive()
			if c.shouldPrintSnapshotOutput() && strings.TrimSpace(fallbackSnapshot) != "" {
				logs = append(logs, fallbackSnapshot)
			}
			if fallbackErr == nil {
				c.FillUsageFromSnapshot(result, fallbackSnapshot)
			}
		}
		result.ScanCompleted = true
		return result, strings.Join(logs, "\n"), nil
	}

	loginRef, err := c.FindLoginRef(snapshotOutput)
	if err != nil {
		directLoginOutput, directLoginErr := c.OpenDirectLoginPage()
		if strings.TrimSpace(directLoginOutput) != "" {
			logs = append(logs, directLoginOutput)
		}
		if directLoginErr != nil {
			return nil, strings.Join(logs, "\n"), fmt.Errorf("%v; direct login fallback failed: %w", err, directLoginErr)
		}
		screenshotPath, screenshotOutput, screenshotErr := c.CaptureScreenshotInCwd()
		if strings.TrimSpace(screenshotOutput) != "" {
			logs = append(logs, screenshotOutput)
		}
		if screenshotErr != nil {
			return nil, strings.Join(logs, "\n"), screenshotErr
		}
		result.EnteredLogin = true
		result.ScreenshotPath = screenshotPath
		result.ScanCompleted = false
		return result, strings.Join(logs, "\n"), nil
	}

	clickOutput, err := c.ClickLoginByRef(loginRef)
	if strings.TrimSpace(clickOutput) != "" {
		logs = append(logs, clickOutput)
	}
	if err != nil {
		return nil, strings.Join(logs, "\n"), err
	}
	result.EnteredLogin = true

	readyOutput, err := c.WaitWindowReady()
	if strings.TrimSpace(readyOutput) != "" {
		logs = append(logs, readyOutput)
	}
	if err != nil {
		return nil, strings.Join(logs, "\n"), err
	}

	screenshotPath, screenshotOutput, err := c.CaptureScreenshotInCwd()
	if strings.TrimSpace(screenshotOutput) != "" {
		logs = append(logs, screenshotOutput)
	}
	if err != nil {
		return nil, strings.Join(logs, "\n"), err
	}
	result.ScreenshotPath = screenshotPath
	result.ScanCompleted = false

	return result, strings.Join(logs, "\n"), nil
}

func (c *AliyunCodingLoginClient) SnapshotInteractive() (string, error) {
	output, err := c.runSessionCommand("snapshot", "-i")
	if err != nil {
		return output, fmt.Errorf("failed to snapshot aliyun home page: %w", err)
	}
	return output, nil
}

func (c *AliyunCodingLoginClient) shouldPrintSnapshotOutput() bool {
	if c.session == nil {
		return false
	}
	return c.session.config.DevMode
}

func (c *AliyunCodingLoginClient) SnapshotAll() (string, error) {
	output, err := c.runSessionCommand("snapshot")
	if err != nil {
		return output, fmt.Errorf("failed to snapshot aliyun page: %w", err)
	}
	return output, nil
}

func (c *AliyunCodingLoginClient) OpenCodingPlanDetail() (string, error) {
	logs := make([]string, 0, 2)
	output, err := c.runSessionCommand("open", codingPlanDetailURL)
	if strings.TrimSpace(output) != "" {
		logs = append(logs, output)
	}
	if err != nil {
		return strings.Join(logs, "\n"), fmt.Errorf("failed to open coding plan detail page: %w", err)
	}
	waitOutput, waitErr := c.WaitWindowReady()
	if strings.TrimSpace(waitOutput) != "" {
		logs = append(logs, waitOutput)
	}
	if waitErr != nil {
		return strings.Join(logs, "\n"), waitErr
	}
	return strings.Join(logs, "\n"), nil
}

func (c *AliyunCodingLoginClient) OpenDirectLoginPage() (string, error) {
	logs := make([]string, 0, 2)
	output, err := c.runSessionCommand("open", aliyunDirectLoginURL)
	if strings.TrimSpace(output) != "" {
		logs = append(logs, output)
	}
	if err != nil {
		return strings.Join(logs, "\n"), fmt.Errorf("failed to open direct login page: %w", err)
	}
	waitOutput, waitErr := c.WaitWindowReady()
	if strings.TrimSpace(waitOutput) != "" {
		logs = append(logs, waitOutput)
	}
	if waitErr != nil {
		return strings.Join(logs, "\n"), waitErr
	}
	return strings.Join(logs, "\n"), nil
}

func (c *AliyunCodingLoginClient) FillUsageFromSnapshot(result *AliyunCodingLoginResult, snapshotOutput string) {
	result.Hours5ResetTime, result.Hours5Usage = extractUsageByMarker(snapshotOutput, "近一周用量")
	result.WeekResetTime, result.WeekUsage = extractUsageByMarker(snapshotOutput, "近一月用量")
	result.MonthResetTime, result.MonthUsage = extractUsageByMarker(snapshotOutput, "套餐专属API Key")
	if result.MonthUsage == "" {
		result.MonthResetTime, result.MonthUsage = extractUsageByMarker(snapshotOutput, "套餐专属API")
	}
	if result.Hours5Usage == "" || result.WeekUsage == "" || result.MonthUsage == "" {
		sequence := extractUsageBySequence(snapshotOutput)
		if len(sequence) >= 3 {
			if result.Hours5Usage == "" {
				result.Hours5ResetTime, result.Hours5Usage = sequence[0][0], sequence[0][1]
			}
			if result.WeekUsage == "" {
				result.WeekResetTime, result.WeekUsage = sequence[1][0], sequence[1][1]
			}
			if result.MonthUsage == "" {
				result.MonthResetTime, result.MonthUsage = sequence[2][0], sequence[2][1]
			}
		}
	}
}

func (c *AliyunCodingLoginClient) IsLoggedInFromSnapshot(snapshotOutput string) bool {
	content := strings.ToLower(snapshotOutput)
	notLoggedInMarkers := []string{
		"登录阿里云",
		"立即登录",
		"快捷注册",
		"sign in",
		"log in",
	}
	for _, marker := range notLoggedInMarkers {
		if strings.Contains(content, strings.ToLower(marker)) {
			return false
		}
	}

	loggedInMarkers := []string{
		"退出登录",
		"sign out",
		"账号中心",
		"accesskey",
	}
	for _, marker := range loggedInMarkers {
		if strings.Contains(content, strings.ToLower(marker)) {
			return true
		}
	}

	return false
}

func (c *AliyunCodingLoginClient) FindLoginRef(snapshotOutput string) (string, error) {
	loginRefs := extractLoginRefCandidates(snapshotOutput)
	if len(loginRefs) == 0 {
		return "", fmt.Errorf("failed to find login ref from snapshot")
	}
	return "@" + loginRefs[0], nil
}

func (c *AliyunCodingLoginClient) ClickLoginByRef(loginRef string) (string, error) {
	loginRef = strings.TrimSpace(loginRef)
	if loginRef == "" {
		return "", fmt.Errorf("login ref is required")
	}
	if !strings.HasPrefix(loginRef, "@") {
		loginRef = "@" + loginRef
	}

	logs := make([]string, 0, 2)
	scrollOutput, scrollErr := c.runSessionCommand("scrollintoview", loginRef)
	if strings.TrimSpace(scrollOutput) != "" {
		logs = append(logs, scrollOutput)
	}
	if scrollErr != nil {
		clickOutput, clickErr := c.runSessionCommand("click", loginRef)
		if strings.TrimSpace(clickOutput) != "" {
			logs = append(logs, clickOutput)
		}
		if clickErr != nil {
			return strings.Join(logs, "\n"), fmt.Errorf("failed to click login ref %s: %w", loginRef, clickErr)
		}
		return strings.Join(logs, "\n"), nil
	}

	clickOutput, clickErr := c.runSessionCommand("click", loginRef)
	if strings.TrimSpace(clickOutput) != "" {
		logs = append(logs, clickOutput)
	}
	if clickErr != nil {
		return strings.Join(logs, "\n"), fmt.Errorf("failed to click login ref %s: %w", loginRef, clickErr)
	}
	return strings.Join(logs, "\n"), nil
}

func (c *AliyunCodingLoginClient) WaitWindowReady() (string, error) {
	output, err := c.runSessionCommand("wait", "10000")
	if err != nil {
		return output, fmt.Errorf("failed to wait 10s: %w", err)
	}
	return output, nil
}

func (c *AliyunCodingLoginClient) CaptureScreenshotInCwd() (string, string, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return "", "", fmt.Errorf("failed to get working directory: %w", err)
	}

	fileName := "aliyu-login.png"
	filePath := filepath.Join(workingDir, fileName)

	output, err := c.runSessionCommand("screenshot", filePath)
	if err != nil {
		return "", output, fmt.Errorf("failed to save login screenshot: %w", err)
	}
	return filePath, output, nil
}

func (c *AliyunCodingLoginClient) detectLoggedIn() (bool, string, error) {
	output, err := c.SnapshotInteractive()
	if err != nil {
		return false, output, err
	}
	content := strings.ToLower(output)
	notLoggedInMarkers := []string{
		"登录阿里云",
		"立即登录",
		"快捷注册",
		"sign in",
		"log in",
	}
	for _, marker := range notLoggedInMarkers {
		if strings.Contains(content, strings.ToLower(marker)) {
			return false, output, nil
		}
	}

	loggedInMarkers := []string{
		"退出登录",
		"sign out",
		"账号中心",
		"accesskey",
	}
	for _, marker := range loggedInMarkers {
		if strings.Contains(content, strings.ToLower(marker)) {
			return true, output, nil
		}
	}

	if strings.Contains(content, "登录") {
		return false, output, nil
	}
	return false, output, nil
}

func (c *AliyunCodingLoginClient) clickLoginEntry() (string, error) {
	snapshotOutput, err := c.SnapshotInteractive()
	if err != nil {
		return snapshotOutput, err
	}
	loginRef, err := c.FindLoginRef(snapshotOutput)
	if err != nil {
		return snapshotOutput, err
	}
	clickOutput, clickErr := c.ClickLoginByRef(loginRef)
	if clickErr != nil {
		if strings.TrimSpace(snapshotOutput) == "" {
			return clickOutput, clickErr
		}
		if strings.TrimSpace(clickOutput) == "" {
			return snapshotOutput, clickErr
		}
		return snapshotOutput + "\n" + clickOutput, clickErr
	}
	if strings.TrimSpace(snapshotOutput) == "" {
		return clickOutput, nil
	}
	if strings.TrimSpace(clickOutput) == "" {
		return snapshotOutput, nil
	}
	return snapshotOutput + "\n" + clickOutput, nil
}

func extractLoginRefCandidates(snapshot string) []string {
	lines := strings.Split(snapshot, "\n")
	priority0 := make([]string, 0)
	priority1 := make([]string, 0)
	priority2 := make([]string, 0)
	seen := make(map[string]struct{})

	for _, line := range lines {
		if !strings.Contains(line, "[ref=") || !strings.Contains(line, "登录") {
			continue
		}
		if strings.Contains(line, "合作伙伴") || strings.Contains(line, "管理后台") {
			continue
		}
		matches := loginRefRegex.FindStringSubmatch(line)
		if len(matches) < 2 {
			continue
		}
		ref := matches[1]
		if _, ok := seen[ref]; ok {
			continue
		}
		seen[ref] = struct{}{}

		if strings.Contains(line, "\"登录\"") {
			priority0 = append(priority0, ref)
			continue
		}
		if strings.Contains(line, "\"登录阿里云\"") {
			priority1 = append(priority1, ref)
			continue
		}
		priority2 = append(priority2, ref)
	}

	candidates := make([]string, 0, len(priority0)+len(priority1)+len(priority2))
	candidates = append(candidates, priority0...)
	candidates = append(candidates, priority1...)
	candidates = append(candidates, priority2...)
	return candidates
}

func (c *AliyunCodingLoginClient) runSessionCommand(args ...string) (string, error) {
	commandArgs := append([]string{"--session-name", c.session.config.SessionName}, args...)
	return c.session.run(commandArgs)
}

func extractUsageByMarker(snapshotOutput string, marker string) (string, string) {
	linePattern := regexp.MustCompile(`(\d{4}-\d{2}-\d{2}\s\d{2}:\d{2}:\d{2})\s*重置\s*([0-9]+%)`)
	lines := strings.Split(snapshotOutput, "\n")
	for _, line := range lines {
		if !strings.Contains(line, marker) {
			continue
		}
		matches := linePattern.FindStringSubmatch(line)
		if len(matches) >= 3 {
			return matches[1], matches[2]
		}
	}
	return "", ""
}

func extractUsageBySequence(snapshotOutput string) [][2]string {
	pattern := regexp.MustCompile(`(\d{4}-\d{2}-\d{2}\s\d{2}:\d{2}:\d{2})\s*重置\s*([0-9]+%)`)
	matches := pattern.FindAllStringSubmatch(snapshotOutput, -1)
	result := make([][2]string, 0, 3)
	for _, match := range matches {
		if len(match) < 3 {
			continue
		}
		result = append(result, [2]string{match[1], match[2]})
		if len(result) == 3 {
			break
		}
	}
	return result
}

package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	defaultAgentBrowserPath = "agent-browser"
	defaultSessionName      = "jeayapp"

	envKeyPath        = "AGENT_BROWSER_PATH"
	envKeyDevMode     = "AGENT_BROWSER_DEV_MODE"
	envKeySessionName = "AGENT_BROWSER_SESSION_NAME"
)

type AgentBrowserConfig struct {
	Path        string
	SessionName string
	DevMode     bool
}

type AgentBrowserSession struct {
	config    AgentBrowserConfig
	targetURL string
	opened    bool
}

func parseEnvFile(envFilePath string) map[string]string {
	values := make(map[string]string)
	file, err := os.Open(envFilePath)
	if err != nil {
		return values
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		value = strings.Trim(value, `"'`)
		if key == "" || value == "" {
			continue
		}
		values[key] = value
	}

	return values
}

func envFileCandidates() []string {
	candidates := make([]string, 0, 2)
	seen := make(map[string]struct{})

	workingDir, err := os.Getwd()
	if err == nil {
		workingDirEnvPath := filepath.Join(workingDir, ".env")
		candidates = append(candidates, workingDirEnvPath)
		seen[workingDirEnvPath] = struct{}{}
	}

	executablePath, err := os.Executable()
	if err == nil {
		executableDir := filepath.Dir(executablePath)
		executableEnvPath := filepath.Join(executableDir, ".env")
		if _, ok := seen[executableEnvPath]; !ok {
			candidates = append(candidates, executableEnvPath)
		}
	}

	return candidates
}

func resolveEnvValue(key string) string {
	for _, envPath := range envFileCandidates() {
		values := parseEnvFile(envPath)
		if value, ok := values[key]; ok {
			return value
		}
	}

	return ""
}

func isDevMode(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "1", "true", "yes", "y", "on", "是":
		return true
	default:
		return false
	}
}

func nonEmptyOrDefault(value string, defaultValue string) string {
	if strings.TrimSpace(value) == "" {
		return defaultValue
	}
	return value
}

func loadAgentBrowserConfig() AgentBrowserConfig {
	return AgentBrowserConfig{
		Path:        nonEmptyOrDefault(resolveEnvValue(envKeyPath), defaultAgentBrowserPath),
		SessionName: nonEmptyOrDefault(resolveEnvValue(envKeySessionName), defaultSessionName),
		DevMode:     isDevMode(resolveEnvValue(envKeyDevMode)),
	}
}

func buildOpenArgs(sessionName string, devMode bool, targetURL string) []string {
	args := []string{"--session-name", sessionName}
	if devMode {
		args = append(args, "--headed")
	}
	args = append(args, "open", targetURL)
	return args
}

func buildCloseArgs(sessionName string) []string {
	return []string{"--session-name", sessionName, "close"}
}

func NewAgentBrowserSession(targetURL string) (*AgentBrowserSession, string, error) {
	targetURL = strings.TrimSpace(targetURL)
	if targetURL == "" {
		return nil, "", errors.New("target url is required")
	}

	session := &AgentBrowserSession{
		config:    loadAgentBrowserConfig(),
		targetURL: targetURL,
	}
	output, err := session.Open()
	if err != nil {
		return nil, output, err
	}
	return session, output, nil
}

func (s *AgentBrowserSession) run(args []string) (string, error) {
	cmd := exec.Command(s.config.Path, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("command failed: %w", err)
	}
	return string(output), nil
}

func (s *AgentBrowserSession) Open() (string, error) {
	output, err := s.run(buildOpenArgs(s.config.SessionName, s.config.DevMode, s.targetURL))
	if err != nil {
		return output, err
	}
	s.opened = true
	return output, nil
}

func (s *AgentBrowserSession) Close() (string, error) {
	if !s.opened {
		return "", nil
	}
	output, err := s.run(buildCloseArgs(s.config.SessionName))
	if err != nil {
		return output, err
	}
	s.opened = false
	return output, nil
}

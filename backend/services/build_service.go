package services

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"
)

type BuildResult struct {
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	Output   string `json:"output"`
	ExePath  string `json:"exePath"`
	Progress int    `json:"progress"`
}

type BuildConfig struct {
	AppName     string `json:"appName"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Company     string `json:"company"`
	Copyright   string `json:"copyright"`
	IconPath    string `json:"iconPath"`
	OutputPath  string `json:"outputPath"`
	TargetOS    string `json:"targetOS"`
	TargetArch  string `json:"targetArch"`
	Optimize    bool   `json:"optimize"`
	StripDebug  bool   `json:"stripDebug"`
	UPXCompress bool   `json:"upxCompress"`
	EmbedAssets bool   `json:"embedAssets"`
}

type BuildProgress struct {
	Stage    string `json:"stage"`
	Progress int    `json:"progress"`
	Message  string `json:"message"`
}

type LogCallback func(level string, message string)

type BuildService struct {
	mu           sync.Mutex
	logMu        sync.Mutex
	cmd          *exec.Cmd
	cancel       context.CancelFunc
	isRunning    bool
	building     bool
	outputDir    string
	progressChan chan BuildProgress
	logCallback  LogCallback
}

func NewBuildService() *BuildService {
	return &BuildService{
		outputDir:    filepath.Join(os.TempDir(), "yigou_build"),
		progressChan: make(chan BuildProgress, 100),
	}
}

func (s *BuildService) SetLogCallback(cb LogCallback) {
	s.logMu.Lock()
	defer s.logMu.Unlock()
	s.logCallback = cb
}

func (s *BuildService) log(level string, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	s.logMu.Lock()
	cb := s.logCallback
	s.logMu.Unlock()
	if cb != nil {
		cb(level, msg)
	}
}

func (s *BuildService) sendProgress(stage string, progress int, message string) {
	s.log("info", "[%s] %d%% - %s", stage, progress, message)
	select {
	case s.progressChan <- BuildProgress{Stage: stage, Progress: progress, Message: message}:
	default:
	}
}

func (s *BuildService) GetProgressChan() <-chan BuildProgress {
	return s.progressChan
}

var nonAsciiRegex = regexp.MustCompile(`[^a-zA-Z0-9_.-]`)

func sanitizeModuleName(name string) string {
	name = nonAsciiRegex.ReplaceAllString(name, "_")
	name = strings.Trim(name, "_.-")
	if name == "" {
		name = "yigou_app"
	}
	if name[0] >= '0' && name[0] <= '9' {
		name = "app_" + name
	}
	return strings.ToLower(name)
}

var fullWidthChars = map[rune]rune{
	'\uFF01': '!', '\uFF02': '"', '\uFF03': '#', '\uFF04': '$', '\uFF05': '%',
	'\uFF06': '&', '\uFF07': '\'', '\uFF08': '(', '\uFF09': ')', '\uFF0A': '*',
	'\uFF0B': '+', '\uFF0C': ',', '\uFF0D': '-', '\uFF0E': '.', '\uFF0F': '/',
	'\uFF10': '0', '\uFF11': '1', '\uFF12': '2', '\uFF13': '3', '\uFF14': '4',
	'\uFF15': '5', '\uFF16': '6', '\uFF17': '7', '\uFF18': '8', '\uFF19': '9',
	'\uFF1A': ':', '\uFF1B': ';', '\uFF1C': '<', '\uFF1D': '=', '\uFF1E': '>',
	'\uFF1F': '?', '\uFF20': '@', '\uFF21': 'A', '\uFF22': 'B', '\uFF23': 'C',
	'\uFF24': 'D', '\uFF25': 'E', '\uFF26': 'F', '\uFF27': 'G', '\uFF28': 'H',
	'\uFF29': 'I', '\uFF2A': 'J', '\uFF2B': 'K', '\uFF2C': 'L', '\uFF2D': 'M',
	'\uFF2E': 'N', '\uFF2F': 'O', '\uFF30': 'P', '\uFF31': 'Q', '\uFF32': 'R',
	'\uFF33': 'S', '\uFF34': 'T', '\uFF35': 'U', '\uFF36': 'V', '\uFF37': 'W',
	'\uFF38': 'X', '\uFF39': 'Y', '\uFF3A': 'Z', '\uFF3B': '[', '\uFF3C': '\\',
	'\uFF3D': ']', '\uFF3E': '^', '\uFF3F': '_', '\uFF40': '`', '\uFF41': 'a',
	'\uFF42': 'b', '\uFF43': 'c', '\uFF44': 'd', '\uFF45': 'e', '\uFF46': 'f',
	'\uFF47': 'g', '\uFF48': 'h', '\uFF49': 'i', '\uFF4A': 'j', '\uFF4B': 'k',
	'\uFF4C': 'l', '\uFF4D': 'm', '\uFF4E': 'n', '\uFF4F': 'o', '\uFF50': 'p',
	'\uFF51': 'q', '\uFF52': 'r', '\uFF53': 's', '\uFF54': 't', '\uFF55': 'u',
	'\uFF56': 'v', '\uFF57': 'w', '\uFF58': 'x', '\uFF59': 'y', '\uFF5A': 'z',
	'\uFF5B': '{', '\uFF5C': '|', '\uFF5D': '}', '\uFF5E': '~', '\u3000': ' ',
}

func convertFullWidthToHalfWidth(s string) string {
	var result strings.Builder
	for _, r := range s {
		if replacement, ok := fullWidthChars[r]; ok {
			result.WriteRune(replacement)
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

func cleanBuildDir(buildDir string) {
	os.Remove(filepath.Join(buildDir, "go.mod"))
	os.Remove(filepath.Join(buildDir, "go.sum"))
	os.Remove(filepath.Join(buildDir, "main.go"))
}

func hideWindowCmd(name string, arg ...string) *exec.Cmd {
	cmd := exec.Command(name, arg...)
	if runtime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{
			HideWindow:    true,
			CreationFlags: 0x08000000,
		}
	}
	return cmd
}

func (s *BuildService) CompileAndRun(goCode string, projectName string, projectDir string, debugMode bool) *BuildResult {
	s.mu.Lock()

	for s.building {
		s.mu.Unlock()
		time.Sleep(100 * time.Millisecond)
		s.mu.Lock()
	}

	if s.isRunning {
		s.stopLocked()
		time.Sleep(500 * time.Millisecond)
	}

	s.building = true
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		s.building = false
		s.mu.Unlock()
	}()

	modeLabel := "正常模式"
	if debugMode {
		modeLabel = "调试模式"
	}
	s.log("info", "========== 开始编译运行 (%s) ==========", modeLabel)
	s.log("info", "项目名称: %s", projectName)

	moduleName := sanitizeModuleName(projectName)

	buildDir := filepath.Join(s.outputDir, moduleName)
	if projectDir != "" {
		buildDir = filepath.Join(projectDir, "build")
	}
	if err := os.MkdirAll(buildDir, 0755); err != nil {
		s.log("error", "创建构建目录失败: %v", err)
		return &BuildResult{
			Success: false,
			Message: fmt.Sprintf("创建构建目录失败: %v", err),
		}
	}

	cleanBuildDir(buildDir)
	goCode = convertFullWidthToHalfWidth(goCode)

	s.log("info", "正在写入源代码...")
	mainFile := filepath.Join(buildDir, "main.go")
	if err := os.WriteFile(mainFile, []byte(goCode), 0644); err != nil {
		s.log("error", "写入代码文件失败: %v", err)
		return &BuildResult{
			Success: false,
			Message: fmt.Sprintf("写入代码文件失败: %v", err),
		}
	}

	s.log("info", "正在初始化Go模块: %s", moduleName)
	initCmd := hideWindowCmd("go", "mod", "init", moduleName)
	initCmd.Dir = buildDir
	initOutput, initErr := initCmd.CombinedOutput()
	if initErr != nil {
		s.log("error", "初始化模块失败: %v\n%s", initErr, string(initOutput))
		return &BuildResult{
			Success: false,
			Message: fmt.Sprintf("初始化模块失败: %v", initErr),
			Output:  string(initOutput),
		}
	}

	s.log("info", "正在配置运行时环境...")
	s.setupRuntime(buildDir)
	s.fixGoMod(buildDir, moduleName)

	s.log("info", "正在整理依赖 (go mod tidy)...")
	tidyCmd := hideWindowCmd("go", "mod", "tidy")
	tidyCmd.Dir = buildDir
	tidyOutput, tidyErr := tidyCmd.CombinedOutput()
	if tidyErr != nil {
		s.log("warn", "依赖整理警告: %v\n%s", tidyErr, string(tidyOutput))
	} else {
		s.log("info", "依赖整理完成")
	}

	exePath := filepath.Join(buildDir, moduleName+".exe")
	s.log("info", "正在编译程序...")
	buildArgs := []string{"build", "-o", exePath}
	if !debugMode {
		buildArgs = append(buildArgs, "-ldflags", "-s -w -H windowsgui")
	} else {
		buildArgs = append(buildArgs, "-gcflags", "all=-N -l")
	}
	buildArgs = append(buildArgs, mainFile)
	buildCmd := hideWindowCmd("go", buildArgs...)
	buildCmd.Dir = buildDir
	buildOutput, err := buildCmd.CombinedOutput()

	if err != nil {
		errMsg := string(buildOutput)
		errMsg = s.translateError(errMsg)
		s.log("error", "编译失败:\n%s", errMsg)
		return &BuildResult{
			Success: false,
			Message: fmt.Sprintf("编译失败: %v", err),
			Output:  errMsg + "\n" + string(tidyOutput),
		}
	}

	s.log("info", "编译成功！正在启动程序...")

	_, cancel := context.WithCancel(context.Background())

	runCmd := exec.Command(exePath)
	runCmd.Dir = buildDir

	stdoutPipe, _ := runCmd.StdoutPipe()
	stderrPipe, _ := runCmd.StderrPipe()

	if err := runCmd.Start(); err != nil {
		cancel()
		s.log("error", "启动程序失败: %v", err)
		return &BuildResult{
			Success: false,
			Message: fmt.Sprintf("启动程序失败: %v", err),
			ExePath: exePath,
		}
	}

	s.mu.Lock()
	s.cancel = cancel
	s.cmd = runCmd
	s.isRunning = true
	s.mu.Unlock()

	var outputBuilder strings.Builder
	var outputMu sync.Mutex

	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := stdoutPipe.Read(buf)
			if n > 0 {
				data := string(buf[:n])
				outputMu.Lock()
				outputBuilder.WriteString(data)
				outputMu.Unlock()
				s.log("output", "%s", strings.TrimRight(data, "\r\n"))
			}
			if err != nil {
				break
			}
		}
	}()

	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := stderrPipe.Read(buf)
			if n > 0 {
				data := string(buf[:n])
				outputMu.Lock()
				outputBuilder.WriteString(data)
				outputMu.Unlock()
				s.log("warn", "%s", strings.TrimRight(data, "\r\n"))
			}
			if err != nil {
				break
			}
		}
	}()

	go func() {
		runCmd.Wait()
		s.mu.Lock()
		s.isRunning = false
		s.cmd = nil
		s.mu.Unlock()
		s.log("info", "程序已退出")
	}()

	time.Sleep(500 * time.Millisecond)

	outputMu.Lock()
	collected := outputBuilder.String()
	outputMu.Unlock()

	return &BuildResult{
		Success: true,
		Message: "编译运行成功",
		Output:  collected,
		ExePath: exePath,
	}
}

func (s *BuildService) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.stopLocked()
}

func (s *BuildService) stopLocked() {
	s.building = false
	if s.cancel != nil {
		s.cancel()
	}
	if s.cmd != nil && s.cmd.Process != nil {
		s.cmd.Process.Kill()
		s.log("info", "程序已被终止")
	}
	s.isRunning = false
	s.cmd = nil
	s.cancel = nil
}

func (s *BuildService) IsRunning() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.isRunning
}

func (s *BuildService) BuildExe(goCode string, projectName string, outputPath string) *BuildResult {
	s.mu.Lock()
	for s.building {
		s.mu.Unlock()
		time.Sleep(100 * time.Millisecond)
		s.mu.Lock()
	}
	s.building = true
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		s.building = false
		s.mu.Unlock()
	}()

	s.log("info", "========== 开始打包构建 ==========")
	s.log("info", "项目名称: %s", projectName)
	s.log("info", "输出路径: %s", outputPath)

	moduleName := sanitizeModuleName(projectName)
	buildDir := filepath.Join(s.outputDir, moduleName+"_build")
	if err := os.MkdirAll(buildDir, 0755); err != nil {
		s.log("error", "创建构建目录失败: %v", err)
		return &BuildResult{
			Success: false,
			Message: fmt.Sprintf("创建构建目录失败: %v", err),
		}
	}

	cleanBuildDir(buildDir)
	goCode = convertFullWidthToHalfWidth(goCode)

	s.log("info", "正在写入源代码...")
	mainFile := filepath.Join(buildDir, "main.go")
	if err := os.WriteFile(mainFile, []byte(goCode), 0644); err != nil {
		s.log("error", "写入代码文件失败: %v", err)
		return &BuildResult{
			Success: false,
			Message: fmt.Sprintf("写入代码文件失败: %v", err),
		}
	}

	s.log("info", "正在初始化Go模块: %s", moduleName)
	initCmd := hideWindowCmd("go", "mod", "init", moduleName)
	initCmd.Dir = buildDir
	initOutput, initErr := initCmd.CombinedOutput()
	if initErr != nil {
		s.log("error", "初始化模块失败: %v\n%s", initErr, string(initOutput))
		return &BuildResult{
			Success: false,
			Message: fmt.Sprintf("初始化模块失败: %v", initErr),
			Output:  string(initOutput),
		}
	}

	s.log("info", "正在配置运行时环境...")
	s.setupRuntime(buildDir)
	s.fixGoMod(buildDir, moduleName)

	s.log("info", "正在整理依赖 (go mod tidy)...")
	tidyCmd := hideWindowCmd("go", "mod", "tidy")
	tidyCmd.Dir = buildDir
	tidyOutput, tidyErr := tidyCmd.CombinedOutput()
	if tidyErr != nil {
		s.log("warn", "依赖整理警告: %v\n%s", tidyErr, string(tidyOutput))
	} else {
		s.log("info", "依赖整理完成")
	}

	if outputPath == "" {
		outputPath = filepath.Join(buildDir, moduleName+".exe")
	}

	outputDir := filepath.Dir(outputPath)
	os.MkdirAll(outputDir, 0755)

	s.log("info", "正在编译程序 (go build -ldflags \"-s -w -H windowsgui\")...")
	buildCmd := hideWindowCmd("go", "build", "-ldflags", "-s -w -H windowsgui", "-o", outputPath, mainFile)
	buildCmd.Dir = buildDir
	buildOutput, err := buildCmd.CombinedOutput()

	if err != nil {
		errMsg := string(buildOutput)
		errMsg = s.translateError(errMsg)
		s.log("error", "编译失败:\n%s", errMsg)
		return &BuildResult{
			Success: false,
			Message: fmt.Sprintf("打包失败: %v", err),
			Output:  errMsg,
		}
	}

	s.log("info", "打包成功！输出: %s", outputPath)
	return &BuildResult{
		Success: true,
		Message: "打包成功",
		Output:  string(buildOutput),
		ExePath: outputPath,
	}
}

func (s *BuildService) translateError(errMsg string) string {
	translations := map[string]string{
		"undefined":             "未定义",
		"cannot use":            "不能使用",
		"declared and not used": "已声明但未使用",
		"imported and not used": "已导入但未使用",
		"syntax error":          "语法错误",
		"undefined:":            "未定义:",
		"not enough arguments":  "参数不足",
		"too many arguments":    "参数过多",
		"cannot convert":        "无法转换",
		"type mismatch":         "类型不匹配",
	}

	result := errMsg
	for eng, chn := range translations {
		result = strings.ReplaceAll(result, eng, chn)
	}

	return result
}

func (s *BuildService) BuildWithConfig(goCode string, config BuildConfig) *BuildResult {
	s.sendProgress("初始化", 5, "正在准备构建环境...")

	if config.AppName == "" {
		config.AppName = "app"
	}
	if config.TargetOS == "" {
		config.TargetOS = runtime.GOOS
	}
	if config.TargetArch == "" {
		config.TargetArch = runtime.GOARCH
	}

	moduleName := sanitizeModuleName(config.AppName)
	buildDir := filepath.Join(s.outputDir, moduleName+"_build_"+time.Now().Format("20060102150405"))
	if err := os.MkdirAll(buildDir, 0755); err != nil {
		return &BuildResult{
			Success: false,
			Message: fmt.Sprintf("创建构建目录失败: %v", err),
		}
	}
	defer os.RemoveAll(buildDir)

	s.sendProgress("准备", 10, "正在写入源代码...")

	mainFile := filepath.Join(buildDir, "main.go")
	if err := os.WriteFile(mainFile, []byte(goCode), 0644); err != nil {
		return &BuildResult{
			Success: false,
			Message: fmt.Sprintf("写入代码文件失败: %v", err),
		}
	}

	s.sendProgress("准备", 15, "正在初始化Go模块...")

	initCmd := hideWindowCmd("go", "mod", "init", moduleName)
	initCmd.Dir = buildDir
	initCmd.Run()

	s.setupRuntime(buildDir)
	s.fixGoMod(buildDir, moduleName)

	s.sendProgress("准备", 20, "正在整理依赖...")

	tidyCmd := hideWindowCmd("go", "mod", "tidy")
	tidyCmd.Dir = buildDir
	tidyCmd.Run()

	s.sendProgress("编译", 30, "正在编译程序...")

	outputPath := config.OutputPath
	if outputPath == "" {
		homeDir, _ := os.UserHomeDir()
		outputDir := filepath.Join(homeDir, "Desktop")
		outputPath = filepath.Join(outputDir, config.AppName+".exe")
	}

	outputDir := filepath.Dir(outputPath)
	os.MkdirAll(outputDir, 0755)

	buildArgs := []string{"build"}
	if config.Optimize {
		buildArgs = append(buildArgs, "-ldflags", "-s -w -H windowsgui")
	}
	if config.StripDebug {
		if len(buildArgs) > 1 && buildArgs[1] == "-ldflags" {
			buildArgs[2] += " -s -H windowsgui"
		} else {
			buildArgs = append(buildArgs, "-ldflags", "-s -H windowsgui")
		}
	}
	if len(buildArgs) <= 1 {
		buildArgs = append(buildArgs, "-ldflags", "-H windowsgui")
	}
	buildArgs = append(buildArgs, "-o", outputPath, mainFile)

	buildCmd := hideWindowCmd("go", buildArgs...)
	buildCmd.Dir = buildDir
	buildCmd.Env = append(os.Environ(),
		"GOOS="+config.TargetOS,
		"GOARCH="+config.TargetArch,
	)

	s.sendProgress("编译", 50, "正在编译...")

	buildOutput, err := buildCmd.CombinedOutput()

	if err != nil {
		errMsg := string(buildOutput)
		errMsg = s.translateError(errMsg)
		return &BuildResult{
			Success: false,
			Message: fmt.Sprintf("编译失败: %v", err),
			Output:  errMsg,
		}
	}

	s.sendProgress("编译", 70, "编译完成，正在处理...")

	if config.UPXCompress {
		s.sendProgress("压缩", 80, "正在使用UPX压缩...")
		upxCmd := hideWindowCmd("upx", "--best", outputPath)
		upxCmd.Run()
	}

	s.sendProgress("完成", 100, "构建完成！")

	fileInfo, _ := os.Stat(outputPath)
	fileSize := int64(0)
	if fileInfo != nil {
		fileSize = fileInfo.Size()
	}

	successMsg := fmt.Sprintf("构建成功！\n目标平台: %s/%s\n输出路径: %s\n文件大小: %.2f MB",
		config.TargetOS, config.TargetArch, outputPath, float64(fileSize)/(1024*1024))

	return &BuildResult{
		Success: true,
		Message: successMsg,
		Output:  string(buildOutput),
		ExePath: outputPath,
	}
}

func (s *BuildService) GetSupportedPlatforms() []map[string]string {
	return []map[string]string{
		{"os": "windows", "arch": "amd64", "label": "Windows 64位"},
		{"os": "windows", "arch": "386", "label": "Windows 32位"},
		{"os": "windows", "arch": "arm64", "label": "Windows ARM64"},
		{"os": "linux", "arch": "amd64", "label": "Linux 64位"},
		{"os": "linux", "arch": "386", "label": "Linux 32位"},
		{"os": "linux", "arch": "arm64", "label": "Linux ARM64"},
		{"os": "darwin", "arch": "amd64", "label": "macOS Intel"},
		{"os": "darwin", "arch": "arm64", "label": "macOS Apple Silicon"},
	}
}

func (s *BuildService) Cleanup() {
	s.Stop()
	os.RemoveAll(s.outputDir)
}

func (s *BuildService) setupRuntime(buildDir string) error {
	yigouDir := filepath.Join(buildDir, "yigou")
	if err := os.MkdirAll(yigouDir, 0755); err != nil {
		return fmt.Errorf("创建yigou目录失败: %v", err)
	}

	runtimeFile := filepath.Join(yigouDir, "runtime.go")
	if err := os.WriteFile(runtimeFile, []byte(yigouRuntimeSource), 0644); err != nil {
		return fmt.Errorf("写入运行时文件失败: %v", err)
	}

	goModFile := filepath.Join(yigouDir, "go.mod")
	goModContent := "module yigou\n\ngo 1.25.4\n"
	if err := os.WriteFile(goModFile, []byte(goModContent), 0644); err != nil {
		return fmt.Errorf("写入yigou go.mod失败: %v", err)
	}

	return nil
}

func (s *BuildService) fixGoMod(buildDir string, moduleName string) error {
	modFile := filepath.Join(buildDir, "go.mod")
	content, err := os.ReadFile(modFile)
	if err != nil {
		return err
	}

	modContent := string(content)
	if !strings.Contains(modContent, "replace yigou") {
		modContent += "\nreplace yigou => ./yigou\n"
	}

	if !strings.Contains(modContent, "require yigou") {
		modContent = strings.Replace(modContent,
			"module "+moduleName,
			"module "+moduleName+"\n\nrequire yigou v0.0.0",
			1)
	}

	if !strings.Contains(modContent, "github.com/jchv/go-webview2") {
		modContent += "\nrequire github.com/jchv/go-webview2 v0.0.0-20250406165304-0bcfea011047\n"
	}
	if !strings.Contains(modContent, "golang.org/x/sys") {
		modContent += "\nrequire golang.org/x/sys v0.30.0\n"
	}

	return os.WriteFile(modFile, []byte(modContent), 0644)
}

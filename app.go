package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"yigou-ide/backend/core"
	"yigou-ide/backend/models"
	"yigou-ide/backend/services"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx            context.Context
	fileService    *services.FileService
	codeConverter  *core.CodeConverter
	buildService   *services.BuildService
	projectService *services.ProjectService
	aiService      *services.AIService
}

func NewApp() *App {
	return &App{
		fileService:    services.NewFileService(),
		codeConverter:  core.NewCodeConverter(),
		buildService:   services.NewBuildService(),
		projectService: services.NewProjectService(),
		aiService:      services.NewAIService(),
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.buildService.SetLogCallback(func(level string, message string) {
		runtime.EventsEmit(ctx, "buildLog", map[string]string{
			"level":   level,
			"message": message,
		})
	})
}

func (a *App) shutdown(ctx context.Context) {
	a.buildService.Cleanup()
	a.aiService.Shutdown()
}

func (a *App) Greet(name string) string {
	return "欢迎使用易狗 IDE， " + name + "！"
}

func (a *App) SaveFile(path string, data string) error {
	return a.fileService.SaveEgoFile(path, data)
}

func (a *App) LoadFile(path string) (string, error) {
	return a.fileService.LoadEgoFile(path)
}

func (a *App) ExportGoCode(egoData string) (string, error) {
	return a.fileService.ExportToGoCode(egoData)
}

func (a *App) OpenDialog() (string, error) {
	file, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "打开项目文件",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "易狗项目文件 (*.ego)",
				Pattern:     "*.ego",
			},
		},
	})
	if err != nil {
		return "", err
	}
	return file, nil
}

func (a *App) SaveDialog() (string, error) {
	file, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           "保存项目文件",
		DefaultFilename: "我的项目.ego",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "易狗项目文件 (*.ego)",
				Pattern:     "*.ego",
			},
		},
	})
	if err != nil {
		return "", err
	}
	return file, nil
}

func (a *App) ConvertChineseToGo(chineseCode string) (string, error) {
	return a.codeConverter.ConvertChineseToGo(chineseCode)
}

func (a *App) ConvertCodeRowsToGo(codeRows []core.CodeRowData) string {
	return a.codeConverter.ConvertCodeRowsToGo(codeRows)
}

func (a *App) GetFunctionSuggestions(prefix string) []core.ChineseFunction {
	return a.codeConverter.GetFunctionSuggestions(prefix)
}

func (a *App) GetAllFunctions() []core.ChineseFunction {
	return a.codeConverter.GetAllFunctions()
}

func (a *App) ValidateCode(chineseCode string) []string {
	return a.codeConverter.ValidateCode(chineseCode)
}

func (a *App) GetCategories() []string {
	return a.codeConverter.GetCategories()
}

func (a *App) GetTemplates() []core.CodeTemplate {
	return a.codeConverter.GetTemplates()
}

func (a *App) GetFunctionsByCategory(category string) []core.ChineseFunction {
	return a.codeConverter.GetFunctionsByCategory(category)
}

func (a *App) CompileAndRun(goCode string, projectName string, projectDir string, debugMode bool) *services.BuildResult {
	return a.buildService.CompileAndRun(goCode, projectName, projectDir, debugMode)
}

func (a *App) StopProgram() {
	a.buildService.Stop()
}

func (a *App) IsProgramRunning() bool {
	return a.buildService.IsRunning()
}

func (a *App) BuildExe(goCode string, projectName string, outputPath string) *services.BuildResult {
	return a.buildService.BuildExe(goCode, projectName, outputPath)
}

func (a *App) BuildWithConfig(goCode string, config services.BuildConfig) *services.BuildResult {
	return a.buildService.BuildWithConfig(goCode, config)
}

func (a *App) GetSupportedPlatforms() []map[string]string {
	return a.buildService.GetSupportedPlatforms()
}

func (a *App) SelectSavePath() (string, error) {
	file, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           "选择打包输出路径",
		DefaultFilename: "程序.exe",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "可执行文件 (*.exe)",
				Pattern:     "*.exe",
			},
		},
	})
	if err != nil {
		return "", err
	}
	return file, nil
}

func (a *App) CreateProject(name, description, projectType, author, path string) *models.Project {
	project, err := a.projectService.CreateProject(name, description, projectType, author, path)
	if err != nil {
		runtime.LogError(a.ctx, "创建项目失败: "+err.Error())
		return nil
	}
	return project
}

func (a *App) OpenProject(id int64) *models.Project {
	project, err := a.projectService.OpenProject(id)
	if err != nil {
		runtime.LogError(a.ctx, "打开项目失败: "+err.Error())
		return nil
	}
	return project
}

func (a *App) DeleteProject(id int64) error {
	return a.projectService.DeleteProject(id)
}

func (a *App) GetAllProjects() []*models.Project {
	return a.projectService.GetAllProjects()
}

func (a *App) GetRecentProjects(limit int) []*models.Project {
	return a.projectService.GetRecentProjects(limit)
}

func (a *App) GetFavoriteProjects() []*models.Project {
	return a.projectService.GetFavoriteProjects()
}

func (a *App) ToggleFavorite(id int64) bool {
	result, _ := a.projectService.ToggleFavorite(id)
	return result
}

func (a *App) UpdateProject(id int64, name, description string) error {
	return a.projectService.UpdateProject(id, name, description)
}

func (a *App) GetProjectTemplates() []models.ProjectTemplate {
	return a.projectService.GetTemplates()
}

func (a *App) SearchProjects(query string) []*models.Project {
	return a.projectService.SearchProjects(query)
}

func (a *App) SelectFolder() (string, error) {
	folder, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "选择项目保存目录",
	})
	if err != nil {
		return "", err
	}
	return folder, nil
}

func (a *App) AIChat(messages []services.ChatMessage) *services.ChatResponse {
	resp, err := a.aiService.Chat(messages)
	if err != nil {
		return &services.ChatResponse{
			Error: err.Error(),
			Done:  true,
		}
	}
	return resp
}

func (a *App) AIComplete(prefix string, suffix string, language string, maxTokens int) *services.CompletionResponse {
	req := services.CompletionRequest{
		Prefix:    prefix,
		Suffix:    suffix,
		Language:  language,
		MaxTokens: maxTokens,
	}
	resp, err := a.aiService.Complete(req)
	if err != nil {
		return &services.CompletionResponse{
			Error: err.Error(),
		}
	}
	return resp
}

func (a *App) AIUpdateConfig(config services.AIConfig) {
	a.aiService.UpdateConfig(config)
}

func (a *App) AIGetConfig() services.AIConfig {
	return a.aiService.GetConfig()
}

func (a *App) DownloadModelFile(modelId string, downloadUrl string) error {
	modelsDir := filepath.Join(filepath.Dir(os.Args[0]), ".models")
	if err := os.MkdirAll(modelsDir, 0755); err != nil {
		return fmt.Errorf("创建模型目录失败: %w", err)
	}

	fileName := modelId + ".gguf"
	savePath := filepath.Join(modelsDir, fileName)

	if _, err := os.Stat(savePath); err == nil {
		runtime.EventsEmit(a.ctx, "modelDownloadProgress", map[string]interface{}{
			"modelId":  modelId,
			"progress": 100,
			"status":   "installed",
		})
		return nil
	}

	runtime.EventsEmit(a.ctx, "modelDownloadProgress", map[string]interface{}{
		"modelId":  modelId,
		"progress": 0,
		"status":   "downloading",
	})

	req, err := http.NewRequestWithContext(a.ctx, "GET", downloadUrl, nil)
	if err != nil {
		return fmt.Errorf("创建下载请求失败: %w", err)
	}
	req.Header.Set("User-Agent", "yigou-ide/1.0")

	client := &http.Client{
		Timeout: 30 * time.Minute,
	}
	resp, err := client.Do(req)
	if err != nil {
		runtime.EventsEmit(a.ctx, "modelDownloadProgress", map[string]interface{}{
			"modelId":  modelId,
			"progress": 0,
			"status":   "error",
		})
		return fmt.Errorf("下载失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		runtime.EventsEmit(a.ctx, "modelDownloadProgress", map[string]interface{}{
			"modelId":  modelId,
			"progress": 0,
			"status":   "error",
		})
		return fmt.Errorf("下载失败，HTTP状态码: %d", resp.StatusCode)
	}

	tmpFile := savePath + ".tmp"
	out, err := os.Create(tmpFile)
	if err != nil {
		return fmt.Errorf("创建临时文件失败: %w", err)
	}
	defer out.Close()

	contentLength := resp.ContentLength
	buf := make([]byte, 32*1024)
	var downloaded int64

	for {
		nr, readErr := resp.Body.Read(buf)
		if nr > 0 {
			nw, writeErr := out.Write(buf[:nr])
			if writeErr != nil {
				os.Remove(tmpFile)
				return fmt.Errorf("写入文件失败: %w", writeErr)
			}
			if nw != nr {
				os.Remove(tmpFile)
				return fmt.Errorf("写入不完整: 期望 %d 实际 %d", nr, nw)
			}
			downloaded += int64(nw)

			if contentLength > 0 {
				progress := int(downloaded * 100 / contentLength)
				runtime.EventsEmit(a.ctx, "modelDownloadProgress", map[string]interface{}{
					"modelId":  modelId,
					"progress": progress,
					"status":   "downloading",
				})
			}
		}
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			os.Remove(tmpFile)
			runtime.EventsEmit(a.ctx, "modelDownloadProgress", map[string]interface{}{
				"modelId":  modelId,
				"progress": 0,
				"status":   "error",
			})
			return fmt.Errorf("读取下载数据失败: %w", readErr)
		}
	}

	if err := os.Rename(tmpFile, savePath); err != nil {
		os.Remove(tmpFile)
		return fmt.Errorf("移动文件失败: %w", err)
	}

	runtime.EventsEmit(a.ctx, "modelDownloadProgress", map[string]interface{}{
		"modelId":  modelId,
		"progress": 100,
		"status":   "installed",
	})

	return nil
}

func (a *App) GetModelDownloadStatus(modelId string) map[string]interface{} {
	modelsDir := filepath.Join(filepath.Dir(os.Args[0]), ".models")
	savePath := filepath.Join(modelsDir, modelId+".gguf")

	if _, err := os.Stat(savePath); err == nil {
		return map[string]interface{}{
			"modelId": modelId,
			"status":  "installed",
			"path":    savePath,
		}
	}

	return map[string]interface{}{
		"modelId": modelId,
		"status":  "not_installed",
		"path":    savePath,
	}
}

func (a *App) GetModelsDir() string {
	modelsDir := filepath.Join(filepath.Dir(os.Args[0]), ".models")
	os.MkdirAll(modelsDir, 0755)
	return modelsDir
}

func (a *App) ReadEgoMemory(projectPath string) (string, error) {
	return a.fileService.ReadEgoMemory(projectPath)
}

func (a *App) WriteEgoMemory(projectPath string, content string) error {
	return a.fileService.WriteEgoMemory(projectPath, content)
}

func (a *App) InitEgoMemory(projectPath string, projectName string) error {
	return a.fileService.InitEgoMemory(projectPath, projectName)
}

func (a *App) AppendEgoMemory(projectPath string, section string, content string) error {
	return a.fileService.AppendEgoMemory(projectPath, section, content)
}

func (a *App) GetEgoMemorySummary(projectPath string) string {
	return a.fileService.GetEgoMemorySummary(projectPath)
}

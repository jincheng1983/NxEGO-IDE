package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"yigou-ide/backend/core"
)

type UIComponent struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	X        int    `json:"x"`
	Y        int    `json:"y"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Text     string `json:"text"`
	Visible  bool   `json:"visible"`
	Color    string `json:"color"`
	FontSize int    `json:"fontSize"`
}

type WindowProps struct {
	Title           string `json:"title"`
	Width           int    `json:"width"`
	Height          int    `json:"height"`
	WindowType      string `json:"windowType"`
	BackgroundColor string `json:"backgroundColor"`
	Resizable       bool   `json:"resizable"`
	ShowInTaskbar   bool   `json:"showInTaskbar"`
	ShowIcon        bool   `json:"showIcon"`
	IconPath        string `json:"iconPath"`
	MinWidth        int    `json:"minWidth"`
	MinHeight       int    `json:"minHeight"`
	MaxWidth        int    `json:"maxWidth"`
	MaxHeight       int    `json:"maxHeight"`
	Opacity         int    `json:"opacity"`
	AlwaysOnTop     bool   `json:"alwaysOnTop"`
	BorderStyle     string `json:"borderStyle"`
	ShowMenuBar     bool   `json:"showMenuBar"`
	ShowMinBtn      bool   `json:"showMinBtn"`
	ShowMaxBtn      bool   `json:"showMaxBtn"`
	ShowCloseBtn    bool   `json:"showCloseBtn"`
	TitleBarHeight  int    `json:"titleBarHeight"`
}

type EgoFile struct {
	Version     string        `json:"version"`
	AppName     string        `json:"appName"`
	Components  []UIComponent `json:"components"`
	CodeRows    []CodeRow     `json:"codeRows"`
	CodeContent string        `json:"codeContent"`
	WindowProps WindowProps   `json:"windowProps"`
}

type CodeRow struct {
	Event string `json:"event"`
	Code  string `json:"code"`
}

type FileService struct{}

func NewFileService() *FileService {
	return &FileService{}
}

func (s *FileService) SaveEgoFile(path string, data string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(data), 0644)
}

func (s *FileService) LoadEgoFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

var egoMemoryTemplate = `# 项目记忆文件

> 此文件记录项目的重要信息，AI 助手会参考此文件来理解项目上下文。
> 你可以随时更新此文件，AI 会自动读取。

## 项目概述
- 项目名称: {{.ProjectName}}
- 创建时间: {{.CreatedAt}}

## 项目结构

## 重要决策

## 当前状态

## AI 指令
`

func (s *FileService) ReadEgoMemory(projectPath string) (string, error) {
	memoryPath := filepath.Join(projectPath, ".ego.md")
	data, err := os.ReadFile(memoryPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	return string(data), nil
}

func (s *FileService) WriteEgoMemory(projectPath string, content string) error {
	memoryPath := filepath.Join(projectPath, ".ego.md")
	dir := filepath.Dir(memoryPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(memoryPath, []byte(content), 0644)
}

func (s *FileService) InitEgoMemory(projectPath string, projectName string) error {
	memoryPath := filepath.Join(projectPath, ".ego.md")
	if _, err := os.Stat(memoryPath); err == nil {
		return nil
	}
	content := strings.ReplaceAll(egoMemoryTemplate, "{{.ProjectName}}", projectName)
	content = strings.ReplaceAll(content, "{{.CreatedAt}}", fmt.Sprintf("%s", "项目创建时自动生成"))
	return s.WriteEgoMemory(projectPath, content)
}

func (s *FileService) AppendEgoMemory(projectPath string, section string, content string) error {
	existing, err := s.ReadEgoMemory(projectPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if existing == "" {
		existing = egoMemoryTemplate
	}

	appended := strings.TrimRight(existing, "\n") + "\n\n## " + section + "\n" + content + "\n"
	return s.WriteEgoMemory(projectPath, appended)
}

func (s *FileService) GetEgoMemorySummary(projectPath string) string {
	content, err := s.ReadEgoMemory(projectPath)
	if err != nil || content == "" {
		return "暂无项目记忆"
	}
	lines := strings.Split(content, "\n")
	summaryLines := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if len(trimmed) > 120 {
			trimmed = trimmed[:120] + "..."
		}
		summaryLines = append(summaryLines, trimmed)
	}
	if len(summaryLines) > 30 {
		summaryLines = summaryLines[:30]
		summaryLines = append(summaryLines, "...(已截断)")
	}
	return strings.Join(summaryLines, "\n")
}

func (s *FileService) ExportToGoCode(egoData string) (string, error) {
	var ego EgoFile
	if err := json.Unmarshal([]byte(egoData), &ego); err != nil {
		return "", err
	}

	converter := core.NewCodeConverter()

	hasComponents := len(ego.Components) > 0

	if !hasComponents {
		var sb strings.Builder
		var bodyBuilder strings.Builder
		needsYigou := false

		bodyBuilder.WriteString(fmt.Sprintf("\tfmt.Println(\"程序: %s 启动\")\n", ego.AppName))
		for _, row := range ego.CodeRows {
			if row.Code == "" {
				continue
			}
			bodyBuilder.WriteString(fmt.Sprintf("\t// %s\n", row.Event))
			goCode, err := converter.ConvertChineseToGo(row.Code)
			if err != nil {
				bodyBuilder.WriteString(fmt.Sprintf("\t// 转换错误: %v\n", err))
				continue
			}
			lines := strings.Split(strings.TrimSpace(goCode), "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line != "" {
					bodyBuilder.WriteString(fmt.Sprintf("\t%s\n", line))
					if strings.Contains(line, "yigou.") {
						needsYigou = true
					}
				}
			}
		}
		bodyBuilder.WriteString("\tfmt.Println(\"程序运行完毕，按回车键退出...\")\n")
		bodyBuilder.WriteString("\tfmt.Scanln()\n")

		sb.WriteString("package main\n\n")
		if needsYigou {
			sb.WriteString("import (\n")
			sb.WriteString("\t\"fmt\"\n")
			sb.WriteString("\t\"yigou\"\n")
			sb.WriteString(")\n\n")
		} else {
			sb.WriteString("import \"fmt\"\n\n")
		}
		sb.WriteString("func main() {\n")
		sb.WriteString(bodyBuilder.String())
		sb.WriteString("}\n")
		return sb.String(), nil
	}

	return s.generateWinGUI(ego, converter)
}

type eventHandler struct {
	HandlerName string
	Code        string
	CompName    string
	EventName   string
}

func (s *FileService) generateWinGUI(ego EgoFile, converter *core.CodeConverter) (string, error) {
	var sb strings.Builder

	winTitle := ego.AppName
	winWidth := 800
	winHeight := 600
	if ego.WindowProps.Title != "" {
		winTitle = ego.WindowProps.Title
	}
	if ego.WindowProps.Width > 0 {
		winWidth = ego.WindowProps.Width
	}
	if ego.WindowProps.Height > 0 {
		winHeight = ego.WindowProps.Height
	}

	var handlers []eventHandler
	for i, row := range ego.CodeRows {
		if strings.TrimSpace(row.Code) == "" {
			continue
		}
		compName := ""
		eventName := ""
		parts := strings.SplitN(row.Event, " ", 2)
		if len(parts) >= 2 {
			compName = strings.TrimSpace(parts[0])
			eventName = strings.TrimSpace(parts[1])
		} else {
			compName = row.Event
			eventName = "被单击"
		}

		goCode, err := converter.ConvertChineseToGo(row.Code)
		if err != nil {
			continue
		}
		handlerName := fmt.Sprintf("egoEvt_%d", i)
		handlers = append(handlers, eventHandler{
			HandlerName: handlerName,
			Code:        goCode,
			CompName:    compName,
			EventName:   eventName,
		})
	}

	sb.WriteString("package main\n\n")
	sb.WriteString("import (\n")
	sb.WriteString("\t\"fmt\"\n")
	sb.WriteString("\t\"log\"\n")
	sb.WriteString("\t\"yigou\"\n")
	sb.WriteString("\t\"github.com/jchv/go-webview2\"\n")
	sb.WriteString(")\n\n")

	for i := range handlers {
		funcBody := s.extractFuncBody(handlers[i].Code)
		handlers[i].Code = funcBody
	}

	for _, h := range handlers {
		sb.WriteString(fmt.Sprintf("func %s() {\n", h.HandlerName))
		lines := strings.Split(strings.TrimSpace(h.Code), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" {
				sb.WriteString(fmt.Sprintf("\t%s\n", line))
			}
		}
		sb.WriteString("}\n\n")
	}

	sb.WriteString("func egoDispatchEvent(compIdx int, eventName string, compName string) {\n")
	sb.WriteString("\tyigou.PrintLine(\"[\" + compName + \"] \" + eventName)\n")
	if len(handlers) > 0 {
		sb.WriteString("\tkey := compName + \"_\" + eventName\n")
		sb.WriteString("\tswitch key {\n")
		for _, h := range handlers {
			eventKey := h.CompName + "_" + h.EventName
			sb.WriteString(fmt.Sprintf("\tcase \"%s\":\n", eventKey))
			sb.WriteString(fmt.Sprintf("\t\t%s()\n", h.HandlerName))
		}
		sb.WriteString("\t}\n")
	}
	sb.WriteString("}\n\n")

	sb.WriteString("var uiHTML = `\n")
	sb.WriteString(s.generateHTML(ego))
	sb.WriteString("`\n\n")

	sb.WriteString("func main() {\n")
	sb.WriteString(fmt.Sprintf("\tw := webview2.NewWithOptions(webview2.WebViewOptions{\n"))
	sb.WriteString("\t\tDebug:     false,\n")
	sb.WriteString("\t\tAutoFocus: true,\n")
	sb.WriteString("\t\tWindowOptions: webview2.WindowOptions{\n")
	sb.WriteString(fmt.Sprintf("\t\t\tTitle:  \"%s\",\n", winTitle))
	sb.WriteString(fmt.Sprintf("\t\t\tWidth:  %d,\n", winWidth))
	sb.WriteString(fmt.Sprintf("\t\t\tHeight: %d,\n", winHeight))
	sb.WriteString("\t\t\tCenter: true,\n")
	sb.WriteString("\t\t},\n")
	sb.WriteString("\t})\n")
	sb.WriteString("\tif w == nil {\n")
	sb.WriteString("\t\tlog.Fatalln(\"Failed to load webview.\")\n")
	sb.WriteString("\t}\n")
	sb.WriteString("\tdefer w.Destroy()\n")

	sb.WriteString("\tw.Bind(\"egoDispatchEvent\", egoDispatchEvent)\n")
	for _, h := range handlers {
		sb.WriteString(fmt.Sprintf("\tw.Bind(\"%s\", %s)\n", h.HandlerName, h.HandlerName))
	}

	sb.WriteString("\tw.SetHtml(uiHTML)\n")
	sb.WriteString("\tw.Run()\n")
	sb.WriteString("}\n")

	return sb.String(), nil
}

func (s *FileService) extractFuncBody(code string) string {
	code = strings.TrimSpace(code)

	funcStart := strings.Index(code, "{")
	if funcStart == -1 {
		return code
	}

	bodyStart := funcStart + 1
	bodyEnd := strings.LastIndex(code, "}")

	if bodyEnd > bodyStart {
		return strings.TrimSpace(code[bodyStart:bodyEnd])
	}
	return strings.TrimSpace(code[bodyStart:])
}

func (s *FileService) generateHTML(ego EgoFile) string {
	var sb strings.Builder

	bgColor := "#f0f0f0"
	if ego.WindowProps.BackgroundColor != "" {
		bgColor = ego.WindowProps.BackgroundColor
	}

	sb.WriteString("<!DOCTYPE html>\n")
	sb.WriteString("<html lang=\"zh-CN\">\n<head>\n")
	sb.WriteString("<meta charset=\"UTF-8\">\n")
	sb.WriteString(fmt.Sprintf("<title>%s</title>\n", ego.AppName))
	sb.WriteString("<style>\n")
	sb.WriteString("* { margin: 0; padding: 0; box-sizing: border-box; }\n")
	sb.WriteString(fmt.Sprintf("body { font-family: 'Microsoft YaHei', 'Segoe UI', sans-serif; background: %s; overflow: hidden; width: 100vw; height: 100vh; position: relative; }\n", bgColor))
	sb.WriteString(".win-button { position: absolute; padding: 6px 16px; border: 1px solid #0078d4; background: #0078d4; color: #fff; border-radius: 4px; cursor: pointer; font-size: 14px; }\n")
	sb.WriteString(".win-button:hover { background: #106ebe; }\n")
	sb.WriteString(".win-label { position: absolute; font-size: 14px; color: #333; }\n")
	sb.WriteString(".win-input { position: absolute; padding: 4px 8px; border: 1px solid #ccc; border-radius: 4px; font-size: 14px; background: #fff; }\n")
	sb.WriteString(".win-input:focus { border-color: #0078d4; outline: none; }\n")
	sb.WriteString(".win-checkbox { position: absolute; display: flex; align-items: center; gap: 6px; font-size: 14px; cursor: pointer; }\n")
	sb.WriteString(".win-radio { position: absolute; display: flex; align-items: center; gap: 6px; font-size: 14px; cursor: pointer; }\n")
	sb.WriteString(".win-group { position: absolute; border: 1px solid #ccc; border-radius: 4px; padding: 16px 8px 8px; }\n")
	sb.WriteString(".win-group-title { position: absolute; top: -8px; left: 12px; background: #f0f0f0; padding: 0 4px; font-size: 13px; color: #666; }\n")
	sb.WriteString(".win-combo { position: absolute; padding: 4px 8px; border: 1px solid #ccc; border-radius: 4px; font-size: 14px; background: #fff; }\n")
	sb.WriteString(".win-textarea { position: absolute; padding: 4px 8px; border: 1px solid #ccc; border-radius: 4px; font-size: 14px; background: #fff; resize: none; }\n")
	sb.WriteString(".win-textarea:focus { border-color: #0078d4; outline: none; }\n")
	sb.WriteString(".win-listbox { position: absolute; border: 1px solid #ccc; border-radius: 4px; background: #fff; overflow-y: auto; font-size: 14px; }\n")
	sb.WriteString(".win-listbox-item { padding: 4px 8px; cursor: pointer; }\n")
	sb.WriteString(".win-listbox-item:hover { background: #e8f0fe; }\n")
	sb.WriteString(".win-listbox-item.selected { background: #0078d4; color: #fff; }\n")
	sb.WriteString(".win-progress { position: absolute; height: 20px; border-radius: 10px; background: #e0e0e0; overflow: hidden; }\n")
	sb.WriteString(".win-progress-bar { height: 100%; background: #0078d4; border-radius: 10px; transition: width 0.3s; }\n")
	sb.WriteString(".win-slider { position: absolute; width: 100%; }\n")
	sb.WriteString(".win-table { position: absolute; border-collapse: collapse; background: #fff; border: 1px solid #ccc; border-radius: 4px; overflow: hidden; }\n")
	sb.WriteString(".win-table th { background: #f5f5f5; padding: 6px 12px; border: 1px solid #ddd; font-size: 13px; text-align: left; }\n")
	sb.WriteString(".win-table td { padding: 6px 12px; border: 1px solid #ddd; font-size: 13px; }\n")
	sb.WriteString(".win-tab { position: absolute; }\n")
	sb.WriteString(".win-tab-header { display: flex; border-bottom: 1px solid #ccc; }\n")
	sb.WriteString(".win-tab-btn { padding: 6px 16px; cursor: pointer; border: 1px solid #ccc; border-bottom: none; background: #f5f5f5; font-size: 13px; border-radius: 4px 4px 0 0; }\n")
	sb.WriteString(".win-tab-btn.active { background: #fff; border-bottom: 1px solid #fff; }\n")
	sb.WriteString(".win-tab-content { border: 1px solid #ccc; border-top: none; padding: 8px; background: #fff; min-height: 100px; }\n")
	sb.WriteString(".win-link { position: absolute; color: #0078d4; cursor: pointer; font-size: 14px; text-decoration: underline; }\n")
	sb.WriteString(".win-link:hover { color: #106ebe; }\n")
	sb.WriteString(".win-divider { position: absolute; height: 1px; background: #ccc; }\n")
	sb.WriteString(".win-container { position: absolute; border: 1px solid #ddd; border-radius: 4px; background: #fafafa; }\n")
	sb.WriteString(".win-image { position: absolute; border: 1px solid #ccc; background: #fff; display: flex; align-items: center; justify-content: center; color: #999; font-size: 12px; }\n")
	sb.WriteString("</style>\n</head>\n<body>\n")

	for i, comp := range ego.Components {
		w := comp.Width
		h := comp.Height
		if w == 0 {
			w = 100
		}
		if h == 0 {
			h = 28
		}
		text := comp.Text
		if text == "" {
			text = comp.Name
		}

		idx := i
		compName := comp.Name
		if compName == "" {
			compName = comp.Text
		}
		if compName == "" {
			compName = comp.Type + "_" + comp.ID
		}

		switch comp.Type {
		case "button":
			sb.WriteString(fmt.Sprintf("<div class=\"win-button\" data-comp-idx=\"%d\" data-comp-name=\"%s\" style=\"left:%dpx;top:%dpx;width:%dpx;height:%dpx;display:flex;align-items:center;justify-content:center;\">%s</div>\n",
				idx, compName, comp.X, comp.Y, w, h, text))
		case "label":
			sb.WriteString(fmt.Sprintf("<div class=\"win-label\" data-comp-idx=\"%d\" data-comp-name=\"%s\" style=\"left:%dpx;top:%dpx;width:%dpx;height:%dpx;display:flex;align-items:center;\">%s</div>\n",
				idx, compName, comp.X, comp.Y, w, h, text))
		case "input":
			sb.WriteString(fmt.Sprintf("<input class=\"win-input\" data-comp-idx=\"%d\" data-comp-name=\"%s\" style=\"left:%dpx;top:%dpx;width:%dpx;height:%dpx;\" placeholder=\"%s\" />\n",
				idx, compName, comp.X, comp.Y, w, h, text))
		case "checkbox":
			sb.WriteString(fmt.Sprintf("<label class=\"win-checkbox\" data-comp-idx=\"%d\" data-comp-name=\"%s\" style=\"left:%dpx;top:%dpx;width:%dpx;height:%dpx;\"><input type=\"checkbox\" />%s</label>\n",
				idx, compName, comp.X, comp.Y, w, h, text))
		case "radio":
			sb.WriteString(fmt.Sprintf("<label class=\"win-radio\" data-comp-idx=\"%d\" data-comp-name=\"%s\" style=\"left:%dpx;top:%dpx;width:%dpx;height:%dpx;\"><input type=\"radio\" />%s</label>\n",
				idx, compName, comp.X, comp.Y, w, h, text))
		case "group":
			sb.WriteString(fmt.Sprintf("<fieldset class=\"win-group\" style=\"left:%dpx;top:%dpx;width:%dpx;height:%dpx;\"><legend class=\"win-group-title\">%s</legend></fieldset>\n",
				comp.X, comp.Y, w, h, text))
		case "combo":
			sb.WriteString(fmt.Sprintf("<select class=\"win-combo\" data-comp-idx=\"%d\" data-comp-name=\"%s\" style=\"left:%dpx;top:%dpx;width:%dpx;height:%dpx;\"><option>%s</option></select>\n",
				idx, compName, comp.X, comp.Y, w, h, text))
		case "textarea":
			sb.WriteString(fmt.Sprintf("<textarea class=\"win-textarea\" data-comp-idx=\"%d\" data-comp-name=\"%s\" style=\"left:%dpx;top:%dpx;width:%dpx;height:%dpx;\" placeholder=\"%s\"></textarea>\n",
				idx, compName, comp.X, comp.Y, w, h, text))
		case "listbox":
			sb.WriteString(fmt.Sprintf("<div class=\"win-listbox\" data-comp-idx=\"%d\" data-comp-name=\"%s\" style=\"left:%dpx;top:%dpx;width:%dpx;height:%dpx;\"><div class=\"win-listbox-item\">%s</div></div>\n",
				idx, compName, comp.X, comp.Y, w, h, text))
		case "progress":
			sb.WriteString(fmt.Sprintf("<div class=\"win-progress\" data-comp-idx=\"%d\" data-comp-name=\"%s\" style=\"left:%dpx;top:%dpx;width:%dpx;height:%dpx;\"><div class=\"win-progress-bar\" style=\"width:30%%;\"></div></div>\n",
				idx, compName, comp.X, comp.Y, w, h))
		case "slider":
			sb.WriteString(fmt.Sprintf("<input type=\"range\" class=\"win-slider\" data-comp-idx=\"%d\" data-comp-name=\"%s\" style=\"left:%dpx;top:%dpx;width:%dpx;\" min=\"0\" max=\"100\" value=\"50\" />\n",
				idx, compName, comp.X, w))
		case "table":
			sb.WriteString(fmt.Sprintf("<table class=\"win-table\" style=\"left:%dpx;top:%dpx;width:%dpx;height:%dpx;\"><tr><th>列1</th><th>列2</th><th>列3</th></tr><tr><td>数据</td><td>数据</td><td>数据</td></tr></table>\n",
				comp.X, comp.Y, w, h))
		case "tab":
			sb.WriteString(fmt.Sprintf("<div class=\"win-tab\" data-comp-idx=\"%d\" data-comp-name=\"%s\" style=\"left:%dpx;top:%dpx;width:%dpx;height:%dpx;\"><div class=\"win-tab-header\"><div class=\"win-tab-btn active\">%s</div></div><div class=\"win-tab-content\">选项卡内容</div></div>\n",
				idx, compName, comp.X, comp.Y, w, h, text))
		case "link":
			sb.WriteString(fmt.Sprintf("<a class=\"win-link\" data-comp-idx=\"%d\" data-comp-name=\"%s\" style=\"left:%dpx;top:%dpx;\" href=\"#\">%s</a>\n",
				idx, compName, comp.X, comp.Y, text))
		case "divider":
			sb.WriteString(fmt.Sprintf("<hr class=\"win-divider\" style=\"left:%dpx;top:%dpx;width:%dpx;\" />\n",
				comp.X, comp.Y, w))
		case "container":
			sb.WriteString(fmt.Sprintf("<div class=\"win-container\" style=\"left:%dpx;top:%dpx;width:%dpx;height:%dpx;\">%s</div>\n",
				comp.X, comp.Y, w, h, text))
		case "image":
			sb.WriteString(fmt.Sprintf("<div class=\"win-image\" data-comp-idx=\"%d\" data-comp-name=\"%s\" style=\"left:%dpx;top:%dpx;width:%dpx;height:%dpx;\">%s</div>\n",
				idx, compName, comp.X, comp.Y, w, h, text))
		case "datepicker":
			sb.WriteString(fmt.Sprintf("<input type=\"date\" class=\"win-input\" data-comp-idx=\"%d\" data-comp-name=\"%s\" style=\"left:%dpx;top:%dpx;width:%dpx;height:%dpx;\" />\n",
				idx, compName, comp.X, comp.Y, w, h))
		}
	}

	sb.WriteString("<script>\n")
	sb.WriteString("function egoEmit(idx, name, evt, val) {\n")
	sb.WriteString("  try {\n")
	sb.WriteString("    if (window.egoDispatchEvent) {\n")
	sb.WriteString("      window.egoDispatchEvent(idx, evt, name);\n")
	sb.WriteString("    }\n")
	sb.WriteString("  } catch(e) { console.log('ego event err:', e); }\n")
	sb.WriteString("}\n")
	sb.WriteString("document.querySelectorAll('[data-comp-idx]').forEach(el => {\n")
	sb.WriteString("  const idx = parseInt(el.getAttribute('data-comp-idx'));\n")
	sb.WriteString("  const name = el.getAttribute('data-comp-name') || '';\n")
	sb.WriteString("  const tag = el.tagName.toLowerCase();\n")
	sb.WriteString("  const cls = el.className || '';\n")
	sb.WriteString("  if (cls.includes('win-button') || cls.includes('win-link') || cls.includes('win-image')) {\n")
	sb.WriteString("    el.addEventListener('click', (e) => { e.preventDefault(); egoEmit(idx, name, '被单击', ''); });\n")
	sb.WriteString("  } else if (tag === 'input' && el.type === 'range') {\n")
	sb.WriteString("    el.addEventListener('change', () => egoEmit(idx, name, '位置被改变', el.value));\n")
	sb.WriteString("    el.addEventListener('input', () => egoEmit(idx, name, '位置被改变', el.value));\n")
	sb.WriteString("  } else if (tag === 'input' || tag === 'textarea') {\n")
	sb.WriteString("    el.addEventListener('change', () => egoEmit(idx, name, '内容被改变', el.value));\n")
	sb.WriteString("    el.addEventListener('blur', () => egoEmit(idx, name, '失去焦点', el.value));\n")
	sb.WriteString("    el.addEventListener('focus', () => egoEmit(idx, name, '获得焦点', ''));\n")
	sb.WriteString("  } else if (tag === 'select') {\n")
	sb.WriteString("    el.addEventListener('change', () => egoEmit(idx, name, '列表项被选择', el.value));\n")
	sb.WriteString("  } else if (el.querySelector('input[type=\"checkbox\"]') || el.querySelector('input[type=\"radio\"]')) {\n")
	sb.WriteString("    const inp = el.querySelector('input');\n")
	sb.WriteString("    if (inp) inp.addEventListener('change', () => egoEmit(idx, name, '被单击', inp.checked));\n")
	sb.WriteString("  } else if (cls.includes('win-tab')) {\n")
	sb.WriteString("    el.querySelectorAll('.win-tab-btn').forEach((btn, ti) => {\n")
	sb.WriteString("      btn.addEventListener('click', () => {\n")
	sb.WriteString("        btn.parentElement.querySelectorAll('.win-tab-btn').forEach(b => b.classList.remove('active'));\n")
	sb.WriteString("        btn.classList.add('active');\n")
	sb.WriteString("        egoEmit(idx, name, '选项卡被切换', ti);\n")
	sb.WriteString("      });\n")
	sb.WriteString("    });\n")
	sb.WriteString("  } else if (cls.includes('win-listbox')) {\n")
	sb.WriteString("    el.querySelectorAll('.win-listbox-item').forEach((item, li) => {\n")
	sb.WriteString("      item.addEventListener('click', () => egoEmit(idx, name, '列表项被选择', li));\n")
	sb.WriteString("    });\n")
	sb.WriteString("  }\n")
	sb.WriteString("});\n")
	sb.WriteString("</script>\n")
	sb.WriteString("</body>\n</html>\n")

	return sb.String()
}

func sanitizeIdent(name string) string {
	result := strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			return r
		}
		return '_'
	}, name)
	if len(result) == 0 {
		return "ctrl"
	}
	if result[0] >= '0' && result[0] <= '9' {
		result = "c" + result
	}
	return result
}

func (s *FileService) GetComponents(egoData string) ([]UIComponent, error) {
	var ego EgoFile
	if err := json.Unmarshal([]byte(egoData), &ego); err != nil {
		return nil, err
	}
	return ego.Components, nil
}

func (s *FileService) SetComponents(egoData string, components []UIComponent) (string, error) {
	var ego EgoFile
	if err := json.Unmarshal([]byte(egoData), &ego); err != nil {
		return "", err
	}
	ego.Components = components
	data, err := json.MarshalIndent(ego, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (s *FileService) GetCodeRows(egoData string) ([]CodeRow, error) {
	var ego EgoFile
	if err := json.Unmarshal([]byte(egoData), &ego); err != nil {
		return nil, err
	}
	return ego.CodeRows, nil
}

func (s *FileService) SetCodeRows(egoData string, codeRows []CodeRow) (string, error) {
	var ego EgoFile
	if err := json.Unmarshal([]byte(egoData), &ego); err != nil {
		return "", err
	}
	ego.CodeRows = codeRows
	data, err := json.MarshalIndent(ego, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (s *FileService) GetAppName(egoData string) (string, error) {
	var ego EgoFile
	if err := json.Unmarshal([]byte(egoData), &ego); err != nil {
		return "", err
	}
	return ego.AppName, nil
}

func (s *FileService) SetAppName(egoData string, appName string) (string, error) {
	var ego EgoFile
	if err := json.Unmarshal([]byte(egoData), &ego); err != nil {
		return "", err
	}
	ego.AppName = appName
	data, err := json.MarshalIndent(ego, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (s *FileService) GetVersion(egoData string) (string, error) {
	var ego EgoFile
	if err := json.Unmarshal([]byte(egoData), &ego); err != nil {
		return "", err
	}
	return ego.Version, nil
}

func (s *FileService) SetVersion(egoData string, version string) (string, error) {
	var ego EgoFile
	if err := json.Unmarshal([]byte(egoData), &ego); err != nil {
		return "", err
	}
	ego.Version = version
	data, err := json.MarshalIndent(ego, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

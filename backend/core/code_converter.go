package core

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

type ChineseFunction struct {
	ChineseName string   `json:"chineseName"`
	GoFunction  string   `json:"goFunction"`
	Description string   `json:"description"`
	Category    string   `json:"category"`
	Params      []string `json:"params"`
	ReturnType  string   `json:"returnType"`
}

type CodeTemplate struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Code        string `json:"code"`
}

var FunctionMap = map[string]ChineseFunction{
	// ========== 输出 ==========
	"输出.打印":  {ChineseName: "输出.打印", GoFunction: "yigou.Print", Description: "输出内容到控制台", Category: "输出", Params: []string{"内容"}, ReturnType: ""},
	"输出.打印行": {ChineseName: "输出.打印行", GoFunction: "yigou.PrintLine", Description: "输出一行内容到控制台", Category: "输出", Params: []string{"内容"}, ReturnType: ""},
	"输出.格式化": {ChineseName: "输出.格式化", GoFunction: "yigou.FormatString", Description: "格式化字符串", Category: "输出", Params: []string{"格式", "参数..."}, ReturnType: "文本型"},
	"调试输出":   {ChineseName: "调试输出", GoFunction: "yigou.PrintLine", Description: "输出调试信息到IDE控制台", Category: "输出", Params: []string{"输出内容"}, ReturnType: ""},

	// ========== 窗口/交互 ==========
	"窗口.弹出提示": {ChineseName: "窗口.弹出提示", GoFunction: "yigou.MsgBox", Description: "弹出提示信息", Category: "窗口", Params: []string{"提示内容"}, ReturnType: ""},
	"窗口.弹出询问": {ChineseName: "窗口.弹出询问", GoFunction: "yigou.ConfirmBox", Description: "弹出询问对话框，返回是/否", Category: "窗口", Params: []string{"询问内容"}, ReturnType: "逻辑型"},
	"窗口.弹出输入": {ChineseName: "窗口.弹出输入", GoFunction: "yigou.InputBox", Description: "弹出输入对话框", Category: "窗口", Params: []string{"提示内容", "默认值"}, ReturnType: "文本型"},
	"窗口.关闭":   {ChineseName: "窗口.关闭", GoFunction: "yigou.Exit", Description: "关闭程序", Category: "窗口", Params: []string{"退出码"}, ReturnType: ""},

	// ========== 文件 ==========
	"文件.读取文本":   {ChineseName: "文件.读取文本", GoFunction: "yigou.ReadFile", Description: "读取文本文件内容", Category: "文件", Params: []string{"文件路径"}, ReturnType: "文本型"},
	"文件.写入文本":   {ChineseName: "文件.写入文本", GoFunction: "yigou.WriteFile", Description: "写入文本内容到文件", Category: "文件", Params: []string{"文件路径", "文本内容"}, ReturnType: "逻辑型"},
	"文件.是否存在":   {ChineseName: "文件.是否存在", GoFunction: "yigou.FileExists", Description: "判断文件是否存在", Category: "文件", Params: []string{"文件路径"}, ReturnType: "逻辑型"},
	"文件.删除":     {ChineseName: "文件.删除", GoFunction: "yigou.DeleteFile", Description: "删除文件", Category: "文件", Params: []string{"文件路径"}, ReturnType: "逻辑型"},
	"文件.重命名":    {ChineseName: "文件.重命名", GoFunction: "yigou.RenameFile", Description: "重命名文件", Category: "文件", Params: []string{"原路径", "新路径"}, ReturnType: "逻辑型"},
	"文件.复制":     {ChineseName: "文件.复制", GoFunction: "yigou.CopyFile", Description: "复制文件", Category: "文件", Params: []string{"源路径", "目标路径"}, ReturnType: "逻辑型"},
	"文件.获取大小":   {ChineseName: "文件.获取大小", GoFunction: "yigou.FileSize", Description: "获取文件大小（字节）", Category: "文件", Params: []string{"文件路径"}, ReturnType: "整数型"},
	"文件.获取修改时间": {ChineseName: "文件.获取修改时间", GoFunction: "yigou.FileModTime", Description: "获取文件最后修改时间", Category: "文件", Params: []string{"文件路径"}, ReturnType: "文本型"},
	"文件.创建目录":   {ChineseName: "文件.创建目录", GoFunction: "yigou.MkdirAll", Description: "创建多级目录", Category: "文件", Params: []string{"目录路径"}, ReturnType: "逻辑型"},
	"文件.删除目录":   {ChineseName: "文件.删除目录", GoFunction: "yigou.RemoveAll", Description: "删除目录及所有内容", Category: "文件", Params: []string{"目录路径"}, ReturnType: "逻辑型"},
	"文件.列出文件":   {ChineseName: "文件.列出文件", GoFunction: "yigou.ReadDir", Description: "列出目录中的文件", Category: "文件", Params: []string{"目录路径"}, ReturnType: "文本型数组"},

	// ========== 字符串 ==========
	"字符串.截取":   {ChineseName: "字符串.截取", GoFunction: "yigou.Substring", Description: "截取字符串的一部分", Category: "字符串", Params: []string{"原字符串", "起始位置", "截取长度"}, ReturnType: "文本型"},
	"字符串.长度":   {ChineseName: "字符串.长度", GoFunction: "yigou.StrLen", Description: "获取字符串长度", Category: "字符串", Params: []string{"字符串"}, ReturnType: "整数型"},
	"字符串.查找":   {ChineseName: "字符串.查找", GoFunction: "yigou.StrIndex", Description: "查找子串位置", Category: "字符串", Params: []string{"原字符串", "查找内容"}, ReturnType: "整数型"},
	"字符串.替换":   {ChineseName: "字符串.替换", GoFunction: "yigou.StrReplaceAll", Description: "替换字符串中的内容", Category: "字符串", Params: []string{"原字符串", "旧内容", "新内容"}, ReturnType: "文本型"},
	"字符串.分割":   {ChineseName: "字符串.分割", GoFunction: "yigou.StrSplit", Description: "按分隔符分割字符串", Category: "字符串", Params: []string{"原字符串", "分隔符"}, ReturnType: "文本型数组"},
	"字符串.合并":   {ChineseName: "字符串.合并", GoFunction: "yigou.StrJoin", Description: "合并字符串数组", Category: "字符串", Params: []string{"字符串数组", "分隔符"}, ReturnType: "文本型"},
	"字符串.转大写":  {ChineseName: "字符串.转大写", GoFunction: "yigou.StrToUpper", Description: "转换为大写", Category: "字符串", Params: []string{"字符串"}, ReturnType: "文本型"},
	"字符串.转小写":  {ChineseName: "字符串.转小写", GoFunction: "yigou.StrToLower", Description: "转换为小写", Category: "字符串", Params: []string{"字符串"}, ReturnType: "文本型"},
	"字符串.去空格":  {ChineseName: "字符串.去空格", GoFunction: "yigou.StrTrimSpace", Description: "去除首尾空格", Category: "字符串", Params: []string{"字符串"}, ReturnType: "文本型"},
	"字符串.是否包含": {ChineseName: "字符串.是否包含", GoFunction: "yigou.StrContains", Description: "判断是否包含子串", Category: "字符串", Params: []string{"原字符串", "子串"}, ReturnType: "逻辑型"},

	// ========== 数值 ==========
	"数值.转文本":  {ChineseName: "数值.转文本", GoFunction: "yigou.Itoa", Description: "将整数转换为文本", Category: "数值", Params: []string{"整数"}, ReturnType: "文本型"},
	"数值.转整数":  {ChineseName: "数值.转整数", GoFunction: "yigou.Atoi", Description: "将文本转换为整数", Category: "数值", Params: []string{"文本"}, ReturnType: "整数型"},
	"数值.转小数":  {ChineseName: "数值.转小数", GoFunction: "yigou.ParseFloat", Description: "将文本转换为小数", Category: "数值", Params: []string{"文本"}, ReturnType: "小数型"},
	"数值.取整":   {ChineseName: "数值.取整", GoFunction: "yigou.Floor", Description: "向下取整", Category: "数值", Params: []string{"数值"}, ReturnType: "小数型"},
	"数值.四舍五入": {ChineseName: "数值.四舍五入", GoFunction: "yigou.Round", Description: "四舍五入取整", Category: "数值", Params: []string{"数值"}, ReturnType: "小数型"},
	"数值.取绝对值": {ChineseName: "数值.取绝对值", GoFunction: "yigou.Abs", Description: "取绝对值", Category: "数值", Params: []string{"数值"}, ReturnType: "小数型"},
	"数值.取最大值": {ChineseName: "数值.取最大值", GoFunction: "yigou.Max", Description: "取两个数中的最大值", Category: "数值", Params: []string{"数值1", "数值2"}, ReturnType: "小数型"},
	"数值.取最小值": {ChineseName: "数值.取最小值", GoFunction: "yigou.Min", Description: "取两个数中的最小值", Category: "数值", Params: []string{"数值1", "数值2"}, ReturnType: "小数型"},
	"数值.取随机数": {ChineseName: "数值.取随机数", GoFunction: "yigou.RandomInt", Description: "取0到N-1的随机整数", Category: "数值", Params: []string{"最大值"}, ReturnType: "整数型"},
	"数值.求平方根": {ChineseName: "数值.求平方根", GoFunction: "yigou.Sqrt", Description: "求平方根", Category: "数值", Params: []string{"数值"}, ReturnType: "小数型"},
	"数值.求幂":   {ChineseName: "数值.求幂", GoFunction: "yigou.Pow", Description: "求x的y次方", Category: "数值", Params: []string{"底数", "指数"}, ReturnType: "小数型"},

	// ========== 系统 ==========
	"延时":         {ChineseName: "延时", GoFunction: "yigou.Sleep", Description: "程序延时指定毫秒数", Category: "系统", Params: []string{"毫秒数"}, ReturnType: ""},
	"系统.获取当前时间":  {ChineseName: "系统.获取当前时间", GoFunction: "yigou.NowTime", Description: "获取当前系统时间", Category: "系统", Params: []string{}, ReturnType: "文本型"},
	"系统.执行命令":    {ChineseName: "系统.执行命令", GoFunction: "yigou.RunCommand", Description: "执行系统命令", Category: "系统", Params: []string{"命令", "参数..."}, ReturnType: "文本型"},
	"系统.获取环境变量":  {ChineseName: "系统.获取环境变量", GoFunction: "yigou.GetEnv", Description: "获取环境变量值", Category: "系统", Params: []string{"变量名"}, ReturnType: "文本型"},
	"系统.设置环境变量":  {ChineseName: "系统.设置环境变量", GoFunction: "yigou.SetEnv", Description: "设置环境变量", Category: "系统", Params: []string{"变量名", "值"}, ReturnType: ""},
	"系统.获取命令行参数": {ChineseName: "系统.获取命令行参数", GoFunction: "yigou.GetArgs", Description: "获取命令行参数", Category: "系统", Params: []string{}, ReturnType: "文本型数组"},
	"系统.获取当前目录":  {ChineseName: "系统.获取当前目录", GoFunction: "yigou.Getwd", Description: "获取当前工作目录", Category: "系统", Params: []string{}, ReturnType: "文本型"},
	"系统.设置当前目录":  {ChineseName: "系统.设置当前目录", GoFunction: "yigou.Chdir", Description: "设置当前工作目录", Category: "系统", Params: []string{"目录路径"}, ReturnType: "逻辑型"},
	"系统.打开文件夹":   {ChineseName: "系统.打开文件夹", GoFunction: "yigou.OpenDir", Description: "打开指定文件夹", Category: "系统", Params: []string{"文件夹路径"}, ReturnType: ""},

	// ========== 网络 ==========
	"网络.发送请求":   {ChineseName: "网络.发送请求", GoFunction: "yigou.HTTPGet", Description: "发送HTTP GET请求", Category: "网络", Params: []string{"网址"}, ReturnType: "文本型"},
	"网络.发送POST": {ChineseName: "网络.发送POST", GoFunction: "yigou.HTTPPost", Description: "发送HTTP POST请求", Category: "网络", Params: []string{"网址", "数据类型", "数据"}, ReturnType: "文本型"},
	"网络.下载文件":   {ChineseName: "网络.下载文件", GoFunction: "yigou.DownloadFile", Description: "下载文件到本地", Category: "网络", Params: []string{"网址", "保存路径"}, ReturnType: "逻辑型"},
	"网络.编码URL":  {ChineseName: "网络.编码URL", GoFunction: "yigou.URLEncode", Description: "URL编码", Category: "网络", Params: []string{"文本"}, ReturnType: "文本型"},
	"网络.解码URL":  {ChineseName: "网络.解码URL", GoFunction: "yigou.URLDecode", Description: "URL解码", Category: "网络", Params: []string{"文本"}, ReturnType: "文本型"},

	// ========== JSON ==========
	"JSON.编码": {ChineseName: "JSON.编码", GoFunction: "yigou.JSONMarshal", Description: "将数据编码为JSON文本", Category: "JSON", Params: []string{"数据"}, ReturnType: "文本型"},
	"JSON.解码": {ChineseName: "JSON.解码", GoFunction: "yigou.JSONUnmarshal", Description: "将JSON文本解码为数据", Category: "JSON", Params: []string{"JSON文本", "目标变量"}, ReturnType: ""},

	// ========== 加密 ==========
	"加密.MD5":      {ChineseName: "加密.MD5", GoFunction: "yigou.MD5Hash", Description: "计算MD5哈希值", Category: "加密", Params: []string{"文本"}, ReturnType: "文本型"},
	"加密.SHA256":   {ChineseName: "加密.SHA256", GoFunction: "yigou.SHA256Hash", Description: "计算SHA256哈希值", Category: "加密", Params: []string{"文本"}, ReturnType: "文本型"},
	"加密.Base64编码": {ChineseName: "加密.Base64编码", GoFunction: "yigou.Base64Encode", Description: "Base64编码", Category: "加密", Params: []string{"数据"}, ReturnType: "文本型"},
	"加密.Base64解码": {ChineseName: "加密.Base64解码", GoFunction: "yigou.Base64Decode", Description: "Base64解码", Category: "加密", Params: []string{"文本"}, ReturnType: "文本型"},
}

var CodeTemplates = []CodeTemplate{
	{
		Name:        "如果-否则",
		Description: "条件判断结构",
		Category:    "流程控制",
		Code:        "如果(条件) {\n    // 条件成立时执行\n} 否则 {\n    // 条件不成立时执行\n}",
	},
	{
		Name:        "循环-计次",
		Description: "计次循环",
		Category:    "流程控制",
		Code:        "循环(计次变量, 起始值, 结束值, 步长) {\n    // 循环体\n}",
	},
	{
		Name:        "循环-判断",
		Description: "条件循环",
		Category:    "流程控制",
		Code:        "循环(条件) {\n    // 条件成立时循环\n}",
	},
	{
		Name:        "变量-定义",
		Description: "定义变量",
		Category:    "变量",
		Code:        "变量 变量名 = 初始值",
	},
	{
		Name:        "函数-定义",
		Description: "定义函数",
		Category:    "函数",
		Code:        "函数 函数名(参数1, 参数2) {\n    // 函数体\n    返回 返回值\n}",
	},
	{
		Name:        "异常-捕获",
		Description: "异常捕获处理",
		Category:    "异常处理",
		Code:        "尝试 {\n    // 可能出错的代码\n} 捕获(错误信息) {\n    // 错误处理\n}",
	},
}

type CodeConverter struct{}

func NewCodeConverter() *CodeConverter {
	return &CodeConverter{}
}

func (c *CodeConverter) ConvertChineseToGo(chineseCode string) (string, error) {
	result := chineseCode

	type replacement struct {
		placeholder string
		goFunc      string
		chineseName string
	}

	var replacements []replacement
	for chineseName, fn := range FunctionMap {
		placeholder := fmt.Sprintf("__GO_FUNC_%d__", len(replacements))
		replacements = append(replacements, replacement{
			placeholder: placeholder,
			goFunc:      fn.GoFunction,
			chineseName: chineseName,
		})
	}

	sort.Slice(replacements, func(i, j int) bool {
		return len(replacements[i].chineseName) > len(replacements[j].chineseName)
	})

	for _, r := range replacements {
		result = strings.ReplaceAll(result, r.chineseName+"(", r.placeholder+"(")
	}

	for _, r := range replacements {
		result = strings.ReplaceAll(result, r.placeholder, r.goFunc)
	}

	return result, nil
}

func (c *CodeConverter) ConvertCodeRowsToGo(codeRows []CodeRowData) string {
	var sb strings.Builder
	var bodyBuilder strings.Builder
	needsYigou := false

	for _, row := range codeRows {
		if row.Code == "" {
			continue
		}
		bodyBuilder.WriteString(fmt.Sprintf("\t// %s\n", row.Event))
		goCode, _ := c.ConvertChineseToGo(row.Code)
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

	sb.WriteString("package main\n\n")
	if needsYigou {
		sb.WriteString("import (\n")
		sb.WriteString("\t\"yigou\"\n")
		sb.WriteString(")\n\n")
	}
	sb.WriteString("func main() {\n")
	sb.WriteString(bodyBuilder.String())
	sb.WriteString("}\n")

	return sb.String()
}

func (c *CodeConverter) GetFunctionSuggestions(prefix string) []ChineseFunction {
	var suggestions []ChineseFunction
	lowerPrefix := strings.ToLower(prefix)

	for _, fn := range FunctionMap {
		lowerName := strings.ToLower(fn.ChineseName)
		if strings.HasPrefix(lowerName, lowerPrefix) ||
			strings.Contains(lowerName, lowerPrefix) {
			suggestions = append(suggestions, fn)
		}
	}

	sort.Slice(suggestions, func(i, j int) bool {
		return suggestions[i].ChineseName < suggestions[j].ChineseName
	})

	if len(suggestions) > 20 {
		suggestions = suggestions[:20]
	}

	return suggestions
}

func (c *CodeConverter) GetFunctionsByCategory(category string) []ChineseFunction {
	var functions []ChineseFunction
	for _, fn := range FunctionMap {
		if fn.Category == category {
			functions = append(functions, fn)
		}
	}
	sort.Slice(functions, func(i, j int) bool {
		return functions[i].ChineseName < functions[j].ChineseName
	})
	return functions
}

func (c *CodeConverter) GetAllFunctions() []ChineseFunction {
	var functions []ChineseFunction
	for _, fn := range FunctionMap {
		functions = append(functions, fn)
	}
	sort.Slice(functions, func(i, j int) bool {
		return functions[i].ChineseName < functions[j].ChineseName
	})
	return functions
}

func (c *CodeConverter) GetCategories() []string {
	categorySet := make(map[string]bool)
	for _, fn := range FunctionMap {
		categorySet[fn.Category] = true
	}
	var categories []string
	for cat := range categorySet {
		categories = append(categories, cat)
	}
	sort.Strings(categories)
	return categories
}

func (c *CodeConverter) GetTemplates() []CodeTemplate {
	return CodeTemplates
}

func (c *CodeConverter) ValidateCode(chineseCode string) []string {
	var errors []string

	re := regexp.MustCompile(`([\p{Han}]+\.?[\p{Han}]*)\(([^)]*)\)`)
	matches := re.FindAllStringSubmatch(chineseCode, -1)

	for _, match := range matches {
		funcName := match[1]
		args := match[2]

		if fn, exists := FunctionMap[funcName]; exists {
			if args != "" {
				argList := strings.Split(args, ",")
				expectedCount := len(fn.Params)
				actualCount := len(argList)
				if actualCount < expectedCount {
					errors = append(errors, fmt.Sprintf("函数 %s 参数不足：需要%d个参数，提供了%d个", funcName, expectedCount, actualCount))
				} else if actualCount > expectedCount {
					errors = append(errors, fmt.Sprintf("函数 %s 参数过多：需要%d个参数，提供了%d个", funcName, expectedCount, actualCount))
				}
			} else if len(fn.Params) > 0 {
				errors = append(errors, fmt.Sprintf("函数 %s 缺少参数：需要%d个参数", funcName, len(fn.Params)))
			}
		} else {
			errors = append(errors, fmt.Sprintf("未定义的函数: %s", funcName))
		}
	}

	return errors
}

type CodeRowData struct {
	Event string `json:"event"`
	Code  string `json:"code"`
}

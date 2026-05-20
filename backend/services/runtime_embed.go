package services

const yigouRuntimeSource = `package yigou

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
	"time"
)

func MsgBox(msg string) {
	user32 := syscall.NewLazyDLL("user32.dll")
	messageBox := user32.NewProc("MessageBoxW")
	title, _ := syscall.UTF16PtrFromString("提示")
	text, _ := syscall.UTF16PtrFromString(msg)
	messageBox.Call(0, uintptr(unsafe.Pointer(text)), uintptr(unsafe.Pointer(title)), 0x40)
}

func ConfirmBox(msg string) bool {
	user32 := syscall.NewLazyDLL("user32.dll")
	messageBox := user32.NewProc("MessageBoxW")
	title, _ := syscall.UTF16PtrFromString("确认")
	text, _ := syscall.UTF16PtrFromString(msg)
	ret, _, _ := messageBox.Call(0, uintptr(unsafe.Pointer(text)), uintptr(unsafe.Pointer(title)), 0x04|0x20)
	return ret == 6
}

func InputBox(prompt string, defaultVal string) string {
	fmt.Print("📝 " + prompt)
	if defaultVal != "" {
		fmt.Print(" [" + defaultVal + "]")
	}
	fmt.Print(": ")
	var input string
	fmt.Scanln(&input)
	if input == "" {
		return defaultVal
	}
	return input
}

func Print(args ...interface{}) {
	fmt.Print(args...)
}

func PrintLine(args ...interface{}) {
	fmt.Println(args...)
}

func FormatString(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}

func Sleep(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

func NowTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func ReadFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("读取文件失败:", err)
		return ""
	}
	return string(data)
}

func WriteFile(path string, content string) bool {
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		fmt.Println("写入文件失败:", err)
		return false
	}
	return true
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func DeleteFile(path string) bool {
	err := os.Remove(path)
	if err != nil {
		fmt.Println("删除文件失败:", err)
		return false
	}
	return true
}

func RenameFile(oldPath, newPath string) bool {
	err := os.Rename(oldPath, newPath)
	if err != nil {
		fmt.Println("重命名失败:", err)
		return false
	}
	return true
}

func CopyFile(src, dst string) bool {
	sourceFile, err := os.Open(src)
	if err != nil {
		fmt.Println("打开源文件失败:", err)
		return false
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		fmt.Println("创建目标文件失败:", err)
		return false
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		fmt.Println("复制文件失败:", err)
		return false
	}
	return true
}

func FileSize(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		return -1
	}
	return info.Size()
}

func FileModTime(path string) string {
	info, err := os.Stat(path)
	if err != nil {
		return ""
	}
	return info.ModTime().Format("2006-01-02 15:04:05")
}

func MkdirAll(path string) bool {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		fmt.Println("创建目录失败:", err)
		return false
	}
	return true
}

func RemoveAll(path string) bool {
	err := os.RemoveAll(path)
	if err != nil {
		fmt.Println("删除目录失败:", err)
		return false
	}
	return true
}

func ReadDir(path string) []string {
	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("读取目录失败:", err)
		return nil
	}
	var names []string
	for _, entry := range entries {
		names = append(names, entry.Name())
	}
	return names
}

func Substring(s string, start, length int) string {
	runes := []rune(s)
	if start < 0 {
		start = 0
	}
	if start >= len(runes) {
		return ""
	}
	end := start + length
	if end > len(runes) {
		end = len(runes)
	}
	return string(runes[start:end])
}

func StrLen(s string) int {
	return len([]rune(s))
}

func StrIndex(s, substr string) int {
	return strings.Index(s, substr)
}

func StrReplaceAll(s, old, new string) string {
	return strings.ReplaceAll(s, old, new)
}

func StrSplit(s, sep string) []string {
	return strings.Split(s, sep)
}

func StrJoin(elems []string, sep string) string {
	return strings.Join(elems, sep)
}

func StrToUpper(s string) string {
	return strings.ToUpper(s)
}

func StrToLower(s string) string {
	return strings.ToLower(s)
}

func StrTrimSpace(s string) string {
	return strings.TrimSpace(s)
}

func StrContains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func Itoa(n int) string {
	return strconv.Itoa(n)
}

func Atoi(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		fmt.Println("转换整数失败:", err)
		return 0
	}
	return n
}

func ParseFloat(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		fmt.Println("转换小数失败:", err)
		return 0
	}
	return f
}

func Floor(f float64) float64 {
	return math.Floor(f)
}

func Round(f float64) float64 {
	return math.Round(f)
}

func Abs(f float64) float64 {
	return math.Abs(f)
}

func Max(a, b float64) float64 {
	return math.Max(a, b)
}

func Min(a, b float64) float64 {
	return math.Min(a, b)
}

func RandomInt(n int) int {
	return rand.Intn(n)
}

func Sqrt(f float64) float64 {
	return math.Sqrt(f)
}

func Pow(x, y float64) float64 {
	return math.Pow(x, y)
}

func SortStrings(slice []string) []string {
	sorted := make([]string, len(slice))
	copy(sorted, slice)
	sort.Strings(sorted)
	return sorted
}

func SortInts(slice []int) []int {
	sorted := make([]int, len(slice))
	copy(sorted, slice)
	sort.Ints(sorted)
	return sorted
}

func ReverseStrings(slice []string) []string {
	reversed := make([]string, len(slice))
	for i, v := range slice {
		reversed[len(slice)-1-i] = v
	}
	return reversed
}

func HTTPGet(urlStr string) string {
	resp, err := http.Get(urlStr)
	if err != nil {
		fmt.Println("HTTP请求失败:", err)
		return ""
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应失败:", err)
		return ""
	}
	return string(body)
}

func HTTPPost(urlStr, contentType, data string) string {
	resp, err := http.Post(urlStr, contentType, strings.NewReader(data))
	if err != nil {
		fmt.Println("HTTP POST请求失败:", err)
		return ""
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应失败:", err)
		return ""
	}
	return string(body)
}

func DownloadFile(urlStr, savePath string) bool {
	resp, err := http.Get(urlStr)
	if err != nil {
		fmt.Println("下载失败:", err)
		return false
	}
	defer resp.Body.Close()

	dir := filepath.Dir(savePath)
	os.MkdirAll(dir, 0755)

	file, err := os.Create(savePath)
	if err != nil {
		fmt.Println("创建文件失败:", err)
		return false
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Println("保存文件失败:", err)
		return false
	}
	return true
}

func URLEncode(s string) string {
	return url.QueryEscape(s)
}

func URLDecode(s string) string {
	decoded, err := url.QueryUnescape(s)
	if err != nil {
		return s
	}
	return decoded
}

func JSONMarshal(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		fmt.Println("JSON编码失败:", err)
		return ""
	}
	return string(data)
}

func JSONUnmarshal(data string, v interface{}) error {
	return json.Unmarshal([]byte(data), v)
}

func MD5Hash(s string) string {
	hash := md5.Sum([]byte(s))
	return fmt.Sprintf("%x", hash)
}

func SHA256Hash(s string) string {
	hash := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", hash)
}

func Base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func Base64Decode(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		fmt.Println("Base64解码失败:", err)
		return nil
	}
	return data
}

func RunCommand(name string, args ...string) string {
	cmd := exec.Command(name, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: 0x08000000,
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("执行命令失败:", err)
	}
	return string(output)
}

func GetEnv(key string) string {
	return os.Getenv(key)
}

func SetEnv(key, value string) {
	os.Setenv(key, value)
}

func GetArgs() []string {
	return os.Args[1:]
}

func Getwd() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	return dir
}

func Chdir(path string) bool {
	err := os.Chdir(path)
	if err != nil {
		fmt.Println("切换目录失败:", err)
		return false
	}
	return true
}

func OpenDir(path string) {
	var cmd *exec.Cmd
	if _, err := os.Stat("C:\\Windows\\explorer.exe"); err == nil {
		cmd = exec.Command("explorer", path)
	} else {
		cmd = exec.Command("open", path)
	}
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: 0x08000000,
	}
	cmd.Start()
}

func Exit(code int) {
	os.Exit(code)
}
`

package utils

import (
	"fmt"
	"os"
	"runtime"
)

// 颜色常量
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	ColorBold   = "\033[1m"
)

// isColorSupported 检查是否支持颜色输出
func isColorSupported() bool {
	if runtime.GOOS == "windows" {
		return false // Windows cmd 默认不支持 ANSI 颜色
	}

	term := os.Getenv("TERM")
	return term != "" && term != "dumb"
}

// colorize 为文本添加颜色
func Colorize(text, color string) string {
	if !isColorSupported() {
		return text
	}
	return color + text + ColorReset
}

// colorizeMethod 为 HTTP 方法添加对应颜色
func ColorizeMethod(method string) string {
	if !isColorSupported() {
		return method
	}

	colors := map[string]string{
		"GET":    ColorGreen,
		"POST":   ColorBlue,
		"PUT":    ColorYellow,
		"DELETE": ColorRed,
		"PATCH":  ColorPurple,
		"HEAD":   ColorCyan,
	}

	if color, ok := colors[method]; ok {
		return Colorize(fmt.Sprintf("%-6s", method), color)
	}
	return fmt.Sprintf("%-6s", method)
}

// colorizeMode 为 Gin 模式添加颜色
func ColorizeMode(mode string) string {
	switch mode {
	case "debug":
		return Colorize(mode, ColorYellow)
	case "release":
		return Colorize(mode, ColorGreen)
	case "test":
		return Colorize(mode, ColorBlue)
	default:
		return mode
	}
}

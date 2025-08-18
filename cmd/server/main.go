package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"actiondelta/internal/config"
	"actiondelta/internal/indexer"
	"actiondelta/internal/repository"
	"actiondelta/internal/router"
	"actiondelta/internal/utils"
)

func main() {
    // 设置生产环境日志
    logger, _ := zap.NewProduction()
    defer logger.Sync()
    zap.ReplaceGlobals(logger)

    // 打印启动横幅
    printBanner()

    // 加载配置
    printStep("📋 Loading configuration...")
    if err := config.Load(); err != nil {
        zap.L().Fatal("failed to load config", zap.Error(err))
    }
    printSuccess("Configuration loaded successfully")

    // 初始化数据库
    printStep("🎯 Connecting to MongoDB...")
    if err := repository.InitMongo(context.Background()); err != nil {
        zap.L().Fatal("failed to init mongo", zap.Error(err))
    }
    defer repository.CloseMongo(context.Background())
    printSuccess("MongoDB connected successfully")
    zap.L().Info("database connected",
        zap.String("database", "actiondelta"),
        zap.String("status", "connected"))

    // 确保索引
    printStep("📊 Ensuring database indexes...")
    if err := indexer.EnsureAllIndexes(context.Background()); err != nil {
        zap.L().Fatal("failed to ensure indexes", zap.Error(err))
    }
    printSuccess("Database indexes ensured")

    // 创建路由
    printStep("🛣️  Setting up routes...")
    r := router.New()

    // 美化显示路由信息
    printRoutes(r)

    // 创建服务器
    srv := &http.Server{
        Addr:              fmt.Sprintf(":%d", config.C.Server.Port),
        Handler:           r,
        ReadTimeout:       15 * time.Second,
        ReadHeaderTimeout: 10 * time.Second,
        WriteTimeout:      30 * time.Second,
        IdleTimeout:       60 * time.Second,
    }

    // 启动服务器
    go func() {
        printServerInfo(config.C.Server.Port, gin.Mode())
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            zap.L().Fatal("http server error", zap.Error(err))
        }
    }()

    // 等待关闭信号
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    printStep("⏹️  Shutting down server...")

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    if err := srv.Shutdown(ctx); err != nil {
        zap.L().Error("server shutdown error", zap.Error(err))
    }

    printSuccess("Server stopped gracefully")
    zap.L().Info("server stopped")
}

// printBanner 打印启动横幅
func printBanner() {
    banner := `
╔══════════════════════════════════════╗
║            🎭 actiondelta API           ║
║              v1.0.0                  ║
║         Built with ❤️  in Go          ║
╚══════════════════════════════════════╝`

    fmt.Println(utils.Colorize(banner, utils.ColorCyan))
    fmt.Println()
}

// printStep 打印步骤信息
func printStep(message string) {
    fmt.Printf("%s %s\n", utils.Colorize("▶", utils.ColorBlue), message)
}

// printSuccess 打印成功信息
func printSuccess(message string) {
    fmt.Printf("%s %s\n", utils.Colorize("✅", utils.ColorGreen), utils.Colorize(message, utils.ColorGreen))
}

// printServerInfo 打印服务器启动信息
func printServerInfo(port int, mode string) {
    fmt.Println()
    fmt.Println(utils.Colorize("🚀 Server Information", utils.ColorGreen))
    fmt.Println(strings.Repeat("─", 40))
    fmt.Printf("   ├─ %s %d\n", utils.Colorize("Port:", utils.ColorWhite), port)
    fmt.Printf("   ├─ %s %s\n", utils.Colorize("Mode:", utils.ColorWhite), utils.ColorizeMode(mode))
    fmt.Printf("   ├─ %s %s\n", utils.Colorize("Time:", utils.ColorWhite), time.Now().Format("15:04:05"))
    fmt.Printf("   └─ %s %s\n", utils.Colorize("Status:", utils.ColorWhite), utils.Colorize("Running", utils.ColorGreen))
    fmt.Println()

    zap.L().Info("server started",
        zap.Int("port", port),
        zap.String("mode", mode),
        zap.Time("start_time", time.Now()))
}

// printRoutes 美化打印路由信息
func printRoutes(r *gin.Engine) {
    routes := r.Routes()

    if len(routes) == 0 {
        return
    }

    // 按功能分组路由
    groups := groupRoutes(routes)

    fmt.Printf("%s Found %d routes\n", utils.Colorize("📋", utils.ColorYellow), len(routes))
    fmt.Println()

    // 打印每个分组
    for _, groupName := range []string{
        "🏥 Health Check",
        "🔐 Authentication",
        "👤 User Management",
        "👥 Relations",
        "💬 Messaging",
        "🏠 Rooms",
        "📁 File Operations",
        "📄 Static Files",
    } {
        if routeList, exists := groups[groupName]; exists && len(routeList) > 0 {
            fmt.Printf("%s\n", utils.Colorize(groupName, utils.ColorCyan))
            fmt.Println(strings.Repeat("─", 50))

            // 排序路由
            sort.Slice(routeList, func(i, j int) bool {
                return routeList[i].Path < routeList[j].Path
            })

            for _, route := range routeList {
                fmt.Printf("  %-10s %s\n",
                    utils.ColorizeMethod(route.Method),
                    route.Path)
            }
            fmt.Println()
        }
    }
}

// groupRoutes 按功能对路由进行分组
func groupRoutes(routes gin.RoutesInfo) map[string]gin.RoutesInfo {
    groups := make(map[string]gin.RoutesInfo)

    for _, route := range routes {
        var groupName string

        switch {
        case strings.Contains(route.Path, "/healthz"):
            groupName = "🏥 Health Check"
        case strings.Contains(route.Path, "/static"):
            groupName = "📄 Static Files"
        case strings.Contains(route.Path, "/auth") ||
             strings.Contains(route.Path, "/login") ||
             strings.Contains(route.Path, "/send_code"):
            groupName = "🔐 Authentication"
        case strings.Contains(route.Path, "/user"):
            groupName = "👤 User Management"
        case strings.Contains(route.Path, "/relation"):
            groupName = "👥 Relations"
        case strings.Contains(route.Path, "/message"):
            groupName = "💬 Messaging"
        case strings.Contains(route.Path, "/room"):
            groupName = "🏠 Rooms"
        case strings.Contains(route.Path, "/file"):
            groupName = "📁 File Operations"
        default:
            groupName = "🔧 Others"
        }

        groups[groupName] = append(groups[groupName], route)
    }

    return groups
}

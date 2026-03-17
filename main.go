package main

import (
	"context"
	"embed"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"fg-abyss/internal/app/handlers"
	"fg-abyss/internal/app/services"
	"fg-abyss/internal/infrastructure/database"
	"fg-abyss/internal/infrastructure/repositories"
	"fg-abyss/internal/plugin"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// getAppDirectory 获取应用程序所在目录
func getAppDirectory() string {
	exePath, err := os.Executable()
	if err != nil {
		log.Printf("Warning: Could not get executable path: %v", err)
		return "."
	}
	return filepath.Dir(exePath)
}

//go:embed all:frontend/dist
var assets embed.FS

// WindowCreateEvent is the event data for creating a new window
type WindowCreateEvent struct {
	Title  string `json:"title"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
}

func init() {
	// Register a custom event whose associated data type is string.
	application.RegisterEvent[string]("time")
	application.RegisterEvent[WindowCreateEvent]("createWindow")
}

func main() {
	// 初始化数据库
	log.Println("=== Initializing Database ===")
	dbInstance, err := database.Init(&database.Config{
		Path: "data/app.db",
	})
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	log.Println("=== Database initialized successfully ===")

	// 创建仓储实例
	log.Println("=== Creating Repositories ===")
	projectRepo := repositories.NewProjectRepository(dbInstance)
	webshellRepo := repositories.NewWebShellRepository(dbInstance)
	settingRepo := repositories.NewSettingRepository(dbInstance)
	log.Println("=== Repositories created ===")

	// 创建服务实例
	log.Println("=== Creating Services ===")
	appService := services.NewAppService(dbInstance, projectRepo, webshellRepo)
	projectService := services.NewProjectService(projectRepo)
	webshellService := services.NewWebShellService(webshellRepo)
	settingService := services.NewSettingService(settingRepo)
	log.Println("=== Services created ===")

	// 初始化设置（自动创建表并填充默认值）
	log.Println("=== Initializing Settings ===")
	defaultSettings := settingService.GetDefaultSettings()
	if err := settingService.InitializeDefaults(defaultSettings); err != nil {
		log.Printf("Warning: Failed to initialize default settings: %v", err)
	} else {
		log.Println("=== Default settings initialized successfully ===")
	}

	// 加载设置到内存
	if err := settingService.Initialize(); err != nil {
		log.Printf("Warning: Failed to load settings: %v", err)
	} else {
		log.Println("=== Settings loaded to memory successfully ===")
	}

	// 创建处理器实例
	log.Println("=== Creating Handlers ===")
	systemHandler := handlers.NewSystemHandler(appService)
	projectHandler := handlers.NewProjectHandler(projectService)
	webshellHandler := handlers.NewWebShellHandler(webshellService)
	settingHandler := handlers.NewSettingHandler(settingService)
	connectionHandler := handlers.NewConnectionHandler()
	commandHandler := handlers.NewCommandHandler()
	fileHandler := handlers.NewFileHandler()
	payloadHandler := handlers.NewPayloadHandler()
	databaseHandler := handlers.NewDatabaseHandler()
	batchHandler := handlers.NewBatchHandler()
	proxyHandler := handlers.NewProxyHandler()
	encryptionHandler := handlers.NewEncryptionHandler()
	auditHandler := handlers.NewAuditHandler()

	// 初始化插件系统
	log.Println("=== Initializing Plugin System ===")
	appDir := getAppDirectory()
	pluginDir := filepath.Join(appDir, "plugins")
	dataDir := filepath.Join(appDir, "data", "plugins")

	pluginLoader := plugin.NewPluginLoader(pluginDir, dataDir, "1.0.0")
	if err := pluginLoader.Initialize(context.Background()); err != nil {
		log.Printf("Warning: Failed to initialize plugin system: %v", err)
	} else {
		log.Println("=== Plugin system initialized successfully ===")
	}

	pluginHandler := handlers.NewPluginHandler(pluginLoader)
	log.Println("=== Handlers created ===")

	// 加载前端资源
	log.Println("Loading frontend assets...")
	frontendAssets, err := fs.Sub(assets, "frontend/dist")
	if err != nil {
		log.Fatalf("Failed to load frontend assets: %v", err)
	}

	// 调试：列出嵌入的文件
	log.Println("=== Embedded files in frontend/dist ===")
	entries, err := fs.ReadDir(frontendAssets, ".")
	if err != nil {
		log.Printf("ERROR: Could not read embedded files: %v", err)
	} else {
		for _, entry := range entries {
			if entry.IsDir() {
				log.Printf("  [DIR]  %s", entry.Name())
			} else {
				info, _ := entry.Info()
				log.Printf("  [FILE] %s (%d bytes)", entry.Name(), info.Size())
			}
		}
	}

	// 创建 Wails 应用
	log.Println("Creating Wails application...")
	app := application.New(application.Options{
		Name:        "FG-ABYSS",
		Description: "FG-ABYSS Application",
		Services: []application.Service{
			application.NewService(systemHandler),
			application.NewService(projectHandler),
			application.NewService(webshellHandler),
			application.NewService(settingHandler),
			application.NewService(connectionHandler),
			application.NewService(commandHandler),
			application.NewService(fileHandler),
			application.NewService(payloadHandler),
			application.NewService(databaseHandler),
			application.NewService(batchHandler),
			application.NewService(proxyHandler),
			application.NewService(encryptionHandler),
			application.NewService(auditHandler),
			application.NewService(pluginHandler),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(frontendAssets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})
	log.Println("Wails application created successfully")

	// 创建窗口
	log.Println("Creating window...")
	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "FG-ABYSS",
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/index.html",
		Frameless:        true,
		Width:            1600,
		Height:           900,
		MinWidth:         1500,
		MinHeight:        900,
	})
	log.Println("Window created successfully")

	// 启动应用
	log.Println("Starting application...")
	if err := app.Run(); err != nil {
		log.Printf("Application error: %v", err)
	}
}

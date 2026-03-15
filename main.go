package main

import (
	"embed"
	"io/fs"
	"log"

	"fg-abyss/internal/app/handlers"
	"fg-abyss/internal/app/services"
	"fg-abyss/internal/infrastructure/database"
	"fg-abyss/internal/infrastructure/repositories"

	"github.com/wailsapp/wails/v3/pkg/application"
)

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
	log.Println("=== Repositories created ===")

	// 创建服务实例
	log.Println("=== Creating Services ===")
	appService := services.NewAppService(dbInstance, projectRepo, webshellRepo)
	projectService := services.NewProjectService(projectRepo)
	webshellService := services.NewWebShellService(webshellRepo, projectRepo)
	log.Println("=== Services created ===")

	// 创建处理器实例
	log.Println("=== Creating Handlers ===")
	systemHandler := handlers.NewSystemHandler(appService)
	projectHandler := handlers.NewProjectHandler(projectService)
	webshellHandler := handlers.NewWebShellHandler(webshellService)
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

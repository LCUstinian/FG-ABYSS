package main

import (
	"embed"
	"io/fs"
	"log"
	"time"

	"fg-abyss/backend/db"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// Wails uses Go's `embed` package to embed the frontend files into the binary.
// Any files in the frontend/dist folder will be embedded into the binary and
// made available to the frontend.
// See https://pkg.go.dev/embed for more information.

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
	// This is not required, but the binding generator will pick up registered events
	// and provide a strongly typed JS/TS API for them.
	application.RegisterEvent[string]("time")
	application.RegisterEvent[WindowCreateEvent]("createWindow")
}

// main function serves as the application's entry point. It initializes the application, creates a window,
// and starts a goroutine that emits a time-based event every second. It subsequently runs the application and
// logs any error that might occur.

//go:generate go tool golang.org/x/tools/cmd/stringer -type=WindowCreateEvent

func main() {
	// 初始化数据库
	dbInstance, err := db.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 获取底层的 sql.DB 实例以检查连接
	sqlDB, err := dbInstance.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}

	// 检查连接
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("=== Database initialized successfully ===")

	// 创建 App 实例
	appInstance := NewApp(dbInstance)
	log.Println("=== App instance created ===")

	// 使用 fs.Sub 提取前端资源子目录
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

	// Create a new Wails application
	log.Println("Creating Wails application...")
	app := application.New(application.Options{
		Name:        "FG-ABYSS",
		Description: "FG-ABYSS Application",
		Services: []application.Service{
			application.NewService(appInstance),
		},
		Assets: application.AssetOptions{
			// 使用 AssetFileServerFS 处理嵌入的前端资源
			// 注意：frontendAssets 已经是 fs.Sub 提取的子目录，直接指向 dist 目录
			Handler: application.AssetFileServerFS(frontendAssets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})
	log.Println("Wails application created successfully")

	// Create a new window with the necessary options.
	log.Println("Creating window...")
	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "FG-ABYSS",
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		// URL 设置为空字符串，让 Wails 自动加载 index.html
		// 或者显式设置为 "/index.html"
		URL:       "/index.html",
		Frameless: true,
		Width:     1600,
		Height:    900,
		MinWidth:  1500,
		MinHeight: 900,
	})
	log.Println("Window created successfully")

	// Create a goroutine that emits an event containing the current time every second.
	// The frontend can listen to this event and update the UI accordingly.
	go func() {
		for {
			now := time.Now().Format(time.RFC1123)
			app.Event.Emit("time", now)
			time.Sleep(time.Second)
		}
	}()

	// Listen for createWindow events from the frontend
	app.Event.On("createWindow", func(event *application.CustomEvent) {
		// Handle the event
		app.Window.NewWithOptions(application.WebviewWindowOptions{
			Title:            "New Window",
			BackgroundColour: application.NewRGB(27, 38, 54),
			URL:              "about:blank",
			Width:            800,
			Height:           600,
			X:                100,
			Y:                100,
		})
	})

	// Run the application. This blocks until the application has been exited.
	err = app.Run()

	// If an error occurred while running the application, log it and exit.
	if err != nil {
		log.Fatal(err)
	}
}

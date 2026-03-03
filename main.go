package main

import (
	"embed"
	_ "embed"
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

	log.Println("Database initialized successfully (Pure Go SQLite)")

	// Create a new Wails application by providing the necessary options.
	// Variables 'Name' and 'Description' are for application metadata.
	// 'Assets' configures the asset server with the 'FS' variable pointing to the frontend files.
	// 'Bind' is a list of Go struct instances. The frontend has access to the methods of these instances.
	// 'Mac' options tailor the application when running an macOS.
	app := application.New(application.Options{
		Name:        "FG-ABYSS",
		Description: "A demo of using raw HTML & CSS",
		Services: []application.Service{
			application.NewService(&GreetService{}),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	// Create a new window with the necessary options.
	// 'Title' is the title of the window.
	// 'Mac' options tailor the window when running on macOS.
	// 'BackgroundColour' is the background colour of the window.
	// 'URL' is the URL that will be loaded into the webview.
	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "FG-ABYSS",
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/",
		Frameless:        true,
		Width:            1600,
		Height:           900,
		MinWidth:         1200,
		MinHeight:        800,
	})

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

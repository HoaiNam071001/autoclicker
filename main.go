package main

import (
	"embed"
	"os"
	"os/signal"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"

	"autoclicker/internal/cli"
	"autoclicker/internal/winfloat"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// `autoclicker run …` (bước 22): gửi lệnh sang app đang mở rồi thoát – không mở giao diện
	// `autoclicker --mcp` (bước 31): MCP qua stdio cho AI – chuyển tiếp sang app; stdin / stdout là pipe của AI, không gắn console
	if cli.IsMCP(os.Args[1:]) {
		os.Exit(cli.NewClient(nil).MCP(os.Stdin, os.Stdout, os.Stderr))
	}
	if args := os.Args[1:]; cli.IsCommand(args) {
		cli.AttachConsole()
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt)
		os.Exit(cli.Main(args, cli.NewClient(sig), os.Stdout, os.Stderr))
	}

	winfloat.PrepareLaunch() // trước wails.Run – thanh mini / overlay nổi được trên app toàn màn hình (macOS)

	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "autoclicker",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		// Frameless + nền trong suốt cố định từ đầu (Wails v2 không đổi được lúc chạy) để overlay nhìn xuyên được.
		// Mỗi màn tự vẽ nền; MainScreen có thanh tiêu đề riêng để kéo cửa sổ.
		Frameless:        true,
		BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 0},
		Mac:              &mac.Options{WebviewIsTransparent: true},
		// Windows: chỉ WebviewIsTransparent thì cửa sổ Win32 vẫn đặc – cần WindowIsTranslucent (không bitmap nền) và
		// BackdropType None (không Mica / Acrylic) để nhìn xuyên đúng alpha của trang.
		Windows:       &windows.Options{WebviewIsTransparent: true, WindowIsTranslucent: true, BackdropType: windows.None},
		OnStartup:     app.startup,
		OnShutdown:    app.shutdown,
		OnBeforeClose: app.beforeClose,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

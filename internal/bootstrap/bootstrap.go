package bootstrap

import (
	"fmt"

	ChatHandler "websocket-chat/internal/app/chat/interface/rest"
	ChatUsecase "websocket-chat/internal/app/chat/usecase"
	"websocket-chat/internal/domain/env"
	"websocket-chat/internal/infra/fiber"
)

func Start() error {
	// 1. Muat konfigurasi environment
	config, err := env.New()
	if err != nil {
		return err
	}

	// 2. Inisialisasi Usecase & jalankan broadcast event loop di goroutine latar belakang
	chatUsecase := ChatUsecase.NewChatUsecase()
	go chatUsecase.Run()

	// 3. Inisialisasi Fiber app
	app := fiber.New(config)

	// 4. Daftarkan router handler
	ChatHandler.NewChatHandler(app, chatUsecase)

	// 5. Jalankan server listen
	return app.Listen(fmt.Sprintf(":%d", config.AppPort))
}

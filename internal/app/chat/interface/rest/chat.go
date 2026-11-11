package rest

import (
	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"

	"websocket-chat/internal/app/chat/usecase"
)

// WsConfig mendefinisikan konfigurasi koneksi WebSocket untuk Fiber
var WsConfig = websocket.Config{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	Origins:         []string{"*"},
}

// ChatHandler menangani interaksi rute HTTP & upgrade WebSocket
type ChatHandler struct {
	chatUsecase usecase.ChatUsecaseItf
}

// NewChatHandler mendaftarkan endpoint WebSocket ke router Fiber
func NewChatHandler(routerGroup fiber.Router, chatUsecase usecase.ChatUsecaseItf) {
	handler := &ChatHandler{
		chatUsecase: chatUsecase,
	}
	_ = handler

	// =========================================================================
	// TODO [Lab 1 - Task 1]: Daftarkan Route Group "/ws" & Endpoint Echo
	// 1. Buat route group "/ws" dengan middleware upgrade:
	//    wsGroup := routerGroup.Group("/ws", handler.UpgradeMiddleware)
	// 2. Daftarkan endpoint GET "/echo" menggunakan websocket.New:
	//    wsGroup.Get("/echo", websocket.New(handler.HandleEcho, WsConfig))
	// =========================================================================

	// =========================================================================
	// TODO [Lab 2 - Task 6]: Daftarkan rute multi-user chat room (/ws/chat)
	// wsGroup.Get("/chat", websocket.New(handler.HandleChat, WsConfig))
	// =========================================================================
}

// UpgradeMiddleware memastikan request HTTP memiliki header Upgrade: websocket
func (h *ChatHandler) UpgradeMiddleware(c *fiber.Ctx) error {
	// =========================================================================
	// TODO [Lab 1 - Task 2]: Periksa apakah request adalah upgrade WebSocket
	// 1. Periksa header dengan: websocket.IsWebSocketUpgrade(c)
	// 2. Jika valid (true): lanjutkan request dengan return c.Next()
	// 3. Jika bukan (false): tolak dengan return fiber.ErrUpgradeRequired
	// =========================================================================

	return fiber.ErrUpgradeRequired
}

// HandleEcho menangani Lab 1: Echo Server (/ws/echo)
func (h *ChatHandler) HandleEcho(c *websocket.Conn) {
	log.Println("[ECHO] Client baru terhubung untuk uji coba Echo")
	defer c.Close()

	// =========================================================================
	// TODO [Lab 1 - Task 3]: Loop baca pesan & balas kembali (Echo)
	// 1. Jalankan loop tak terbatas: for { ... }
	// 2. Baca pesan masuk: messageType, message, err := c.ReadMessage()
	// 3. Jika err != nil: log disconnect & break loop
	// 4. Balas langsung ke client:
	//    echoReply := append([]byte("[Echo]: "), message...)
	//    if err := c.WriteMessage(messageType, echoReply); err != nil { break }
	// =========================================================================
}

// HandleChat menangani Lab 2 & 3: Multi-User Chat Room (/ws/chat?user=nama)
func (h *ChatHandler) HandleChat(c *websocket.Conn) {
	username := c.Query("user")
	if username == "" {
		username = "Anonymous"
	}

	// =========================================================================
	// TODO [Lab 2 - Task 7]: Inisialisasi Client & Daftarkan ke Hub Usecase
	// 1. Buat instance client: client := NewClientConn(h.chatUsecase, c, username)
	// 2. Registrasi client ke Hub: h.chatUsecase.Register(client)
	// 3. Jalankan WritePump di background goroutine: go client.WritePump()
	// 4. Jalankan ReadPump secara blocking: client.ReadPump()
	// =========================================================================
}

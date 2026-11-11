package usecase

import (
	"websocket-chat/internal/domain/entity"
)

// ChatUsecaseItf adalah kontrak interface untuk usecase pengelolaan chat room & broadcast
type ChatUsecaseItf interface {
	Register(client entity.Client)
	Unregister(client entity.Client)
	Broadcast(msg entity.Message)
	Run()
}

// ChatUsecase mengelola state koneksi aktif dan mendistribusikan pesan secara aman (single goroutine select loop)
type ChatUsecase struct {
	clients    map[entity.Client]bool
	register   chan entity.Client
	unregister chan entity.Client
	broadcast  chan entity.Message
}

// NewChatUsecase membuat instance baru ChatUsecase
func NewChatUsecase() ChatUsecaseItf {
	return &ChatUsecase{
		clients:    make(map[entity.Client]bool),
		register:   make(chan entity.Client),
		unregister: make(chan entity.Client),
		broadcast:  make(chan entity.Message),
	}
}

func (u *ChatUsecase) Register(client entity.Client) {
	u.register <- client
}

func (u *ChatUsecase) Unregister(client entity.Client) {
	u.unregister <- client
}

func (u *ChatUsecase) Broadcast(msg entity.Message) {
	u.broadcast <- msg
}

// Run menjalankan event loop utama dalam satu goroutine untuk menghindari data race pada map clients
func (u *ChatUsecase) Run() {
	for {
		select {
		// =========================================================================
		// TODO [Lab 2 - Task 1]: Tangani pendaftaran client baru (case client := <-u.register)
		// 1. Simpan client ke dalam map: u.clients[client] = true
		// 2. Log client baru yang terhubung
		// 3. Buat notifikasi sistem:
		//    joinMsg := entity.NewSystemMessage(fmt.Sprintf("%s bergabung ke dalam chat!", client.Username()))
		// 4. Siarkan notifikasi ke seluruh client: u.broadcastToAll(joinMsg)
		// =========================================================================
		case client := <-u.register:
			_ = client

		// =========================================================================
		// TODO [Lab 2 - Task 2]: Tangani pelepasan client yang keluar (case client := <-u.unregister)
		// 1. Periksa apakah client ada di u.clients:
		//    if _, ok := u.clients[client]; ok { ... }
		// 2. Hapus dari map: delete(u.clients, client)
		// 3. Tutup resource client: client.Close()
		// 4. Buat notifikasi sistem:
		//    leaveMsg := entity.NewSystemMessage(fmt.Sprintf("%s telah meninggalkan chat", client.Username()))
		// 5. Siarkan notifikasi ke seluruh client tersisa: u.broadcastToAll(leaveMsg)
		// =========================================================================
		case client := <-u.unregister:
			_ = client

		// =========================================================================
		// TODO [Lab 2 - Task 3]: Tangani broadcast pesan chat (case msg := <-u.broadcast)
		// Siarkan pesan ke seluruh client aktif:
		// u.broadcastToAll(msg)
		// =========================================================================
		case msg := <-u.broadcast:
			_ = msg
		}
	}
}

// broadcastToAll menyalurkan pesan ke seluruh channel send client aktif
func (u *ChatUsecase) broadcastToAll(msg entity.Message) {
	for client := range u.clients {
		select {
		case client.SendChannel() <- msg:
		default:
			// Buffer client penuh atau client unresponsive, putuskan koneksi untuk mencegah lagging
			delete(u.clients, client)
			client.Close()
		}
	}
}

package rest

import (
	"sync"
	"time"

	"github.com/gofiber/contrib/websocket"

	"websocket-chat/internal/app/chat/usecase"
	"websocket-chat/internal/domain/entity"
)

const (
	// Waktu batas toleransi penulisan pesan ke socket TCP
	writeWait = 10 * time.Second

	// Waktu batas toleransi menunggu pesan Pong dari client sebelum koneksi dianggap mati
	pongWait = 60 * time.Second

	// Interval pengiriman Ping ke client (harus lebih kecil dari pongWait)
	pingPeriod = (pongWait * 9) / 10

	// Ukuran maksimal payload pesan yang diizinkan (bytes)
	maxMessageSize = 1024
)

// ClientConn mengimplementasikan kontrak interface entity.Client
type ClientConn struct {
	chatUsecase usecase.ChatUsecaseItf
	conn        *websocket.Conn
	username    string
	send        chan entity.Message
	closeOnce   sync.Once
}

// NewClientConn menginisialisasi koneksi client baru
func NewClientConn(chatUsecase usecase.ChatUsecaseItf, conn *websocket.Conn, username string) *ClientConn {
	return &ClientConn{
		chatUsecase: chatUsecase,
		conn:        conn,
		username:    username,
		send:        make(chan entity.Message, 256),
	}
}

func (c *ClientConn) Username() string {
	return c.username
}

func (c *ClientConn) SendChannel() chan entity.Message {
	return c.send
}

func (c *ClientConn) Close() {
	c.closeOnce.Do(func() {
		close(c.send)
	})
}

// ReadPump adalah goroutine pembaca tunggal dari koneksi WebSocket menuju Usecase
func (c *ClientConn) ReadPump() {
	defer func() {
		c.chatUsecase.Unregister(c)
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))

	// =========================================================================
	// TODO [Lab 3 - Task 1]: Pasang Pong Handler untuk mendeteksi heartbeat
	// Setiap kali frame Pong diterima dari client, perpanjang batas waktu pembacaan:
	// c.conn.SetPongHandler(func(string) error {
	//     _ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	//     return nil
	// })
	// =========================================================================

	for {
		// =====================================================================
		// TODO [Lab 2 - Task 4]: Baca pesan dari WebSocket frame & kirim ke Hub
		// 1. Baca frame websocket:
		//    _, messageBytes, err := c.conn.ReadMessage()
		//    if err != nil { break }
		// 2. Bungkus payload pesan:
		//    msg := entity.NewChatMessage(c.username, string(messageBytes))
		// 3. Kirim ke Hub untuk disiarkan ke semua client:
		//    c.chatUsecase.Broadcast(msg)
		// =====================================================================
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

// WritePump adalah goroutine penulis tunggal dari channel send menuju koneksi TCP client
func (c *ClientConn) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		// =====================================================================
		// TODO [Lab 2 - Task 5]: Tulis pesan dari channel send ke WebSocket client
		// 1. Terima pesan dari channel send: case msg, ok := <-c.send:
		// 2. Set deadline penulisan: c.conn.SetWriteDeadline(time.Now().Add(writeWait))
		// 3. Jika !ok (channel ditutup oleh Hub):
		//    _ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
		//    return
		// 4. Buka writer TextMessage:
		//    w, err := c.conn.NextWriter(websocket.TextMessage)
		//    if err != nil { return }
		// 5. Marshal pesan ke JSON & kirim ke socket:
		//    msgBytes, _ := json.Marshal(msg)
		//    w.Write(msgBytes)
		//    w.Close()
		// =====================================================================
		case msg, ok := <-c.send:
			_ = msg
			if !ok {
				return
			}

		// =====================================================================
		// TODO [Lab 3 - Task 2]: Kirim heartbeat frame Ping secara periodik
		// 1. Setiap ticker.C berdetak: case <-ticker.C:
		// 2. Set deadline penulisan: c.conn.SetWriteDeadline(time.Now().Add(writeWait))
		// 3. Kirim frame Ping ke client:
		//    if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil { return }
		// =====================================================================
		case <-ticker.C:
		}
	}
}

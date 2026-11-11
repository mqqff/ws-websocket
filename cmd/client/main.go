package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type incomingMessage struct {
	Type      string    `json:"type"`
	Sender    string    `json:"sender"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

func main() {
	username := flag.String("user", "budi", "Username anda di chat room")
	mode := flag.String("mode", "chat", "Pilihan mode koneksi: 'chat' atau 'echo'")
	serverHost := flag.String("host", "localhost:8080", "Host server WebSocket")
	flag.Parse()

	var path string
	var rawQuery string
	if *mode == "echo" {
		path = "/ws/echo"
	} else {
		path = "/ws/chat"
		rawQuery = "user=" + url.QueryEscape(*username)
	}

	u := url.URL{Scheme: "ws", Host: *serverHost, Path: path, RawQuery: rawQuery}
	targetURL := u.String()

	fmt.Println("=======================================================")
	fmt.Printf(" Menghubungkan ke %s ...\n", targetURL)

	conn, resp, err := websocket.DefaultDialer.Dial(targetURL, nil)
	if err != nil {
		fmt.Printf(" Gagal terhubung ke server (%v)\n", err)
		if resp != nil {
			fmt.Printf(" Status HTTP: %s\n", resp.Status)
		}
		fmt.Println(" Tips: Pastikan server sudah dijalankan dengan: go run cmd/api/main.go")
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println(" Terhubung! Status HTTP 101 Switching Protocols.")
	fmt.Println("=======================================================")
	if *mode == "echo" {
		fmt.Println(" Mode: ECHO TEST (Lab 1)")
		fmt.Println(" Ketik teks apa saja, server akan membalas dengan [Echo] ...")
	} else {
		fmt.Printf(" Mode: LIVE CHAT (User: %s) [Lab 2 & 3]\n", *username)
		fmt.Println(" Ketik pesan lalu tekan Enter untuk mengirim.")
	}
	fmt.Println(" Ketik '/exit' atau tekan Ctrl+C untuk keluar.")
	fmt.Println("=======================================================")

	// Channel untuk menangani sinyal interrupt (Ctrl+C)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	done := make(chan struct{})

	// Goroutine pembaca pesan dari server
	go func() {
		defer close(done)
		for {
			_, messageBytes, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
					fmt.Printf("\n[Koneksi terputus]: %v\n", err)
				} else {
					fmt.Println("\n[Koneksi ditutup]")
				}
				return
			}

			// Parse pesan chat JSON
			var msg incomingMessage
			if err := json.Unmarshal(messageBytes, &msg); err == nil && msg.Content != "" {
				timeStr := msg.Timestamp.Format("15:04:05")
				if msg.Type == "system" {
					fmt.Printf("[%s] [SYSTEM] %s\n", timeStr, msg.Content)
				} else {
					fmt.Printf("[%s] %s: %s\n", timeStr, msg.Sender, msg.Content)
				}
			} else {
				// Tampilkan pesan teks mentah (misal respon echo)
				fmt.Printf("%s\n", string(messageBytes))
			}
		}
	}()

	// Scanner input keyboard
	scanner := bufio.NewScanner(os.Stdin)
	inputChan := make(chan string)
	go func() {
		for scanner.Scan() {
			inputChan <- scanner.Text()
		}
	}()

	for {
		select {
		case <-done:
			return

		case text := <-inputChan:
			text = strings.TrimSpace(text)
			if text == "" {
				continue
			}
			if text == "/exit" {
				fmt.Println("Keluar dari chat...")
				_ = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Keluar"))
				return
			}

			err := conn.WriteMessage(websocket.TextMessage, []byte(text))
			if err != nil {
				log.Printf("Gagal mengirim pesan: %v\n", err)
				return
			}

		case <-interrupt:
			fmt.Println("\nMenerima sinyal interrupt, menutup koneksi...")
			_ = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Client terminated"))
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

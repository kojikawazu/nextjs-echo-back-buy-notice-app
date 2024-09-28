package websocket

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

// WebSocketアップグレード用の設定
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// セキュリティのため、特定のオリジンを許可する
		origin := r.Header.Get("Origin")
		if origin == "http://localhost:3000" || origin == "http://yourdomain.com" {
			return true
		}
		return false
	},
}

// クライアントを保持するためのマップとロック
var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan []byte)
var mutex = &sync.Mutex{}

// WebSocketハンドラー
func HandleWebSocket(c echo.Context) error {
	// WebSocket接続をアップグレード
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	// クライアントをマップに追加
	mutex.Lock()
	clients[ws] = true
	mutex.Unlock()

	log.Println("WebSocket connection established")

	// クライアントからのメッセージを受信
	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Println("WebSocket read error:", err)
			mutex.Lock()
			delete(clients, ws)
			mutex.Unlock()
			ws.Close()
			break
		}
		log.Printf("Received message: %s", msg)

		// 受信したメッセージをブロードキャストチャンネルに送信
		broadcast <- msg
	}

	log.Println("WebSocket connection closed")
	return nil
}

func HandleMessages() {
	for {
		// ブロードキャストされたメッセージを受信
		msg := <-broadcast
		// すべてのクライアントにメッセージを送信
		mutex.Lock()
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Printf("WebSocket write error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
		mutex.Unlock()
	}
}

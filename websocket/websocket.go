package websocket

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

var (
	ctx = context.Background()

	// Redisクライアントの設定
	rdb = redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_URL"), // Docker ComposeでRedisコンテナを利用
	})

	// 許可するオリジンを環境変数から取得
	allowedOrigins = strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")

	// WebSocketアップグレード用の設定
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			// 許可されたオリジンと照合
			for _, allowedOrigin := range allowedOrigins {
				if origin == strings.TrimSpace(allowedOrigin) {
					return true
				}
			}
			return false
		},
	}

	// クライアントを保持するためのマップとロック
	clients = make(map[*websocket.Conn]bool)
	mutex   = &sync.Mutex{}
)

// WebSocketハンドラー
func HandleWebSocket(c echo.Context) error {
	log.Println("WebSocket connection requested")

	// WebSocket接続をアップグレード
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Printf("Failed to upgrade to WebSocket: %v", err)
		return err
	}
	log.Println("WebSocket connection upgraded")
	defer ws.Close()

	// クライアントをマップに追加
	mutex.Lock()
	clients[ws] = true
	mutex.Unlock()

	log.Println("WebSocket connection established")

	// クライアントが切断されたときに、マップから削除する
	defer func() {
		mutex.Lock()
		delete(clients, ws)
		mutex.Unlock()
		log.Println("WebSocket connection closed")
	}()

	// クライアントからのメッセージを受信
	for {
		log.Println("Waiting for message...")

		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}
		log.Println("Message received")

		// メッセージのタイプで処理を分岐
		var messageData map[string]interface{}
		if err := json.Unmarshal(msg, &messageData); err != nil {
			log.Printf("Invalid message format: %v", err)
			continue
		}

		// メッセージタイプが存在するか確認
		messageType, ok := messageData["type"].(string)
		if !ok {
			log.Println("Message type missing or invalid")
			continue
		}

		log.Println("Message type:", messageType)
		// メッセージタイプによって処理を分岐
		switch messageType {
		case "debug":
			// デバッグメッセージの処理
			handleDebugMessage(messageData)
		case "reservation_notification":
			// 予約通知の処理
			handleReservationNotification(messageData)
		default:
			log.Printf("Unknown message type: %s", messageType)
		}
	}

	log.Println("WebSocket connection closed")
	return nil
}

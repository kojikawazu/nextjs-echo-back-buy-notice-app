package websocket

import (
	"context"
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
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}

		log.Printf("Received message: %s", msg)

		// Redisにメッセージをパブリッシュ
		err = rdb.Publish(ctx, "websocket-channel", msg).Err()
		if err != nil {
			log.Printf("Failed to publish message to Redis: %v", err)
		}
	}

	log.Println("WebSocket connection closed")
	return nil
}

// Redisからのメッセージをすべてのクライアントにブロードキャスト
func HandleMessages() {
	log.Println("Starting to broadcast messages from Redis")
	// Redis Pub/Sub をサブスクライブ
	pubsub := rdb.Subscribe(ctx, "websocket-channel")
	defer pubsub.Close()

	for {
		// Redisからのメッセージを受信
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			log.Printf("Failed to receive message from Redis: %v", err)
			continue
		}

		log.Printf("Broadcasting message from Redis: %s", msg.Payload)

		// すべてのWebSocketクライアントにメッセージを送信
		mutex.Lock()
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(msg.Payload))
			if err != nil {
				log.Printf("WebSocket write error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
		mutex.Unlock()
	}
}

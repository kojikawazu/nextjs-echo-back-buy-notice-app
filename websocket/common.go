package websocket

import (
	"context"
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

// Redisに通知をパブリッシュする関数
func PublishToRedis(channel, message string) error {
	log.Println("Publishing message to Redis")

	// Redisにメッセージをパブリッシュ
	err := rdb.Publish(context.Background(), channel, message).Err()
	if err != nil {
		log.Printf("Failed to publish message to Redis: %v", err)
		return err
	}

	log.Println("Published message to Redis successfully")
	return nil
}

// ブロードキャストメッセージ
func broadcastMessage(messageType string, content string) {
	log.Println("Broadcasting message to all clients")

	mutex.Lock()
	defer mutex.Unlock()

	// すべてのクライアントにメッセージを送信
	for client := range clients {
		log.Println("Sending message to client")

		msg := map[string]string{
			"type":    messageType,
			"content": content,
		}

		// メッセージをJSON形式に変換
		messageJSON, _ := json.Marshal(msg)
		// クライアントにメッセージを送信
		err := client.WriteMessage(websocket.TextMessage, messageJSON)
		if err != nil {
			// エラーが発生した場合、クライアントをクローズし、クライアントリストから削除
			log.Printf("WebSocket write error: %v", err)
			client.Close()
			delete(clients, client)
		}
	}

	log.Println("Broadcasted message to all clients")
}

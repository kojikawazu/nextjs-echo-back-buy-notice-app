package websocket

import (
	"log"
)

// Redisからのメッセージをすべてのクライアントにブロードキャスト
func HandleMessages() {
	log.Println("Starting to broadcast messages from Redis")

	// Redis Pub/Sub をサブスクライブ。複数チャンネルを指定。
	pubsub := rdb.Subscribe(ctx, "reservation-notifications", "debug-channel")
	defer pubsub.Close()

	for {
		log.Println("Waiting for messages from Redis")

		// Redisからのメッセージを受信
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			log.Printf("Failed to receive message from Redis: %v", err)
			continue
		}
		log.Printf("Received message from Redis: %s", msg.Payload)

		// チャンネルに応じて処理を分岐
		switch msg.Channel {
		case "reservation-notifications":
			log.Printf("Broadcasting reservation notification: %s", msg.Payload)
			// WebSocketクライアントにメッセージを送信
			broadcastMessage("reservation_notification", msg.Payload)
		case "debug-channel":
			log.Printf("Broadcasting debug message: %s", msg.Payload)
			// WebSocketクライアントにメッセージを送信
			broadcastMessage("debug", msg.Payload)
		default:
			log.Printf("Unknown channel: %s", msg.Channel)
		}

		log.Println("Message broadcasted to all clients")
	}
}

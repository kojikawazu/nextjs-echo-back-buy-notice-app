package websocket

import (
	"encoding/json"
	"log"
)

// 動作確認用のデバッグメッセージ処理
func handleDebugMessage(messageData map[string]interface{}) {
	log.Println("Handling debug message")
	message := messageData["message"].(string)
	log.Printf("Debug message: %s", message)

	// Redisにデバッグメッセージをパブリッシュ
	err := rdb.Publish(ctx, "debug-channel", message).Err()
	if err != nil {
		log.Printf("Failed to publish debug message to Redis: %v", err)
	}

	log.Println("Published debug message to Redis")
}

// 予約通知メッセージ処理
func handleReservationNotification(messageData map[string]interface{}) {
	log.Println("Handling reservation notification")

	// messageData を JSON に変換
	messageJSON, err := json.Marshal(messageData)
	if err != nil {
		log.Printf("Failed to marshal reservation notification: %v", err)
		return
	}

	// Redisに予約通知メッセージをパブリッシュ
	err = rdb.Publish(ctx, "reservation-channel", messageJSON).Err()
	if err != nil {
		log.Printf("Failed to publish reservation notification to Redis: %v", err)
	}

	log.Println("Published reservation notification to Redis")
}

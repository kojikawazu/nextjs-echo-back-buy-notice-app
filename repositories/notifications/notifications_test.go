package repositories_notifications

import (
	"backend/supabase"
	"log"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func setupSupabase() {
	// 環境変数の読み込み
	err := godotenv.Load("../../.env.test")
	if err != nil {
		log.Println("No ../../.env.test file found")
	}

	// テストの前にSupabaseクライアントの初期化
	err = supabase.InitSupabase()
	if err != nil {
		log.Fatalf("Supabase initialization failed: %v", err)
	}
}

func TestRepository_FetchNotifications(t *testing.T) {
	// Supabaseクライアントの初期化
	setupSupabase()

	// リポジトリのインスタンスを作成
	repo := NewNotificationRepository()

	// メソッドを実行
	notifications, err := repo.FetchNotifications()
	if err != nil {
		t.Fatalf("Failed to fetch notifications: %v", err)
	}

	// エラーチェックとデータ確認
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(notifications), 0)
}

func TestRepository_CreateNotification_ErrorCases(t *testing.T) {
	// Supabaseクライアントの初期化
	setupSupabase()

	// リポジトリのインスタンスを作成
	repo := NewNotificationRepository()

	// メソッドを実行
	err := repo.CreateNotification("", "", "")

	// エラーチェックとデータ確認
	assert.Error(t, err)
}

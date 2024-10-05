package repositories_reservations

import (
	"backend/supabase"
	"log"
	"os"
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

func TestRepository_FetchReservations(t *testing.T) {
	// Supabaseクライアントの初期化
	setupSupabase()

	// リポジトリのインスタンスを作成
	repo := NewReservationRepository()

	// メソッドを実行
	reservations, err := repo.FetchReservations()
	if err != nil {
		t.Fatalf("Failed to fetch reservations: %v", err)
	}

	// エラーチェックとデータ確認
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(reservations), 0)
}

func TestRepository_FetchReservationById(t *testing.T) {
	// Supabaseクライアントの初期化
	setupSupabase()

	// リポジトリのインスタンスを作成
	repo := NewReservationRepository()

	// メソッドを実行
	reservation, err := repo.FetchReservationById("eb0abf4d-82f7-44c8-b0c3-b33a5f680756")

	// エラーチェックとデータ確認
	assert.NoError(t, err)
	assert.NotNil(t, reservation)
}

func TestRepository_FetchReservationById_NoData(t *testing.T) {
	// Supabaseクライアントの初期化
	setupSupabase()

	// リポジトリのインスタンスを作成
	repo := NewReservationRepository()

	// 環境変数
	testID := os.Getenv("TEST_RESERVATION_ID")

	// メソッドを実行
	reservation, err := repo.FetchReservationById(testID)

	// エラーチェックとデータ確認
	assert.Error(t, err)
	assert.Nil(t, reservation)
}

func TestRepository_FetchReservationByUserId(t *testing.T) {
	// Supabaseクライアントの初期化
	setupSupabase()

	// リポジトリのインスタンスを作成
	repo := NewReservationRepository()

	// 環境変数
	testID := os.Getenv("TEST_USER_ID")

	// メソッドを実行
	reservation, err := repo.FetchReservationByUserId(testID)

	// エラーチェックとデータ確認
	assert.NoError(t, err)
	assert.NotNil(t, reservation)
}

func TestRepository_FetchReservationByUserId_NoData(t *testing.T) {
	// Supabaseクライアントの初期化
	setupSupabase()

	// リポジトリのインスタンスを作成
	repo := NewReservationRepository()

	// メソッドを実行
	reservation, err := repo.FetchReservationByUserId("99")

	// エラーチェックとデータ確認
	assert.Error(t, err)
	assert.Nil(t, reservation)
}

func TestRepository_CreateReservation(t *testing.T) {
	// Supabaseクライアントの初期化
	setupSupabase()

	// リポジトリのインスタンスを作成
	repo := NewReservationRepository()

	// メソッドを実行
	reservationId, err := repo.CreateReservation("", "", 0, "", "")

	// エラーチェックとデータ確認
	assert.Error(t, err)
	assert.Empty(t, reservationId)
}

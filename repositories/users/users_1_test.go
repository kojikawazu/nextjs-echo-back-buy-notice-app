package repositories_users

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

func TestRepository_FetchUsers(t *testing.T) {
	// Supabaseクライアントの初期化
	setupSupabase()

	// リポジトリのインスタンスを作成
	repo := NewUserRepository()

	// メソッドを実行
	users, err := repo.FetchUsers()
	if err != nil {
		t.Fatalf("Failed to fetch users: %v", err)
	}

	// エラーチェックとデータ確認
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(users), 0)
}

func TestRepository_FetchUserByEmailAndPassword(t *testing.T) {
	// Supabaseクライアントの初期化
	setupSupabase()

	// リポジトリのインスタンスを作成
	repo := NewUserRepository()

	// テスト用の環境変数を取得
	testName := os.Getenv("TEST_USER_NAME")
	testEmail := os.Getenv("TEST_USER_EMAIL")
	testPasswd := os.Getenv("TEST_USER_PASSWD")

	// メソッドを実行
	user, err := repo.FetchUserByEmailAndPassword(testEmail, testPasswd)
	if err != nil {
		t.Fatalf("Failed to fetch user: %v", err)
	}

	// エラーチェックとデータ確認
	assert.NoError(t, err)
	assert.Equal(t, testName, user.Name)
	assert.Equal(t, testEmail, user.Email)
}

func TestRepository_FetchUserByEmailAndPassword_ErrorCases(t *testing.T) {
	// Supabaseクライアントの初期化
	setupSupabase()

	// リポジトリのインスタンスを作成
	repo := NewUserRepository()

	// メソッドを実行
	user, err := repo.FetchUserByEmailAndPassword("", "")

	// エラーチェックとデータ確認
	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestRepository_FetchUserById(t *testing.T) {
	// Supabaseクライアントの初期化
	setupSupabase()

	// リポジトリのインスタンスを作成
	repo := NewUserRepository()

	// テスト用の環境変数を取得
	testUserId := os.Getenv("TEST_USER_ID")
	testName := os.Getenv("TEST_USER_NAME")
	testEmail := os.Getenv("TEST_USER_EMAIL")

	// メソッドを実行
	user, err := repo.FetchUserById(testUserId)
	if err != nil {
		t.Fatalf("Failed to fetch user: %v", err)
	}

	// エラーチェックとデータ確認
	assert.NoError(t, err)
	assert.Equal(t, testName, user.Name)
	assert.Equal(t, testEmail, user.Email)
}

func TestRepository_FetchUserById_ErrorCases(t *testing.T) {
	// Supabaseクライアントの初期化
	setupSupabase()

	// リポジトリのインスタンスを作成
	repo := NewUserRepository()

	// メソッドを実行
	user, err := repo.FetchUserById("")

	// エラーチェックとデータ確認
	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestRepository_FetchUserByEmail(t *testing.T) {
	// Supabaseクライアントの初期化
	setupSupabase()

	// リポジトリのインスタンスを作成
	repo := NewUserRepository()

	// テスト用の環境変数を取得
	testName := os.Getenv("TEST_USER_NAME")
	testEmail := os.Getenv("TEST_USER_EMAIL")

	// メソッドを実行
	user, err := repo.FetchUserByEmail(testEmail)
	if err != nil {
		t.Fatalf("Failed to fetch user: %v", err)
	}

	// エラーチェックとデータ確認
	assert.NoError(t, err)
	assert.Equal(t, testName, user.Name)
	assert.Equal(t, testEmail, user.Email)
}

func TestRepository_FetchUserByEmail_ErrorCases(t *testing.T) {
	// Supabaseクライアントの初期化
	setupSupabase()

	// リポジトリのインスタンスを作成
	repo := NewUserRepository()

	// メソッドを実行
	user, err := repo.FetchUserByEmail("")

	// エラーチェックとデータ確認
	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestRepository_CreateUser_ErrorCases(t *testing.T) {
	// Supabaseクライアントの初期化
	setupSupabase()

	// リポジトリのインスタンスを作成
	repo := NewUserRepository()

	// メソッドを実行
	err := repo.CreateUser("", "", "")

	// エラーチェックとデータ確認
	assert.Error(t, err)
}

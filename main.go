package main

import (
	"backend/auth"
	"backend/handlers"
	"backend/supabase"
	"backend/websocket"
	"strings"

	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// 環境変数の読み込み
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")

	// Supabaseクライアントの初期化
	err = supabase.InitSupabase()
	if err != nil {
		log.Fatalf("Supabase initialization failed: %v", err)
	}
	// テストクエリの実行
	err = supabase.TestQuery()
	if err != nil {
		log.Fatalf("Test query failed: %v", err)
	}

	e := echo.New()

	// ミドルウェアの設定
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// CORSを有効化
	// AllowCredentialsをtrueに設定すると、クライアント側でwithCredentialsをtrueに設定する必要がある
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     strings.Split(allowedOrigins, ","),
		AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))

	// APIエンドポイントの設定
	e.GET("/api/users", handlers.GetUsers)
	e.POST("/api/user", handlers.GetUserByEmailAndPassword)
	e.POST("/api/user/add", handlers.AddUser)

	e.GET("/api/reservations", handlers.GetReservations)
	e.GET("/api/reservations/:user_id", handlers.GetReservationByUserId)
	e.POST("/api/reservation", handlers.AddReservation)

	e.GET("/api/notifications", handlers.GetNotifications)
	e.POST("/api/notification", handlers.AddNotification)

	e.POST("/api/login", auth.Login)
	e.GET("/api/auth/check", auth.CheckAuth)

	// WebSocketエンドポイントの設定
	e.GET("/ws", websocket.HandleWebSocket)
	// メッセージをブロードキャストするためのゴルーチン
	go websocket.HandleMessages()

	// ヘルスチェックエンドポイントの追加
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Service is running")
	})

	// シグナルハンドラーの設定
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("Shutting down server...")

		// Echoサーバーのシャットダウン
		if err := e.Close(); err != nil {
			log.Printf("Echo shutdown failed: %v", err)
		}

		// Supabaseコネクションプールのクローズ
		supabase.ClosePool()
	}()

	// サーバーの起動
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Echo server failed: %v", err)
	}
}

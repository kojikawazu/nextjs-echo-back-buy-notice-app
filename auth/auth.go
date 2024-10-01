package auth

import (
	"backend/services"
	"backend/utils"
	"log"
	"net/http"
	"net/mail"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

var JwtKey = []byte(os.Getenv("JWT_SECRET_KEY")) // 環境変数から読み込む

func init() {
	if len(JwtKey) == 0 {
		log.Fatal("JWT_SECRET_KEY is not set in the environment")
	}
}

// ユーザー情報のペイロード
type Claims struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// ログインエンドポイント（JWTトークンの発行）
func Login(c echo.Context) error {
	utils.LogInfo(c, "Logging in...")

	// JSONのリクエストボディからemailとpasswordを取得
	type RequestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// リクエストボディをバインド
	var reqBody RequestBody
	if err := c.Bind(&reqBody); err != nil {
		utils.LogError(c, "Failed to bind request body: "+err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// バリデーション：emailとpasswordが空でないことを確認
	if reqBody.Email == "" || reqBody.Password == "" {
		utils.LogError(c, "Email and password are required")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Email and password are required",
		})
	}
	// バリデーション：emailが有効な形式であることを確認
	if _, err := mail.ParseAddress(reqBody.Email); err != nil {
		utils.LogError(c, "Invalid email format: "+err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid email format",
		})
	}
	utils.LogInfo(c, "Email and password are valid")

	// サービス層からユーザーデータを取得
	user, err := services.FetchUserByEmailAndPassword(reqBody.Email, reqBody.Password)
	if err != nil {
		utils.LogError(c, "Error fetching user: "+err.Error())
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "User not found",
		})
	}

	// 認証成功
	utils.LogInfo(c, "User authenticated successfully:"+user.Email)

	// JWTトークンの作成
	expirationTime := time.Now().Add(1 * time.Hour) // トークンの有効期限を1時間に設定
	claims := &Claims{
		UserID:   user.ID,
		Email:    user.Email, // 取得したユーザー情報を使う
		Username: user.Name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		utils.LogError(c, "Could not create JWT token: "+err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Could not create token"})
	}
	utils.LogInfo(c, "JWT token created successfully")

	// HTTP-onlyクッキーにトークンをセット
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = tokenString
	cookie.Expires = expirationTime
	cookie.HttpOnly = true
	c.SetCookie(cookie)

	utils.LogInfo(c, "JWT token set in HTTP-only cookie")
	return c.JSON(http.StatusOK, map[string]string{"message": "Login successful"})
}

// 認証確認エンドポイント
func CheckAuth(c echo.Context) error {
	utils.LogInfo(c, "Checking authentication...")

	// クッキーからJWTトークンを取得
	cookie, err := c.Cookie("token")
	if err != nil {
		utils.LogError(c, "Token not found in cookies")
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Token not found"})
	}
	tokenString := cookie.Value

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})

	if err != nil {
		utils.LogError(c, "Failed to parse token: "+err.Error())
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid token"})
	}

	if !token.Valid {
		utils.LogError(c, "Invalid token")
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid token"})
	}

	// 認証成功
	utils.LogInfo(c, "Authentication successful for user: "+claims.Email)
	return c.JSON(http.StatusOK, map[string]string{
		"message":  "Authenticated",
		"user_id":  claims.UserID,
		"username": claims.Username,
		"email":    claims.Email,
	})
}

// ログアウトエンドポイント
func Logout(c echo.Context) error {
	utils.LogInfo(c, "Logging out...")

	// クッキーを削除するために、空のトークンと過去の有効期限を設定
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = ""
	cookie.Expires = time.Unix(0, 0) // 有効期限を過去に設定して削除
	cookie.HttpOnly = true
	c.SetCookie(cookie)

	utils.LogInfo(c, "User logged out and token removed from cookie")
	return c.JSON(http.StatusOK, map[string]string{"message": "Logout successful"})
}

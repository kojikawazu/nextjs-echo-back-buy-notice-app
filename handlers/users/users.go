package handlers_users

import (
	services_users "backend/services/users"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	UserService services_users.UserService
}

// コンストラクタ
func NewUserHandler(userService services_users.UserService) *UserHandler {
	return &UserHandler{
		UserService: userService,
	}
}

// 全ユーザーを取得し、JSON形式で返すハンドラー
// ユーザー取得に失敗した場合、500エラーを返す。
func (h *UserHandler) GetUsers(c echo.Context) error {
	log.Println("Fetching users...")

	// サービス層でユーザー一覧を取得
	users, err := h.UserService.FetchUsers()
	if err != nil {
		log.Printf("Error fetching users from Supabase: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch users",
		})
	}

	log.Println("Fetched users successfully")
	return c.JSON(http.StatusOK, users)
}

// リクエストボディで指定されたemailとpasswordでユーザーを取得する。
// 有効なemailフォーマットかをチェックし、データベースに該当ユーザーがいない場合、404エラーを返す。
func (h *UserHandler) GetUserByEmailAndPassword(c echo.Context) error {
	log.Println("Fetching user by email and password...")

	// JSONのリクエストボディからemailとpasswordを取得
	type RequestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// リクエストボディをバインド
	var reqBody RequestBody
	if err := c.Bind(&reqBody); err != nil {
		log.Printf("Failed to bind request body: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// サービス層からユーザーデータを取得
	user, err := h.UserService.FetchUserByEmailAndPassword(reqBody.Email, reqBody.Password)
	if err != nil {
		switch err.Error() {
		case "email and password are required":
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Email and password are required",
			})
		case "invalid email format":
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid email format",
			})
		case "user not found":
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "User not found",
			})
		default:
			log.Printf("Error fetching user: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to fetch user",
			})
		}
	}

	log.Println("Fetched user successfully")
	return c.JSON(http.StatusOK, user)
}

// 新しいユーザーを追加するハンドラー
// ユーザーが既に存在する場合、409 Conflictを返し、存在しない場合は新規作成する。
func (h *UserHandler) AddUser(c echo.Context) error {
	log.Println("Creating new user...")

	// リクエストボディからデータを取得
	type RequestBody struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// リクエストボディをバインド
	var reqBody RequestBody
	if err := c.Bind(&reqBody); err != nil {
		log.Printf("Failed to bind request body: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// 新規ユーザーを作成
	err := h.UserService.CreateUser(reqBody.Name, reqBody.Email, reqBody.Password)
	if err != nil {
		switch err.Error() {
		case "name, email and password are required":
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Name, email and password are required",
			})
		case "invalid email format":
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid email format",
			})
		case "user already exists":
			return c.JSON(http.StatusConflict, map[string]string{
				"error": "User already exists",
			})
		case "error creating user":
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Error creating user",
			})
		default:
			log.Printf("Failed to create user: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to create user",
			})
		}
	}

	log.Println("User created successfully")
	return c.JSON(http.StatusCreated, map[string]string{
		"message": "User created successfully",
	})
}

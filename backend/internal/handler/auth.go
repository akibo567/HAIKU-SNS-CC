package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/goshichigo/backend/internal/middleware"
	"github.com/goshichigo/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userRepo         *repository.UserRepository
	jwtSecret        string
	jwtRefreshSecret string
}

func NewAuthHandler(userRepo *repository.UserRepository, jwtSecret, jwtRefreshSecret string) *AuthHandler {
	return &AuthHandler{
		userRepo:         userRepo,
		jwtSecret:        jwtSecret,
		jwtRefreshSecret: jwtRefreshSecret,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username    string `json:"username"`
		Email       string `json:"email"`
		Password    string `json:"password"`
		DisplayName string `json:"displayName"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "リクエストが不正です")
		return
	}

	if req.Username == "" || req.Email == "" || req.Password == "" || req.DisplayName == "" {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "必須フィールドが不足しています")
		return
	}

	existing, _ := h.userRepo.FindByUsername(r.Context(), req.Username)
	if existing != nil {
		writeError(w, http.StatusConflict, "CONFLICT", "このユーザー名は既に使われています")
		return
	}

	existingEmail, _ := h.userRepo.FindByEmail(r.Context(), req.Email)
	if existingEmail != nil {
		writeError(w, http.StatusConflict, "CONFLICT", "このメールアドレスは既に使われています")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "サーバーエラーが発生しました")
		return
	}

	user, err := h.userRepo.Create(r.Context(), req.Username, req.Email, string(hash), req.DisplayName)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "ユーザー作成に失敗しました")
		return
	}

	accessToken, refreshToken, err := h.generateTokens(user.ID, user.Username)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "トークン生成に失敗しました")
		return
	}

	h.setRefreshCookie(w, refreshToken)
	h.storeRefreshToken(r, user.ID, refreshToken)

	writeJSON(w, http.StatusCreated, map[string]any{
		"data": map[string]any{
			"accessToken": accessToken,
			"user":        userResponse(user),
		},
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "リクエストが不正です")
		return
	}

	user, err := h.userRepo.FindByEmail(r.Context(), req.Email)
	if err != nil || user == nil {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "メールアドレスまたはパスワードが正しくありません")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "メールアドレスまたはパスワードが正しくありません")
		return
	}

	accessToken, refreshToken, err := h.generateTokens(user.ID, user.Username)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "トークン生成に失敗しました")
		return
	}

	h.setRefreshCookie(w, refreshToken)
	h.storeRefreshToken(r, user.ID, refreshToken)

	writeJSON(w, http.StatusOK, map[string]any{
		"data": map[string]any{
			"accessToken": accessToken,
			"user":        userResponse(user),
		},
	})
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "リフレッシュトークンがありません")
		return
	}

	tokenHash := hashToken(cookie.Value)
	userID, err := h.userRepo.FindRefreshToken(r.Context(), tokenHash)
	if err != nil || userID == "" {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "リフレッシュトークンが無効です")
		return
	}

	user, err := h.userRepo.FindByID(r.Context(), userID)
	if err != nil || user == nil {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "ユーザーが見つかりません")
		return
	}

	// rotate token
	h.userRepo.DeleteRefreshToken(r.Context(), tokenHash)
	accessToken, refreshToken, err := h.generateTokens(user.ID, user.Username)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "トークン生成に失敗しました")
		return
	}

	h.setRefreshCookie(w, refreshToken)
	h.storeRefreshToken(r, user.ID, refreshToken)

	writeJSON(w, http.StatusOK, map[string]any{
		"data": map[string]any{
			"accessToken": accessToken,
		},
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie("refresh_token"); err == nil {
		tokenHash := hashToken(cookie.Value)
		h.userRepo.DeleteRefreshToken(r.Context(), tokenHash)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
		Path:     "/",
	})

	writeJSON(w, http.StatusOK, map[string]any{"data": map[string]string{"message": "ログアウトしました"}})
}

func (h *AuthHandler) generateTokens(userID, username string) (string, string, error) {
	now := time.Now()

	accessClaims := jwt.MapClaims{
		"sub":      userID,
		"username": username,
		"exp":      now.Add(15 * time.Minute).Unix(),
		"iat":      now.Unix(),
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(h.jwtSecret))
	if err != nil {
		return "", "", err
	}

	refreshClaims := jwt.MapClaims{
		"sub": userID,
		"exp": now.Add(7 * 24 * time.Hour).Unix(),
		"iat": now.Unix(),
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(h.jwtRefreshSecret))
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (h *AuthHandler) setRefreshCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    token,
		HttpOnly: true,
		Path:     "/api/v1/auth",
		MaxAge:   7 * 24 * 60 * 60,
		SameSite: http.SameSiteLaxMode,
	})
}

func (h *AuthHandler) storeRefreshToken(r *http.Request, userID, token string) {
	hash := hashToken(token)
	expiresAt := time.Now().Add(7 * 24 * time.Hour).Format(time.RFC3339)
	h.userRepo.StoreRefreshToken(r.Context(), userID, hash, expiresAt)
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func userResponse(u *repository.User) map[string]any {
	bio := ""
	if u.Bio != nil {
		bio = *u.Bio
	}
	return map[string]any{
		"id":          u.ID,
		"username":    u.Username,
		"displayName": u.DisplayName,
		"bio":         bio,
		"createdAt":   u.CreatedAt,
	}
}

func GetUserIDFromContext(r *http.Request) string {
	return middleware.GetUserID(r)
}

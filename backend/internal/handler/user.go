package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/goshichigo/backend/internal/middleware"
	"github.com/goshichigo/backend/internal/repository"
)

type UserHandler struct {
	userRepo  *repository.UserRepository
	haikuRepo *repository.HaikuRepository
}

func NewUserHandler(userRepo *repository.UserRepository, haikuRepo *repository.HaikuRepository) *UserHandler {
	return &UserHandler{userRepo: userRepo, haikuRepo: haikuRepo}
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	user, err := h.userRepo.FindByUsername(r.Context(), username)
	if err != nil || user == nil {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "ユーザーが見つかりません")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": userResponse(user)})
}

func (h *UserHandler) GetPosts(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	user, err := h.userRepo.FindByUsername(r.Context(), username)
	if err != nil || user == nil {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "ユーザーが見つかりません")
		return
	}

	cursor := r.URL.Query().Get("cursor")
	limit := 20
	if l := r.URL.Query().Get("limit"); l != "" {
		if n, err := strconv.Atoi(l); err == nil && n > 0 && n <= 50 {
			limit = n
		}
	}

	posts, err := h.haikuRepo.ListByUserID(r.Context(), user.ID, cursor, limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "投稿の取得に失敗しました")
		return
	}

	currentUserID := middleware.GetUserID(r)
	likedMap := map[string]bool{}
	if currentUserID != "" && len(posts) > 0 {
		ids := make([]string, len(posts))
		for i, p := range posts {
			ids[i] = p.ID
		}
		likedMap, _ = h.haikuRepo.LikedPostIDs(r.Context(), currentUserID, ids)
	}

	var nextCursor string
	if len(posts) == limit {
		nextCursor = posts[len(posts)-1].ID
	}

	items := make([]map[string]any, len(posts))
	for i, p := range posts {
		items[i] = haikuResponse(p, likedMap[p.ID])
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"data": items,
		"meta": map[string]any{"cursor": nextCursor, "hasNext": nextCursor != ""},
	})
}

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	var req struct {
		DisplayName string  `json:"displayName"`
		Bio         *string `json:"bio"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "リクエストが不正です")
		return
	}

	if req.DisplayName == "" {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "表示名は必須です")
		return
	}

	user, err := h.userRepo.UpdateProfile(r.Context(), userID, req.DisplayName, req.Bio)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "プロフィール更新に失敗しました")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"data": userResponse(user)})
}

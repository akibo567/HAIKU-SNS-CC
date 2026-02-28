package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/goshichigo/backend/internal/middleware"
	"github.com/goshichigo/backend/internal/mora"
	"github.com/goshichigo/backend/internal/repository"
)

type HaikuHandler struct {
	haikuRepo *repository.HaikuRepository
	replyRepo *repository.ReplyRepository
}

func NewHaikuHandler(haikuRepo *repository.HaikuRepository, replyRepo *repository.ReplyRepository) *HaikuHandler {
	return &HaikuHandler{haikuRepo: haikuRepo, replyRepo: replyRepo}
}

func (h *HaikuHandler) ListTimeline(w http.ResponseWriter, r *http.Request) {
	cursor := r.URL.Query().Get("cursor")
	limit := 20
	if l := r.URL.Query().Get("limit"); l != "" {
		if n, err := strconv.Atoi(l); err == nil && n > 0 && n <= 50 {
			limit = n
		}
	}

	posts, err := h.haikuRepo.List(r.Context(), cursor, limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "タイムラインの取得に失敗しました")
		return
	}

	userID := middleware.GetUserID(r)
	likedMap := map[string]bool{}
	if userID != "" && len(posts) > 0 {
		ids := make([]string, len(posts))
		for i, p := range posts {
			ids[i] = p.ID
		}
		likedMap, _ = h.haikuRepo.LikedPostIDs(r.Context(), userID, ids)
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
		"meta": map[string]any{
			"cursor":  nextCursor,
			"hasNext": nextCursor != "",
		},
	})
}

func (h *HaikuHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	var req struct {
		Ku1 string `json:"ku1"`
		Ku2 string `json:"ku2"`
		Ku3 string `json:"ku3"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "リクエストが不正です")
		return
	}

	if err := mora.ValidateHaiku(req.Ku1, req.Ku2, req.Ku3); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_MORA_COUNT", err.Error())
		return
	}

	post, err := h.haikuRepo.Create(r.Context(), userID, req.Ku1, req.Ku2, req.Ku3)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "俳句の投稿に失敗しました")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"data": haikuResponse(*post, false),
	})
}

func (h *HaikuHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	post, err := h.haikuRepo.FindByID(r.Context(), id)
	if err != nil || post == nil {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "俳句が見つかりません")
		return
	}

	userID := middleware.GetUserID(r)
	liked := false
	if userID != "" {
		liked, _ = h.haikuRepo.IsLikedByUser(r.Context(), userID, id)
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"data": haikuResponse(*post, liked),
	})
}

func (h *HaikuHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID := middleware.GetUserID(r)

	deleted, err := h.haikuRepo.Delete(r.Context(), id, userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "削除に失敗しました")
		return
	}
	if !deleted {
		writeError(w, http.StatusForbidden, "FORBIDDEN", "この俳句を削除する権限がありません")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"data": map[string]string{"message": "削除しました"}})
}

func (h *HaikuHandler) Like(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID := middleware.GetUserID(r)

	if err := h.haikuRepo.AddLike(r.Context(), userID, id); err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "いいねに失敗しました")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"data": map[string]string{"message": "いいねしました"}})
}

func (h *HaikuHandler) Unlike(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID := middleware.GetUserID(r)

	if err := h.haikuRepo.RemoveLike(r.Context(), userID, id); err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "いいね取消に失敗しました")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"data": map[string]string{"message": "いいねを取り消しました"}})
}

func (h *HaikuHandler) ListReplies(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "id")

	replies, err := h.replyRepo.ListByPostID(r.Context(), postID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "返句の取得に失敗しました")
		return
	}

	items := make([]map[string]any, len(replies))
	for i, rep := range replies {
		items[i] = replyResponse(rep)
	}

	writeJSON(w, http.StatusOK, map[string]any{"data": items})
}

func (h *HaikuHandler) CreateReply(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "id")
	userID := middleware.GetUserID(r)

	// 親投稿の存在確認
	post, err := h.haikuRepo.FindByID(r.Context(), postID)
	if err != nil || post == nil {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "俳句が見つかりません")
		return
	}

	var req struct {
		Ku1 string `json:"ku1"`
		Ku2 string `json:"ku2"`
		Ku3 string `json:"ku3"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "リクエストが不正です")
		return
	}

	if err := mora.ValidateHaiku(req.Ku1, req.Ku2, req.Ku3); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_MORA_COUNT", err.Error())
		return
	}

	reply, err := h.replyRepo.Create(r.Context(), postID, userID, req.Ku1, req.Ku2, req.Ku3)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "返句の投稿に失敗しました")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{"data": replyResponse(*reply)})
}

func replyResponse(rep repository.Reply) map[string]any {
	return map[string]any{
		"id":        rep.ID,
		"postId":    rep.PostID,
		"ku1":       rep.Ku1,
		"ku2":       rep.Ku2,
		"ku3":       rep.Ku3,
		"createdAt": rep.CreatedAt,
		"author": map[string]any{
			"id":          rep.UserID,
			"username":    rep.Username,
			"displayName": rep.DisplayName,
		},
	}
}

func haikuResponse(p repository.HaikuPost, likedByMe bool) map[string]any {
	return map[string]any{
		"id":        p.ID,
		"ku1":       p.Ku1,
		"ku2":       p.Ku2,
		"ku3":       p.Ku3,
		"likeCount": p.LikeCount,
		"likedByMe": likedByMe,
		"createdAt": p.CreatedAt,
		"author": map[string]any{
			"id":          p.UserID,
			"username":    p.Username,
			"displayName": p.DisplayName,
		},
	}
}

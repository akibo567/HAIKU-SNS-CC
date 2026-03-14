package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/goshichigo/backend/internal/middleware"
	"github.com/goshichigo/backend/internal/repository"
)

type stubHaikuRepo struct {
	createCalled   bool
	findByIDResult *repository.HaikuPost
}

func (s *stubHaikuRepo) List(ctx context.Context, cursor string, limit int) ([]repository.HaikuPost, error) {
	return nil, nil
}

func (s *stubHaikuRepo) Create(ctx context.Context, userID, ku1, ku2, ku3 string) (*repository.HaikuPost, error) {
	s.createCalled = true
	return &repository.HaikuPost{
		ID:          "post-1",
		UserID:      userID,
		Username:    "user",
		DisplayName: "User",
		Ku1:         ku1,
		Ku2:         ku2,
		Ku3:         ku3,
		CreatedAt:   time.Now(),
	}, nil
}

func (s *stubHaikuRepo) FindByID(ctx context.Context, id string) (*repository.HaikuPost, error) {
	return s.findByIDResult, nil
}

func (s *stubHaikuRepo) Delete(ctx context.Context, id, userID string) (bool, error) {
	return false, nil
}

func (s *stubHaikuRepo) AddLike(ctx context.Context, userID, postID string) error {
	return nil
}

func (s *stubHaikuRepo) RemoveLike(ctx context.Context, userID, postID string) error {
	return nil
}

func (s *stubHaikuRepo) IsLikedByUser(ctx context.Context, userID, postID string) (bool, error) {
	return false, nil
}

func (s *stubHaikuRepo) LikedPostIDs(ctx context.Context, userID string, postIDs []string) (map[string]bool, error) {
	return map[string]bool{}, nil
}

type stubReplyRepo struct {
	createCalled bool
}

func (s *stubReplyRepo) Create(ctx context.Context, postID, userID, ku1, ku2, ku3 string) (*repository.Reply, error) {
	s.createCalled = true
	return &repository.Reply{
		ID:          "reply-1",
		PostID:      postID,
		UserID:      userID,
		Username:    "user",
		DisplayName: "User",
		Ku1:         ku1,
		Ku2:         ku2,
		Ku3:         ku3,
		CreatedAt:   time.Now(),
	}, nil
}

func (s *stubReplyRepo) ListByPostID(ctx context.Context, postID string) ([]repository.Reply, error) {
	return nil, nil
}

func TestCreateRejectsInvalidHaikuBeforeRepository(t *testing.T) {
	haikuRepo := &stubHaikuRepo{}
	replyRepo := &stubReplyRepo{}
	h := NewHaikuHandler(haikuRepo, replyRepo)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/posts", strings.NewReader(`{"ku1":"あいうえ","ku2":"あいうえおかき","ku3":"あいうえお"}`))
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, "user-1"))
	rec := httptest.NewRecorder()

	h.Create(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
	if haikuRepo.createCalled {
		t.Fatal("expected repository Create not to be called")
	}

	var body struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Error.Code != "INVALID_MORA_COUNT" {
		t.Fatalf("error code = %q, want %q", body.Error.Code, "INVALID_MORA_COUNT")
	}
}

func TestCreateReplyRejectsInvalidHaikuBeforeRepository(t *testing.T) {
	haikuRepo := &stubHaikuRepo{
		findByIDResult: &repository.HaikuPost{
			ID:          "post-1",
			UserID:      "author-1",
			Username:    "author",
			DisplayName: "Author",
			Ku1:         "あいうえお",
			Ku2:         "あいうえおかき",
			Ku3:         "あいうえお",
			CreatedAt:   time.Now(),
		},
	}
	replyRepo := &stubReplyRepo{}
	h := NewHaikuHandler(haikuRepo, replyRepo)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/posts/post-1/replies", strings.NewReader(`{"ku1":"あいうえ","ku2":"あいうえおかき","ku3":"あいうえお"}`))
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, "user-1"))
	req = withURLParam(req, "id", "post-1")
	rec := httptest.NewRecorder()

	h.CreateReply(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
	if replyRepo.createCalled {
		t.Fatal("expected repository Create not to be called")
	}

	var body struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Error.Code != "INVALID_MORA_COUNT" {
		t.Fatalf("error code = %q, want %q", body.Error.Code, "INVALID_MORA_COUNT")
	}
}

func withURLParam(req *http.Request, key, value string) *http.Request {
	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add(key, value)
	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
}

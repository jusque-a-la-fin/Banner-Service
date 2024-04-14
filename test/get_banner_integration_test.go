// Тест на получение баннера
package integration_test

import (
	"banner/configs"
	"banner/internal/banner"
	"banner/internal/datastore"
	"banner/internal/handlers"
	"banner/internal/session"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
)

// Тест корректного ответа: 200
func TestGetBannerOK(t *testing.T) {
	if err := configs.Init(); err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/user_banner?tag_id=2&feature_id=1&&use_last_revision=false", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("token", "user_token")

	sess := &session.Session{
		ParticipantToken: "user_token",
		IsAdmin:          false,
	}

	ctx := session.ContextWithSession(req.Context(), sess)
	req = req.WithContext(ctx)

	mdtb, err := datastore.CreateNewMainDB()
	if err != nil {
		log.Fatalf("error connecting to banner_service database: %#v", err)
	}
	rdc := datastore.CreateNewRedisClient()

	rtr := setupRouter(mdtb, rdc)

	rr := httptest.NewRecorder()

	rtr.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d; got %d", http.StatusOK, rr.Code)
	}

	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
	}

	var body map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Errorf("Expected JSON response, got error: %s", err)
	}
}

func setupRouter(db *sql.DB, rd *redis.Client) http.Handler {
	bannerHandler := &handlers.BannerHandler{
		BannerRepo: banner.NewDBRepo(db, rd),
	}

	rtr := mux.NewRouter()
	rtr.HandleFunc("/user_banner", bannerHandler.GetBanner)
	return rtr
}

// Тест некорректных данных: 400
func TestGetBannerBadRequest(t *testing.T) {
	if err := configs.Init(); err != nil {
		t.Fatal(err)
	}

	mdtb, err := datastore.CreateNewMainDB()
	if err != nil {
		log.Fatalf("error connecting to banner_service database: %#v", err)
	}
	rdc := datastore.CreateNewRedisClient()

	rtr := setupRouter(mdtb, rdc)

	req, err := http.NewRequest("GET", "/user_banner?tag_id=2&feature_id=#&use_last_revision=1", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("token", "user_token")

	sess := &session.Session{
		ParticipantToken: "user_token",
		IsAdmin:          false,
	}

	ctx := session.ContextWithSession(req.Context(), sess)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	rtr.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d; got %d", http.StatusBadRequest, rr.Code)
	}

	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
	}

	var body map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Errorf("Expected JSON response, got error: %s", err)
	}
}

// Тест: пользователь не авторизован: 401
func TestGetBannerUnauthorized(t *testing.T) {
	if err := configs.Init(); err != nil {
		t.Fatal(err)
	}

	tdtb, err := datastore.CreateNewTokenDB()
	if err != nil {
		log.Fatalf("error connecting to token_storage database: %#v", err)
	}
	snm := session.NewSessionsManager(tdtb)

	req := httptest.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	req.Header.Set("token", "user_token1")

	w := httptest.NewRecorder()

	snm.CheckToken(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d; got %d", http.StatusUnauthorized, w.Code)
	}
}

// Тест: пользователь не имеет доступа, 403
func TestGetBannerForbidden(t *testing.T) {
	if err := configs.Init(); err != nil {
		t.Fatal(err)
	}

	mdtb, err := datastore.CreateNewMainDB()
	if err != nil {
		log.Fatalf("error connecting to banner_service database: %#v", err)
	}
	rdc := datastore.CreateNewRedisClient()

	rtr := setupRouter(mdtb, rdc)

	req, err := http.NewRequest("GET", "/user_banner?tag_id=3&feature_id=7&&use_last_revision=true", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("token", "user_token")

	sess := &session.Session{
		ParticipantToken: "user_token",
		IsAdmin:          false,
	}

	ctx := session.ContextWithSession(req.Context(), sess)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	rtr.ServeHTTP(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Errorf("Expected status code %d; got %d", http.StatusForbidden, rr.Code)
	}
}

// Тест: баннер не найден, 404
func TestGetBannerNotFound(t *testing.T) {
	if err := configs.Init(); err != nil {
		t.Fatal(err)
	}

	mdtb, err := datastore.CreateNewMainDB()
	if err != nil {
		log.Fatalf("error connecting to banner_service database: %#v", err)
	}
	rdc := datastore.CreateNewRedisClient()

	rtr := setupRouter(mdtb, rdc)

	req, err := http.NewRequest("GET", "/user_banner?feature_id=1&tag_id=20&use_last_revision=true", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("token", "user_token")

	sess := &session.Session{
		ParticipantToken: "user_token",
		IsAdmin:          false,
	}

	ctx := session.ContextWithSession(req.Context(), sess)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	rtr.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d; got %d", http.StatusNotFound, rr.Code)
	}
}

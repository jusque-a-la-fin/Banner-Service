package integration_test

import (
	"banner/configs"
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

// Тестирование создания нового баннера
func TestCreateNewBannerOK(t *testing.T) {
	if err := configs.Init(); err != nil {
		t.Fatal(err)
	}

	requestBody := []byte(`{"tag_ids": [2, 3, 5], "feature_id": 3, "content": {"title": "title_14", "text": "text14", "url": "url16"}, "is_active": true}`)

	req, err := http.NewRequest("POST", "http://localhost:8080/banner", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("token", "admin_token")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	// если он баннер уже был добавлен, то вернется ошибка 500
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status code %d; got %d", http.StatusCreated, resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
	}

	var body map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Errorf("Failed to parse JSON response body: %s", err)
	}
}

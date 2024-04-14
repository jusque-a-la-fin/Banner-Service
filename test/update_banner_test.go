package integration_test

import (
	"banner/configs"
	"bytes"
	"net/http"
	"testing"
)

// Тестирование обновления баннера
func TestUpdateBannerOK(t *testing.T) {
	if err := configs.Init(); err != nil {
		t.Fatal(err)
	}

	requestBody := []byte(`{"tag_ids": [2, 3], "feature_id": 3, "content": {"title": "sotle_17", "text": "sotext18", "url": "sourl19"}, "is_active": true}`)

	req, err := http.NewRequest("PATCH", "http://localhost:8080/banner/1", bytes.NewBuffer(requestBody))
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

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d; got %d", http.StatusOK, resp.StatusCode)
	}
}

package integration_test

import (
	"banner/configs"
	"bytes"
	"net/http"
	"testing"
)

// Тестирование удаления баннера
func TestDeleteBannerOK(t *testing.T) {
	if err := configs.Init(); err != nil {
		t.Fatal(err)
	}

	var requestBody []byte

	req, err := http.NewRequest("DELETE", "http://localhost:8080/banner/4", bytes.NewBuffer(requestBody))
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

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code %d; got %d", http.StatusOK, resp.StatusCode)
	}
}

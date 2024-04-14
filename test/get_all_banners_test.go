package integration_test

import (
	"banner/configs"
	"net/http"
	"testing"
)

func TestGetAllBannersOK(t *testing.T) {
	if err := configs.Init(); err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", "http://localhost:8080/banner?feature_id=3&limit=5&offset=0", nil)
	if err != nil {
		t.Fatal(err)
	}

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

	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
	}

}

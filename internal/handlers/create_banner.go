package handlers

import (
	"banner/internal/banner"
	"encoding/json"
	"log"
	"net/http"
)

type BannerRequest struct {
	TagIDs    []int                  `json:"tag_ids"`
	FeatureID *int                   `json:"feature_id"`
	Content   map[string]interface{} `json:"content"`
	IsActive  *bool                  `json:"is_active"`
}

// CreateNewBanner создает новый баннер
func (hnd *BannerHandler) CreateNewBanner(wrt http.ResponseWriter, rqt *http.Request) {
	isAdmin := checkIfAdmin(wrt, rqt)
	if !isAdmin {
		return
	}

	ban := obtainBanner(wrt, rqt)
	if ban == nil {
		return
	}

	banID, err := hnd.BannerRepo.InsertNewBannerIntoDB(*ban)
	if err != nil {
		log.Printf("error inserting new banner into database: %#v", err)
		err := errorResponse{Error: "Внутренняя ошибка сервера"}
		SendJSON(wrt, err, http.StatusInternalServerError)
		return
	}

	response := struct {
		BannerID int `json:"banner_id"`
	}{
		BannerID: *banID,
	}
	SendJSON(wrt, response, http.StatusCreated)
}

// obtainBanner получает баннер из тела запроса
func obtainBanner(wrt http.ResponseWriter, rqt *http.Request) *banner.Banner {
	var banReq BannerRequest
	err := json.NewDecoder(rqt.Body).Decode(&banReq)
	if err != nil {
		log.Printf("error decoding json from request body: %#v", err)
		err := errorResponse{Error: "Некорректные данные"}
		SendJSON(wrt, err, http.StatusBadRequest)
		return nil
	}

	ban := banner.Banner{}
	ban.TagsIDs = banReq.TagIDs
	ban.FeatureID = banReq.FeatureID
	if banReq.Content != nil {
		ban.Content, err = json.Marshal(banReq.Content)
		if err != nil {
			log.Printf("error marshaling: %#v", err)
			err := errorResponse{Error: "Внутренняя ошибка сервера"}
			SendJSON(wrt, err, http.StatusInternalServerError)
			return nil
		}
	}
	ban.IsActive = banReq.IsActive
	return &ban
}

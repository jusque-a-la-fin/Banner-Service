package handlers

import (
	"banner/internal/session"
	"log"
	"net/http"
)

// GetAllBanners получает все баннеры c фильтрацией по фиче и/или тегу
func (hnd *BannerHandler) GetAllBanners(wrt http.ResponseWriter, rqt *http.Request) {
	isAdmin := checkIfAdmin(wrt, rqt)
	if !isAdmin {
		return
	}

	featureID := rqt.URL.Query().Get("feature_id")
	tagID := rqt.URL.Query().Get("tag_id")
	offset := rqt.URL.Query().Get("offset")
	limit := rqt.URL.Query().Get("limit")
	bans, err := hnd.BannerRepo.GetAllBannersFromDB(featureID, tagID, offset, limit)
	if err != nil {
		log.Printf("error getting all banners from database: %#v", err)
		err := errorResponse{Error: "Внутренняя ошибка сервера"}
		SendJSON(wrt, err, http.StatusInternalServerError)
		return
	}
	SendJSON(wrt, bans, http.StatusOK)
}

// checkIfAdmin проверяет, является ли пользователь админом
func checkIfAdmin(wrt http.ResponseWriter, rqt *http.Request) bool {
	sess, err := session.SessionFromContext(rqt.Context())
	if err != nil || sess == nil {
		log.Printf("error getting session info: %#v", err)
		err := errorResponse{Error: "Внутренняя ошибка сервера"}
		SendJSON(wrt, err, http.StatusInternalServerError)
		return false
	}

	if !sess.IsAdmin {
		log.Printf("%s", errNoAccess)
		http.Error(wrt, "Пользователь не имеет доступа", http.StatusForbidden)
		return false
	}
	return true
}

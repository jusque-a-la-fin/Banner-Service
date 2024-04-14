package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

const bannerNotFound = "error: banner with this id hasn't been found"

// UpdateBanner обновляет содержимое баннера
func (hnd *BannerHandler) UpdateBanner(wrt http.ResponseWriter, rqt *http.Request) {
	isAdmin := checkIfAdmin(wrt, rqt)
	if !isAdmin {
		return
	}
	vars := mux.Vars(rqt)
	bannerID := vars["id"]
	banID, err := strconv.Atoi(bannerID)
	if err != nil {
		log.Printf("error converting string to int: %#v", err)
		err := errorResponse{Error: "Внутренняя ошибка сервера"}
		SendJSON(wrt, err, http.StatusInternalServerError)
		return
	}

	ban := obtainBanner(wrt, rqt)
	if ban == nil {
		return
	}

	err = hnd.BannerRepo.UpdateBannerInDB(banID, *ban)
	if err != nil && err.Error() != bannerNotFound {
		log.Printf("%#v", err)
		err := errorResponse{Error: "Внутренняя ошибка сервера"}
		SendJSON(wrt, err, http.StatusInternalServerError)
		return
	}

	if err != nil && err.Error() == bannerNotFound {
		log.Printf("%#v", err)
		http.Error(wrt, "Баннер не найден", http.StatusNotFound)
		return
	}
	wrt.WriteHeader(http.StatusOK)
}

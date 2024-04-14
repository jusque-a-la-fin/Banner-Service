package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// DeleteBanner удаляет баннер
func (hnd *BannerHandler) DeleteBanner(wrt http.ResponseWriter, rqt *http.Request) {
	isAdmin := checkIfAdmin(wrt, rqt)
	if !isAdmin {
		return
	}
	vars := mux.Vars(rqt)
	bannerID := vars["id"]
	banID, err := strconv.Atoi(bannerID)
	if err != nil {
		log.Printf("error converting string to int: %#v", err)
		err := errorResponse{Error: "Некорректные данные"}
		SendJSON(wrt, err, http.StatusBadRequest)
		return
	}

	err = hnd.BannerRepo.DeleteBannerFromDB(banID)
	if err != nil && err.Error() != bannerNotFound {
		log.Printf("%#v", err)
		err := errorResponse{Error: "Внутренняя ошибка сервера"}
		SendJSON(wrt, err, http.StatusInternalServerError)
		return
	}

	if err != nil && err.Error() == bannerNotFound {
		log.Printf("%#v", err)
		http.Error(wrt, "Баннер для тэга не найден", http.StatusNotFound)
		return
	}

	wrt.WriteHeader(http.StatusNoContent)
}

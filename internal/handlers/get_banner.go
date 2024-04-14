package handlers

import (
	"banner/internal/banner"
	"banner/internal/session"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

const (
	errUserNotFound = "error: value for this user token hasn't been found"
	errOutdated     = "banner info is out of date"
	errNotFound     = "error: banner hasn't been found"
	errNoAccess     = "error: user doesn't have access"
)

type BannerHandler struct {
	BannerRepo banner.BannerRepo
}

type errorResponse struct {
	Error string `json:"error"`
}

// GetBanner получает баннер для пользователя
func (hnd *BannerHandler) GetBanner(wrt http.ResponseWriter, rqt *http.Request) {
	featureID := rqt.URL.Query().Get("feature_id")
	tagID := rqt.URL.Query().Get("tag_id")
	useLastRevision := rqt.URL.Query().Get("use_last_revision")
	token := rqt.Header.Get("token")
	isValid := validateUserParams(token, featureID, tagID, useLastRevision)
	if !isValid {
		log.Printf("error: incorrect request parameters")
		err := errorResponse{Error: "Некорректные данные"}
		SendJSON(wrt, err, http.StatusBadRequest)
		return
	}

	sess, err := session.SessionFromContext(rqt.Context())
	if err != nil || sess == nil {
		log.Printf("error getting session info: %#v", err)
		err := errorResponse{Error: "Внутренняя ошибка сервера"}
		SendJSON(wrt, err, http.StatusInternalServerError)
		return
	}

	uLR, err := strconv.ParseBool(useLastRevision)
	if err != nil {
		log.Printf("error converting string to bool: %#v", err)
		err := errorResponse{Error: "Внутренняя ошибка сервера"}
		SendJSON(wrt, err, http.StatusInternalServerError)
		return
	}

	var ban *banner.Banner
	if !uLR {
		ban, err = hnd.BannerRepo.GetBannerFromCache(rqt.Context(), token, featureID, tagID, sess.IsAdmin)
	}

	if err != nil && err.Error() != errUserNotFound && err.Error() != errOutdated && err.Error() != errNotFound {
		log.Printf("%#v", err)
		err := errorResponse{Error: "Внутренняя ошибка сервера"}
		SendJSON(wrt, err, http.StatusInternalServerError)
		return
	}

	if ban == nil {
		ban, err = hnd.BannerRepo.GetBannerFromDB(featureID, tagID, sess.IsAdmin)
		if err != nil {
			switch err.Error() {
			case errNotFound:
				log.Printf("%s", err)
				http.Error(wrt, "Баннер для не найден", http.StatusNotFound)
			case errNoAccess:
				log.Printf("%s", err)
				http.Error(wrt, "Пользователь не имеет доступа", http.StatusForbidden)
			default:
				log.Printf("%#v", err)
				err := errorResponse{Error: "Внутренняя ошибка сервера"}
				SendJSON(wrt, err, http.StatusInternalServerError)
				return
			}
			return
		}

		err = hnd.BannerRepo.SetBannerInCache(rqt.Context(), *ban, token, featureID, tagID)
		if err != nil {
			log.Printf("%#v", err)
			err := errorResponse{Error: "Внутренняя ошибка сервера"}
			SendJSON(wrt, err, http.StatusInternalServerError)
			return
		}
	}

	wrt.Header().Set("Content-Type", "application/json")
	wrt.WriteHeader(http.StatusOK)
	if _, err := wrt.Write(ban.Content); err != nil {
		log.Printf("error writing response: %#v", err)
		err := errorResponse{Error: "Внутренняя ошибка сервера"}
		SendJSON(wrt, err, http.StatusInternalServerError)
		return
	}
}

// validateUserParams проверяет корректность параметров, переданных пользователем
func validateUserParams(token, featureID, tagID, useLastRevision string) bool {
	if token == "" {
		return false
	}
	_, err := strconv.Atoi(featureID)
	if err != nil {
		return false
	}

	_, err = strconv.Atoi(tagID)
	if err != nil {
		return false
	}

	if useLastRevision != "true" && useLastRevision != "false" {
		return false
	}
	return true
}

// SendJSON сериализует данные и посылает JSON
func SendJSON(wrt http.ResponseWriter, data interface{}, statusCode int) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("error marshaling JSON: %#v", err)
		http.Error(wrt, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	wrt.Header().Set("Content-Type", "application/json")

	wrt.WriteHeader(statusCode)
	if _, err := wrt.Write(jsonData); err != nil {
		log.Printf("error writing response: %#v", err)
		http.Error(wrt, "Внутренняя ошибка сервера", http.StatusInternalServerError)
	}
}

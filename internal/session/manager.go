package session

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type SessionsManager struct {
	sessionsDB *sql.DB
}

func NewSessionsManager(sessDB *sql.DB) *SessionsManager {
	return &SessionsManager{
		sessionsDB: sessDB,
	}
}

// CheckToken проверяет, находится ли полученный токен в базах данных токенов зарегистрированных пользователей
func (snm *SessionsManager) CheckToken(wrt http.ResponseWriter, rqt *http.Request) *Session {
	tokenValue := rqt.Header.Get("token")
	var isUser bool
	err := snm.sessionsDB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE user_token = $1)", tokenValue).Scan(&isUser)
	if err != nil {
		log.Printf("error checking if the token belongs to user: %#v", err)
		http.Error(wrt, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return nil
	}

	var isAdmin bool
	if !isUser {
		err = snm.sessionsDB.QueryRow("SELECT EXISTS(SELECT 1 FROM admins WHERE admin_token = $1)", tokenValue).Scan(&isAdmin)
		if err != nil {
			log.Printf("error checking if the token belongs to admin: %#v", err)
			http.Error(wrt, "Внутренняя ошибка сервера", http.StatusInternalServerError)
			return nil
		}
		if !isAdmin {
			log.Printf("the token is invalid. User is unauthorized")
			http.Error(wrt, "Пользователь не авторизован", http.StatusUnauthorized)
			return nil
		}
	}

	sess := &Session{}
	sess.ParticipantToken = tokenValue
	sess.IsAdmin = isAdmin
	return sess
}

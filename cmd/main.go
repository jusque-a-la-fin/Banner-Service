package main

import (
	"banner/configs"
	"banner/internal/banner"
	"banner/internal/datastore"
	"banner/internal/handlers"
	"banner/internal/middleware"
	"banner/internal/session"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func main() {
	if errConf := configs.Init(); errConf != nil {
		log.Fatalf("config error: %#v", errConf)
	}
	tdb, err := datastore.CreateNewTokenDB()
	if err != nil {
		log.Fatalf("error connecting to token_storage database: %#v", err)
	}

	sm := session.NewSessionsManager(tdb)

	mdb, err := datastore.CreateNewMainDB()
	if err != nil {
		log.Fatalf("error connecting to banner_service database: %#v", err)
	}

	rdc := datastore.CreateNewRedisClient()
	bannerRepo := banner.NewDBRepo(mdb, rdc)
	bannerHandler := &handlers.BannerHandler{
		BannerRepo: bannerRepo,
	}

	rtr := mux.NewRouter()
	rtr.HandleFunc("/user_banner", bannerHandler.GetBanner).Methods("GET")
	rtr.HandleFunc("/banner", bannerHandler.GetAllBanners).Methods("GET")
	rtr.HandleFunc("/banner", bannerHandler.CreateNewBanner).Methods("POST")
	rtr.HandleFunc("/banner/{id}", bannerHandler.UpdateBanner).Methods("PATCH")
	rtr.HandleFunc("/banner/{id}", bannerHandler.DeleteBanner).Methods("DELETE")
	mux := middleware.Auth(sm, rtr)
	port := viper.GetString("port")
	err = http.ListenAndServe(port, mux)
	if err != nil {
		log.Fatalf("ListenAndServe error: %#v", err)
	}
}

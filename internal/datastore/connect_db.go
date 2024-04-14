package datastore

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

func CreateNewTokenDB() (*sql.DB, error) {
	postgresHost := os.Getenv("POSTGRES_HOST1")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB1")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", postgresHost, port, user, password, dbname)
	dtb, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error from `Open` function, package `sql`: %#v", err)
	}
	return dtb, nil
}

func CreateNewMainDB() (*sql.DB, error) {
	postgresHost := os.Getenv("POSTGRES_HOST2")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB2")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", postgresHost, port, user, password, dbname)
	dtb, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error from `Open` function, package `sql`: %#v", err)
	}
	return dtb, nil
}

func CreateNewRedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: "",
		DB:       0,
	})
	return rdb
}

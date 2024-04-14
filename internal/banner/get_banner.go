package banner

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type BannerDBRepository struct {
	dtb *sql.DB
	rdb *redis.Client
}

const (
	// numberOfMatchedValues - это число параметров(feature_id и tag_id), по которым определяется, является ли баннер искомым
	numberOfMatchedValues = 2
	timeLayout            = "2006-01-02 15:04:05.999999999 -0700 MST"
)

func NewDBRepo(sdb *sql.DB, rdc *redis.Client) *BannerDBRepository {
	return &BannerDBRepository{dtb: sdb, rdb: rdc}
}

// GetBannerFromDB получает баннер либо из кэша, либо из базы данных
func (repo *BannerDBRepository) GetBannerFromDB(featureID, tagID string, isAdmin bool) (*Banner, error) {
	featureValue, err := repo.getFeatureID(featureID)
	if err != nil {
		return nil, err
	}

	// Так как одна и та же комбинация из feature_id и tag_id может принадлежать разным баннерам(с разными JSON-ами),
	// а пользователь должен получить один баннер, то будет выбран последний добавленный из подходящих баннеров
	query := `
		SELECT banners.id, content, is_active
		FROM banners
		INNER JOIN banners_tags ON banners.id = banners_tags.banner_id
		WHERE banners.feature_id = $1 AND banners_tags.tag_id = $2
		ORDER BY banners.id DESC
		LIMIT 1;
	`

	banner := &Banner{}
	var isActive bool
	var contentBytes []byte
	err = repo.dtb.QueryRow(query, featureValue, tagID).Scan(new(interface{}), &contentBytes, &isActive)
	if contentBytes == nil {
		return nil, fmt.Errorf("error: banner hasn't been found")
	}

	if err != nil {
		return nil, fmt.Errorf("error while getting the banner info from database: %#v", err)
	}

	if !isActive && !isAdmin {
		return nil, fmt.Errorf("error: user doesn't have access")
	}

	if !json.Valid(contentBytes) {
		return nil, fmt.Errorf("error: incorrect banner info")
	}

	valid, err := CheckFields(contentBytes)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, fmt.Errorf("error: incorrect banner info")
	}

	banner.Content = contentBytes
	banner.IsActive = &isActive
	return banner, nil
}

// getFeatureID получает id для feature_id
func (repo *BannerDBRepository) getFeatureID(featureID string) (string, error) {
	query := `SELECT id FROM features WHERE feature_id = $1;`
	stmt, err := repo.dtb.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var id string
	err = stmt.QueryRow(featureID).Scan(&id)
	switch {
	case err == sql.ErrNoRows:
		return "", fmt.Errorf("error: feature_id hasn't been found")
	case err != nil:
		return "", fmt.Errorf("error: getting row from features table")
	default:
		return id, nil
	}
}

// CheckFields проверяет обязательные поля для JSON
func CheckFields(contentBytes []byte) (bool, error) {
	var content map[string]interface{}
	err := json.Unmarshal(contentBytes, &content)
	if err != nil {
		return false, fmt.Errorf("error unmarshaling banner info: %#v", err)
	}

	_, okTitle := content["title"]
	_, okText := content["text"]
	_, okURL := content["url"]
	if okTitle && okText && okURL {
		return true, nil
	}
	return false, nil
}

// GetBannerFromCache получает баннер из кэша
func (repo *BannerDBRepository) GetBannerFromCache(ctx context.Context, token, featureID, tagID string, isAdmin bool) (*Banner, error) {
	values, err := repo.rdb.SMembers(ctx, token).Result()
	if err != nil {
		return nil, fmt.Errorf("error getting values from set: %#v", err)
	}

	if len(values) == 0 {
		return nil, fmt.Errorf("error: value for this user token hasn't been found")
	}

	// Поиск баннера
	// Так как в кэше могут оказаться разные баннеры(с разными JSON-ами), соответствующие одной и той же комбинации из feature_id и tag_id,
	// а пользователь должен получить один баннер, то будет выбран последний из найденных баннеров
	var banners []Banner
	for _, val := range values {
		parts := strings.Split(val, "&")
		count := 0
		var creationDate time.Time
		banner := Banner{}
		for _, part := range parts {
			keyValue := strings.SplitN(part, "=", 2)
			if len(keyValue) != 2 {
				continue
			}

			key := strings.TrimSpace(keyValue[0])
			value := strings.TrimSpace(keyValue[1])
			switch key {
			case "feature_id":
				if value == featureID {
					count++
				}
			case "tag_id":
				if value == tagID {
					count++
				}
			case "updated_at":
				creationDate, err = time.Parse(timeLayout, value)
				if err != nil {
					return nil, fmt.Errorf("error parsing time: %#v", err)
				}
				currentTime := time.Now()
				if currentTime.Sub(creationDate) >= 5*time.Minute {
					err := repo.rdb.SRem(ctx, token, val).Err()
					if err != nil {
						return nil, fmt.Errorf("error deleting value from set: %#v", err)
					}
					return nil, fmt.Errorf("banner info is out of date")
				}

			case "is_active":
				if !isAdmin && value == "false" {
					return nil, fmt.Errorf("error: user doesn't have access")
				}
			case "banner_info":
				if count == numberOfMatchedValues {
					banner.Content = []byte(value)
					banners = append(banners, banner)
				}
			}
		}
	}

	if len(banners) == 0 {
		return nil, fmt.Errorf("error: banner hasn't been found")
	}
	banner := banners[len(banners)-1]
	return &banner, nil
}

// SetBannerInCache добавляет в кэш новый баннер
func (repo *BannerDBRepository) SetBannerInCache(ctx context.Context, banner Banner, token, featureID, tagID string) error {
	isActiveStr := strconv.FormatBool(*banner.IsActive)
	// время записи баннера в кэш
	updatedAt := time.Now().Format(timeLayout)
	content := banner.Content
	value := fmt.Sprintf("feature_id=%s&tag_id=%s&updated_at=%s&is_active=%s&banner_info=%s", featureID, tagID, updatedAt, isActiveStr, content)
	err := repo.rdb.SAdd(ctx, token, value).Err()
	if err != nil {
		return fmt.Errorf("error adding to set: %#v", err)
	}
	return nil
}

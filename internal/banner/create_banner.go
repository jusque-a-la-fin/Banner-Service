package banner

import (
	"fmt"
	"time"
)

// InsertNewBannerIntoDB добавляет новый баннер в базу данных
func (repo *BannerDBRepository) InsertNewBannerIntoDB(ban Banner) (*int, error) {
	currentTime := time.Now().Format(timeLayout)
	createdAt := currentTime
	updatedAt := currentTime

	// запись баннера в базу данных
	var bannerID int
	err := repo.dtb.QueryRow(
		"INSERT INTO banners (feature_id, content, created_at, updated_at, is_active) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		ban.FeatureID, ban.Content, createdAt, updatedAt, ban.IsActive,
	).Scan(&bannerID)
	if err != nil {
		return nil, fmt.Errorf("error inserting a banner into database: %#v", err)
	}

	queryPart := "INSERT INTO banners_tags (banner_id, tag_id) VALUES%s;"
	errMesg := "error inserting tags into database"
	err = repo.insertTagsIDs(bannerID, ban.TagsIDs, queryPart, errMesg)
	if err != nil {
		return nil, err
	}
	return &bannerID, nil
}

// insertTagsIDs записывает теги, связанные с баннером, в базу данных
func (repo *BannerDBRepository) insertTagsIDs(bannerID int, tagsIDs []int, queryPart, errMesg string) error {
	values := ""
	for idx := range tagsIDs {
		values = fmt.Sprintf("%s (%d, $%d),", values, bannerID, idx+1)
	}

	values = formatStr(values)
	query := fmt.Sprintf("%s;", queryPart)
	query = fmt.Sprintf(query, values)

	tagsIDS := make([]interface{}, len(tagsIDs))
	for idx, tagID := range tagsIDs {
		tagsIDS[idx] = tagID
	}

	_, err := repo.dtb.Exec(query, tagsIDS...)
	if err != nil {
		return fmt.Errorf("%s: %#v", errMesg, err)
	}
	return nil
}

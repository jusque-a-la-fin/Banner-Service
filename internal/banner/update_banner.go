package banner

import (
	"fmt"
	"time"
)

// UpdateBannerInDB обновляет содержимое баннера
func (repo *BannerDBRepository) UpdateBannerInDB(bannerID int, ban Banner) error {
	query := "SELECT EXISTS(SELECT 1 FROM banners WHERE id = $1);"
	var exists bool
	err := repo.dtb.QueryRow(query, bannerID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error while checking if banner with this id exists, %#v", err)
	}

	if !exists {
		return fmt.Errorf("error: banner with this id hasn't been found")
	}

	query = "UPDATE banners SET"
	var updateValues []interface{}
	if ban.FeatureID != nil {
		query = fmt.Sprintf("%s feature_id = $1", query)
		updateValues = append(updateValues, *ban.FeatureID)
	}
	if ban.Content != nil {
		query = fmt.Sprintf("%s, content = $2", query)
		updateValues = append(updateValues, ban.Content)
	}
	if ban.IsActive != nil {
		query = fmt.Sprintf("%s, is_active = $3", query)
		updateValues = append(updateValues, *ban.IsActive)
	}

	query = fmt.Sprintf("%s, updated_at = $4", query)
	if len(updateValues) == 0 {
		return fmt.Errorf("no columns to update")
	}

	currentTime := time.Now().Format(timeLayout)
	updatedAt := currentTime
	updateValues = append(updateValues, updatedAt)
	query = fmt.Sprintf("%s WHERE id = $5;", query)
	updateValues = append(updateValues, bannerID)

	_, err = repo.dtb.Exec(query, updateValues...)
	if err != nil {
		return fmt.Errorf("error while updating the banner: %#v", err)
	}

	if len(ban.TagsIDs) != 0 {
		queryPart := "INSERT INTO banners_tags (banner_id, tag_id) VALUES%s ON CONFLICT (banner_id, tag_id) DO NOTHING;"
		errMesg := "error while updating tag_ids"
		err := repo.insertTagsIDs(bannerID, ban.TagsIDs, queryPart, errMesg)
		if err != nil {
			return err
		}
	}
	return nil
}

func formatStr(values string) string {
	runes := []rune(values)
	runes = runes[:len(runes)-1]
	return string(runes)
}

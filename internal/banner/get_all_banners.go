package banner

import (
	"fmt"
	"strconv"
	"strings"
)

// GetAllBannersFromDB получает все баннеры из базы данных
func (repo *BannerDBRepository) GetAllBannersFromDB(featureID, tagID, offset, limit string) ([]Banner, error) {
	// формирование запроса
	query := `SELECT b.id, STRING_AGG(bt.tag_id::TEXT, ',') AS tag_ids, b.feature_id, b.content, b.is_active, b.created_at, b.updated_at
	          FROM banners b
			  JOIN banners_tags bt ON b.id = bt.banner_id`

	args := make([]interface{}, 0)
	if tagID != "" && featureID == "" {
		query = fmt.Sprintf("%s\nWHERE bt.tag_id = $1", query)
		args = append(args, tagID)
	}

	if featureID != "" && tagID == "" {
		query = fmt.Sprintf("%s\nWHERE b.feature_id = $1", query)
		args = append(args, featureID)
	}

	if tagID != "" && featureID != "" {
		query = fmt.Sprintf("%s\nWHERE bt.tag_id = $1 AND b.feature_id = $2", query)
		args = append(args, tagID, featureID)
	}

	query = fmt.Sprintf("%s\nGROUP BY b.id, b.feature_id, b.content, b.is_active, b.created_at, b.updated_at", query)
	query = fmt.Sprintf("%s\nORDER BY b.id", query)
	if offset != "" {
		query = fmt.Sprintf("%s %s", query, fmt.Sprintf("OFFSET %s", offset))
	}

	if limit != "" {
		query = fmt.Sprintf("%s %s", query, fmt.Sprintf("LIMIT %s", limit))
	}
	query = fmt.Sprintf("%s;", query)

	rows, err := repo.dtb.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error selecting all banners from database: %#v", err)
	}
	defer rows.Close()

	bans := make([]Banner, 0)
	for rows.Next() {
		ban := Banner{}
		tagsStr := ""
		if err = rows.Scan(&ban.ID, &tagsStr, &ban.FeatureID, &ban.Content, &ban.IsActive, &ban.CreatedAt, &ban.UpdatedAt); err != nil {
			return nil, fmt.Errorf("error scanning info of all banners: %#v", err)
		}

		tags := strings.Split(tagsStr, ",")
		for _, tag := range tags {
			var tagInt int
			tagInt, err = strconv.Atoi(tag)
			if err != nil {
				return nil, fmt.Errorf("error converting string to int: %#v", err)
			}
			ban.TagsIDs = append(ban.TagsIDs, tagInt)
		}
		bans = append(bans, ban)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error processing result of database query: %#v", err)
	}
	return bans, nil
}

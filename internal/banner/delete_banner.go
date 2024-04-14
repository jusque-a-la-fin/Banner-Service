package banner

import "fmt"

// DeleteBannerFromDB удаляет баннер из базы данных
func (repo *BannerDBRepository) DeleteBannerFromDB(bannerID int) error {
	// проверка, есть ли баннер с таким id
	query := "SELECT EXISTS(SELECT 1 FROM banners WHERE id = $1);"
	var exists bool
	err := repo.dtb.QueryRow(query, bannerID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error while checking if banner with this id exists: %#v", err)
	}

	if !exists {
		return fmt.Errorf("error: banner with this id hasn't been found")
	}

	// удаление тегов, связанных с баннером, из базы данных
	_, err = repo.dtb.Exec("DELETE FROM banners_tags WHERE banner_id = $1", bannerID)
	if err != nil {
		return fmt.Errorf("error while deleting from banners_tags %#v", err)
	}

	// удаление баннера из базы данных
	_, err = repo.dtb.Exec("DELETE FROM banners WHERE id = $1", bannerID)
	if err != nil {
		return fmt.Errorf("error while deleting from banners%#v", err)
	}
	return nil
}

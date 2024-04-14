package banner

import (
	"context"
)

type Banner struct {
	ID        *int
	TagsIDs   []int
	FeatureID *int
	Content   []byte
	IsActive  *bool
	CreatedAt string
	UpdatedAt string
}

type BannerRepo interface {
	GetBannerFromDB(featureID, tagID string, isAdmin bool) (*Banner, error)
	GetBannerFromCache(ctx context.Context, token, featureID, tagID string, isAdmin bool) (*Banner, error)
	SetBannerInCache(ctx context.Context, banner Banner, token, featureID, tagID string) error
	GetAllBannersFromDB(featureID, tagID, offset, limit string) ([]Banner, error)
	InsertNewBannerIntoDB(ban Banner) (*int, error)
	UpdateBannerInDB(bannerID int, ban Banner) error
	DeleteBannerFromDB(bannerID int) error
}

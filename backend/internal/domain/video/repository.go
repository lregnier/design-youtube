package video

import "context"

type VideoRepository interface {
	Save(ctx context.Context, v *Video) error
	FindByID(ctx context.Context, id VideoID) (*Video, error)
	ListReady(ctx context.Context) ([]*Video, error)
}

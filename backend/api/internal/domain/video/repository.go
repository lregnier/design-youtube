package video

import "context"

type VideoRepository interface {
	Save(ctx context.Context, v *Video) error
	FindByID(ctx context.Context, id VideoID) (*Video, error)
	List(ctx context.Context) ([]*Video, error)
}

package application

import (
	"context"

	"github.com/lregnier/design-youtube/api/internal/domain/video"
)

type EventPublisher interface {
	Publish(ctx context.Context, event video.DomainEvent) error
}

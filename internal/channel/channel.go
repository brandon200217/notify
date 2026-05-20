package channel

import (
	"context"

	"github.com/brandon200217/NOTIFY/internal/models"
)

type Channel interface {
	Send(ctx context.Context, req *models.NotifyRequest) error
}

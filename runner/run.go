package runner

import (
	"context"
)

type TaskRunner interface {
	Run(ctx context.Context) error
}

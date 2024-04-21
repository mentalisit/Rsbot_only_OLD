package storage

import (
	"Rsbot_only/internal/models"
	"context"
)

type Timers interface {
	UpdateMitutsQueue(ctx context.Context, name, CorpName string) models.Sborkz
	MinusMin(ctx context.Context) []models.Sborkz
}
type TimeDeleteMessage interface {
	TimerDeleteMessage() []models.Timer
	TimerInsert(c models.Timer)
}

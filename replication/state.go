package replication

import (
	"context"

	"github.com/jackc/pglogrepl"
)

type ReplicationState interface {
	GetLsn(ctx context.Context) (pglogrepl.LSN, error)
	SetLsn(ctx context.Context, lsn pglogrepl.LSN) error
}

type memoryReplicationState struct {
	LSN pglogrepl.LSN `json:"lsn"`
}

func NewMemoryReplicationState() ReplicationState {
	return &memoryReplicationState{
		LSN: pglogrepl.LSN(0),
	}
}

// GetLsn implements ReplicationState.
func (m *memoryReplicationState) GetLsn(ctx context.Context) (pglogrepl.LSN, error) {
	return m.LSN, nil
}

// SetLsn implements ReplicationState.
func (m *memoryReplicationState) SetLsn(ctx context.Context, lsn pglogrepl.LSN) error {
	m.LSN = lsn
	return nil
}

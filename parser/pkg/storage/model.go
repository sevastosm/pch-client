package storage

import (
	"github.com/sermojohn/postgres-client/pkg/domain"
	"time"
)

type SummaryRow struct {
	Server    domain.IXPServer
	Summary   domain.BGPSummary
	UpdatedAt time.Time
}

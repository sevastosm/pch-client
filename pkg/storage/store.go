package storage

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/sermojohn/postgres-client/pkg/config"
	"github.com/sermojohn/postgres-client/pkg/domain"
)

const (
	InsertStatement = "insert into IXP_MONITORING_DATA_V2 (IXP, COUNTRY, CITY, RS_LOCAL_ASN, NUMBER_OF_RIB_ENTRIES, NUMBER_OF_PEERS, TOTAL_NUMBER_OF_NEIGHBORS, UPDATED) values ($1, $2, $3, $4, $5, $6, $7, $8);"
	UpdateStatement = "update IXP_MONITORING_DATA_V2 set RS_LOCAL_ASN=$1, NUMBER_OF_RIB_ENTRIES=$2, NUMBER_OF_PEERS=$3, TOTAL_NUMBER_OF_NEIGHBORS=$4, UPDATED=$5 where IXP=$6 AND COUNTRY=$7 AND CITY=$8;"
	SelectStatement = "select IXP, COUNTRY, CITY, RS_LOCAL_ASN, NUMBER_OF_RIB_ENTRIES, NUMBER_OF_PEERS, TOTAL_NUMBER_OF_NEIGHBORS, UPDATED from IXP_MONITORING_DATA_V2;"
)

type Store interface {
	InsertSummary(server domain.IXPServer, summary domain.BGPSummary) error
	UpdateSummary(server domain.IXPServer, summary domain.BGPSummary) error
	UpsertSummary(server domain.IXPServer, summary domain.BGPSummary) error
	SelectAllSummaries() ([]SummaryRow, error)
}

type store struct {
	db *sql.DB
}

func (s *store) InsertSummary(server domain.IXPServer, summary domain.BGPSummary) error {
	_, err := s.db.Exec(InsertStatement, server.IXP, server.Country, server.City,
		summary.LocalASNumber, summary.RIBEntries, summary.NumberOfPeers, summary.TotalNumberOfNeighbors, time.Now())
	if err != nil {
		return err
	}
	return nil
}

func (s *store) UpdateSummary(server domain.IXPServer, summary domain.BGPSummary) error {
	_, err := s.db.Exec(UpdateStatement, summary.LocalASNumber, summary.RIBEntries, summary.NumberOfPeers, summary.TotalNumberOfNeighbors,
		time.Now(), server.IXP, server.Country, server.City)
	if err != nil {
		return err
	}
	return nil
}

func (s *store) UpsertSummary(server domain.IXPServer, summary domain.BGPSummary) error {
	res, err := s.db.Exec(UpdateStatement, summary.LocalASNumber, summary.RIBEntries, summary.NumberOfPeers, summary.TotalNumberOfNeighbors,
		time.Now(), server.IXP, server.Country, server.City)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 || err != nil {
		_, err = s.db.Exec(InsertStatement, server.IXP, server.Country, server.City,
			summary.LocalASNumber, summary.RIBEntries, summary.NumberOfPeers, summary.TotalNumberOfNeighbors, time.Now())
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *store) SelectAllSummaries() ([]SummaryRow, error) {
	rows, err := s.db.Query(SelectStatement)
	if err != nil {
		return nil, err
	}

	var summaryRows []SummaryRow
	for ; rows.Next(); {
		row := SummaryRow{}
		err := rows.Scan(&row.Server.IXP, &row.Server.Country, &row.Server.Country, &row.Summary.LocalASNumber,
			&row.Summary.RIBEntries, &row.Summary.NumberOfPeers, &row.Summary.TotalNumberOfNeighbors,
			&row.Updated)
		if err != nil {
			return nil, err
		}
		summaryRows = append(summaryRows, row)
	}

	return summaryRows, nil
}

func New(dbc config.DBConfig) (Store, error) {
	pgConn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbc.Host, dbc.Port, dbc.User, dbc.Password, dbc.Name)

	db, err := sql.Open("postgres", pgConn)
	if err != nil {
		return nil, err
	}

	return &store{db: db}, nil
}

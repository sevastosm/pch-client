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
	InsertStatement = "insert into IXP_SERVER_DATA (IXP, COUNTRY, CITY, PROTOCOL, RS_LOCAL_ASN, NUMBER_OF_RIB_ENTRIES, NUMBER_OF_PEERS, TOTAL_NUMBER_OF_NEIGHBORS, UPDATED_AT) values ($1, $2, $3, $4, $5, $6, $7, $8, $9);"
	UpdateStatement = "update IXP_SERVER_DATA set RS_LOCAL_ASN=$1, NUMBER_OF_RIB_ENTRIES=$2, NUMBER_OF_PEERS=$3, TOTAL_NUMBER_OF_NEIGHBORS=$4, UPDATED_AT=$5 where IXP=$6 AND COUNTRY=$7 AND CITY=$8 AND PROTOCOL=$9;"
	SelectStatement = "select IXP, COUNTRY, CITY, PROTOCOL, RS_LOCAL_ASN, NUMBER_OF_RIB_ENTRIES, NUMBER_OF_PEERS, TOTAL_NUMBER_OF_NEIGHBORS, UPDATED_AT from IXP_SERVER_DATA;"
)

type Store interface {
	UpsertSummary(server domain.IXPServer, summary domain.BGPSummary) error
	SelectAllSummaries() ([]SummaryRow, error)
}

type store struct {
	db *sql.DB
}

func New(dbc config.DBConfig) (Store, error) {
	pgConn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbc.Host, dbc.Port, dbc.User, dbc.Password, dbc.Database)

	db, err := sql.Open("postgres", pgConn)
	if err != nil {
		return nil, err
	}

	s := &store{db: db}
	_, err = s.SelectAllSummaries()
	if err != nil {
		return nil, err
	}

	return &store{db: db}, nil
}

func (s *store) UpsertSummary(server domain.IXPServer, summary domain.BGPSummary) error {
	res, err := s.db.Exec(UpdateStatement, summary.LocalASNumber, summary.RIBEntries, summary.NumberOfPeers, summary.TotalNumberOfNeighbors,
		time.Now(), server.IXP, server.Country, server.City, server.Protocol)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		_, err = s.db.Exec(InsertStatement, server.IXP, server.Country, server.City, server.Protocol,
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
		err := rows.Scan(&row.Server.IXP, &row.Server.Country, &row.Server.City, &row.Server.Protocol,
			&row.Summary.LocalASNumber, &row.Summary.RIBEntries, &row.Summary.NumberOfPeers,
			&row.Summary.TotalNumberOfNeighbors, &row.UpdatedAt)
		if err != nil {
			return nil, err
		}
		summaryRows = append(summaryRows, row)
	}

	return summaryRows, nil
}

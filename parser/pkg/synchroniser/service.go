package synchroniser

import (
	"github.com/sermojohn/postgres-client/pkg/parser"
	"github.com/sermojohn/postgres-client/pkg/storage"
	"log"
)

// Synchroniser triggers parser to fetch data and storage component to store
type Synchroniser interface {
	UpdateSummaries()
}

func New(store storage.Store, parser parser.PCHParser) Synchroniser {
	return &synchroniser{
		store:  store,
		parser: parser,
	}
}

type synchroniser struct {
	store  storage.Store
	parser parser.PCHParser
}

func (upd *synchroniser) UpdateSummaries() {
	err := upd.parser.FetchSummaries(func(resp *parser.FetchResponse) {
		if err2 := upd.store.UpsertSummary(resp.Server.IXPServer, resp.Summary); err2 != nil {
			log.Printf("[synchroniser] failed to update BGP summary for %s, error: %v\n", resp.Server.IXPServer.IXP, err2)
		}
	})
	if err != nil {
		log.Printf("[synchroniser] failed to fetch BGP summaries, error: %v\n", err)
	}
}

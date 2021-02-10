package updater

import (
	"github.com/sermojohn/postgres-client/pkg/parser"
	"github.com/sermojohn/postgres-client/pkg/storage"
	"log"
)

type Updater interface {
	UpdateSummaries() error
}

func New(store storage.Store, parser parser.IXPParser) Updater {
	return &updater{
		store:  store,
		parser: parser,
	}
}

type updater struct {
	store  storage.Store
	parser parser.IXPParser
}

func (upd *updater) UpdateSummaries() error {
	err := upd.parser.ForEachSummary(func(resp *parser.FetchResponse) error {
		if err2 := upd.store.UpsertSummary(resp.Server.IXPServer, resp.Summary); err2 != nil {
			log.Printf("[updater] failed to update BGP summary for %s, error: %v\n", resp.Server.IXPServer.IXP, err2)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

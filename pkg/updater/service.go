package updater

import (
	"github.com/sermojohn/postgres-client/pkg/config"
	"github.com/sermojohn/postgres-client/pkg/parser"
	"github.com/sermojohn/postgres-client/pkg/storage"
	"log"
)

type Updater interface {
	UpdateSummaries() error
}

func New(store storage.Store, parser parser.IXPParser, config config.UpdateConfig) Updater {
	return &updater{
		store:  store,
		parser: parser,
		config: config,
	}
}

type updater struct {
	store  storage.Store
	parser parser.IXPParser
	config config.UpdateConfig
}

func (upd *updater) UpdateSummaries() error {
	err := upd.parser.ForEachSummary(upd.config.AmountOfServers, upd.config.ParserRateLimitDelayMillis, func(resp *parser.FetchResponse) error {
		if err := upd.store.UpsertSummary(resp.Server.IXPServer, resp.Summary); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	log.Printf("[updater] updated %d BGP summaries", upd.config.AmountOfServers)
	return nil
}

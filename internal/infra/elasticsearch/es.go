package elasticsearch

import (
	"shop/internal/config"

	"github.com/elastic/go-elasticsearch/v8"
	"go.uber.org/zap"
)

func NewElasticsearch(cfg *config.Config, logger *zap.Logger) (*elasticsearch.Client, error) {
	esCfg := elasticsearch.Config{
		Addresses: cfg.Elasticsearch.Addresses,
		Username:  cfg.Elasticsearch.Username,
		Password:  cfg.Elasticsearch.Password,
	}

	es, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		logger.Error("failed to create elasticsearch client", zap.Error(err))
		return nil, err
	}

	// Ping to check connection (optional but recommended)
	res, err := es.Info()
	if err != nil {
		logger.Error("failed to connect to elasticsearch", zap.Error(err))
		return nil, err
	}
	defer res.Body.Close()

	return es, nil
}

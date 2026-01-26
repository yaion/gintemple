package idgen

import (
	"shop/internal/config"

	"github.com/bwmarrin/snowflake"
	"go.uber.org/zap"
)

// IDGenerator is an interface for generating unique IDs
type IDGenerator interface {
	GenerateID() int64
	GenerateStringID() string
}

type snowflakeGenerator struct {
	node *snowflake.Node
}

// NewIDGenerator creates a new Snowflake ID generator
func NewIDGenerator(cfg *config.Config, logger *zap.Logger) (IDGenerator, error) {
	node, err := snowflake.NewNode(cfg.Server.NodeID)
	if err != nil {
		logger.Error("failed to create snowflake node", zap.Error(err))
		return nil, err
	}

	return &snowflakeGenerator{node: node}, nil
}

func (g *snowflakeGenerator) GenerateID() int64 {
	return g.node.Generate().Int64()
}

func (g *snowflakeGenerator) GenerateStringID() string {
	return g.node.Generate().String()
}

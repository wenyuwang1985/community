package snowflake

import (
	"fmt"

	sf "github.com/bwmarrin/snowflake"
)

var node *sf.Node

func Init(nodeID int64) error {
	n, err := sf.NewNode(nodeID)
	if err != nil {
		return fmt.Errorf("Snowflake 初始化失败: %w", err)
	}
	node = n
	return nil
}

func Generate() int64 {
	return node.Generate().Int64()
}

package uuid

import "github.com/bwmarrin/snowflake"

var node *snowflake.Node

func init() {
	var err error
	node, err = snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}
}

func GenerateID() int64 {
	id := node.Generate()
	return int64(id)
}

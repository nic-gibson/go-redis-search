package main

import (
	"context"
	"fmt"

	"github.com/nic-gibson/go-redis-search/ftsearch"
	"github.com/redis/go-redis/v9"
)

func main() {

	client := ftsearch.NewClient(&redis.Options{})

	createCmd := client.CreateIndex(context.Background(), "simple", ftsearch.NewIndexOptions().AddSchemaAttribute(ftsearch.TextAttribute{
		Name:  "foo",
		Alias: "bar",
	}))

	_, _ = createCmd.Result()
	cmd := createCmd.String()
	fmt.Println(cmd)
}

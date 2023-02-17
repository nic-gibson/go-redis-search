package ftsearch_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/redis/go-redis/v9"

	"github.com/nic-gibson/go-redis-search/ftsearch"
)

var client *ftsearch.Client
var ctx = context.Background()

var _ = Describe("Create", func() {

	BeforeEach(func() {
		// Requires redis on localhost:6379 with search module!
		client = ftsearch.NewClient(&redis.Options{})
		Expect(client.Ping(ctx).Err()).NotTo(HaveOccurred())
		Expect(client.FlushDB(ctx).Err()).NotTo(HaveOccurred())
	})

	It("can build the simplest index", func() {
		createCmd := client.CreateIndex(ctx, "simple", ftsearch.NewIndexOptions().AddSchemaAttribute(ftsearch.TextAttribute{
			Name:  "foo",
			Alias: "bar",
		}))
		Expect(createCmd.Err()).NotTo(HaveOccurred())
		Expect(createCmd.String()).To(Equal("ft.create simple on hash score 1 schema foo as bar text: true"))
	})

	It("can build a hash index with options", func() {
		createCmd := client.CreateIndex(ctx, "withoptions", ftsearch.NewIndexOptions().
			AddPrefix("account:").
			WithMaxTextFields().
			WithScore(0.5).
			WithLanguage("spanish").
			AddSchemaAttribute(ftsearch.TextAttribute{
				Name:  "foo",
				Alias: "bar",
			}))
		Expect(createCmd.Err()).NotTo(HaveOccurred())
		Expect(createCmd.String()).To(Equal("ft.create withoptions on hash prefix 1 account: language spanish score 0.5 maxtextfields schema foo as bar text: true"))
	})

	It("can build a hash index with multiple schema entries", func() {
		createCmd := client.CreateIndex(ctx, "multiattrib", ftsearch.NewIndexOptions().
			AddSchemaAttribute(ftsearch.TextAttribute{
				Name:  "texttest",
				Alias: "xxtext",
			}).
			AddSchemaAttribute(ftsearch.NumericAttribute{
				Name:     "numtest",
				Sortable: true,
			}))
		Expect(createCmd.Err()).NotTo(HaveOccurred())
		Expect(createCmd.String()).To(Equal("ft.create multiattrib on hash score 1 schema texttest as xxtext text numtest numeric sortable: true"))
	})

})

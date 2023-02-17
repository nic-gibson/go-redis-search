// query provides an interface to RedisSearch's query functionality.
package ftsearch

import (
	"context"
	"fmt"
	"math"
)

type QueryOptions struct {
	Index        string
	NoContent    bool
	Verbatim     bool
	NoStopWords  bool
	WithScores   bool
	InOrder      bool
	ExplainScore bool
	Limit        *queryLimit
	ReturnFields []string
	Filters      queryFilterList
	InKeys       []string
	InFields     []string
	Language     string
	Slop         int32
	Summarize    *querySummarize
	HighLight    *queryHighlight
}

const (
	noSlop                   = -100 // impossible value for slop to indicate none set
	defaultOffset            = 0    // default first value for return offset
	defaultLimit             = 10   // default number of results to return
	defaultSumarizeSeparator = "..."
	defaultSummarizeLen      = 20
	defaultSummarizeFrags    = 3
)

type QueryResult struct {
	Score       float64
	Value       map[string]string
	Explanation []interface{}
}

type QueryResults struct {
	Count int64
	Data  map[string]QueryResult
}

func (c *Client) FTSearch(ctx context.Context, index string, query string, qry *QueryOptions) *QueryCmd {
	serialized := qry.serialize()
	args := []interface{}{"ft.search", index, query}
	args = append(args, serialized...)
	return NewQueryCmd(ctx, args...)
}

/******************************************************************************
* Functions operating on the query struct itself							  *
******************************************************************************/

// NewQuery creates a new query with defaults set
func NewQuery() *QueryOptions {
	return &QueryOptions{
		Limit: DefaultQueryLimit(),
		Slop:  noSlop,
	}
}

// String returns the serialized query as a single string. Any quoting
// required to use it in redis-cli is not done.
func (q *QueryOptions) String() string {
	return fmt.Sprintf("%v", q.serialize())
}

// WithIndex sets the index to be search on a query, returning the
// udpated query for chaining
func (q *QueryOptions) WithIndex(index string) *QueryOptions {
	q.Index = index
	return q
}

// WithLimit adds a limit to a query, returning the Query with
// the limit added (to allow chaining)
func (q *QueryOptions) WithLimit(first int64, num int64) *QueryOptions {
	q.Limit = NewQueryLimit(first, num)
	return q
}

// WithReturnFields sets the return fields, replacing any which
// might currently be set, returning the updated qry.
func (q *QueryOptions) WithReturnFields(fields []string) *QueryOptions {
	q.ReturnFields = fields
	return q
}

// AddReturnField appends a single field to the return fields,
// returning the updated query
func (q *QueryOptions) AddReturnField(field string) *QueryOptions {
	q.ReturnFields = append(q.ReturnFields, field)
	return q
}

// WithFilters sets the filters, replacing any which might
// be currently set, returning the updated query
func (q *QueryOptions) WithFilters(filters []*queryFilter) *QueryOptions {
	q.Filters = filters
	return q
}

// WithFilters sets the filters, replacing any which might
// be currently set, returning the updated query
func (q *QueryOptions) AddFilter(filter *queryFilter) *QueryOptions {
	q.Filters = append(q.Filters, filter)
	return q
}

// WithInKeys sets the keys to be searched, limiting the search
// to only these keys. The updated query is returned.
func (q *QueryOptions) WithInKeys(keys []string) *QueryOptions {
	q.InKeys = keys
	return q
}

// AddKey adds a single key to the keys to be searched, limiting the search
// to only these keys. The updated query is returned.
func (q *QueryOptions) AddKey(key string) *QueryOptions {
	q.InKeys = append(q.InKeys, key)
	return q
}

// WithInKeys sets the fields to be searched, limiting the search
// to only these fields. The updated query is returned.
func (q *QueryOptions) WithInFields(fields []string) *QueryOptions {
	q.InFields = fields
	return q
}

// AddField adds a single field to the fields to be searched in, limiting the search
// to only these fields. The updated query is returned.
func (q *QueryOptions) AddField(field string) *QueryOptions {
	q.InFields = append(q.InFields, field)
	return q
}

// WithSummarize sets the Summarize member of the query, returning the updated query.
func (q *QueryOptions) WithSummarize(s *querySummarize) *QueryOptions {
	q.Summarize = s
	return q
}

// WithHighlight sets the Highlight member of the query, returning the updated query.
func (q *QueryOptions) WithHighlight(h *queryHighlight) *QueryOptions {
	q.HighLight = h
	return q
}

// serialize converts a query struct to a slice of  interface{}
// ready for execution against Redis
func (q *QueryOptions) serialize() []interface{} {
	var args = []interface{}{}

	if q.NoContent {
		args = append(args, "NOCONTENT")
	}

	if q.Verbatim {
		args = append(args, "VERBATIM")
	}

	if q.NoStopWords {
		args = append(args, "NOSTOPWORDS")
	}

	if q.WithScores {
		args = append(args, "WITHSCORES")
	}

	args = append(args, q.Filters.serialize()...)
	args = append(args, serializeCountedArgs("RETURN", false, q.ReturnFields)...)

	if q.Summarize != nil {
		args = append(args, q.Summarize.serialize()...)
	}

	if q.HighLight != nil {
		args = append(args, q.HighLight.serialize()...)
	}

	args = append(args, q.serializeSlop()...)

	if q.InOrder {
		args = append(args, "INORDER")
	}

	args = append(args, q.serializeLanguage()...)
	args = append(args, serializeCountedArgs("INKEYS", false, q.InKeys)...)
	args = append(args, serializeCountedArgs("INFIELDS", false, q.InFields)...)

	if q.ExplainScore {
		args = append(args, "EXPLAINSCORE")
	}

	if q.Limit != nil {
		args = append(args, q.Limit.serialize()...)
	}

	return args
}

// resultSize uses the query to work out how many entries
// in the query raw results slice are used per result.
func (q *QueryOptions) resultSize() int {
	count := 2 // default to 2 - key and value

	if q.WithScores { // one more if returning scores
		count += 1
	}

	if q.NoContent { // one less if not content
		count -= 1
	}

	if q.ExplainScore { // one more if explaining
		count += 1
	}

	return count
}

func (q *QueryOptions) serializeSlop() []interface{} {
	if q.Slop != noSlop {
		return []interface{}{"SLOP", q.Slop}
	} else {
		return nil
	}
}

func (q *QueryOptions) serializeLanguage() []interface{} {
	if q.Language != "" {
		return []interface{}{"LANGUAGE", q.Language}
	} else {
		return nil
	}
}

/******************************************************************************
* Public utilities                                                            *
******************************************************************************/

// FilterValue formats a value for use in a filter and returns it
func FilterValue(val float64, exclusive bool) interface{} {
	prefix := ""
	if exclusive {
		prefix = "("
	}

	if math.IsInf(val, -1) {
		return prefix + "-inf"
	} else if math.IsInf(val, 1) {
		return prefix + "+inf"
	} else {
		return fmt.Sprintf("%s%f", prefix, val)
	}
}

/******************************************************************************
* Internal utilities                                                          *
******************************************************************************/

func toMap(input []interface{}) map[string]string {

	results := make(map[string]string, len(input)/2)
	key := ""
	for i := 0; i < len(input); i += 2 {
		key = input[i].(string)
		value := input[i+1].(string)
		results[key] = value
	}
	return results
}

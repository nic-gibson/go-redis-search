// query provides an interface to RedisSearch's query functionality.
package ftsearch

import (
	"context"
	"fmt"
	"math"
	"strconv"

	"github.com/go-redis/redis/v8"
)

type countedArgs []string

type query struct {
	Index        string
	QueryString  string
	NoContent    bool
	Verbatim     bool
	NoStopWords  bool
	WithScores   bool
	InOrder      bool
	ExplainScore bool
	Limit        *queryLimit
	ReturnFields countedArgs
	Filters      queryFilterList
	InKeys       countedArgs
	InFields     countedArgs
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

func (c *Client) Search(ctx context.Context, qry *query) (*QueryResults, error) {

	serialized := qry.serialize()
	cmd := redis.NewSliceCmd(ctx, serialized...)
	if err := c.client.Process(ctx, cmd); err != nil {
		return nil, err
	} else if rawResults, err := cmd.Result(); err != nil {
		return nil, err
	} else {
		resultSize := qry.resultSize()
		resultCount := (len(rawResults) - 1) / resultSize
		results := QueryResults{
			Count: rawResults[0].(int64),
			Data:  make(map[string]QueryResult, resultCount),
		}

		for i := 1; i < len(rawResults); i += resultSize {
			j := 0
			var score float64 = 0

			key := rawResults[i+j].(string)
			j++

			if qry.WithScores {
				score, _ = strconv.ParseFloat(rawResults[i+j].(string), 64)
				j++
			}

			result := QueryResult{
				Score: score,
			}

			if !qry.NoContent {
				result.Value = toMap(rawResults[i+j].([]interface{}))
				j++
			}

			results.Data[key] = result

		}
		return &results, nil
	}
}

/******************************************************************************
* Functions operating on the query struct itself							  *
******************************************************************************/

// NewQuery creates a new query with defaults set
func NewQuery() *query {
	return &query{
		Limit: DefaultQueryLimit(),
		Slop:  noSlop,
	}
}

// String returns the serialized query as a single string. Any quoting
// required to use it in redis-cli is not done.
func (q *query) String() string {
	return fmt.Sprintf("%v", q.serialize())
}

// WithQueryString sets the raw query string on a Query, returning
// the updated query for chaining.
func (q *query) WithQueryString(queryString string) *query {
	q.QueryString = queryString
	return q
}

// WithIndex sets the index to be search on a query, returning the
// udpated query for chaining
func (q *query) WithIndex(index string) *query {
	q.Index = index
	return q
}

// WithLimit adds a limit to a query, returning the Query with
// the limit added (to allow chaining)
func (q *query) WithLimit(first int64, num int64) *query {
	q.Limit = NewQueryLimit(first, num)
	return q
}

// WithReturnFields sets the return fields, replacing any which
// might currently be set, returning the updated qry.
func (q *query) WithReturnFields(fields []string) *query {
	q.ReturnFields = fields
	return q
}

// AddReturnField appends a single field to the return fields,
// returning the updated query
func (q *query) AddReturnField(field string) *query {
	q.ReturnFields = append(q.ReturnFields, field)
	return q
}

// WithFilters sets the filters, replacing any which might
// be currently set, returning the updated query
func (q *query) WithFilters(filters []*queryFilter) *query {
	q.Filters = filters
	return q
}

// WithFilters sets the filters, replacing any which might
// be currently set, returning the updated query
func (q *query) AddFilter(filter *queryFilter) *query {
	q.Filters = append(q.Filters, filter)
	return q
}

// WithInKeys sets the keys to be searched, limiting the search
// to only these keys. The updated query is returned.
func (q *query) WithInKeys(keys []string) *query {
	q.InKeys = keys
	return q
}

// AddKey adds a single key to the keys to be searched, limiting the search
// to only these keys. The updated query is returned.
func (q *query) AddKey(key string) *query {
	q.InKeys = append(q.InKeys, key)
	return q
}

// WithInKeys sets the fields to be searched, limiting the search
// to only these fields. The updated query is returned.
func (q *query) WithInFields(fields []string) *query {
	q.InFields = fields
	return q
}

// AddField adds a single field to the fields to be searched in, limiting the search
// to only these fields. The updated query is returned.
func (q *query) AddField(field string) *query {
	q.InFields = append(q.InFields, field)
	return q
}

// WithSummarize sets the Summarize member of the query, returning the updated query.
func (q *query) WithSummarize(s *querySummarize) *query {
	q.Summarize = s
	return q
}

// WithHighlight sets the Highlight member of the query, returning the updated query.
func (q *query) WithHighlight(h *queryHighlight) *query {
	q.HighLight = h
	return q
}

// serialize converts a query struct to a slice of  interface{}
// ready for execution against Redis
func (q *query) serialize() []interface{} {
	var args = []interface{}{"FT.SEARCH", q.Index, q.QueryString}

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
	args = append(args, q.ReturnFields.serialize("RETURN")...)

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
	args = append(args, q.InKeys.serialize("INKEYS")...)
	args = append(args, q.InFields.serialize("INFIELDS")...)

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
func (q *query) resultSize() int {
	count := 2 // default to 2 - key and value

	if q.WithScores {
		count += 1
	}

	if q.NoContent {
		count -= 1
	}

	if q.ExplainScore {
		count += 1
	}

	return count
}

func (q *query) serializeSlop() []interface{} {
	if q.Slop != noSlop {
		return []interface{}{"SLOP", q.Slop}
	} else {
		return nil
	}
}

func (q *query) serializeLanguage() []interface{} {
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

func (c countedArgs) serialize(name string) []interface{} {
	if len(c) > 0 {
		result := make([]interface{}, 2+len(c))

		result[0] = name
		result[1] = len(c)
		for pos, val := range c {
			result[pos+2] = val
		}

		return result
	} else {
		return nil
	}
}

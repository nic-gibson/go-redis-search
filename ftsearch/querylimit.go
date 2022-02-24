package ftsearch

/******************************************************************************
* Functions operating on QueryLimit structs                                   *
******************************************************************************/

// queryLimit defines the results by offset and number.
type queryLimit struct {
	First int64
	Num   int64
}

// NewQueryLimit returns an initialized QueryLimit struct
func NewQueryLimit(first int64, num int64) *queryLimit {
	return &queryLimit{First: first, Num: num}
}

// DefaultQueryLimit returns an initialzied QueryLimit struct with the
// default limit range
func DefaultQueryLimit() *queryLimit {
	return NewQueryLimit(defaultOffset, defaultLimit)
}

// Serialize the limit for output
func (ql *queryLimit) serialize() []interface{} {
	if ql.First == defaultOffset && ql.Num == defaultLimit {
		return nil
	} else {
		return []interface{}{"LIMIT", ql.First, ql.Num}
	}
}

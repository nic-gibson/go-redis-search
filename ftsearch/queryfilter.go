package ftsearch

// Functions operating on QueryFilter structs                                  *

import "math"

type queryFilter struct {
	Attribute string
	Min       interface{} // either a numeric value or +inf, -inf or "(" followed by numeric
	Max       interface{} // as above
}

type queryFilterList []*queryFilter

func NewQueryFilter(attribute string) *queryFilter {
	qf := &queryFilter{Attribute: attribute}
	return qf.
		WithMinInclusive(math.Inf(-1)).
		WithMaxInclusive(math.Inf(1))
}

// WithMinInclusive sets an inclusive minimum for the query filter value and
// returns it
func (qf *queryFilter) WithMinInclusive(val float64) *queryFilter {
	qf.Min = FilterValue(val, false)
	return qf
}

// WithMaxInclusive sets an inclusive maximum for the query filter value and
// returns it
func (qf *queryFilter) WithMaxInclusive(val float64) *queryFilter {
	qf.Max = FilterValue(val, false)
	return qf
}

// WithMinExclusive sets an exclusive minimum for the query filter value and
// returns it
func (qf *queryFilter) WithMinExclusive(val float64) *queryFilter {
	qf.Min = FilterValue(val, true)
	return qf
}

// WithMaxExclusive sets an exclusive maximum for the query filter value and
// returns it
func (qf *queryFilter) WithMaxExclusive(val float64) *queryFilter {
	qf.Max = FilterValue(val, true)
	return qf
}

// serialize converts a filter list to an array of interface{} objects for execution
func (q queryFilterList) serialize() []interface{} {
	if len(q) > 0 {
		var args []interface{}
		for _, arg := range q {
			args = append(args, "FILTER", arg.Attribute, arg.Min, arg.Max)
		}
		return args
	} else {
		return nil
	}
}

package ftsearch

// Functions and structs used to set up summarization and highlighting.

type querySummarize struct {
	Fields    countedArgs
	Frags     int32
	Len       int32
	Separator string
}

func DefaultQuerySummarize() *querySummarize {
	return &querySummarize{
		Separator: defaultSumarizeSeparator,
		Len:       defaultSummarizeLen,
		Frags:     defaultSummarizeFrags,
	}
}

func NewQuerySummarize() *querySummarize {
	return &querySummarize{}
}

func NewQueryHighlight() *queryHighlight {
	return &queryHighlight{}
}

// WithLen sets the length of the query summarization fragment (in words)
// The modified struct is returned to support chaining
func (s *querySummarize) WithLen(len int32) *querySummarize {
	s.Len = len
	return s
}

// WithFrags sets the number of the fragements to create and return
// The modified struct is returned to support chaining
func (s *querySummarize) WithFrags(n int32) *querySummarize {
	s.Frags = n
	return s
}

// WithSeparator sets the fragment separator to be used.
// The modified struct is returned to support chaining
func (s *querySummarize) WithSeparator(sep string) *querySummarize {
	s.Separator = sep
	return s
}

// WithFields sets the fields to be summarized. Leaving it empty
// (the default) will cause all fields to be summarized
// The modified struct is returned to support chaining
func (s *querySummarize) WithFields(fields []string) *querySummarize {
	s.Fields = fields
	return s
}

// AddField adds a new field to the list of those to be summarised.
// The modified struct is returned to support chaining
func (s *querySummarize) AddField(field string) *querySummarize {
	s.Fields = append(s.Fields, field)
	return s
}

// serialize prepares the summarisation to be passed to Redis.
func (s *querySummarize) serialize() []interface{} {
	args := []interface{}{"SUMMARIZE"}
	args = append(args, s.Fields.serialize("FIELDS")...)
	args = append(args, "FRAGS", s.Frags)
	args = append(args, "LEN", s.Len)
	args = append(args, "SEPARATOR", s.Separator)
	return args
}

// queryHighlight allows the user to define optional query highlighting
type queryHighlight struct {
	Fields   countedArgs
	OpenTag  string
	CloseTag string
}

// WithFields sets the fields to be highlighting. Leaving it empty
// (the default) will cause all fields to be highlighted
// The modified struct is returned to support chaining
func (h *queryHighlight) WithFields(fields []string) *queryHighlight {
	h.Fields = fields
	return h
}

// AddField adds a new field to the list of those to be highlighted.
// The modified struct is returned to support chaining
func (h *queryHighlight) AddField(field string) *queryHighlight {
	h.Fields = append(h.Fields, field)
	return h
}

// SetTags sets the start and end tags. Both must be non empty or
// both empty. This is not enforced in this code to keep the API consistent
// but will lead to a Redis error if not set correctly.
func (h *queryHighlight) SetTags(open string, close string) *queryHighlight {
	h.OpenTag = open
	h.CloseTag = close
	return h
}

// serialize prepares the highlighting to be passed to Redis.
func (h *queryHighlight) serialize() []interface{} {
	args := []interface{}{"HIGHLIGHT"}
	args = append(args, h.Fields.serialize("FIELDS")...)
	if h.OpenTag != "" || h.CloseTag != "" {
		args = append(args, "TAGS", h.OpenTag, h.CloseTag)
	}
	return args
}

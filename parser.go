package gotemplate

import (
	"errors"
	"fmt"
	"io"
	"os"

	csvmap "github.com/recursionpharma/go-csv-map"
	"github.com/shoobyban/mxj"
	"github.com/shoobyban/slog"
)

// ParserFunc is to parse a []byte into an interface{}
type ParserFunc func(io.Reader) (interface{}, error)

// Parser is the main type
type Parser struct {
	parsers map[string]ParserFunc
}

// NewParser defines a new parser
func NewParser() *Parser {
	return &Parser{
		parsers: map[string]ParserFunc{
			"xml": func(content io.Reader) (interface{}, error) {
				return mxj.NewMapXmlReader(content)
			},
			"json": func(content io.Reader) (interface{}, error) {
				return mxj.NewMapJsonReader(content)
			},
			"csv": func(content io.Reader) (interface{}, error) {
				r := csvmap.NewReader(content)
				r.Reader.LazyQuotes = true
				var err error
				r.Columns, err = r.ReadHeader()
				if err != nil {
					slog.Errorf("Error reading csv header %v", err)
				}
				return r.ReadAll()
			},
		},
	}
}

// RegisterParser registers or overrides a format parser func. Indices are lower case.
func (l *Parser) RegisterParser(format string, parser ParserFunc) {
	l.parsers[format] = parser
}

// ReadStruct reads from given file, parsing into structure
func (l *Parser) ReadStruct(filename, format string) (interface{}, error) {
	f, err := os.Open(filename)
	if err != nil {
		slog.Infof("Can't open file %s", filename)
		return nil, err
	}
	defer f.Close()
	return l.ParseStruct(f, format)
}

// ParseStruct parses byte slice into map or slice
func (l *Parser) ParseStruct(content io.Reader, format string) (interface{}, error) {
	var out interface{}
	var err error
	if parser, ok := l.parsers[format]; ok {
		out, err = parser(content)
	} else {
		return nil, errors.New("Unknown file")
	}
	if err != nil {
		return nil, fmt.Errorf("Can't parse %s: %v", format, err)
	}
	return out, nil
}

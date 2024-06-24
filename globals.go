// globals.go
package main

import (
	"github.com/blevesearch/bleve/v2"
	"github.com/spf13/afero"
)

var (
	devMode     bool     // If enabled, removes the database and search index on shutdown
	osFS        afero.Fs = afero.NewOsFs()
	chatTurn             = 1
	sqliteDB    *SQLiteDB
	searchIndex bleve.Index
)

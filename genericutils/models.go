package genericutils

import (
	"google.golang.org/api/docs/v1"
	"google.golang.org/api/sheets/v4"
)

// struct representing an object containing all Google Services
// to be passed between functions that interact with the google APIs
type GoogleServices struct {
	Docs   *docs.Service
	Sheets *sheets.Service
}

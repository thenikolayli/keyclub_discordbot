package eventutils

import "strings"

// https://docs.google.com/document/d/id-example/edit?tab=t.0
// extracts a Google docs id from the url
// splits it at /d/ and gets the second part, then splits it at /edit and gets the first one
func DocsUrlToId(url string) string {
	id := strings.Split(url, "/d/")[1]
	id = strings.Split(id, "/edit")[0]
	return id
}

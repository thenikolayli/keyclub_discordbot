package eventutils

import "strings"

// https://docs.google.com/document/d/id-example/edit?tab=t.0
// https://docs.google.com/document/u/0/d/1x8B8h9ZFNIUcartcK7JLDUHjmnMTu62LP8hNzK82xgI/mobilebasic
// extracts a Google docs id from the url
// splits it at /d/ and gets the second part, then splits it at /edit and gets the first one
func DocsUrlToId(url string) string {
	id := strings.Split(url, "/d/")[1]
	id = strings.Split(id, "/edit")[0]
	id = strings.Split(id, "/mobilebasic")[0]
	return id
}

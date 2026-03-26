package genericutils

import (
	"fmt"
	"time"
)

// just so i don't have to rewrite this in every footer
func GetFormattedLastUpdated(hoursLastUpdated time.Time) string {
	return fmt.Sprintf("Last updated: %v", hoursLastUpdated.Format("Jan 2 2006 15:04:05"))
}

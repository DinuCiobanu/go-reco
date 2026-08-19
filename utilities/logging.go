package utilities

import (
	"fmt"
	"time"
)

func LogFetch(entity string, interval time.Duration) {
	fmt.Printf("[%s] fetching %s, refetching in %s\n", time.Now().Format(time.RFC3339), entity, interval)
}

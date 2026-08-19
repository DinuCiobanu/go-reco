package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"resty.dev/v3"
)

func GetWithRetry(req *resty.Request, url string) (*resty.Response, error) {
	for {
		resp, err := req.Get(url)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode() != http.StatusTooManyRequests {
			return resp, nil
		}

		wait := 1
		if seconds, err := strconv.Atoi(resp.Header().Get("Retry-After")); err == nil {
			wait = seconds
		}

		fmt.Printf("rate limited, retrying in %ds\n", wait)
		time.Sleep(time.Duration(wait) * time.Second)
	}
}

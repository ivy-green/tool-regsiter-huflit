package huflit

import (
	"time"

	"github.com/imroc/req/v3"
)

const (
	userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36"
	timeout   = 2
)

type HuflitScraper struct {
	client *req.Client
}

func NewHuflitScraper() *HuflitScraper {
	return &HuflitScraper{
		client: requestClient(),
	}
}

func requestClient() *req.Client {
	client := req.C().
		SetUserAgent(userAgent).
		SetTimeout(timeout * time.Second).
		SetCommonRetryCount(10).
		SetCommonRetryCondition(func(resp *req.Response, err error) bool {
			return err != nil || resp.GetStatusCode() != 200
		})

	return client
}

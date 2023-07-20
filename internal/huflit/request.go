package huflit

import (
	"github.com/imroc/req/v3"
)

func (scraper *HuflitScraper) httpGet(url string, headers map[string]string) (*req.Response, error) {
	request := scraper.client.R()
	if headers != nil {
		request.SetHeaders(headers)
	}

	resp, err := request.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.IsErrorState() {
		return nil, ErrResponse(resp)
	}

	return resp, nil
}

func (scraper *HuflitScraper) httpPost(url string, data map[string]string, headers map[string]string) (*req.Response, error) {
	request := scraper.client.R()
	if headers != nil {
		request.SetHeaders(headers)
	}
	if data != nil {
		request.SetFormData(data)
	}

	resp, err := request.Post(url)
	if err != nil {
		return nil, err
	}

	if resp.IsErrorState() {
		return nil, ErrResponse(resp)
	}

	return resp, nil
}

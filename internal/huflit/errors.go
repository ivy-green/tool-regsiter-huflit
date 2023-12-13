package huflit

import (
	"errors"
	"fmt"

	"github.com/imroc/req/v3"
)

var ErrLoginFailed = errors.New("login failed")
var ErrCannotGetInfo = errors.New("get student info failed")

func ErrResponse(resp *req.Response) error {
	return fmt.Errorf("dump response: %v and statusCode : %d ", resp.String(), resp.StatusCode)
}

package huflit

import (
	"errors"
	"fmt"

	"github.com/imroc/req/v3"
)

var ErrLoginFailed = errors.New("username or password invalid")
var ErrCannotGetInfo = errors.New("cannot get info")

func ErrResponse(resp *req.Response) error {
	return fmt.Errorf("dump response: %v and statusCode : %d ", resp.String(), resp.StatusCode)
}

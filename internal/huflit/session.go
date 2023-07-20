package huflit

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type StudentInfo struct {
	Name string
	Code string
}

func (scraper *HuflitScraper) Login(username, password string) (*StudentInfo, error) {
	data := map[string]string{
		"txtTaiKhoan": username,
		"txtMatKhau":  password,
	}
	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
		"Origin":       "https://portal.huflit.edu.vn",
		"Referer":      "https://portal.huflit.edu.vn/Login",
	}
	resp, err := scraper.httpPost("https://portal.huflit.edu.vn/Login", data, headers)
	if err != nil {
		return nil, err
	}

	document, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	loginMessage := strings.TrimSpace(document.Find(".loginbox-forgot").Text())
	if loginMessage == "Tên đăng nhập hoặc mật khẩu không chính xác" {
		return nil, ErrLoginFailed
	}

	info := strings.TrimSpace(document.Find(`.stylecolor  a[data-toggle="dropdown"]`).Text())
	if info == "" {
		return nil, ErrCannotGetInfo
	}

	parseInfo := strings.Split(info, "|")

	return &StudentInfo{
		Code: parseInfo[0],
		Name: parseInfo[1],
	}, nil
}

func (scraper *HuflitScraper) GetSessionDKMH() error {
	_, err := scraper.httpGet("https://portal.huflit.edu.vn/Home/DangKyHocPhan", nil)
	if err != nil {
		return err
	}

	return nil
}

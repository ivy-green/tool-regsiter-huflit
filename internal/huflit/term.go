package huflit

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Term struct {
	Code      string
	Name      string
	Credit    float64
	RequestId string
	Subjects  []Subject
}

func (scraper *HuflitScraper) GetTerms() ([]Term, error) {
	resp, err := scraper.httpGet("https://dkmh.huflit.edu.vn/DangKyHocPhan/DanhSachHocPhan?typeId=KH&id=", nil)
	if err != nil {
		return nil, err
	}

	document, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	terms := make([]Term, 0)

	document.Find("#DanhSachLop").Each(func(i int, selection *goquery.Selection) {
		selection.Find("tr").Each(func(x int, tr *goquery.Selection) {
			if x < 2 {
				return
			}
			data := strings.Split(strings.TrimSpace(tr.Text()), "\n")
			if len(data) == 9 {
				credit, err := strconv.ParseFloat(strings.TrimSpace(data[3]), 64)
				if err != nil {
					return
				}
				link := ""

				tr.Find(`td[style="text-align:right"]`).Each(func(_ int, s *goquery.Selection) {
					val, exists := s.Find("a").Attr("href")
					if exists {
						link = val
					}
				})
				if link == "" {
					return
				}
				terms = append(terms, Term{
					Code:      strings.TrimSpace(data[1]),
					Name:      strings.TrimSpace(data[2]),
					Credit:    credit,
					RequestId: getRequestIdFromLink(link),
				})
			}
		})
	})

	return terms, nil
}

func getRequestIdFromLink(input string) string {
	//input := `javascript:GetClassStudyUnit('nu/Oq2Y7lSgnAK3c7GZvuQ==','Tư tưởng Hồ Chí Minh','KH', '')`
	regex := `'([^']*)'`
	re := regexp.MustCompile(regex)
	match := re.FindStringSubmatch(input)

	if len(match) > 1 {
		extractedValue := match[1]
		return extractedValue
	}

	return ""
}

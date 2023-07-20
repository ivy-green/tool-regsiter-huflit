package huflit

import (
	"errors"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Subject struct {
	Code      string
	RequestId string
}

type RegisterResp struct {
	State bool        `json:"State"`
	Obj   string      `json:"Obj"`
	Obj1  interface{} `json:"Obj1"`
	Obj2  interface{} `json:"Obj2"`
	Msg   string      `json:"Msg"`
}

func (scraper *HuflitScraper) GetClassStudyUnit(requestId string, registType string) ([]Subject, error) {
	url := fmt.Sprintf("https://dkmh.huflit.edu.vn/DangKyHocPhan/DanhSachLopHocPhan?id=%v&registType=%v", requestId, registType)
	resp, err := scraper.httpGet(url, nil)
	if err != nil {
		return nil, err
	}

	document, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	subjects := make([]Subject, 0)
	document.Find("tr").Each(func(_ int, selection *goquery.Selection) {
		val, ok := selection.Find(`input[name="theory"]`).Attr("id")
		if !ok {
			return
		}

		subject := Subject{
			RequestId: val,
		}
		selection.Find("td[style=\"text-align:center\"]").Each(func(i int, data *goquery.Selection) {
			switch i {
			case 2:
				subject.Code = strings.TrimSpace(data.Text())
				return
			}
		})
		subjects = append(subjects, subject)
	})

	var practicalSubjects map[string]Subject
	document.Find(".tr-no-hover").Each(func(_ int, tr *goquery.Selection) {
		tr.Find("tr").Each(func(_ int, trClass *goquery.Selection) {
			subject := Subject{}
			trClass.Find("td").Each(func(i int, td *goquery.Selection) {
				if i == 1 {
					subject.Code = strings.TrimSpace(td.Text())
				}
				val, exists := td.Find(`.classCheckChon`).Attr("id")
				if exists {
					subject.RequestId = val
				}
			})
			if subject.Code != "" {
				practicalSubjects[subject.Code] = subject
			}
		})
	})

	for _, subject := range practicalSubjects {
		subjects = append(subjects, subject)
	}

	return subjects, nil
}

func (scraper *HuflitScraper) RegisterSubject(firstRequestId, secondRequestId string) (*RegisterResp, error) {
	if secondRequestId != "" {
		secondRequestId += "|"
	}
	// https://dkmh.huflit.edu.vn/DangKyHocPhan/RegistUpdateScheduleStudyUnit?Hide=iqhcfeZlCkP2wJbeIWovVGXmVDNWkd1hLthhZ1A9lJUpGiEEStVARS1H+meXmaEY26MSLfpCAig=|&ScheduleStudyUnitOld=&acceptConflict=
	url := fmt.Sprintf("https://dkmh.huflit.edu.vn/DangKyHocPhan/RegistUpdateScheduleStudyUnit?Hide=%v|%v&ScheduleStudyUnitOld=&acceptConflict=", firstRequestId, secondRequestId)
	resp, err := scraper.httpGet(url, nil)
	if err != nil {
		return nil, err
	}

	if resp.IsErrorState() {
		return nil, errors.New("")
	}

	var registerResp RegisterResp
	if err := resp.UnmarshalJson(&registerResp); err != nil {
		return nil, err
	}

	return &registerResp, nil
}

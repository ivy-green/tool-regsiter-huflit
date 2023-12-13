package huflit

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Subject struct {
	Code      string
	RequestId string
}

type RegisterResp struct {
	State bool        `json:"State"`
	Obj   interface{} `json:"Obj"`
	Obj1  interface{} `json:"Obj1"`
	Obj2  interface{} `json:"Obj2"`
	Msg   string      `json:"Msg"`
}

func (scraper *HuflitScraper) GetClassStudyUnit(register Register, requestId string, registType string) ([]Subject, error) {
	url := fmt.Sprintf("https://dkmh.huflit.edu.vn/DangKyHocPhan/DanhSachLopHocPhan?id=%v&registType=%v", requestId, registType)
	resp, err := scraper.httpGet(url, nil)
	if err != nil {
		return nil, err
	}

	writeHTMLToFile("data.html", resp.Bytes())
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

	practicalSubjects := make(map[string]Subject, 0)

	selector := fmt.Sprintf("#tr-of-%v", register.FirstCode)
	document.Find(selector).Each(func(_ int, trClass *goquery.Selection) {

		trClass.Find("tr").Each(func(i int, tdElement *goquery.Selection) {
			code := tdElement.Find("td[style=\"text-align:center\"]").Eq(1).Text()
			if code != "" {
				val, exists := tdElement.Find(`.classCheckChon`).Attr("id")
				if exists {
					//subject.RequestId = val
					practicalSubjects[code] = Subject{
						Code:      code,
						RequestId: val,
					}
				}
			}
		})
		//trClass.Find("td").Each(func(i int, td *goquery.Selection) {
		//	log.Println(td.Text())
		//	if i == 1 {
		//		subject.Code = strings.TrimSpace(td.Text())
		//	}
		//	val, exists := td.Find(`.classCheckChon`).Attr("id")
		//	if exists {
		//		subject.RequestId = val
		//	}
		//})
		//if subject.Code != "" {
		//	practicalSubjects[subject.Code] = subject
		//}
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
	//log.Println(url)
	//return nil, nil
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

func writeHTMLToFile(fileName string, content []byte) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(content)
	if err != nil {
		return err
	}

	return nil
}

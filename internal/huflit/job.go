package huflit

import (
	"sync"

	"github.com/rs/zerolog/log"
)

func (scraper *HuflitScraper) StartJob(username, password string, registers []Register) {
	//_, err := scraper.Login(username, password)
	//if err != nil {
	//	log.Error().Err(err).Msg("cannot login")
	//	return
	//}
	//if err := scraper.GetSessionDKMH(); err != nil {
	//	log.Error().Err(err).Msg("cannot get session")
	//	return
	//}

	scraper.client.SetCommonHeader("Cookie", "ASP.NET_SessionId=ujgct3didhrl0oumqcbpckme; User=19DH110082; UserPW=0596F7A388BA87782E99AB3CB983ADF4; UserID=19DH110082")
	terms, err := scraper.GetTerms("NKH")
	if err != nil {
		log.Error().Err(err).Msg("cannot fetch terms")
		return
	}

	var wg sync.WaitGroup
	for _, term := range terms {
		for _, register := range registers {
			if term.Code == register.Code {
				wg.Add(1)
				go func(register Register, requestId string) {
					defer wg.Done()
					scraper.fetchSubjectAndRegister(register, requestId, "NKH")
				}(register, term.RequestId)
			}
		}
	}
	wg.Wait()
}

func (scraper *HuflitScraper) fetchSubjectAndRegister(register Register, requestId string, registType string) {
	subjects, err := scraper.GetClassStudyUnit(register, requestId, registType)
	if err != nil {
		log.Error().Err(err).Msg("cannot fetch class")
		return
	}

	firstRequestId := ""
	secondRequestId := ""
	for _, subject := range subjects {
		if subject.Code == register.FirstCode {
			firstRequestId = subject.RequestId
		}
		if register.SecondCode != "" && subject.Code == register.SecondCode {
			secondRequestId = subject.RequestId
		}
	}

	registerResp, err := scraper.RegisterSubject(firstRequestId, secondRequestId)
	if err != nil {
		log.Error().Err(err).Msg("cannot register")
		return
	}

	if registerResp.State == true {
		log.Info().Str("name", register.Name).Str("id", register.Code).Str("msg", registerResp.Msg).Msg("register successfully")
	} else {
		log.Info().Str("name", register.Name).Str("id", register.Code).Str("msg", registerResp.Msg).Msg("register failed")
	}

}

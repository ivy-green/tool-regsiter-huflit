package huflit

import (
	"sync"

	"github.com/rs/zerolog/log"
)

func (scraper *HuflitScraper) StartJob(username, password string, registers []Register) {
	_, err := scraper.Login(username, password)
	if err != nil {
		log.Error().Err(err).Msg("cannot login")
		return
	}
	if err := scraper.GetSessionDKMH(); err != nil {
		log.Error().Err(err).Msg("cannot get session")
		return
	}

	//scraper.client.SetCommonHeader("Cookie", "User=21DH110592; UserPW=A5AD3206605D5494DDD2D66E53B97814; ASP.NET_SessionId=tot2qkt2hbbankvb41rdg5h2; UserID=21DH110592")
	terms, err := scraper.GetTerms()
	if err != nil {
		log.Error().Err(err).Msg("cannot fetch terms")
		return
	}

	var wg sync.WaitGroup
	for _, term := range terms {
		for _, register := range registers {
			if term.Code == register.Code {
				wg.Add(1)
				go func() {
					defer wg.Done()
					scraper.fetchSubjectAndRegister(register, term.RequestId)
				}()
			}
		}
	}
	wg.Wait()
}

func (scraper *HuflitScraper) fetchSubjectAndRegister(register Register, requestId string) {
	subjects, err := scraper.GetClassStudyUnit(requestId, "KH")
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
		log.Info().Str("name", register.Name).Str("id", register.Code).Msg("register successfully")
	} else {
		log.Info().Str("name", register.Name).Str("id", register.Code).Str("msg", registerResp.Msg).Msg("register failed")
	}

}

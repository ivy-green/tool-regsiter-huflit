package huflit

import (
	"errors"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/rs/zerolog/log"
)

func (scraper *HuflitScraper) StartJob(workerId int, username, password string, registers []Register, typeId string) error {
	if err := backoff.Retry(func() error {
		_, err := scraper.Login(username, password)
		if err != nil {
			log.Error().Err(err).Int("worker_id", workerId).Msg("cannot login")
			return errors.New("cannot login")
		}
		log.Info().Msg("login successfully")
		return nil
	}, backoff.NewConstantBackOff(5*time.Millisecond)); err != nil {
		return err
	}

	if err := backoff.Retry(func() error {
		if err := scraper.GetSessionDKMH(); err != nil {
			log.Error().Err(err).Int("worker_id", workerId).Msg("cannot get session")
			return errors.New("cannot get session")
		}
		return nil
	}, backoff.NewConstantBackOff(5*time.Millisecond)); err != nil {
		return err
	}

	terms, err := backoff.RetryWithData(func() ([]Term, error) {
		terms, err := scraper.GetTerms(typeId)
		if err != nil {
			log.Error().Err(err).Int("worker_id", workerId).Msg("cannot fetch terms")
			return nil, errors.New("cannot fetch terms")
		}

		return terms, nil
	}, backoff.NewConstantBackOff(5*time.Millisecond))
	if err != nil {
		return err
	}

	c := make(chan Register, len(registers))

	for _, term := range terms {
		for _, register := range registers {
			if term.Code == register.Code {
				register.RequestId = term.RequestId
				c <- register
			}
		}
	}

	for register := range c {
		go func(register Register) {
			defer func() {
				if r := recover(); r != nil {
					log.Info().Interface("val", r).Msg("Recovered")
					c <- register
				}
			}()

			if err := scraper.fetchSubjectAndRegister(register, register.RequestId, typeId); err != nil {
				log.Error().Err(err).Str("name", register.Name).Msg("retry !!!")
				c <- register
			}
		}(register)
	}

	return nil
}

func (scraper *HuflitScraper) fetchSubjectAndRegister(register Register, requestId string, registType string) error {
	log.Info().Str("name", register.Name).Str("id", register.Code).Msg("running !!!")

	subjects, err := backoff.RetryWithData(func() ([]Subject, error) {
		subjects, err := scraper.GetClassStudyUnit(register, requestId, registType)
		if err != nil {
			log.Error().Err(err).Msg("cannot fetch class")
			return nil, errors.New("cannot fetch class")
		}
		return subjects, nil
	}, backoff.NewConstantBackOff(5*time.Millisecond))
	if err != nil {
		return err
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

	registerResp, err := backoff.RetryWithData(func() (*RegisterResp, error) {
		registerResp, err := scraper.RegisterSubject(firstRequestId, secondRequestId)
		if err != nil {
			log.Error().Err(err).Msg("cannot register")
			return nil, errors.New("cannot register")
		}

		return registerResp, nil
	}, backoff.NewConstantBackOff(5*time.Millisecond))
	if err != nil {
		return err
	}

	if !registerResp.State {
		log.Error().Str("name", register.Name).
			Str("id", register.Code).
			Str("resp_message", registerResp.Msg).
			Msg("register failed")

		return errors.New("register failed")
	}

	log.Info().Str("name", register.Name).
		Str("id", register.Code).
		Str("resp_message", registerResp.Msg).
		Msg("register successfully")
	return nil
}

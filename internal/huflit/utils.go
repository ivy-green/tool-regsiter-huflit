package huflit

import (
	"github.com/rs/zerolog/log"
)

func Recovery() {
	if r := recover(); r != nil {
		log.Info().Interface("val", r).Msg("Recovered")
	}
}

package main

import (
	"bufio"
	"flag"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/whoant/subject-huflit/internal/huflit"
)

var (
	username        string
	password        string
	numberOfWorkers int
	filepath        string
)

func main() {
	flag.StringVar(&username, "username", "", "Username")
	flag.StringVar(&password, "password", "", "Password")
	flag.StringVar(&filepath, "file", "", "File path")
	flag.IntVar(&numberOfWorkers, "workers", 5, "Number of workers")
	flag.Parse()
	registers := readSubjectsFromInput(filepath)
	for _, register := range registers {
		log.Info().Str("code", register.Code).Str("name", register.Name).Msg("!!!")
	}

	c := make(chan int, numberOfWorkers)
	for i := 1; i <= numberOfWorkers; i++ {
		c <- i
	}

	for workerId := range c {
		go func(workerId int) {
			defer func() {
				if r := recover(); r != nil {
					log.Info().Interface("val", r).Msg("Recovered")
					c <- workerId
				}
			}()
			client := huflit.NewHuflitScraper()
			if err := client.StartJob(workerId, username, password, registers, "KH"); err != nil {
				log.Info().Int("worker_id", workerId).Msg("retry !!")
				c <- workerId
			}
		}(workerId)
	}
}

func readSubjectsFromInput(filepath string) []huflit.Register {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot read input")
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	res := make([]huflit.Register, 0)
	for scanner.Scan() {
		line := scanner.Text()
		//code|code_detail|name
		a := strings.Split(line, "|")
		res = append(res, huflit.Register{
			Code:       a[0],
			FirstCode:  a[1],
			SecondCode: a[2],
			Name:       a[3],
		})
	}

	return res
}

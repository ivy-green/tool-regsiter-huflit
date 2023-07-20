package main

import (
	"github.com/whoant/subject-huflit/internal/huflit"
)

func main() {
	username := "21DH110592"
	password := "9.G2U2WryL_Q2Jj"
	client := huflit.NewHuflitScraper()
	var registers []huflit.Register
	registers = append(registers, huflit.Register{
		Code:       "1230723",
		Name:       "Đồ án phần mềm",
		FirstCode:  "231123072310",
		SecondCode: "231123072337",
	})
	client.StartJob(username, password, registers)

}

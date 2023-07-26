package main

import (
	"github.com/whoant/subject-huflit/internal/huflit"
)

func main() {
	username := ""
	password := ""
	client := huflit.NewHuflitScraper()
	var registers []huflit.Register
	registers = append(registers, huflit.Register{
		Code:       "1230723",
		Name:       "Đồ án phần mềm",
		FirstCode:  "231123072313",
		SecondCode: "231123072344",
	})
	registers = append(registers, huflit.Register{
		Code:      "1010462",
		Name:      "Chủ nghĩa xã hội khoa học",
		FirstCode: "231101046218",
	})
	client.StartJob(username, password, registers)

}

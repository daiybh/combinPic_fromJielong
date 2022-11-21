package main

import (
	"fmt"
	"log"

	"github.com/xuri/excelize/v2"
)

func readExcelmain() {
	file, err := excelize.OpenFile("test.xlsx")
	if err != nil {
		log.Fatal(err)
	}
	c1, err := file.GetCellValue("Sheet1", "A2")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(c1)
	c2, err := file.GetCellValue("Sheet1", "A3")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(c2)
}

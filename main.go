package main

import (
	"log"
	"os"

	"github.com/xuri/excelize/v2"
)

func main() {
	ep := os.Getenv("EXCEL_PATH")
	if ep == "" {
		log.Fatalln("EXCEL_PATH not defined.")
	}
	xlsx, err := excelize.OpenFile(ep)
	if err != nil {
		log.Fatal("ERROR", err.Error())
	}

}

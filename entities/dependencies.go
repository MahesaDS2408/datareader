package entities

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
)

type Dependencies struct {
	Config *YodelConfig
	Router *chi.Mux
	csvMap *CSVMapWithMutex
}

func (d Dependencies) GetCSVMap() *CSVMap {
	if d.csvMap != nil {
		if d.csvMap.Map != nil {
			return &d.csvMap.Map
		}
	}
	return nil
}

func (d *Dependencies) SetCSVMap() {
	d.csvMap = &CSVMapWithMutex{
		Map: CSVMap{},
	}
	files, err := ioutil.ReadDir(d.Config.FolderIndukan)
	if err != nil {
		log.Fatal(err)
	}

	eg := &errgroup.Group{}
	for _, file := range files {
		if !file.IsDir() {
			eg.Go(func() error {
				path := fmt.Sprintf("%s/%s", d.Config.FolderIndukan, file.Name())
				csv, err := NewCSVBlockFromFile(path)
				fileNameNoExt := strings.Split(file.Name(), ".")[0]
				d.csvMap.Add(fileNameNoExt, csv)
				return err
			})
		}
	}
	if err := eg.Wait(); err != nil {
		log.Fatal("Error", err)
	}
}

type HTTPResponse struct {
	Code     int         `json:"status_code"`
	Messsage string      `json:"status_msg"`
	Content  interface{} `json:"content"`
}

func (d Dependencies) Run() {
	fmt.Println("Loading all csv")
	d.SetCSVMap()

	d.Router = NewHttpRoute(&d)

	fmt.Printf(
		"Datareader use %s config and online at %d \n",
		viper.ConfigFileUsed(),
		d.Config.NomorPort,
	)

	fmt.Println("Support me by sponsoring on GitHub at https://github.com/frederett")
	fmt.Println("and/or give star to https://github.com/MahesaDS2408/datareader.")
	http.ListenAndServe(fmt.Sprintf(":%d", d.Config.NomorPort), d.Router)
}

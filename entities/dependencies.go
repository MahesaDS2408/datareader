package entities

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
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

func NewHttpRoute(d *Dependencies) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/{alamat}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		if d == nil {
			err := HTTPResponse{
				Code:     500,
				Messsage: "CSV Map not loaded, check your config or your csvs folder",
				Content:  "",
			}
			msg, _ := json.Marshal(err)
			w.Write(msg)
			return
		}
		csvs := d.GetCSVMap()
		if csvs == nil {
			err := HTTPResponse{
				Code:     500,
				Messsage: "CSV Map not loaded, check your config or your csvs folder",
				Content:  "",
			}
			msg, _ := json.Marshal(err)
			w.Write(msg)
			return
		}
		alamat := chi.URLParam(r, "alamat")
		if csv := (*csvs)[alamat]; csv != nil {
			queryString := r.URL.Query().Get("allowed")
			if queryString != "" {
				data := FilterField(csv.Data, strings.Split(queryString, ","))
				jsonCsv, _ := json.Marshal(data)
				w.Write(jsonCsv)
			}

			return
		}

		err := HTTPResponse{
			Code:     404,
			Messsage: "CSV block not found for that address",
			Content:  "",
		}
		msg, _ := json.Marshal(err)
		w.Write(msg)
	})

	return r
}

func (d Dependencies) Run() {
	fmt.Println("Loading all csv")
	d.SetCSVMap()

	d.Router = NewHttpRoute(&d)

	fmt.Println(
		fmt.Sprintf(
			"Datareader use %s config and online at %d",
			viper.ConfigFileUsed(),
			d.Config.NomorPort,
		),
	)

	fmt.Println("Support me by sponsoring on GitHub at https://github.com/frederett")
	fmt.Println("and/or give star to https://github.com/MahesaDS2408/datareader.")
	http.ListenAndServe(fmt.Sprintf(":%d", d.Config.NomorPort), d.Router)
}

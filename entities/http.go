package entities

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func NewHttpRoute(d *Dependencies) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		err := HTTPResponse{
			Code:     200,
			Messsage: fmt.Sprintf("File OK!, check filename on %s to access them", d.Config.FolderIndukan),
			Content:  "OK!",
		}
		msg, _ := json.Marshal(err)
		w.Write(msg)
		return
	})

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
				return
			}

			jsonCsv, _ := json.Marshal(csv.Data)
			w.Write(jsonCsv)
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

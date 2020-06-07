// go build -ldflags "-s -w" -o index.cgi cgi.go

package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"net/http"
	"net/http/cgi"
	"sort"
	"strings"
	"wota/sotauploader/db"
	"wota/sotauploader/utils"
)

var err error
var errs strings.Builder
var debugIn = true
var debugDb = false

type WotaOption struct {
	Id   int
	Ref  string
	Name string
}

type CreateFormData struct {
	User        string
	Errors      string
	WotaOptions []WotaOption
	Disabled    bool
}

func main() {
	utils.SetDebugInput(debugIn)
	if err := cgi.Serve(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errs.WriteString("")

		header := w.Header()
		header.Set("Content-Type", "text/html; charset=utf-8")

		authCookie := utils.FindAuthCookie(r.Cookies())

		var dbAllGood = debugDb
		if !debugDb {
			db.Init()
			err = db.Open(db.CmsDb)
			if err == nil {
				err = db.Open(db.WotaDb)
				if err == nil {
					_, err = db.LoadSummits()
					if err == nil {
						dbAllGood = true
					}
				}
			}
			if err != nil {
				errs.WriteString(err.Error())
			}
		}

		if dbAllGood {
			var user, err = db.FindUser(authCookie)
			if user == "" {
				errs.WriteString("You must be logged in to use the SOTA CSV Uploader")
			}

			formData := getFormData(user, errs.String(), getWotaOptions())
			tmpl, err := template.ParseFiles("../templates/sota-upload-form.html")

			if err != nil {
				errs.WriteString(err.Error())
			} else {
				err = tmpl.Execute(w, formData)
				if err != nil {
					errs.WriteString(err.Error())
				}
			}
		}
		if !debugDb {
			db.CloseAll()
		}
	})); err != nil {
		fmt.Println(err)
	}
}

func getFormData(user string, errs string, options []WotaOption) CreateFormData {
	var form CreateFormData
	form.User = user
	form.Errors = errs
	form.WotaOptions = options
	form.Disabled = errs != ""
	return form
}

func getWotaOptions() []WotaOption {
	var options []WotaOption
	for _, summit := range db.GetSummits() {
		// Skip 0 'TBC' as it isn't relevant here
		if summit.WotaId == 0 {
			continue
		}

		var summitId string
		if summit.WotaId <= 214 {
			summitId = fmt.Sprintf("LDW-%03d", summit.WotaId)
		} else {
			summitId = fmt.Sprintf("LDO-%03d", summit.WotaId-214)
		}
		var option WotaOption
		option.Id = summit.WotaId
		option.Ref = summitId
		option.Name = summit.Name
		options = append(options, option)
	}
	sort.Slice(options, func(i, j int) bool {
		return options[i].Name < options[j].Name
	})
	return options
}

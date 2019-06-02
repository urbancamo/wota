// go build -ldflags "-s -w" -o index.cgi cgi.go

package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"net/http"
	"net/http/cgi"
	"strings"
	"wota/sotauploader/db"
	"wota/sotauploader/utils"
)

var err error
var errs strings.Builder
var debugIn = false
var debugDb = false

func main() {
	utils.SetDebugInput(debugIn)
	if err := cgi.Serve(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errs.WriteString("")

		header := w.Header()
		header.Set("Content-Type", "text/html; charset=utf-8")

		authCookie := utils.FindAuthCookie(r.Cookies())

		if !debugDb {
			db.Init()
			db.Open(db.CmsDb)
			db.Open(db.WotaDb)
			db.LoadSummits()
		}

		var user, err = db.FindUser(authCookie)
		if user == "" {
			errs.WriteString("You must be logged in to use the SOTA submit form")
		}
		var b []byte
		b, err = ioutil.ReadFile("sota-upload-form.html")
		s := string(b)
		s = strings.ReplaceAll(s, "$USER", user)
		s = strings.ReplaceAll(s, "$ERRORS", errs.String())
		s = strings.ReplaceAll(s, "$OPTIONS", getWotaOptions())
		disabled := "disabled"
		if user != "unknown" {
			disabled = ""
		}
		if user != "" {
			s = strings.ReplaceAll(s, "$DISABLED", disabled)
		}
		w.Write([]byte(s))
		if err != nil {
			//errs.WriteString(err.Error())
		}
		if !debugDb {
			db.CloseAll()
		}
	})); err != nil {
		fmt.Println(err)
	}
}

func getWotaOptions() string {
	var options strings.Builder
	for _, summit := range db.GetSummits() {
		var summitId string
		if summit.WotaId <= 214 {
			summitId = fmt.Sprintf("LDW-%03d", summit.WotaId)
		} else {
			summitId = fmt.Sprintf("LDO-%03d", summit.WotaId-214)
		}
		options.WriteString(fmt.Sprintf("<option id=\"%s\">%s: %s</option>", summit.WotaId, summitId, summit.Name))
	}
	return options.String()
}

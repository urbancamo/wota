// go build -ldflags "-s -w" -o index.cgi cgi.go

package main

import (
	"bytes"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cgi"
	"os"
	"strings"
	"time"
	"wota/sotauploader/csv"
	"wota/sotauploader/db"
	"wota/sotauploader/utils"
	"wota/sotautils"
)

var err error
var errs strings.Builder
var debugIn = false
var debugDb = false
var dumpForm = false
var UTF8Bom = []byte{0xEF, 0xBB, 0xBF}

const UNKNOWN = "UNKNOWN"

type CheckFormView struct {
	User        string
	Filename    string
	CsvData     string
	Activations []ActivationsView
	Chases      []ChaseView
	SummitCount int
	Errors      string
}

type ActivationsView struct {
	Id         int
	Date       string
	CallUsed   string
	WotaId     string
	SummitName string
	StnCall    string
	S2S        bool
}

type ChaseView struct {
	Id         int
	Date       string
	WorkedBy   string
	WotaId     string
	SummitName string
	StnCall    string
}

func main() {
	utils.SetDebugInput(debugIn)
	if err := cgi.Serve(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user = UNKNOWN
		var filename = UNKNOWN
		var csvData = UNKNOWN
		var summitCount int

		filename, csvData = getFile(r)

		header := w.Header()
		header.Set("Content-Type", "text/html; charset=utf-8")

		authCookie := utils.FindAuthCookie(r.Cookies())
		err = openDbs()
		if err != nil {
			errs.WriteString(err.Error())
		} else {
			summitCount, err = db.LoadSummits()
			if err != nil {
				errs.WriteString(fmt.Sprintf("%s occurred, summit count is %d", err.Error(), summitCount))
			} else {
				if debugIn {
					user = "M0NOM"
				} else {
					user, err = db.FindUser(authCookie)
					if err != nil {
						errs.WriteString(err.Error())
					}
				}
			}

			// Store a copy of the CSV file uploaded
			logFilename := GetFilename("/home/wotasite/logs/sota-uploader", user)
			err = WriteToFile(logFilename, csvData)
			if err != nil {
				errs.WriteString(err.Error())
			}

			var contacts csv.Contacts
			if user == UNKNOWN {
				errs.WriteString("You must be logged in to use the SOTA submit form\n")
			} else {
				var csvParsed bool
				contacts, csvParsed = csv.ParseCsv(csvData, user)
				if !csvParsed {
					errs.WriteString("Error parsing CSV file - please check content\n")
				}
			}

			if !debugDb {
				err = db.CloseAll()
				if err != nil {
					errs.WriteString(err.Error())
				}
			}

			templateLoaded := false

			if err == nil {
				formData := getFormData(user, filename, csvData, contacts, summitCount, errs.String())
				tmpl, err := template.ParseFiles("../templates/sota-upload-check.html")

				if err == nil {
					err = tmpl.Execute(w, formData)
					if err != nil {
						errs.WriteString(err.Error())
					} else {
						templateLoaded = true
					}
				} else {
					errs.WriteString(err.Error())
				}
			}

			if !templateLoaded {
				s := "<html><body>$ERRORS</body></html>"
				s = strings.ReplaceAll(s, "$ERRORS", errs.String())
				_, _ = w.Write([]byte(s))
			}
		}
	})); err != nil {
		fmt.Println(err)
	}
}

func openDbs() error {
	var err error
	if !debugDb {
		db.Init()
		err = db.Open(db.CmsDb)
		if err != nil {
			errs.WriteString(err.Error())
		}

		err = db.Open(db.WotaDb)
		if err != nil {
			errs.WriteString(err.Error())
		}
	}
	return err
}

func getFormData(user string, filename string, csvData string, contacts csv.Contacts, summitCount int, errors string) CheckFormView {
	var formData CheckFormView
	formData.User = user
	formData.Filename = filename
	formData.CsvData = csvData

	for i, activationContact := range contacts.ActivationContacts {
		formData.Activations = append(formData.Activations, getActivationView(i, activationContact))
	}
	for i, chaserContact := range contacts.ChaserContacts {
		formData.Chases = append(formData.Chases, getChaseView(i, chaserContact))
	}
	formData.SummitCount = summitCount
	if errors != "" {
		formData.Errors = errors
	} else {
		formData.Errors = fmt.Sprintf("No errors with %d summits loaded from database", summitCount)
	}
	return formData
}

func getActivationView(id int, contact csv.ActivationContact) ActivationsView {
	var view ActivationsView
	view.Id = id + 1
	view.Date = string(contact.Date)
	view.CallUsed = contact.CallUsed
	view.WotaId = sotautils.GetWotaRefFromId(contact.WotaId)
	view.SummitName = db.GetSummitName(contact.WotaId)
	view.StnCall = contact.StnCall
	view.S2S = contact.S2S
	return view
}

func getChaseView(id int, contact csv.ChaserContact) ChaseView {
	var view ChaseView
	view.Id = id + 1
	view.Date = string(contact.Date)
	view.WorkedBy = contact.WkdBy
	view.WotaId = sotautils.GetWotaRefFromId(contact.WotaId)
	view.SummitName = db.GetSummitName(contact.WotaId)
	view.StnCall = contact.StnCall
	return view
}

func getFile(r *http.Request) (string, string) {
	filename := "UNKNOWN"
	content := "UNKNOWN"

	if debugIn {
		filename = "test/csv/2019-05-18-04-Seat-Sandal-SOTA.csv"
		fileContentInBytes, err := ioutil.ReadFile(filename)
		if err == nil {
			content = string(fileContentInBytes)
		}
	} else {
		readForm, err := r.MultipartReader()
		if err != nil {
			errs.WriteString(err.Error())
		} else {
			for {
				part, errPart := readForm.NextPart()
				if errPart == io.EOF {
					break
				}
				filename = part.FormName()
				if filename == "filename" {
					filename = part.FileName()

					buf := new(bytes.Buffer)
					_, _ = buf.ReadFrom(part)
					// G8CPZ attempted an upload on 14/07/2019 which failed due to CSV file having a UTF BOM
					// marker at the start. Strip the marker off if found
					byteData := buf.Bytes()
					strippedByteData := bytes.TrimPrefix(byteData, UTF8Bom)
					content = string(strippedByteData)
				}

			}
		}
	}
	return filename, content
}

func GetFilename(basePath string, callSign string) string {
	t := time.Now()
	return fmt.Sprintf("%s/%d-%02d-%02dT%02d%02d%02d-%s.csv", basePath,
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second(), callSign)
}

func WriteToFile(filename string, data string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, data)
	if err != nil {
		return err
	}
	return file.Sync()
}

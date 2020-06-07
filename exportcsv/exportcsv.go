package exportcsv

import (
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"net/http"
	"net/http/cgi"
	"strings"
	"wota/sotauploader/db"
	"wota/sotauploader/utils"
)

const ActivationContacts = "Activation Contacts"
const ChaseContacts = "Chase Contacts"

// {"Index":"1","Date":"2019-05-18 09:12:00","Callsign Used":"M0NOM/P","Summit Id":"LDW-003","Summit Name":"Helvellyn","Contact":"M0OAT/P","Summit to Summit":"NNY"}
type ActivationContact struct {
	Index          string `json:"Index"`
	Date           string `json:"Date"`
	CallsignUsed   string `json:"Callsign Used"`
	SummitId       string `json:"Summit Id"`
	SummitName     string `json:"Summit Name"`
	Contact        string `json:"Contact"`
	SummitToSummit string `json:"Summit to Summit"`
}

// {"Index":"1","Date":"2019-05-18 09:12:00","Callsign Used":"M0NOM/P","Summit":"LDW-129","Summit Name":"Illgill Head","Station Worked":"M0OAT"}
type ChaseContact struct {
	Index         string `json:"Index"`
	Date          string `json:"Date"`
	CallsignUsed  string `json:"Callsign Used"`
	Summit        string `json:"Summit"`
	SummitName    string `json:"Summit Name"`
	StationWorked string `json:"Station Worked"`
}

type ExportResult struct {
	Type    string `json:"Type"`
	Results string `json:"Results"`
	Errors  string `json:"Errors"`
}

var err error
var errs strings.Builder
var debugIn = false
var debugDb = false
var debugMsgs = true

func main() {
	if !debugIn {
		handleCgi()
	} else {
		handleExport()
	}
}

func handleCgi() {
	if err = cgi.Serve(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user = "UNKNOWN"
		var summitCount int
		var results strings.Builder

		var uploadType = "UNDETERMINED"

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			errs.WriteString(err.Error())
		}

		header := w.Header()
		header.Set("Content-Type", "application/json; charset=utf-8")

		authCookie := utils.FindAuthCookie(r.Cookies())

		if !debugDb {
			db.Init()
			err = db.Open(db.CmsDb)
			if err != nil {
				errs.WriteString(err.Error())
			}
			results.WriteString("CMS DB Opened... ")
			err = db.Open(db.WotaDb)
			if err != nil {
				errs.WriteString(err.Error())
			}
			results.WriteString("WOTA DB Opened... ")
			//summitCount, err = db.LoadSummits()

			if err != nil {
				errs.WriteString(fmt.Sprintf("%s occurred, summit count is %d", err.Error(), summitCount))
			} else {
				user, err = db.FindUser(authCookie)
				if err != nil {
					errs.WriteString(err.Error())
				} else {
					results.WriteString(fmt.Sprintf("User identified: %s\n", user))
					// Determine if activation or chase data are required
					jsonContent := string(body)
					if strings.Contains(jsonContent, "\"activation") {
						exportType = ActivationContacts
						results.WriteString(process(user, exportType, body))
					} else {
						exportType = ChaseContacts
						results.WriteString(process(user, exportType, body))
					}
					results.WriteString("Processing Complete\n")
				}
			}
		}

		// Marshal results and errors into a json object for return
		var jsonData []byte
		var jsonUploadData ExportResult
		jsonUploadData.Type = uploadType
		if err != nil {
			jsonUploadData.Errors = err.Error()
		} else {
			jsonUploadData.Errors = "No errors"
		}
		jsonUploadData.Results = results.String()

		jsonData, err = json.Marshal(jsonUploadData)
		if err != nil {
			errs.WriteString(err.Error())
		}

		// Store a copy of the json result
		logFilename := GetFilename("/home/wotasite/logs/sota-uploader", user)
		err = WriteToFile(logFilename, jsonData)
		if err != nil {
			errs.WriteString(err.Error())
		}

		w.Write(jsonData)
		if !debugDb {
			db.CloseAll()
		}

	})); err != nil {
		fmt.Println(err)
	}
}

func process(user string, exportType string, body []byte) string {
	if !debugDb {
		err := db.PrepareStatementsForInsert()
		if err != nil {
			return err.Error()
		}
	}

	if uploadType == ActivationContacts {
		var activationContacts []ActivationContact
		err := json.Unmarshal(body, &activationContacts)
		if err != nil {
			return err.Error()
		}
		return processActivationContacts(user, activationContacts)
	} else if uploadType == ChaseContacts {
		var chaseContacts []ChaseContact
		err := json.Unmarshal(body, &chaseContacts)
		if err != nil {
			return err.Error()
		}
		return processChaseContacts(user, chaseContacts)
	} else {
		errs.WriteString(fmt.Sprintf("Unknown upload type: %s", uploadType))
	}
	return errs.String()
}

func processActivationContacts(user string, activationContacts []ActivationContact) string {
	var errs strings.Builder
	for _, contact := range activationContacts {
		errs.WriteString(insertActivationContact(user, contact))
	}
	return errs.String()
}

func processChaseContacts(user string, chaseContacts []ChaseContact) string {
	var errs strings.Builder
	for _, contact := range chaseContacts {
		errs.WriteString(insertChaseContact(user, contact))
	}
	return errs.String()
}


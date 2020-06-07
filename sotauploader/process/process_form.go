package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cgi"
	"os"
	"strings"
	"time"
	"wota/domain"
	"wota/sotauploader/db"
	"wota/sotauploader/utils"
	"wota/sotautils"
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

type UploadResult struct {
	Type    string `json:"Type"`
	Results string `json:"Results"`
	Errors  string `json:"Errors"`
}

var err error
var errs strings.Builder
var debugIn = false
var debugDb = false

func main() {
	if !debugIn {
		handleCgi()
	} else {
		handleFile()
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
					// Determine if this is activation or chase data
					jsonContent := string(body)
					if strings.Contains(jsonContent, "\"Summit to Summit\":") {
						uploadType = ActivationContacts
						results.WriteString(process(user, uploadType, body))
					} else {
						uploadType = ChaseContacts
						results.WriteString(process(user, uploadType, body))
					}
					results.WriteString("Processing Complete\n")
				}
			}
		}

		// Marshal results and errors into a json object for return
		var jsonData []byte
		var jsonUploadData UploadResult
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

func handleFile() string {
	var user = "M0NOM"
	var fileContentInBytes []byte
	var processErrs string
	fileContentInBytes, err = ioutil.ReadFile("json/activationTest.json")

	if !debugDb {
		db.Init()
		err = db.Open(db.CmsDb)
		if err != nil {
			errs.WriteString(err.Error())
		} else {
			err = db.Open(db.WotaDb)
			if err != nil {
				errs.WriteString(err.Error())
			}
		}
	}

	if err == nil {
		processErrs = process(user, ActivationContacts, fileContentInBytes)
		if processErrs != "" {
			errs.WriteString(processErrs)
		}
	} else {
		errs.WriteString(err.Error())
	}

	fileContentInBytes, err = ioutil.ReadFile("json/chaseTest.json")
	if err == nil {
		processErrs = process(user, ChaseContacts, fileContentInBytes)
		if processErrs != "" {
			errs.WriteString(processErrs)
		}
	} else {
		errs.WriteString(err.Error())
	}
	return errs.String()
}

func process(user string, uploadType string, body []byte) string {
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
		var chaseContacts []domain.ChaseContact
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

func processChaseContacts(user string, chaseContacts []domain.ChaseContact) string {
	var errs strings.Builder
	for _, contact := range chaseContacts {
		errs.WriteString(insertChaseContact(user, contact))
	}
	return errs.String()
}

func insertActivationContact(user string, contact ActivationContact) string {
	if debugDb {
		_, _ = fmt.Printf("INSERT ACTIVATION - callsignUsed: %s, date: %s, contact: %s, summitId: %s, s2s: %s\n", contact.CallsignUsed, contact.Date, contact.Contact, contact.SummitId, contact.SummitToSummit)
	} else {
		count, err := db.InsertActivation(user, sotautils.GetOperatorFromCallsign(contact.CallsignUsed), contact.Date, contact.Contact, sotautils.GetWotaIdFromRef(contact.SummitId), contact.SummitToSummit)
		if count != 1 {
			return fmt.Sprintf("Could not insert: %s", getActivationDebugLine(contact))
		} else if err != nil {
			return err.Error()
		} else {
			return fmt.Sprintf("Inserted Activation Contact - Callsign: %s, Date: %s, Contact: %s, Summit: %s, Summit to Summit: %s\n", contact.CallsignUsed, contact.Date, contact.Contact, contact.SummitId, contact.SummitToSummit)
		}
	}
	return ""
}

func getActivationDebugLine(contact ActivationContact) string {
	return fmt.Sprintf("Activation callsignUsed: %s, date: %s, contact: %s, summitId: %s, s2s: %s\n", contact.CallsignUsed, contact.Date, contact.Contact, contact.SummitId, contact.SummitToSummit)
}

func insertChaseContact(user string, contact domain.ChaseContact) string {
	if debugDb {
		_, _ = fmt.Printf("INSERT CHASE - callsignUsed: %s, date: %s, summit: %s, stationWorked: %s\n", contact.CallsignUsed, contact.Date, contact.Summit, contact.StationWorked)
	} else {
		count, err := db.InsertChase(user, sotautils.GetOperatorFromCallsign(contact.CallsignUsed), contact.Date, sotautils.GetWotaIdFromRef(contact.Summit), contact.StationWorked)
		if count != 1 {
			return fmt.Sprintf("Could not insert: %s", getChaseDebugLine(contact))
		} else if err != nil {
			return err.Error()
		} else {
			return fmt.Sprintf("Inserted Chase Contact - Callsign: %s, Date: %s, Summit: %s, Station Worked: %s\n", contact.CallsignUsed, contact.Date, contact.Summit, contact.StationWorked)
		}
	}
	return ""
}

func getChaseDebugLine(contact domain.ChaseContact) string {
	return fmt.Sprintf("Chase callsignUsed: callsignUsed: %s, date: %s, summit: %s, stationWorked: %s\n", contact.CallsignUsed, contact.Date, contact.Summit, contact.StationWorked)
}

func GetFilename(basepath string, callsign string) string {
	t := time.Now()
	return fmt.Sprintf("%s/%d-%02d-%02dT%02d%02d%02d-%s.json", basepath,
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second(), callsign)
}

func WriteToFile(filename string, data []byte) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, string(data))
	if err != nil {
		return err
	}
	return file.Sync()
}

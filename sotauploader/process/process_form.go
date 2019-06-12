package main

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
	"wota/sotautils"
)

const ActivationContacts = "activationContacts"
const ChaseContacts = "chaseContacts"

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

var err error
var errs strings.Builder
var debugIn = true
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

		query := r.URL.Query()
		uploadType := query.Get("type")
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			errs.WriteString(err.Error())
		}

		header := w.Header()
		header.Set("Content-Type", "text/html; charset=utf-8")

		authCookie := utils.FindAuthCookie(r.Cookies())

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

			//summitCount, err = db.LoadSummits()
			if err != nil {
				errs.WriteString(fmt.Sprintf("%s occurred, summit count is %d", err.Error(), summitCount))
			} else {
				user, err = db.FindUser(authCookie)
				if err != nil {
					errs.WriteString(err.Error())
				} else {
					process(user, uploadType, body)
				}
			}
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
		callsignUsed := strings.ToUpper(contact.CallsignUsed)
		// Sanity check - make sure the contact is relevant to the user attempting to upload the file
		if strings.Contains(callsignUsed, user) {
			errs.WriteString(insertActivationContact(user, contact))
		}
	}
	return errs.String()
}

func processChaseContacts(user string, chaseContacts []ChaseContact) string {
	var errs strings.Builder
	for _, contact := range chaseContacts {
		callsignUsed := strings.ToUpper(contact.CallsignUsed)
		// Sanity check - make sure the contact is relevant to the user attempting to upload the file
		if strings.Contains(callsignUsed, user) {
			errs.WriteString(insertChaseContact(user, contact))
		}
	}
	return errs.String()
}

func insertActivationContact(user string, contact ActivationContact) string {
	if debugDb {
		_, _ = fmt.Printf("INSERT ACTIVATION - callsignUsed: %s, date: %s, contact: %s, summitId: %s, s2s: %s\n", contact.CallsignUsed, contact.Date, contact.Contact, contact.SummitId, contact.SummitToSummit)
	} else {
		count, err := db.InsertActivation(user, user, contact.Date, contact.Contact, sotautils.GetWotaIdFromRef(contact.SummitId), contact.SummitToSummit)
		if count != 1 {
			return fmt.Sprintf("Could not insert: %s", getActivationDebugLine(contact))
		} else if err != nil {
			return err.Error()
		}
	}
	return ""
}

func getActivationDebugLine(contact ActivationContact) string {
	return fmt.Sprintf("Activation callsignUsed: %s, date: %s, contact: %s, summitId: %s, s2s: %s\n", contact.CallsignUsed, contact.Date, contact.Contact, contact.SummitId, contact.SummitToSummit)
}

func insertChaseContact(user string, contact ChaseContact) string {
	if debugDb {
		_, _ = fmt.Printf("INSERT CHASE - callsignUsed: %s, date: %s, summit: %s, stationWorked: %s\n", contact.CallsignUsed, contact.Date, contact.Summit, contact.StationWorked)
	} else {
		count, err := db.InsertChase(user, contact.CallsignUsed, contact.Date, sotautils.GetWotaIdFromRef(contact.Summit), contact.StationWorked)
		if count != 1 {
			return fmt.Sprintf("Could not insert: %s", getChaseDebugLine(contact))
		} else if err != nil {
			return err.Error()
		}
	}
	return ""
}

func getChaseDebugLine(contact ChaseContact) string {
	return fmt.Sprintf("Chase callsignUsed: callsignUsed: %s, date: %s, summit: %s, stationWorked: %s\n", contact.CallsignUsed, contact.Date, contact.Summit, contact.StationWorked)
}

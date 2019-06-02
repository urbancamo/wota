// go build -ldflags "-s -w" -o index.cgi cgi.go

package main

import (
	"bytes"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cgi"
	"strings"
	"wota/sotauploader/csv"
	"wota/sotauploader/db"
	"wota/sotauploader/utils"
	"wota/sotautils"
)

var err error
var errs strings.Builder
var debugIn = false
var debugDb = false
var debugStages = false
var dumpForm = false

func main() {
	utils.SetDebugInput(debugIn)
	if err := cgi.Serve(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user = "UNKNOWN"
		var filename = "UNKNOWN"
		var csvData = "UNKNOWN"
		var summitCount int

		if debugStages {
			errs.WriteString(" Stage 0 ")
		}
		filename, csvData = getFile(r)

		header := w.Header()
		header.Set("Content-Type", "text/html; charset=utf-8")

		if debugStages {
			errs.WriteString("2 ")
		}
		authCookie := utils.FindAuthCookie(r.Cookies())
		if debugStages {
			errs.WriteString("3 ")
		}

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
			summitCount, err = db.LoadSummits()
			if err != nil {
				errs.WriteString(fmt.Sprintf("%s occurred, summit count is %d", err.Error(), summitCount))
			} else {
				user, err = db.FindUser(authCookie)
				if err != nil {
					errs.WriteString(err.Error())
				}
			}
		}
		if debugStages {
			errs.WriteString("4 ")
		}

		var contacts csv.Contacts
		ok := true
		if user == "" {
			errs.WriteString("\nYou must be logged in to use the SOTA submit form")
		} else {
			contacts, ok = csv.ParseCsv(csvData, user)
			if !ok {
				errs.WriteString("Error parsing CSV file - please check content")
			}
		}

		if debugStages {
			errs.WriteString("5 ")
		}
		if !debugDb {
			err = db.CloseAll()
			if err != nil {
				errs.WriteString(err.Error())
			}
		}
		if debugStages {
			errs.WriteString("6 ")
		}

		var fileContentInBytes []byte
		fileContentInBytes, err = ioutil.ReadFile("sota-upload-check.html")
		if err == nil {
			if debugStages {
				errs.WriteString("7 ")
			}
			s := string(fileContentInBytes)
			s = strings.ReplaceAll(s, "$USER", user)
			s = strings.ReplaceAll(s, "$FILENAME", filename)
			s = strings.ReplaceAll(s, "$CONTENT", csvData)
			s = strings.ReplaceAll(s, "$ACTIVATION_COUNT", fmt.Sprintf("%d", len(contacts.ActivationContacts)))
			s = strings.ReplaceAll(s, "$ACTIVATION_CONTACTS", getActivationTable(contacts.ActivationContacts))
			s = strings.ReplaceAll(s, "$CHASER_COUNT", fmt.Sprintf("%d", len(contacts.ChaserContacts)))
			s = strings.ReplaceAll(s, "$CHASER_CONTACTS", getChaseTable(contacts.ChaserContacts, len(contacts.ActivationContacts) > 0))
			errors := fmt.Sprintf("No errors with %d summits loaded from database", summitCount)
			if errs.String() != "" {
				errors = errs.String()
			}
			s = strings.ReplaceAll(s, "$ERRORS", errors)
			disabled := "disabled"
			if user != "unknown" {
				disabled = ""
			}
			if user != "" {
				s = strings.ReplaceAll(s, "$DISABLED", disabled)
			}
			w.Write([]byte(s))
			if debugStages {
				errs.WriteString("8 ")
			}
		} else {
			s := "<html><body>$ERRORS</body></html>"
			s = strings.ReplaceAll(s, "$ERRORS", errs.String())
			_, _ = w.Write([]byte(s))
		}

		if err != nil {
			//errs.WriteString(err.Error())
		}
	})); err != nil {
		fmt.Println(err)
	}
}

func getActivationTable(contacts []csv.ActivationContact) string {
	var out strings.Builder

	if len(contacts) > 0 {
		for id, contact := range contacts {
			contactId := fmt.Sprintf("activationContact-%04d", id)
			out.WriteString("<tr id=\"" + contactId + "\">")
			out.WriteString(tableColumn(id, "date", contact.Date))
			out.WriteString(tableColumn(id, "callUsed", contact.CallUsed+"/P"))
			out.WriteString(tableColumn(id, "wotaId", sotautils.GetWotaRefFromId(contact.WotaId)))
			out.WriteString(tableColumn(id, "summitName", db.GetSummitName(contact.WotaId)))
			out.WriteString(tableColumn(id, "stnCall", contact.StnCall))
			if contact.S2S {
				out.WriteString(tableColumn(id, "s2s", "<input type=\"checkbox\" checked></checked>"))
			} else {
				out.WriteString(tableColumn(id, "s2s", "<input type=\"checkbox\" ></checked>"))
			}
			out.WriteString("</tr>")
		}
		return out.String()
	} else {
		return "<td span=\"5\">No activation contacts detected</td>"
	}
}

func getChaseTable(contacts []csv.ChaserContact, hasActivationContacts bool) string {
	var out strings.Builder

	if len(contacts) > 0 {
		for id, contact := range contacts {
			contactId := fmt.Sprintf("chaserContact-%04d", id)
			out.WriteString("<tr id=\"" + contactId + "\">")
			out.WriteString(tableColumn(id, "date", contact.Date))
			if hasActivationContacts {
				out.WriteString(tableColumn(id, "wkdBy", contact.WkdBy+"/P"))
			} else {
				out.WriteString(tableColumn(id, "wkdBy", contact.WkdBy))
			}
			out.WriteString(tableColumn(id, "wotaId", sotautils.GetWotaRefFromId(contact.WotaId)))
			out.WriteString(tableColumn(id, "summitName", db.GetSummitName(contact.WotaId)))
			out.WriteString(tableColumn(id, "stnCall", contact.StnCall))
			out.WriteString("</tr>")
		}
	} else {
		return "<td span=\"4\">No chaser contacts detected</td>"
	}
	return out.String()
}

func tableColumn(row int, name string, content string) string {
	var id = fmt.Sprintf("%s-%04d", name, row)
	return fmt.Sprintf("<td id='%s'>%s</td>", id, content)
}

func getFile(r *http.Request) (string, string) {
	filename := "UNKNOWN"
	content := "UNKNOWN"

	read_form, err := r.MultipartReader()
	if err != nil {
		errs.WriteString(err.Error())
	} else {
		for {
			part, errPart := read_form.NextPart()
			if errPart == io.EOF {
				break
			}
			filename = part.FormName()
			if filename == "filename" {
				filename = part.FileName()

				buf := new(bytes.Buffer)
				buf.ReadFrom(part)
				content = buf.String()
			}

		}
	}
	return filename, content
}

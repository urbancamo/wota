package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"wota/sotautils"
)

type SotaAlert struct {
	Id                 int
	TimeStamp          string
	DateActivated      string
	AssociationCode    string
	SummitCode         string
	SummitDetails      string
	PosterCallsign     string
	ActivatingCallsign string
	ActivatorName      string
	Frequency          string
	Comments           string
}

type WotaAlert struct {
	Id       int
	WotaId   int
	DateTime string
	Call     string
	FreqMode string
	Comment  string
	PostedBy string
}

var db *sql.DB
var err error
var debugIn = false
var debugDb = false
var debugMsgs = true

func main() {
	if !debugDb {
		db, err = sql.Open("mysql", WotaDb)
		defer db.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
		err = db.Ping()
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	sotautils.LoadSotaWotaMapId()
	var SotaAlerts = getAlerts()
	var wotaSpots = convertSotaToWotaAlerts(SotaAlerts)
	fmt.Printf("\n\n\n%+v", wotaSpots)
	if len(wotaSpots) > 0 && !debugDb {
		updateAlertsInDb(wotaSpots)
	}
}

func getAlerts() []SotaAlert {
	var alerts []SotaAlert
	if debugIn {
		var alert = createDummyAlert()
		alerts = append(alerts, alert)
	} else {
		response, err := http.Get("http://api2.sota.org.uk/api/alerts")
		if err != nil {
			log.Fatal(err)
		} else {
			defer response.Body.Close()
			body, err2 := ioutil.ReadAll(response.Body)
			if err2 != nil {
				fmt.Println("error:", err2)
			}
			if debugMsgs {
				fmt.Printf("Read %d characters from Swagger API\n", len(body))
			}
			err3 := json.Unmarshal(body, &alerts)
			if err2 != nil {
				fmt.Println("error:", err3)
			}
		}
	}
	return alerts
}

func updateAlertsInDb(wotaAlerts []WotaAlert) {
	for _, alert := range wotaAlerts {
		// check to see if we've added the alert already
		id := 0
		err = db.QueryRow("SELECT `id` FROM `alerts` WHERE `datetime` = ? and `call` = ? and `wotaid` = ?",
			alert.DateTime, alert.Call, alert.WotaId).Scan(&id)
		if id == 0 {
			// add alert if not already there
			_, err = db.Query("INSERT INTO `alerts`(`datetime`, `call`, `wotaid`, `freqmode`, `comment`, `postedby`) VALUES (?, ?, ?, ?, ?, ?)",
				alert.DateTime, alert.Call, alert.WotaId, alert.FreqMode, alert.Comment, alert.PostedBy)
			if err != nil {
				fmt.Println("error:", err)
			}
			fmt.Printf("Alert %s %s %d added\n", alert.DateTime, alert.Call, alert.WotaId)
		} else {
			fmt.Printf("Alert %s %s %d already added, ignoring\n", alert.DateTime, alert.Call, alert.WotaId)
		}
	}

}

func convertSotaToWotaAlerts(SotaAlerts []SotaAlert) []WotaAlert {
	var wotaAlerts []WotaAlert

	for _, alert := range SotaAlerts {
		if alert.AssociationCode == "G" && strings.Split(alert.SummitCode, "-")[0] == "LD" {
			if sotautils.GetWotaIdFromSummitCode(alert.SummitCode) != 0 {
				if debugMsgs {
					fmt.Printf("Creating WOTA alert for summit: %s/%s\n", alert.AssociationCode, alert.SummitCode)
				}
				wotaAlerts = append(wotaAlerts, convertSotaToWotaAlert(alert))
			}
		}
	}
	return wotaAlerts
}

func convertSotaToWotaAlert(SotaAlert SotaAlert) WotaAlert {
	var wotaAlert WotaAlert

	wotaAlert.DateTime = strings.Split(strings.ReplaceAll(SotaAlert.DateActivated, "T", " "), ".")[0]

	wotaAlert.Call = SotaAlert.ActivatingCallsign
	wotaAlert.WotaId = sotautils.GetWotaIdFromSummitCode(SotaAlert.SummitCode)

	wotaAlert.FreqMode = SotaAlert.Frequency
	var commentLen = len(SotaAlert.Comments)
	if commentLen > 79 {
		commentLen = 79
	}
	wotaAlert.Comment = SotaAlert.Comments[0:commentLen]
	if len(wotaAlert.Comment) < 79-len("[SOTA>WOTA] ") {
		// add a header if we have enough room
		wotaAlert.Comment = "[SOTA>WOTA] " + SotaAlert.Comments
	}
	wotaAlert.PostedBy = SotaAlert.PosterCallsign
	return wotaAlert
}

func createDummyAlert() SotaAlert {
	var alert SotaAlert
	alert.Id = 12345
	alert.TimeStamp = "2019-06-15T19:06:59.999"
	alert.Comments = "TEST PLEASE IGNORE"
	alert.PosterCallsign = "G1OHH"
	alert.AssociationCode = "G"
	alert.SummitCode = "LD-056"
	alert.ActivatorName = "Mark"
	alert.ActivatingCallsign = "M0NOM/P"
	alert.Frequency = "14.285-FM"
	return alert
}

//{
// "id":161795,
// "timeStamp":"2019-03-11T09:52:04",
// "dateActivated":"2019-08-11T12:00:00",
// "associationCode":"G",
// "summitCode":"LD-018",
// "summitDetails":"Stony Cove Pike, 763m, 6 point(s)",
// " posterCallsign":"G4YTD",
// "activatingCallsign":"G4YTD/P",
// "activatorName":"Tim","frequency":
// "145.500-fm,144.300-ssb,7.118-ssb",
// "comments":"Time appx. Alert when ready. QRO + 5 ele on 2m"}

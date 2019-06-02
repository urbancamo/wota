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

type SotaSpot struct {
	Id                int
	Timestamp         string
	Comments          string
	Callsign          string
	AssociationCode   string
	SummitCode        string
	ActivatorCallsign string
	ActivatorName     string
	Frequency         string
	Mode              string
	SummitDetails     string
	HighlightColor    string
}

type WotaSpot struct {
	Id       int
	DateTime string
	Call     string
	WotaId   int
	FreqMode string
	Comment  string
	Spotter  string
}

var db *sql.DB
var err error
var debugIn = false
var debugDb = false

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
	var sotaSpots = getSpots()
	var wotaSpots = convertSotaToWotaSpots(sotaSpots)
	fmt.Printf("\n\n\n%+v", wotaSpots)
	if len(wotaSpots) > 0 && !debugDb {
		updateSpotsInDb(wotaSpots)
	}
}

func getSpots() []SotaSpot {
	var spots []SotaSpot
	if debugIn {
		var spot = createDummySpot()
		spots = append(spots, spot)
	} else {
		response, err := http.Get("http://api2.sota.org.uk/api/spots/1")
		if err != nil {
			log.Fatal(err)
		} else {
			defer response.Body.Close()
			body, err2 := ioutil.ReadAll(response.Body)
			if err2 != nil {
				fmt.Println("error:", err2)
			}
			err3 := json.Unmarshal(body, &spots)
			if err2 != nil {
				fmt.Println("error:", err3)
			}
		}
	}
	return spots
}

func updateSpotsInDb(wotaSpots []WotaSpot) {
	for _, spot := range wotaSpots {
		// check to see if we've added the spot already
		id := 0
		err = db.QueryRow("SELECT `id` FROM `spots` WHERE `datetime` = ? and `call` = ? and `wotaid` = ?",
			spot.DateTime, spot.Call, spot.WotaId).Scan(&id)
		if id == 0 {
			// add spot if not already there
			_, err = db.Query("INSERT INTO `spots`(`datetime`, `call`, `wotaid`, `freqmode`, `comment`, `spotter`) VALUES (?, ?, ?, ?, ?, ?)",
				spot.DateTime, spot.Call, spot.WotaId, spot.FreqMode, spot.Comment, spot.Spotter)
			if err != nil {
				fmt.Println("error:", err)
			}
			fmt.Printf("Spot %s %s %d added\n", spot.DateTime, spot.Call, spot.WotaId)
		} else {
			fmt.Printf("Spot %s %s %d already added, ignoring\n", spot.DateTime, spot.Call, spot.WotaId)
		}
	}

}

func convertSotaToWotaSpots(sotaSpots []SotaSpot) []WotaSpot {
	var wotaSpots []WotaSpot

	for _, spot := range sotaSpots {
		if spot.AssociationCode == "G" && strings.Split(spot.SummitCode, "-")[0] == "LD" {
			if sotautils.GetWotaIdFromSotaCode(spot.SummitCode) != 0 {
				wotaSpots = append(wotaSpots, convertSotaToWotaSpot(spot))
			}
		}
	}
	return wotaSpots
}

func convertSotaToWotaSpot(sotaSpot SotaSpot) WotaSpot {
	var wotaSpot WotaSpot

	wotaSpot.DateTime = strings.Split(strings.ReplaceAll(sotaSpot.Timestamp, "T", " "), ".")[0]

	wotaSpot.Call = sotaSpot.ActivatorCallsign
	wotaSpot.WotaId = sotautils.GetWotaIdFromSotaCode(sotaSpot.SummitCode)

	//if wotaId <= 214 {
	//	wotaSpot.WotaId = fmt.Sprintf("LDW-%03d", wotaId)
	//} else {
	//	wotaSpot.WotaId = fmt.Sprintf("LDO-%03d", wotaId-214)
	//}
	wotaSpot.FreqMode = sotaSpot.Frequency + "-" + sotaSpot.Mode
	var commentLen = len(sotaSpot.Comments)
	if commentLen > 79 {
		commentLen = 79
	}
	wotaSpot.Comment = sotaSpot.Comments[0:commentLen]
	if len(wotaSpot.Comment) < 79-len("[SOTA>WOTA] ") {
		// add a header if we have enough room
		wotaSpot.Comment = "[SOTA>WOTA] " + sotaSpot.Comments
	}
	wotaSpot.Spotter = sotaSpot.Callsign
	return wotaSpot
}

func createDummySpot() SotaSpot {
	var spot SotaSpot
	spot.Id = 12345
	spot.Timestamp = "2019-05-21T19:06:59.999"
	spot.Comments = "TEST PLEASE IGNORE"
	spot.Callsign = "G1OHH"
	spot.AssociationCode = "G"
	spot.SummitCode = "LD-056"
	spot.ActivatorName = "Mark"
	spot.ActivatorCallsign = "M0NOM/P"
	spot.Frequency = "14.285"
	spot.Mode = "ssb"
	return spot
}

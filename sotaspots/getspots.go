package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
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

type SotaWota struct {
	Sotaid string
	Wotaid string
}

var sotaToWotaJson = []byte(`[
{"sotaid":"000001","wotaid":"000001"},
{"sotaid":"000003","wotaid":"000003"},
{"sotaid":"000004","wotaid":"000004"},
{"sotaid":"000005","wotaid":"000007"},
{"sotaid":"000006","wotaid":"000008"},
{"sotaid":"000007","wotaid":"000013"},
{"sotaid":"000008","wotaid":"000014"},
{"sotaid":"000009","wotaid":"000020"},
{"sotaid":"000010","wotaid":"000022"},
{"sotaid":"000011","wotaid":"000025"},
{"sotaid":"000012","wotaid":"000029"},
{"sotaid":"000013","wotaid":"000030"},
{"sotaid":"000014","wotaid":"000032"},
{"sotaid":"000015","wotaid":"000040"},
{"sotaid":"000017","wotaid":"000049"},
{"sotaid":"000018","wotaid":"000055"},
{"sotaid":"000019","wotaid":"000056"},
{"sotaid":"000020","wotaid":"000063"},
{"sotaid":"000021","wotaid":"000067"},
{"sotaid":"000022","wotaid":"000069"},
{"sotaid":"000023","wotaid":"000082"},
{"sotaid":"000024","wotaid":"000086"},
{"sotaid":"000025","wotaid":"000093"},
{"sotaid":"000026","wotaid":"000104"},
{"sotaid":"000027","wotaid":"000108"},
{"sotaid":"000028","wotaid":"000112"},
{"sotaid":"000029","wotaid":"000129"},
{"sotaid":"000031","wotaid":"000140"},
{"sotaid":"000033","wotaid":"000147"},
{"sotaid":"000034","wotaid":"000151"},
{"sotaid":"000035","wotaid":"000155"},
{"sotaid":"000036","wotaid":"000168"},
{"sotaid":"000037","wotaid":"000173"},
{"sotaid":"000040","wotaid":"000184"},
{"sotaid":"000041","wotaid":"000190"},
{"sotaid":"000042","wotaid":"000196"},
{"sotaid":"000043","wotaid":"000203"},
{"sotaid":"000044","wotaid":"000210"},
{"sotaid":"000047","wotaid":"000211"},
{"sotaid":"000051","wotaid":"000213"},
{"sotaid":"000030","wotaid":"000216"},
{"sotaid":"000032","wotaid":"000219"},
{"sotaid":"000045","wotaid":"000273"},
{"sotaid":"000048","wotaid":"000278"},
{"sotaid":"000050","wotaid":"000281"},
{"sotaid":"000053","wotaid":"000297"},
{"sotaid":"000054","wotaid":"000303"},
{"sotaid":"000055","wotaid":"000317"},
{"sotaid":"000056","wotaid":"000323"}
]`)

var db *sql.DB
var err error
var debugIn = true
var debugDb = true
var sotaWotaIdMap map[string]string

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
	loadSotaWotaMapId()
	var sotaSpots = getSpots()
	var wotaSpots = convertSotaToWotaSpots(sotaSpots)
	fmt.Printf("\n\n\n%+v", wotaSpots)
	if len(wotaSpots) > 0 && !debugDb {
		updateSpotsInDb(wotaSpots)
	}
}

func loadSotaWotaMapId() {
	sotaWotaIdMap = make(map[string]string)
	var sotaWotaData []SotaWota
	json.Unmarshal(sotaToWotaJson, &sotaWotaData)
	for _, sotaWotaId := range sotaWotaData {
		sotaWotaIdMap[sotaWotaId.Sotaid] = sotaWotaId.Wotaid
	}
}

func getSpots() []SotaSpot {
	var spots []SotaSpot
	if debugIn {
		var spot = createDummySpot()
		spots = append(spots, spot)
	} else {
		response, err := http.Get("http://api2.sota.org.uk/api/spots/1/{filter}?filter=all")
		if err != nil {
			log.Fatal(err)
		} else {
			defer response.Body.Close()
			body, err := ioutil.ReadAll(response.Body)
			err2 := json.Unmarshal(body, &spots)
			if err != nil {
				fmt.Println("error:", err2)
			}
		}
	}
	return spots
}

func updateSpotsInDb(wotaSpots []WotaSpot) {
	for _, spot := range wotaSpots {
		// check to see if we've added the spot already
		id := 0
		err = db.QueryRow("SELECT id FROM spots WHERE datetime = ? and call = ? and wotaid = ?",
			spot.DateTime, spot.Call, spot.WotaId).Scan(&id)
		if id == 0 {
			// add spot if not already there
			_, err = db.Query("INSERT INTO `spots`(`datetime`, `call`, `wotaid`, `freqmode`, `comment`, `spotter`) VALUES (?, ?, ?, ?, ?, ?)",
				spot.DateTime, spot.Call, spot.WotaId, spot.FreqMode, spot.Comment, spot.Spotter)
			if err != nil {
				fmt.Println("error:", err)
			}
		}
	}

}

func convertSotaToWotaSpots(sotaSpots []SotaSpot) []WotaSpot {
	var wotaSpots []WotaSpot

	for _, spot := range sotaSpots {
		wotaId := getWotaIdFromSotaId(spot.SummitCode)
		if spot.AssociationCode == "GLD" && wotaId != "" {
			wotaSpots = append(wotaSpots, convertSotaToWotaSpot(spot))
		}
	}
	return wotaSpots
}

func convertSotaToWotaSpot(sotaSpot SotaSpot) WotaSpot {
	var wotaSpot WotaSpot
	wotaSpot.DateTime = sotaSpot.Timestamp
	wotaSpot.Call = sotaSpot.ActivatorCallsign
	wotaSpot.WotaId, _ = strconv.Atoi(getWotaIdFromSotaId(sotaSpot.SummitCode))
	//if wotaId <= 214 {
	//	wotaSpot.WotaId = fmt.Sprintf("LDW-%03d", wotaId)
	//} else {
	//	wotaSpot.WotaId = fmt.Sprintf("LDO-%03d", wotaId-214)
	//}
	wotaSpot.FreqMode = sotaSpot.Frequency + "-" + sotaSpot.Mode
	wotaSpot.Comment = sotaSpot.Comments
	wotaSpot.Spotter = sotaSpot.Callsign
	return wotaSpot
}

func getWotaIdFromSotaId(summitCode string) string {
	sotaSummitNumber, _ := strconv.Atoi(strings.SplitAfter(summitCode, "-")[1])
	summitRef := fmt.Sprintf("%06d", sotaSummitNumber)
	return sotaWotaIdMap[summitRef]
}

func createDummySpot() SotaSpot {
	var spot SotaSpot
	spot.Id = 12345
	spot.Timestamp = "2019-05-21T19:06:00.000"
	spot.Comments = "TEST PLEASE IGNORE"
	spot.Callsign = "G1OHH"
	spot.AssociationCode = "GLD"
	spot.SummitCode = "LD-056"
	spot.ActivatorName = "Mark"
	spot.ActivatorCallsign = "M0NOM/P"
	spot.Frequency = "14.285"
	spot.Mode = "ssb"
	return spot
}

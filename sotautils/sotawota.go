package sotautils

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

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

var sotaWotaIdMap map[int]int

func LoadSotaWotaMapId() {
	sotaWotaIdMap = make(map[int]int)

	var sotaWotaData []SotaWota
	json.Unmarshal(sotaToWotaJson, &sotaWotaData)
	for _, sotaWota := range sotaWotaData {
		var sotaId, _ = strconv.Atoi(sotaWota.Sotaid)
		var wotaId, _ = strconv.Atoi(sotaWota.Wotaid)
		sotaWotaIdMap[sotaId] = wotaId
	}
}

func GetWotaIdFromSotaCode(summitCode string) int {
	if summitCode == "" {
		return 0
	}
	summitParts := strings.Split(summitCode, "-")
	if len(summitParts) != 2 {
		return 0
	}
	if summitParts[0] != "G/LD" {
		return 0
	}

	sotaSummitNumber, _ := strconv.Atoi(strings.SplitAfter(summitCode, "-")[1])
	return sotaWotaIdMap[sotaSummitNumber]
}

func GetWotaIdFromSotaId(sotaId int) int {
	return sotaWotaIdMap[sotaId]
}

func GetWotaRefFromId(wotaId int) string {
	if wotaId > 214 {
		return fmt.Sprintf("LDO-%03d", wotaId-214)
	} else {
		return fmt.Sprintf("LDW-%03d", wotaId)
	}
}

func ConvertSotaDate(sotaDate string, sotaTime string) string {
	dateFields := strings.Split(sotaDate, "/")
	// Just check that we have three fields, if not try a - as the separator
	if len(dateFields) == 1 {
		dateFields = strings.Split(sotaDate, "-")
	}
	day := dateFields[0]
	month := dateFields[1]
	year := dateFields[2]
	if len(year) == 2 {
		year = "20" + year
	}

	var hour, min string
	if len(sotaTime) == 3 {
		hour = sotaTime[0:1]
		min = sotaTime[1:3]
	} else {
		hour = sotaTime[0:2]
		min = sotaTime[2:4]
	}
	return fmt.Sprintf("%04s-%02s-%02s %02s:%02s:00", year, month, day, hour, min)
}

func ConvertSotaYear(sotaDate string) int {
	dateFields := strings.Split(sotaDate, "/")
	// Just check that we have three fields, if not try a - as the separator
	if len(dateFields) == 1 {
		dateFields = strings.Split(sotaDate, "-")
	}
	year, _ := strconv.Atoi(dateFields[2])
	return year + 2000
}

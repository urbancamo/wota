package db

import (
	"database/sql"
	"fmt"
	"strconv"
)

type Summit struct {
	WotaId    int
	SotaId    int
	Book      string
	Name      string
	Height    int64
	Reference string
	HumpId    int
	GridId    string
}

var summitsMap map[int]Summit
var summitList []Summit

func LoadSummits() (int, error) {
	var err error
	err = getSummitsFromDb()
	return len(summitsMap), err
}

func getSummitsFromDb() error {
	summitsMap = make(map[int]Summit)

	rows, err := dbMap[WotaDb].Query("SELECT `wotaid`, `sotaid`, `book`, `name`, `height`, `reference`, `humpid`, `gridid` FROM `summits`")
	if err != nil {
		return err
	}

	for rows.Next() {
		var summit Summit
		var wotaId, sotaId, book, name, reference, humpId, gridId sql.NullString
		var height sql.NullInt64
		err := rows.Scan(&wotaId, &sotaId, &book, &name, &height, &reference, &humpId, &gridId)
		if err != nil {
			return err
		} else {
			summit.WotaId = convertToIntNullIsZero(wotaId)
			summit.SotaId = convertToIntNullIsZero(sotaId)
			summit.Book = convertToStringNullIsBlank(book)
			summit.Name = convertToStringNullIsBlank(name)
			summit.Height = intNullIsZero(height)
			summit.Reference = convertToStringNullIsBlank(reference)
			summit.HumpId = convertToIntNullIsZero(humpId)
			summit.GridId = convertToStringNullIsBlank(gridId)
			summitsMap[summit.WotaId] = summit
			summitList = append(summitList, summit)
		}
	}
	return err
}

func intNullIsZero(val sql.NullInt64) int64 {
	if val.Valid {
		return val.Int64
	}
	return 0
}

func convertToStringNullIsBlank(val sql.NullString) string {
	if val.Valid {
		return val.String
	}
	return ""
}

func convertToIntNullIsZero(val sql.NullString) int {
	if val.Valid {
		rtn, err := strconv.Atoi(val.String)
		if err == nil {
			return rtn
		}
	}
	return 0
}

func GetSummitName(wotaId int) string {
	if val, ok := summitsMap[wotaId]; ok {
		return val.Name
	}
	return fmt.Sprintf("Unknown summit: %d", wotaId)
}

func GetSummits() []Summit {
	return summitList
}

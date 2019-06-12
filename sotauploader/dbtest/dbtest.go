package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	"wota/sotauploader/db"
)

func main() {
	var summitCount int
	var err error
	var errs strings.Builder

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
		if err != nil {
			errs.WriteString(err.Error())
		}
	}
	db.CloseAll()
}

package db

import "database/sql"

var dbMap map[string]*sql.DB

func Init() {
	dbMap = make(map[string]*sql.DB)
}

func FindUser(sessionId string) (user string, err error) {
	err = dbMap[CmsDb].QueryRow("SELECT `username` FROM `cms_module_feusers_loggedin` JOIN `cms_module_feusers_users` users ON userid = id WHERE `sessionid` = ? ",
		sessionId).Scan(&user)
	return user, err
}

func Open(dburl string) error {
	var err error

	dbMap[dburl], err = sql.Open("mysql", dburl)
	if err != nil {
		return err
	}
	err = dbMap[dburl].Ping()
	return err
}

func Close(dburl string) error {
	return dbMap[dburl].Close()
}

func CloseAll() error {
	var err error

	for _, db := range dbMap {
		err = db.Close()
	}
	return err
}

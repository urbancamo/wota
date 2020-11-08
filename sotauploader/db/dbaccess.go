package db

import (
	"database/sql"
	"strings"
	"wota/sotautils"
)

const ActivationInsertSql = "INSERT INTO `activator_log`(`activatedby`, `callused`, `wotaid`, `date`, `time`, `year`, `stncall`, `ucall`, `s2s`, `confirmed`) VALUES (?,?,?,?,?,?,?,?,?,?)"
const ActivationSelectSql = "SELECT COUNT(*) FROM `activator_log` WHERE `activatedby`=? AND `wotaid`=? AND `stncall`=? AND `date`=?"

const ChaseInsertSql = "INSERT INTO `chaser_log`(`wkdby`, `ucall`, `wotaid`, `date`, `time`, `year`, `stncall`) VALUES (?,?,?,?,?,?,?)"
const ChaseSelectSql = "SELECT COUNT(*) FROM `chaser_log` WHERE `wkdby` = ? AND `wotaid` = ? AND `stncall` = ? AND `date` = ?"
const ChaseSetConfirmedSql = "UPDATE `chaser_log` SET `confirmed` = true WHERE `ucall` = ? AND `wotaid` = ? AND `stncall` = ? AND `date` = ?"

const ChaseWorkedSummitSql = "SELECT COUNT(*) FROM `chaser_log` WHERE `wkdby` = ? AND `wotaid` = ?"
const ChaseWorkedActivatorOnSummitSql = "SELECT COUNT(*) FROM `chaser_log` WHERE `wkdby` = ? AND `wotaid` = ? AND `stncall` = ?"
const ChaseCheckSumnitAnnualPointsSql = "SELECT COUNT(*) FROM `chaser_log` WHERE `year` = ? AND `wkdby` = ? AND `wotaid` = ?"
const ChaseCheckSummitActivatorAnnualPointsSql = "SELECT COUNT(*) FROM `chaser_log` WHERE `year` = ? AND `wkdby` = ? AND `wotaid` = ? AND `stncall` = ?"
const ChaseCheckActivationConfirmationSql = "SELECT COUNT(*) FROM `activator_log` WHERE `callused` = ? AND `wotaid` = ? AND `ucall` = ? AND `date` = ?"
const ActivationSetConfirmedSql = "UPDATE `activator_log` SET `confirmed` = true WHERE `callused` = ? AND `wotaid` = ? AND `ucall` = ? AND `date` = ? AND `time` = ?"
const ChaseInsertWithPointsSql = "INSERT into `chaser_log` (`wkdby`, `ucall`, `wotaid`, `date`, `time`, `year`, `stncall`, `points`, `wawpoints`, `points_yr`, `wawpoints_yr`, `confirmed`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)"

const ExportChaserSql = "SELECT `wkdby`,`ucall`,`wotaid`,`date`,`year`,`stncall`,`confirmed` FROM `chaser_log` WHERE `wkdby` = ? ORDER BY `date`"
const ExportActivatorSql = "SELECT `activatedby`,`callused`,`wotaid`,`date`,`year`,`stncall`,`s2s`,`confirmed` FROM `activator_log` WHERE `activatedby` = ? ORDER BY `date`"

// CHASE SQL
//$wawpoints = 1; $points = 1;
//$sql = "SELECT * FROM `chaser_log` WHERE `wkdby` = '".$wkdby."' AND `wotaid` = '".$summitid."'";
//$result = mysql_query($sql,$con);
//if($result && (mysql_num_rows($result) > 0)) $wawpoints = 0;
//$sql = "SELECT * FROM `chaser_log` WHERE `wkdby` = '".$wkdby."' AND `wotaid` = '".$summitid."' AND `stncall` = '".$stncall."'";
//$result = mysql_query($sql,$con);
//if($result && (@mysql_num_rows($result) > 0)) $points = 0;
// check annual points
//$wawpoints_yr = 1; $points_yr = 1;
//$sql = "SELECT * FROM `chaser_log` WHERE `year` = ".$year." AND `wkdby` = '".$wkdby."' AND `wotaid` = '".$summitid."'";
//$result = mysql_query($sql,$con);
//if($result && (mysql_num_rows($result) > 0)) $wawpoints_yr = 0;
//$sql = "SELECT * FROM `chaser_log` WHERE `year` = ".$year." AND `wkdby` = '".$wkdby."' AND `wotaid` = '".$summitid."' AND `stncall` = '".$stncall."'";
//$result = mysql_query($sql,$con);
//if($result && (@mysql_num_rows($result) > 0)) $points_yr = 0;
// check confirmation
//$sql = "SELECT * FROM `activator_log` WHERE `callused` = '".$stncall."' AND `wotaid` = '".$summitid."' AND `ucall` = '".ucall($wkdby)."' AND `date` = '".$date."'";
//$result = mysql_query($sql,$con);
//if($result && (@mysql_num_rows($result) > 0)) {
//$cfmmsg = " confirmed and";
//mysql_query("UPDATE `activator_log` SET `confirmed` = true WHERE `callused` = '".$stncall."' AND `wotaid` = '".$summitid."' AND `ucall` = '".$wkdby."' AND `date` = '".$date."'");
//$cfm = "true";
// update
//} else {
//$cfmmsg = "";
//$cfm = "false";
//}

var dbMap map[string]*sql.DB
var activationInsertSql *sql.Stmt
var activationSetConfirmedSql *sql.Stmt
var chaseInsertSql *sql.Stmt
var chaseInsertWithPointsSql *sql.Stmt
var chaseSetConfirmedSql *sql.Stmt
var chaseCheckActivationConfirmationSql *sql.Stmt

func Init() {
	dbMap = make(map[string]*sql.DB)
}

func FindUser(sessionId string) (user string, err error) {
	err = dbMap[CmsDb].QueryRow("SELECT `username` FROM `cms_module_feusers_loggedin` JOIN `cms_module_feusers_users` users ON userid = id WHERE `sessionid` = ? ",
		sessionId).Scan(&user)
	return strings.ToUpper(user), err
}

func PrepareStatementsForInsert() error {
	var err error
	activationInsertSql, err = dbMap[WotaDb].Prepare(ActivationInsertSql)
	if err != nil {
		return err
	}
	chaseInsertSql, err = dbMap[WotaDb].Prepare(ChaseInsertSql)
	if err != nil {
		return err
	}
	chaseInsertWithPointsSql, err = dbMap[WotaDb].Prepare(ChaseInsertWithPointsSql)
	if err != nil {
		return err
	}
	activationSetConfirmedSql, err = dbMap[WotaDb].Prepare(ActivationSetConfirmedSql)
	if err != nil {
		return err
	}
	chaseCheckActivationConfirmationSql, err = dbMap[WotaDb].Prepare(ChaseCheckActivationConfirmationSql)
	if err != nil {
		return err
	}
	chaseSetConfirmedSql, err = dbMap[WotaDb].Prepare(ChaseSetConfirmedSql)
	return err
}

func confirmChase(user string, summitId int, callUsed string, date string) error {
	_, err := chaseSetConfirmedSql.Exec(user, summitId, callUsed, date, date)
	return err
}

func confirmActivation(user string, summitId int, callUsed string, date string) error {
	_, err := activationSetConfirmedSql.Exec(user, summitId, callUsed, date)
	return err
}

func hasActivation(user string, date string, stnCall string, summitId int) (bool, error) {
	var count int = 0
	err := dbMap[WotaDb].QueryRow(ActivationSelectSql, user, summitId, stnCall, date).Scan(&count)
	return count > 0, err
}

func hasChase(user string, summitId int, callUsed string, date string) (bool, error) {
	var count int = 0
	err := dbMap[WotaDb].QueryRow(ChaseSelectSql, user, summitId, callUsed, date).Scan(&count)
	return count > 0, err
}

func hasWorkedActivatorOnSummit(user string, summitId int, contact string) (bool, error) {
	var count int = 0
	err := dbMap[WotaDb].QueryRow(ChaseWorkedActivatorOnSummitSql, user, summitId, contact).Scan(&count)
	return count > 0, err
}

func hasWorkedSummit(user string, summitId int) (bool, error) {
	var count int = 0
	err := dbMap[WotaDb].QueryRow(ChaseWorkedSummitSql, user, summitId).Scan(&count)
	return count > 0, err
}

func hasChaseSummitAnnualPoints(date string, wkdby string, summitId int) (bool, error) {
	year := date[0:4]
	var count int = 0
	err := dbMap[WotaDb].QueryRow(ChaseCheckSumnitAnnualPointsSql, year, wkdby, summitId).Scan(&count)
	return count > 0, err
}

func hasChaseSummitActivatorAnnualPoints(date string, wkdby string, summitId int, stnCall string) (bool, error) {
	year := date[0:4]
	var count int = 0
	err := dbMap[WotaDb].QueryRow(ChaseCheckSummitActivatorAnnualPointsSql, year, wkdby, summitId, stnCall).Scan(&count)
	return count > 0, err
}

func hasWorkedChaserFromSummit(callUsed string, summitId int, user string, date string) (bool, error) {
	var count int = 0
	err := dbMap[WotaDb].QueryRow(ChaseCheckActivationConfirmationSql, callUsed, summitId, user, date).Scan(&count)
	return count > 0, err
}

func InsertActivation(user string, callsignUsed string, date string, contact string, summitId int, summitToSummit string) (int64, error) {
	s2s := false
	year := date[0:4] // '2020-11-07 10:07:00'
	//  01234567890123456789
	//  00000000001111111111
	time := date[11:19] // '2020-11-07 10:07:00'
	//  01234567890123456789
	//  00000000001111111111

	if summitToSummit == "Y" {
		s2s = true
	}

	// Check whether this activation record is already in the database
	hasAct, err := hasActivation(user, truncateToDay(date), sotautils.GetOperatorFromCallsign(contact), summitId)
	if err != nil {
		return 0, err
	}

	// If it is pretend we've inserted it
	if hasAct {
		return 1, nil
	}

	// See if there is the chase end of this contact, if there is then confirm it now rather than waiting for the
	// batch job to run
	var chase bool
	chase, err = hasChase(contact, summitId, user, truncateToDay(date))
	if err != nil {
		return 1, nil
	}

	// If not insert the record
	// INSERT INTO `activator_log`(`activatedby`, `callused`, `wotaid`, `date`, `time`, `year`, `stncall`, `ucall`, `s2s`)
	result, err := activationInsertSql.Exec(user, callsignUsed, summitId, date, time, year, sotautils.GetOperatorFromCallsign(contact), contact, boolToInt(s2s), boolToInt(chase))
	if err != nil {
		return 0, err
	}

	// confirm the chase now
	if chase {
		err = confirmChase(contact, summitId, user, date)
		if err != nil {
			return 1, nil
		}
	}

	return result.RowsAffected()
}

func InsertChase(user string, callsignUsed string, date string, summitId int, stationWorked string) (int64, error) {
	year := date[0:4]
	time := date[11:19]

	var err error

	// Check whether this chase record is already in the database
	var chaseExists bool
	chaseExists, err = hasChase(user, summitId, stationWorked, truncateToDay(date))
	if err != nil {
		return 1, nil
	}

	if chaseExists {
		return 1, nil
	}

	var workedAllWainwrightsPoints = 1
	var workedSummit bool
	workedSummit, err = hasWorkedSummit(user, summitId)
	if err != nil {
		return 0, err
	}
	if workedSummit {
		workedAllWainwrightsPoints = 0
	}

	var summitPoints = 1
	var workedSummitActivator bool
	workedSummitActivator, err = hasWorkedActivatorOnSummit(user, summitId, stationWorked)
	if err != nil {
		return 0, err
	}
	if workedSummitActivator {
		summitPoints = 0
	}

	var workedAllWainwrightsAnnualPoints = 1
	var summitAnnualPoints = 1

	var hasSummitAnnualPoints bool
	hasSummitAnnualPoints, err = hasChaseSummitAnnualPoints(truncateToDay(date), user, summitId)
	if err != nil {
		return 0, err
	}
	if hasSummitAnnualPoints {
		workedAllWainwrightsAnnualPoints = 0
	}

	var hasSummitActivatorAnnualPoints bool
	hasSummitActivatorAnnualPoints, err = hasChaseSummitActivatorAnnualPoints(truncateToDay(date), user, summitId, stationWorked)
	if err != nil {
		return 0, err
	}
	if hasSummitActivatorAnnualPoints {
		summitAnnualPoints = 0
	}

	var confirmed = false
	confirmed, err = hasWorkedChaserFromSummit(stationWorked, summitId, callsignUsed, truncateToDay(date))
	if err != nil {
		return 0, err
	}

	result, err := chaseInsertWithPointsSql.Exec(user, callsignUsed, summitId, date, time, year, stationWorked, summitPoints, workedAllWainwrightsPoints, summitAnnualPoints, workedAllWainwrightsAnnualPoints, boolToInt(confirmed))
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func Open(dbUrl string) error {
	var err error

	dbMap[dbUrl], err = sql.Open("mysql", dbUrl)
	if err != nil {
		return err
	}
	err = dbMap[dbUrl].Ping()
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

func boolToInt(val bool) int {
	if val {
		return 1
	}
	return 0
}

func truncateToDay(date string) string {
	return date[0:10]
}

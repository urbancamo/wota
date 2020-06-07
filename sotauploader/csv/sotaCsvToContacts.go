package csv

import (
	"encoding/csv"
	"errors"
	"io"
	"regexp"
	"strings"
	"wota/sotautils"
)

// To log an activation we need:
// Date of activation
// Callsign used
// Summit (WOTA)
// List of contacts

type ActivationContact struct {
	ActivatedBy string
	CallUsed    string
	WotaId      int
	Date        string
	Year        int
	StnCall     string
	UCall       string
	S2S         bool
	Errors      string
}

type ChaserContact struct {
	WkdBy   string
	UCall   string
	WotaId  int
	Date    string
	Year    int
	StnCall string
	Errors  string
}

type Contacts struct {
	ActivationContacts []ActivationContact
	ChaserContacts     []ChaserContact
}

// V2,M0NOM/P,G/LD-022,18/05/19,1347,144MHz,FM,M0OAT/P,G/LD-059,Thx for contact from Seat Sandal
type SotaContact struct {
	Version       string
	CallUsed      string
	MySummitId    string
	Date          string
	Time          string
	Frequency     string
	Mode          string
	StnWorked     string
	TheirSummitId string
	Comment       string
	Error         string
}

// Position of fields in the Version 1 SOTA CSV file
const V1_CALL_USED = 0
const V1_DATE = 1
const V1_TIME = 2
const V1_MY_SOTA_ID = 3
const V1_FREQUENCY = 4
const V1_MODE = 5
const V1_STN_WORKED = 6
const V1_COMMENT = 7

// Position of fields in the Version 2+ SOTA CSV file
const V2P_VERSION = 0
const V2P_CALL_USED = 1
const V2P_MY_SOTA_ID = 2
const V2P_DATE = 3
const V2P_TIME = 4
const V2P_FREQUENCY = 5
const V2P_MODE = 6
const V2P_STN_WORKED = 7
const V2P_THEIR_SOTA_ID = 8
const V2P_COMMENT = 9

func ParseCsv(csvData string, operator string) (Contacts, bool) {
	r := csv.NewReader(strings.NewReader(csvData))
	sotautils.LoadSotaWotaMapId()

	sotaContacts := parseCsvForSotaContacts(r)

	// If there any errors in the Sota Contacts, don't proceed

	var contacts Contacts

	ok := Ok(sotaContacts)
	if Ok(sotaContacts) {
		activationContacts := parseSotaContactsForActivationContact(sotaContacts, operator)
		chaserContacts := parseSotaContactsForChaserContacts(sotaContacts, operator)
		contacts.ActivationContacts = activationContacts
		contacts.ChaserContacts = chaserContacts
	}
	return contacts, ok
}

func Ok(sotaContacts []SotaContact) bool {
	ok := true
	for _, contact := range sotaContacts {
		if contact.Error != "" {
			ok = false
			break
		}
	}
	return ok
}

func parseCsvForSotaContacts(r *csv.Reader) []SotaContact {
	var sotaContacts []SotaContact

	for {
		var contact SotaContact

		fields, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			contact.Error = err.Error()
		}

		contact, err = readSotaContactFromRecord(fields)
		if err != nil {
			contact.Error = err.Error()
		}

		sotaContacts = append(sotaContacts, contact)
	}
	return sotaContacts
}

func readSotaContactFromRecord(fields []string) (SotaContact, error) {
	var contact SotaContact
	if len(fields) < 8 {
		return contact, errors.New("CSV file doesn't contain at least eight CSV fields as required")
	}

	// Is this a version 2+ SOTA CSV file?
	v2pFile := strings.ToUpper(fields[V2P_VERSION]) == "V2"

	if v2pFile {
		contact.Version = strings.ToUpper(fields[V2P_VERSION])
		contact.CallUsed = strings.ToUpper(fields[V2P_CALL_USED])
		contact.MySummitId = strings.ToUpper(fields[V2P_MY_SOTA_ID])
		contact.Date = fields[V2P_DATE]
		contact.Time = fields[V2P_TIME]
		contact.Frequency = strings.ToUpper(fields[V2P_FREQUENCY])
		contact.Mode = strings.ToUpper(fields[V2P_MODE])
		contact.StnWorked = strings.ToUpper(fields[V2P_STN_WORKED])
		if len(fields) > 8 {
			contact.TheirSummitId = strings.ToUpper(fields[V2P_THEIR_SOTA_ID])
		}
		if len(fields) > 9 {
			contact.Comment = fields[V2P_COMMENT]
		}
	} else {
		contact.CallUsed = strings.ToUpper(fields[V1_CALL_USED])
		contact.Date = fields[V1_DATE]
		contact.Time = fields[V1_TIME]
		contact.MySummitId = strings.ToUpper(fields[V1_MY_SOTA_ID])
		contact.Frequency = strings.ToUpper(fields[V1_FREQUENCY])
		contact.Mode = strings.ToUpper(fields[V1_MODE])
		contact.StnWorked = strings.ToUpper(fields[V1_STN_WORKED])
		if len(fields) > 7 {
			contact.Comment = fields[V1_COMMENT]
		}

	}
	return contact, nil
}

func parseSotaContactsForActivationContact(contacts []SotaContact, operator string) []ActivationContact {
	var activationContacts []ActivationContact

	for _, contact := range contacts {
		// Only interested for Activation contacts if we are on a WOTA summit
		if checkCallUsedIsUs(contact.CallUsed, operator) && checkSummitIsAWota(contact.MySummitId) {

			// Candidate contact, fill in as we go, might not make it to the end
			var activationContact ActivationContact

			activationContact.WotaId = sotautils.GetWotaIdFromSummitCode(contact.MySummitId)
			activationContact.Date = sotautils.ConvertSotaDate(contact.Date, contact.Time)
			activationContact.Year = sotautils.ConvertSotaYear(contact.Date)

			activationContact.ActivatedBy = operator
			activationContact.CallUsed = contact.CallUsed
			activationContact.UCall = operator
			activationContact.StnCall = contact.StnWorked
			activationContact.UCall = sotautils.GetOperatorFromCallsign(contact.StnWorked)

			// Is this a S2S?
			activationContact.S2S = checkSummitIsAWota(contact.TheirSummitId)
			if !activationContact.S2S {
				// Check the comments for a valid WOTA Id
				activationContact.S2S = checkStringContainsAWotaRef(contact.Comment)
			}

			activationContacts = append(activationContacts, activationContact)
		}

	}
	return activationContacts
}

func checkSummitIsAWota(summitCode string) bool {
	isWota := sotautils.GetWotaIdFromSummitCode(summitCode) != 0

	// Maybe you they have substituted a WOTA code here instead?
	if !isWota {
		if strings.Contains(summitCode, "LDW-") || strings.Contains(summitCode, "LDO-") {
			isWota = sotautils.GetWotaIdFromRef(summitCode) != 0
		}
	}
	return isWota
}

func checkStringContainsAWotaRef(comment string) bool {
	return getWotaRefFromString(comment) != ""
}

func getWotaRefFromString(comment string) string {
	regexp, _ := regexp.Compile("WOTA: LD[WO]-[0-9][0-9][0-9]")
	wotaSubString := regexp.FindString(comment)
	if len(wotaSubString) == 13 {
		wotaRef := wotaSubString[6:13]
		if wotaRef != "" {
			// Look up the reference
			if checkSummitIsAWota(wotaRef) {
				return wotaRef
			}
		}
	}
	return ""
}

// The example of G4YSS using GX0OOO has proved that this check isn't valid
// the WOTA form doesn't do any checks, so I guess we've just got to roll with it here!
func checkCallUsedIsUs(callUsed string, operator string) bool {
	//return strings.Contains(callUsed, operator)
	// ignore any checks and assume activators know what they're doing!
	return true
}

func parseSotaContactsForChaserContacts(contacts []SotaContact, operator string) []ChaserContact {
	var chaserContacts []ChaserContact

	for _, contact := range contacts {
		// Only interested for Chaser contacts if we are on not on a WOTA summit, but the
		// station worked is on a WOTA summit
		// Note that we might be on a SOTA summit however that isn't a WOTA
		if checkCallUsedIsUs(contact.CallUsed, operator) &&
			(checkSummitIsAWota(contact.TheirSummitId) || checkStringContainsAWotaRef(contact.Comment)) {

			// Candidate contact, fill in as we go, might not make it to the end
			var chaserContact ChaserContact
			chaserContact.WotaId = sotautils.GetWotaIdFromSummitCode(contact.TheirSummitId)
			if chaserContact.WotaId == 0 {
				// Their WOTA must be in the comment
				chaserContact.WotaId = sotautils.GetWotaIdFromSummitCode(getWotaRefFromString(contact.Comment))
			}
			chaserContact.Date = sotautils.ConvertSotaDate(contact.Date, contact.Time)
			chaserContact.Year = sotautils.ConvertSotaYear(contact.Date)

			chaserContact.WkdBy = contact.CallUsed
			chaserContact.UCall = operator
			chaserContact.StnCall = sotautils.GetOperatorFromCallsign(contact.StnWorked)

			chaserContacts = append(chaserContacts, chaserContact)
		}

	}
	return chaserContacts
}

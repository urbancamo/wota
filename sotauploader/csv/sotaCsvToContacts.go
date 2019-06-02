package csv

import (
	"encoding/csv"
	"errors"
	"io"
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
	Version     string
	CallUsed    string
	MySotaId    string
	Date        string
	Time        string
	Frequency   string
	Mode        string
	StnWorked   string
	TheirSotaId string
	Comment     string
	Error       string
}

const VERSION = 0
const CALL_USED = 1
const MY_SOTA_ID = 2
const DATE = 3
const TIME = 4
const FREQUENCY = 5
const MODE = 6
const STN_WORKED = 7
const THEIR_SOTA_ID = 8
const COMMENT = 9

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
	contact.Version = strings.ToUpper(fields[VERSION])
	contact.CallUsed = strings.ToUpper(fields[CALL_USED])
	contact.MySotaId = strings.ToUpper(fields[MY_SOTA_ID])
	contact.Date = fields[DATE]
	contact.Time = fields[TIME]
	contact.Frequency = strings.ToUpper(fields[FREQUENCY])
	contact.Mode = strings.ToUpper(fields[MODE])
	contact.StnWorked = strings.ToUpper(fields[STN_WORKED])
	if len(fields) > 8 {
		contact.TheirSotaId = strings.ToUpper(fields[THEIR_SOTA_ID])
	}
	if len(fields) > 9 {
		contact.Comment = fields[COMMENT]
	}
	return contact, nil
}

func parseSotaContactsForActivationContact(contacts []SotaContact, operator string) []ActivationContact {
	var activationContacts []ActivationContact

	for _, contact := range contacts {
		// Only interested for Activation contacts if we are on a WOTA summit
		if checkCallUsedIsUs(contact.CallUsed, operator) && checkSummitIsAWota(contact.MySotaId) {

			// Candidate contact, fill in as we go, might not make it to the end
			var activationContact ActivationContact

			activationContact.WotaId = sotautils.GetWotaIdFromSotaCode(contact.MySotaId)
			activationContact.Date = sotautils.ConvertSotaDate(contact.Date, contact.Time)
			activationContact.Year = sotautils.ConvertSotaYear(contact.Date)

			activationContact.ActivatedBy = operator
			activationContact.CallUsed = operator
			activationContact.UCall = operator
			activationContact.StnCall = contact.StnWorked
			activationContact.UCall = sotautils.GetOperatorFromCallsign(contact.StnWorked)

			// Is this a S2S?
			activationContact.S2S = checkSummitIsAWota(contact.TheirSotaId)

			activationContacts = append(activationContacts, activationContact)
		}

	}
	return activationContacts
}

func checkSummitIsAWota(sotaId string) bool {
	sotaFromWota := sotautils.GetWotaIdFromSotaCode(sotaId) != 0
	// OK, so maybe we're trying to be clever here and have substituted a WOTA reference instead
	// TODO
	return sotaFromWota
}

func checkCallUsedIsUs(callUsed string, operator string) bool {
	return strings.Contains(callUsed, operator)
}

func parseSotaContactsForChaserContacts(contacts []SotaContact, operator string) []ChaserContact {
	var chaserContacts []ChaserContact

	for _, contact := range contacts {
		// Only interested for Chaser contacts if we are on not on a WOTA summit, but the
		// station worked is on a WOTA summit
		// Note that we might be on a SOTA summit however that isn't a WOTA
		if checkCallUsedIsUs(contact.CallUsed, operator) &&
			checkSummitIsAWota(contact.TheirSotaId) {

			// Candidate contact, fill in as we go, might not make it to the end
			var chaserContact ChaserContact
			chaserContact.WotaId = sotautils.GetWotaIdFromSotaCode(contact.TheirSotaId)
			chaserContact.Date = sotautils.ConvertSotaDate(contact.Date, contact.Time)
			chaserContact.Year = sotautils.ConvertSotaYear(contact.Date)

			chaserContact.WkdBy = operator
			chaserContact.UCall = operator
			chaserContact.StnCall = sotautils.GetOperatorFromCallsign(contact.StnWorked)

			chaserContacts = append(chaserContacts, chaserContact)
		}

	}
	return chaserContacts
}

package main

import (
	"fmt"
	"io/ioutil"
	"wota/sotauploader/csv"
)

var testUpload1 = `V2,M0NOM/P,G/LD-022,18/05/19,1347,144MHz,FM,M0OAT/P,G/LD-059,Thx for contact from Seat Sandal
V2,M0NOM/P,G/LD-022,18/05/19,1351,144MHz,FM,G4VPX/P,G/LD-053,Thx for contact from Seat Sandal
V2,M0NOM/P,G/LD-022,18/05/19,1354,144MHz,FM,G1OHH,,Thx for contact from Seat Sandal
V2,M0NOM/P,G/LD-022,18/05/19,1357,144MHz,FM,G1ZAR/P,G/LD-007,Thx for contact from Seat Sandal
V2,M0NOM/P,G/LD-022,18/05/19,1405,144MHz,FM,G4OBK/P,G/LD-017,Thx for contact from Seat Sandal
V2,M0NOM/P,G/LD-022,18/05/19,1425,14MHz,SSB,G0TDM,,Thx for contact from Seat Sandal
V2,M0NOM/P,G/LD-022,18/05/19,1428,144MHz,FM,G4WHA/P,,Latrigg  LDW-206
V2,M0NOM/P,G/LD-022,18/05/19,1432,5MHz,SSB,MI0JLA/P,GI/SW-007,Thx for contact from Seat Sandal
V2,M0NOM/P,G/LD-022,18/05/19,1436,144MHz,FM,G6PJZ/P,G/LD-057,Thx for contact from Seat Sandal`

var testUpload2 = `V2,M0NOM/P,,17/04/2019,1148,144MHz,FM,2E0MIX/P,G/LD-029,WOTA: LDW-129
V2,M0NOM/P,,17/04/2019,1150,144MHz,FM,M0HQO/P,G/LD-029,WOTA: LDW-129`

var testUpload3 = `V2,G4VPX/P,,17/04/2019,1148,144MHz,FM,2E0MIX/P,G/LD-029,WOTA: LDW-129
V2,G4VPX/P,,17/04/2019,1150,144MHz,FM,M0HQO/P,G/LD-029,WOTA: LDW-129`

func main() {
	var contacts csv.Contacts
	var err error

	// broken file resilience
	contacts, err = loadAndParseCsv("testSotaCsvToContacts.go")
	if err != nil {
		fmt.Printf("%s", err.Error())
	}

	contacts, _ = loadAndParseCsv("csv/2018-10-28-Place-Fell-SOTA.csv")
	checkContactNumberIsCorrect("activator", len(contacts.ActivationContacts), 0)
	checkContactNumberIsCorrect("chaser", len(contacts.ChaserContacts), 3)

	contacts, _ = loadAndParseCsv("csv/2018-05-19-Fairfield-SOTA.csv")
	checkContactNumberIsCorrect("activator", len(contacts.ActivationContacts), 0)
	checkContactNumberIsCorrect("chaser", len(contacts.ChaserContacts), 3)

	contacts, _ = loadAndParseCsv("csv/2018-05-04-Gummers-How-SOTA.csv")
	checkContactNumberIsCorrect("activator", len(contacts.ActivationContacts), 0)
	checkContactNumberIsCorrect("chaser", len(contacts.ChaserContacts), 3)

	contacts, _ = loadAndParseCsv("csv/2018-05-03-Lickbarrow-Road.csv")
	checkContactNumberIsCorrect("activator", len(contacts.ActivationContacts), 0)
	checkContactNumberIsCorrect("chaser", len(contacts.ChaserContacts), 3)

	contacts, _ = csv.ParseCsv(testUpload3, "M0NOM")
	checkContactNumberIsCorrect("activator", len(contacts.ActivationContacts), 0)
	checkContactNumberIsCorrect("chaser", len(contacts.ChaserContacts), 0)

	contacts, _ = csv.ParseCsv(testUpload1, "M0NOM")
	checkContactNumberIsCorrect("activator", len(contacts.ActivationContacts), 9)
	checkContactNumberIsCorrect("chaser", len(contacts.ChaserContacts), 0)
	checkS2SNumberIsCorrect(contacts.ActivationContacts, 3)

	contacts, _ = csv.ParseCsv(testUpload2, "M0NOM")
	checkContactNumberIsCorrect("activator", len(contacts.ActivationContacts), 0)
	checkContactNumberIsCorrect("chaser", len(contacts.ChaserContacts), 2)

}

func loadAndParseCsv(filename string) (csv.Contacts, error) {
	var contacts csv.Contacts
	fileContentInBytes, err := ioutil.ReadFile(filename)
	if err == nil {
		contacts, _ = csv.ParseCsv(string(fileContentInBytes), "M0NOM")
	}
	return contacts, err
}

func checkContactNumberIsCorrect(contactType string, actual int, expected int) {
	if actual != expected {
		fmt.Printf("Number of %s contacts is %d not %d as expected\n", contactType, actual, expected)
	}
}

func checkS2SNumberIsCorrect(activationContacts []csv.ActivationContact, expected int) {
	actual := 0
	for _, activationContact := range activationContacts {
		if activationContact.S2S {
			actual++
		}
	}
	if actual != expected {
		fmt.Printf("Number of S2S contacts is %d not %d as expected\n", actual, expected)
	}
}

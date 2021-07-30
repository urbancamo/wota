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

var filenames = [...]string{
	"2018-05-03-Lickbarrow-Road.csv",
	"2018-05-04-Gummers-How-SOTA.csv",
	"2018-05-05-High-Street-SOTA.csv",
	"2018-05-07-Gummers-How-SOTA.csv",
	"2018-05-19-Fairfield-SOTA.csv",
	"2018-05-19-St-Sunday-Crag-SOTA.csv",
	"2018-05-28-Whitbarrow-Scar-SOTA.csv",
	"2018-05-29-Gummers-How-MBARS-SOTA.csv",
	"2018-06-02-Gummers-How-SOTA.csv",
	"2018-06-03-Holme-Fell-SOTA.csv",
	"2018-06-07-Old-Man-Of-Coniston-SOTA2.csv",
	"2018-06-07-Old-Man-Of-Coniston-SOTA.csv",
	"2018-06-23-Red-Screes-SOTA.csv",
	"2018-06-24-Helvellyn-SOTA.csv",
	"2018-07-05-Skiddaw-SOTA.csv",
	"2018-07-26-Pike-of-Blisco-SOTA.csv",
	"2018-07-27-SOTA-Chase-Home-QTH.csv",
	"2018-08-04-Gummers-How-SOTA.csv",
	"2018-08-11-Alto-de-los-Jarales-SOTA-Chase.csv",
	"2018-08-11-Alto-de-los-Jarales-SOTA.csv",
	"2018-08-18-Alto-de-Calar-SOTA.csv",
	"2018-08-22-Veleta-SOTA.csv",
	"2018-09-15-High-Raise-SOTA.csv",
	"2018-09-18-SOTA-Chase-Home-QTH.csv",
	"2018-10-10-SOTA-Chase-Queen-Adelaides-Hill.csv",
	"2018-10-14-Red-Screes-SOTA.csv",
	"2018-10-18-SOTA-Chase-Home-QTH.csv",
	"2018-10-25-SOTA-Chase-Queen-Adelaide-Hill.csv",
	"2018-10-27-Gummers-How-SOTA.csv",
	"2018-10-28-Place-Fell-SOTA.csv",
	"2018-11-03-Gummers-How-SOTA.csv",
	"2018-11-08-Queen-Adelaides-Hill.csv",
	"2018-11-12-Queen-Adelaides-Hill.csv",
	"2018-11-22-Queen-Adelaides-Hill.csv",
	"2018-11-23-Lickbarrow-Road.csv",
	"2018-11-24-Lingmoor-Fell-SOTA.csv",
	"2018-11-25-SOTA-Chase-Penrith.csv",
	"2018-12-09-Seat-Sandal-SOTA.csv",
	"2018-12-16-Gummers-How-SOTA.csv",
	"2018-12-24-Great-Gable.csv",
	"2018-12-29-Gummers-How-SOTA.csv",
	"2019-01-01-Gummers-How-SOTA.csv",
	"2019-01-03-Home-QTH.csv",
	"2019-01-06-Top-of-Selside-SOTA.csv",
	"20190120-Lickbarrow-Road.csv",
	"2019-02-02-Gummers-How-SOTA.csv",
	"2019-02-10-Little-Mell-Fell-SOTA.csv",
	"2019-02-17-Snezhanka-SOTA.csv",
	"2019-02-24-Mechi-chel-SOTA.csv",
	"2019-03-23-Red-Screes-SOTA.csv",
	"2019-03-26-Lickbarrow-SOTA-Chase.csv",
	"2019-04-06-Gummers-How-SOTA.csv",
	"2019-04-17-Brant-Fell-WOTA-Chase.csv",
	"2019-04-18-Red-Screes-SOTA.csv",
	"2019-04-20-Gummers-How-SOTA.csv",
	"2019-04-29-Queen-Adelaides-Hill.csv",
	"2019-04-29-SOTA-Chase-Home-QTH.csv",
	"2019-05-02-Queen-Adelaides-Hill.csv",
	"2019-05-16-Brant-Fell.csv",
	"2019-05-18-02-Helvellyn-SOTA.csv",
	"2019-05-18-03-Walk-between-Helvellyn-and-Seat-Sandal.csv",
	"2019-05-18-04-Seat-Sandal-SOTA.csv",
	"2019-05-19-SOTA-Chase-Home-QTH.csv",
	"2019-05-24-Lickbarrow-Road-SOTA-Chase.csv",
	"M6VMS_log_20181019.csv",
	"SOTA-DB-Extract-MM0NOM_P.csv",
	"WOTA-to-WOTA-in-comment-test.csv",
	"WOTA-to-WOTA-test.csv",
}

func main() {
	var contacts csv.Contacts
	var err error

	contacts, _ = loadAndParseCsv("csv/2021-03-29T153645-G8CPZ.csv")
	checkContactNumberIsCorrect("activator", len(contacts.ActivationContacts), 4)
	checkContactNumberIsCorrect("chaser", len(contacts.ChaserContacts), 0)

	contacts, _ = loadAndParseCsvWithUser("csv/G4YSS_activationID_358852.csv", "GX0OOO")
	contacts, _ = loadAndParseCsvWithUser("csv/G4YSS_activationID_358857.csv", "GX0OOO")

	contacts, _ = loadAndParseCsv("csv/2018-06-07-Old-Man-Of-Coniston-SOTA2.csv")
	checkContactNumberIsCorrect("activator", len(contacts.ActivationContacts), 0)
	checkContactNumberIsCorrect("chaser", len(contacts.ChaserContacts), 0)

	contacts, _ = loadAndParseCsv("csv/WOTA-to-WOTA-in-comment-test.csv")
	checkContactNumberIsCorrect("activator", len(contacts.ActivationContacts), 1)
	checkContactNumberIsCorrect("chaser", len(contacts.ChaserContacts), 1)

	contacts, _ = loadAndParseCsv("csv/WOTA-to-WOTA-test.csv")
	checkContactNumberIsCorrect("activator", len(contacts.ActivationContacts), 1)
	checkContactNumberIsCorrect("chaser", len(contacts.ChaserContacts), 1)

	contacts, _ = loadAndParseCsv("csv/2018-10-28-Place-Fell-SOTA.csv")
	checkContactNumberIsCorrect("activator", len(contacts.ActivationContacts), 64)
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

	// broken file resilience
	contacts, err = loadAndParseCsv("testSotaCsvToContacts.go")
	if err != nil {
		fmt.Printf("%s", err.Error())
	}

	// Run through every file to make sure there are no major errors
	for _, filename := range filenames {
		contacts, err = loadAndParseCsv("csv/" + filename)
		fmt.Printf("Processed %s resulting in %d activator and %d chaser contacts\n", filename, len(contacts.ActivationContacts), len(contacts.ChaserContacts))
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}

func loadAndParseCsv(filename string) (csv.Contacts, error) {
	var contacts csv.Contacts
	fileContentInBytes, err := ioutil.ReadFile(filename)
	if err == nil {
		contacts, _ = csv.ParseCsv(string(fileContentInBytes), "M0NOM")
	}
	return contacts, err
}

func loadAndParseCsvWithUser(filename string, user string) (csv.Contacts, error) {
	var contacts csv.Contacts
	fileContentInBytes, err := ioutil.ReadFile(filename)
	if err == nil {
		contacts, _ = csv.ParseCsv(string(fileContentInBytes), user)
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

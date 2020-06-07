package domain

// {"Index":"1","Date":"2019-05-18 09:12:00","Callsign Used":"M0NOM/P","Summit":"LDW-129","Summit Name":"Illgill Head","Station Worked":"M0OAT"}
type ChaseContact struct {
	Index         string `json:"Index"`
	Date          string `json:"Date"`
	CallsignUsed  string `json:"Callsign Used"`
	Summit        string `json:"Summit"`
	SummitName    string `json:"Summit Name"`
	StationWorked string `json:"Station Worked"`
}


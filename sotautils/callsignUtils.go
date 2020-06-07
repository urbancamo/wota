package sotautils

import "strings"

// this method extracts the amateur radio operator from
// a callsign. So for example /P and /M are stripped off
// and a non-UK operator callsign for example OK/SQ9MDF/P
// would have the prefix stripped
func GetOperatorFromCallsign(callsign string) string {
	// determine how many slashes in the callsign
	callsignParts := strings.Split(callsign, "/")
	var slashCount = len(callsignParts) - 1
	if slashCount == 1 {
		return callsignParts[0]
	} else if slashCount > 1 {
		return callsignParts[1]
	}
	return callsign
}

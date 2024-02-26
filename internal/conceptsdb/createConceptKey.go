package conceptsdb

// CreateConceptKey sorts two strings and joins them with a '+' character.
// This is used to create a "Concept Key", which is used as the key in Redis.
//
// Example:
//   - "Fart" + "Dangerous" -> "Dangerous+Fart"
func CreateConceptKey(str1, str2 string) string {
	if str1 < str2 {
		return str1 + "+" + str2
	}
	return str2 + "+" + str1
}

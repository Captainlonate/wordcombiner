package conceptsdb

// CreateConceptKey creates a key that will be used to store records in Redis.
// All the records in our Redis database will use keys from this function.
// CreateConceptKey sorts two strings and joins them with a '+' character.
// The redis database stores "Concepts", so this function creates "Concept Keys".
//
// Example:
//   - "Fart" + "Dangerous" -> "Dangerous+Fart"
func CreateConceptKey(str1, str2 string) string {
	if str1 < str2 {
		return str1 + "+" + str2
	}
	return str2 + "+" + str1
}

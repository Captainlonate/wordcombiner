package routes

// Any struct that will ever be placed in the `Data` field of the ApiResponse
// must begin with DTO_
type DTO_NewConcept struct {
	Concept     string `json:"concept"`
	IsFirstTime bool   `json:"is_first"`
}

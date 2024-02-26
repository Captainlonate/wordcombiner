package routes

// Data Transfer Object (DTO) that represents the responds from GET /combine
type DTO_NewConcept struct {
	Concept     string `json:"concept"`
	IsFirstTime bool   `json:"is_first"`
}

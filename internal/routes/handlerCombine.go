package routes

import (
	"captainlonate/wordcombiner/internal/conceptsdb"
	oai "captainlonate/wordcombiner/internal/openai"
	"fmt"
	"net/http"
	"time"

	ce "captainlonate/wordcombiner/internal/customError"
)

// Route handler for GET /combine
func combineRouteHandler(w http.ResponseWriter, r *http.Request) {
	// Parse query params
	conceptOne := r.URL.Query().Get("one")
	conceptTwo := r.URL.Query().Get("two")

	conceptKey := conceptsdb.CreateConceptKey(conceptOne, conceptTwo)

	// Check if there is a match in redis already
	newConceptFromRedis, err := conceptsdb.FetchFromRedis(conceptKey)
	if newConceptFromRedis != "" && err == nil {
		processRedisHit(w, conceptKey, conceptOne, conceptTwo, newConceptFromRedis)
		return
	}

	// Otherwise fetch the combination from OpenAI
	chatCompletion, err := oai.GetCombinationFromOpenAI(conceptOne, conceptTwo)
	if err != nil || len(chatCompletion.Choices) == 0 {
		sendJSON(w, apiResponseFailure(ce.OpenAIAPIErrorCode, "Error getting combination from OpenAI"))
		return
	}
	processOpenAIHit(w, conceptKey, conceptOne, conceptTwo, chatCompletion)
}

func processRedisHit(w http.ResponseWriter, conceptKey, conceptOne, conceptTwo, newConcept string) {
	sendJSON(w, apiResponseSuccess(DTO_NewConcept{
		Concept:     newConcept,
		IsFirstTime: false,
	}))

	// Debug logging to terminal
	fmt.Println("===============================================")
	fmt.Println("====== 👨‍🌾🐑🍄 New Request (Redis) 🧺🌾🐐 ======")
	fmt.Println("===============================================")
	fmt.Printf("'%s' + '%s' = '%s'\n", conceptOne, conceptTwo, newConcept)
	fmt.Printf("\tRedis Key: '%s'\n", conceptKey)
	fmt.Println("")
	fmt.Println("")
}

func processOpenAIHit(w http.ResponseWriter, conceptKey string, conceptOne string, conceptTwo string, chatCompletion oai.ChatCompletion) {
	newConcept := chatCompletion.Choices[0].Message.Content

	// Save the combination in Redis for next time
	conceptsdb.InsertToRedis(conceptKey, newConcept)

	// Respond to the user
	sendJSON(w, apiResponseSuccess(DTO_NewConcept{
		Concept:     newConcept,
		IsFirstTime: true,
	}))

	// Debug logging to terminal
	createdTime := time.Unix(chatCompletion.Created, 0)
	fmt.Println("================================================")
	fmt.Println("====== 👨‍🌾🐑🍄 New Request (OpenAI) 🧺🌾🐐 ======")
	fmt.Println("================================================")
	fmt.Printf("'%s' + '%s' = '%s'\n", conceptOne, conceptTwo, newConcept)
	fmt.Printf("\tID: %s\n", chatCompletion.ID)
	fmt.Printf("\tCreated: %s\n", createdTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("\tModel: %s\n", chatCompletion.Model)
	fmt.Printf("\tTotal Tokens Used: %+v\n", chatCompletion.Usage.TotalTokens)
	fmt.Println("")
	fmt.Println("")
}

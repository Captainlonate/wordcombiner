package routes

import (
	"captainlonate/wordcombiner/internal/conceptsdb"
	oai "captainlonate/wordcombiner/internal/openai"
	"fmt"
	"net/http"
	"time"

	ce "captainlonate/wordcombiner/internal/customerror"
)

// This is the route handler, which handles all requests to "GET /combine"
func combineRouteHandler(w http.ResponseWriter, r *http.Request) {
	// Parse query params (upstream middleware will ensure they are set)
	conceptOne := r.URL.Query().Get("one")
	conceptTwo := r.URL.Query().Get("two")

	// The key for the record in Redis (for both fetching and storing)
	conceptKey := conceptsdb.CreateConceptKey(conceptOne, conceptTwo)

	// Check if there is a match in redis already
	newConceptFromRedis, err := conceptsdb.FetchFromRedis(conceptKey)
	if newConceptFromRedis != "" && err == nil {
		sendJSON(w, apiResponseSuccess(DTO_NewConcept{
			Concept:     newConceptFromRedis,
			IsFirstTime: false,
		}))
		return
	}

	// If not in Redis, get the combination from OpenAI
	chatCompletion, err := oai.GetCombinationFromOpenAI(conceptOne, conceptTwo)
	if err != nil || len(chatCompletion.Choices) == 0 {
		sendJSON(w, apiResponseFailure(ce.OpenAIAPIErrorCode, "Error getting combination from OpenAI"))
		return
	}
	processOpenAIHit(w, conceptKey, conceptOne, conceptTwo, chatCompletion)
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
	fmt.Println("====== 👨‍🌾🐑🍄 New Request To OpenAI 🧺🌾🐐 ======")
	fmt.Println("================================================")
	fmt.Printf("'%s' + '%s' = '%s'\n", conceptOne, conceptTwo, newConcept)
	fmt.Printf("\tID: %s\n", chatCompletion.ID)
	fmt.Printf("\tCreated: %s\n", createdTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("\tModel: %s\n", chatCompletion.Model)
	fmt.Printf("\tTotal Tokens Used: %+v\n", chatCompletion.Usage.TotalTokens)
	fmt.Println("")
	fmt.Println("")
}

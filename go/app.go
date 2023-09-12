package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"

	"github.com/boltdb/bolt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

const MaxStringLength = 1000

var (
	seenStringsMutex sync.RWMutex
	seenStrings      = make(map[string]int)
	mostPopularInput string
	longestInput     string

	db         *bolt.DB
	bucketName = []byte("stats")
	dbPath     = "stats.db" // Path to the BoltDB file
)

type StringData struct {
	Input         string `param:"input" query:"input" form:"input" json:"input" xml:"input" validate:"required"`
	Length        int    `json:"length"`
	MostFrequent  string `json:"most_frequent,omitempty"`
	FrequentCount int    `json:"frequent_count,omitempty"`
}

type StatsData struct {
	Inputs          map[string]int `json:"inputs"`
	MostPopular     string         `json:"most_popular,omitempty"`
	LongestInput    string         `json:"longest_input_received,omitempty"`
	LongestInputLen int            `json:"longest_input_len,omitempty"`
}

func init() {
	// Initialize BoltDB and create a bucket for stats
	var err error
	db, err = bolt.Open(dbPath, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Create a bucket for stats if it doesn't exist
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucketName)
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
}

func validateResponse(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := next(c); err != nil {
			return err
		}
		responseData, ok := c.Get("response_data").(StringData)
		if !ok {
			return echo.NewHTTPError(http.StatusInternalServerError, "Invalid response data")
		}
		if err := validateStringData(responseData); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Invalid response data")
		}
		return nil
	}
}

func validateStringData(data StringData) error {
	validate := validator.New()
	return validate.Struct(data)
}

func maxLengthMiddleware(maxLength int) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			input := c.QueryParam("input") // Adjust to match your query parameter name
			if len(input) > maxLength {
				return echo.NewHTTPError(http.StatusBadRequest, "Input string is too long")
			}
			return next(c)
		}
	}
}

func storeStatistics() {
	// Store statistics in BoltDB
	db.Update(func(tx *bolt.Tx) error {
		// Serialize your statistics data (e.g., seen_strings) to JSON
		statisticsJSON, err := json.Marshal(seenStrings)
		if err != nil {
			return err
		}

		// Store the JSON data in the bucket
		if err := tx.Bucket(bucketName).Put([]byte("statistics"), statisticsJSON); err != nil {
			return err
		}

		return nil
	})
}

func storeStatisticsMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Execute the request handler
		err := next(c)

		// After executing the handler, store statistics
		storeStatistics()

		return err
	}
}

func validateRequest(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		data := new(StringData)
		if err := c.Bind(data); err != nil {
			return err
		}
		validate := validator.New()
		if err := validate.Struct(data); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return next(c)
	}
}

func loadStatistics() {
	// Load statistics from BoltDB
	db.View(func(tx *bolt.Tx) error {
		// Retrieve the stored JSON data from the bucket
		statisticsJSON := tx.Bucket(bucketName).Get([]byte("statistics"))

		// Deserialize the JSON data into the seenStrings map
		if err := json.Unmarshal(statisticsJSON, &seenStrings); err != nil {
			return err
		}

		return nil
	})
}

func remember(input string) {
	seenStringsMutex.Lock()
	defer seenStringsMutex.Unlock()

	if seenStrings[input] == 0 {
		seenStrings[input] = 1
	} else {
		seenStrings[input] += 1
	}
	if seenStrings[input] > seenStrings[mostPopularInput] {
		mostPopularInput = input
	}
	if utf8.RuneCountInString(input) > utf8.RuneCountInString(longestInput) {
		longestInput = input
	}
}

func calculateMostFrequent(input string) (string, int) {
	charCount := make(map[rune]int)
	maxChar := ""
	maxCount := 0

	for _, char := range input {
		if unicode.IsSpace(char) || unicode.IsPunct(char) {
			continue // Ignore white space and punctuation
		}

		charCount[char]++
		if charCount[char] > maxCount {
			maxChar = string(char)
			maxCount = charCount[char]
		}
	}

	return maxChar, maxCount
}

func stringinate(c echo.Context) (err error) {
	requestData := new(StringData)
	if err = c.Bind(requestData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	remember(requestData.Input)

	mostFrequentChar, frequentCount := calculateMostFrequent(requestData.Input)
	requestData.Length = len(requestData.Input)
	requestData.MostFrequent = mostFrequentChar
	requestData.FrequentCount = frequentCount
	storeStatistics()

	return c.JSON(http.StatusOK, requestData)
}

func transformText(c echo.Context) (err error) {
	StringData := new(StringData)
	requestData := struct {
		Text           string `json:"text"`
		Transformation string `json:"transformation"`
	}{}

	if err = c.Bind(&requestData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	remember(requestData.Text)

	mostFrequentChar, frequentCount := calculateMostFrequent(requestData.Text)
	StringData.Length = len(requestData.Text)
	StringData.MostFrequent = mostFrequentChar
	StringData.FrequentCount = frequentCount

	storeStatistics()

	var transformedText string

	switch requestData.Transformation {
	case "uppercase":
		transformedText = strings.ToUpper(requestData.Text)
	case "lowercase":
		transformedText = strings.ToLower(requestData.Text)
	case "titlecase":
		transformedText = strings.Title(requestData.Text)
	default:
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid transformation type")
	}

	responseData := struct {
		OriginalText    string `json:"original_text"`
		TransformedText string `json:"transformed_text"`
	}{
		OriginalText:    requestData.Text,
		TransformedText: transformedText,
	}

	return c.JSON(http.StatusOK, responseData)
}

func resetStatistics(c echo.Context) error {
	seenStringsMutex.Lock()
	defer seenStringsMutex.Unlock()

	// Clear in-memory statistics
	seenStrings = make(map[string]int)
	mostPopularInput = ""
	longestInput = ""

	// Clear statistics in BoltDB
	db.Update(func(tx *bolt.Tx) error {
		// Remove the statistics bucket (if it exists) to reset statistics
		tx.DeleteBucket(bucketName)

		// Recreate the statistics bucket for the next test
		_, err := tx.CreateBucketIfNotExists(bucketName)
		return err
	})

	// Return a success response
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Statistics reset successfully",
	})
}

func stats(c echo.Context) (err error) {
	seenStringsMutex.RLock()
	defer seenStringsMutex.RUnlock()

	mostPopular := ""
	mostPopularCount := 0

	for input, count := range seenStrings {
		if count > mostPopularCount {
			mostPopular = input
			mostPopularCount = count
		}
	}

	return c.JSON(http.StatusOK, StatsData{
		Inputs:          seenStrings,
		MostPopular:     mostPopular,
		LongestInput:    longestInput,
		LongestInputLen: utf8.RuneCountInString(longestInput),
	})
}

func main() {
	loadStatistics()
	defer db.Close()
	e := echo.New()

	// Use the storeStatisticsMiddleware before all routes
	e.Use(storeStatisticsMiddleware)

	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, `
			<pre>
			Welcome to the Stringinator 3000 for all of your string manipulation needs.
			GET / - You're already here!
			POST /stringinate - Get all of the info you've ever wanted about a string. Takes JSON of the following form: {"input":"your-string-goes-here"}
			GET /stats - Get statistics about all strings the server has seen, including the longest and most popular strings.
			</pre>
		`)
	})

	e.POST("/stringinate", stringinate, validateRequest, validateResponse)
	e.GET("/stringinate", stringinate, validateRequest, maxLengthMiddleware(MaxStringLength), validateResponse)
	e.POST("/transform", transformText, validateRequest, validateResponse)
	e.GET("/transform", transformText, validateRequest, maxLengthMiddleware(MaxStringLength), validateResponse)
	e.GET("/reset-stats", resetStatistics)
	e.GET("/stats", stats)
	e.Logger.Fatal(e.Start(":1323"))
}

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"log-parser/match"
	"log-parser/parser"
	"time"
)

const (
	dateLayout = "02/01/2006 15:04"
)

func main() {
	now := time.Now()

	matches, err := parser.ParseLog("qgames.log")
	if err != nil {
		log.Fatal(err)
	}

	n := len(matches)

	games := make([]map[string]*match.Match, n)
	matchSummary := make([]map[string]match.Summary, n)
	for i := 0; i < n; i++ {
		key := fmt.Sprintf("game_%d", i+1)

		games[i] = map[string]*match.Match{
			key: matches[i],
		}

		matchSummary[i] = map[string]match.Summary{
			key: {
				KillsByMeans: matches[i].KillsByMeans,
			},
		}
	}

	matchesOutput, err := json.MarshalIndent(games, "", "    ")
	if err != nil {
		log.Fatalf("marshalling json output: %v", err.Error())
	}

	summaryOutput, err := json.MarshalIndent(matchSummary, "", "    ")
	if err != nil {
		log.Fatalf("marshalling json output: %v", err.Error())
	}

	reportTime := time.Now().Format(dateLayout)

	fmt.Printf("Matches Report - %v\n", reportTime)
	fmt.Println(string(matchesOutput))

	fmt.Printf("Deaths by Death cause - %v\n", reportTime)
	fmt.Println(string(summaryOutput))

	fmt.Printf("reports generated in %d ms", time.Since(now).Milliseconds())
}

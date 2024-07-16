package parser

import (
	"bufio"
	"errors"
	"log"
	"log-parser/match"
	"os"
	"sync"
)

func ParseLog(filepath string) ([]*match.Match, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, errors.New("reading the log file")
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			log.Fatalf("closing the file: %s", err.Error())
		}
	}(file)

	digester := LoadLogsDigester()
	matches := make([]*match.Match, 0)

	resultStream := make(chan *match.Match)
	gatheredLinesStream := make(chan []string)
	lines := make([]string, 0)

	go func() {
		defer close(gatheredLinesStream)
		sc := bufio.NewScanner(file)
		for sc.Scan() {
			logLine, matchLastLine := GatherLines(sc.Text())
			if logLine != "" {
				lines = append(lines, logLine)
			}

			if matchLastLine {
				gatheredLinesStream <- lines
				lines = make([]string, 0)
			}
		}

		if len(lines) > 0 {
			gatheredLinesStream <- lines
		}
	}()

	var wg sync.WaitGroup
	go func() {
		for gatheredLines := range gatheredLinesStream {
			wg.Add(1)
			go func(gatheredLines []string) {
				defer wg.Done()

				gameMatch := match.NewMatch()
				for _, line := range gatheredLines {
					err = digester.Handle(line, gameMatch)
					if err != nil {
						log.Fatalf("digesting log file: %v", err.Error())
					}

					if gameMatch.Done {
						resultStream <- gameMatch
					}
				}
			}(gatheredLines)
		}

		wg.Wait()
		close(resultStream)
	}()

	for result := range resultStream {
		matches = append(matches, result)
	}

	return matches, nil
}

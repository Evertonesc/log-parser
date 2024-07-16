package parser

import (
	"bufio"
	"errors"
	"log"
	"log-parser/match"
	"os"
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

	go func() {
		defer close(resultStream)
		for gatheredLine := range gatheredLinesStream {
			gameMatch := match.NewMatch()
			for _, line := range gatheredLine {
				err = digester.Handle(line, gameMatch)
				if err != nil {
					log.Fatalf("digesting log file: %v", err.Error())
				}

				if gameMatch.Done {
					resultStream <- gameMatch
				}
			}
		}
	}()

	for result := range resultStream {
		matches = append(matches, result)
	}

	return matches, nil
}

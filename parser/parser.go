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
			log.Fatal("closing the file")
		}
	}(file)

	digester := LoadLogsDigester()
	gameMatch := match.NewMatch()

	matches := make([]*match.Match, 0)

	sc := bufio.NewScanner(file)
	for sc.Scan() {
		logLine := sc.Text()

		err = digester.Handle(logLine, gameMatch)
		if err != nil {
			log.Fatal("digesting log file")
		}

		if gameMatch.Done {
			matches = append(matches, gameMatch)
			gameMatch = match.NewMatch()
		}
	}

	return matches, nil
}

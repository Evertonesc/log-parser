package parser

import (
	"bufio"
	"errors"
	"log"
	"log-parser/match"
	"os"
)

func parseLog(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return errors.New("reading the log file")
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			log.Fatal("closing the file")
		}
	}(file)

	digester := LoadLogsDigester()
	match := &Match{}

	sc := bufio.NewScanner(file)
	for sc.Scan() {
		logLine := sc.Text()

		err = digester.Handle(logLine, match)
		if err != nil {
			log.Fatal("digesting log file")
		}
	}

	return nil
}

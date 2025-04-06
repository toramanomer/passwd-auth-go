package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func init() {
	filename := ".env"
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Failed to open %s file: %v\n", filename, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		name, value, found := strings.Cut(line, "=")
		if !found {
			log.Fatalf("Invalid line %d in file %s: %s\n", lineNumber, filename, line)
		}

		if err := os.Setenv(name, value); err != nil {
			log.Fatalf("Failed to set %s environment variable: %v\n", name, err)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error scanning %s file: %v\n", filename, err)
	}
}

func main() {

}

package main

import (
	"encoding/csv"
	"log"
	"os"
)

// Recieves a file name as a argument, reads it and returns the output
func readCSV(filename string) [][]string {
	fd, err := os.Open(filename)
	if err != nil {
		log.Printf("Couldn't load csv file \"%v\"\n", filename)
		os.Exit(1)
	}
	csvReader := csv.NewReader(fd)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Printf("Couldn't load csv file \"%v\"\n", filename)
		os.Exit(1)
	}
	return data
}

// Takes a filename as an argument and returns an array of question parsed from the csv file
func loadCSVIntoQuestions(filename string) []question {
	var qs []question
	data := readCSV(filename)
	if data == nil {
		log.Fatalln("Error Parsing Csv File")
		os.Exit(1)
		return nil
	}
	for _, v := range data {
		q := question{
			question: v[0],
			answer:   v[1],
		}
		qs = append(qs, q)
	}
	return qs
}

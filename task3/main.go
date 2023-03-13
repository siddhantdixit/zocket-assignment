package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

type Person struct {
	Name    string
	Age     int
	Country string
}

func main() {
	file, err := os.Open("people.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	var people []Person
	reader.Read()
	for {
		record, err := reader.Read()
		if err != nil {
			break
		}
		person := Person{Name: record[0], Age: parseInt(record[1]), Country: record[2]}
		people = append(people, person)
	}

	fmt.Printf("|%-20s|%-10s|%-20s|\n", "Name", "Age", "Country")
	fmt.Println("----------------------------------------------")
	for _, person := range people {
		fmt.Printf("|%-20s|%-10d|%-20s|\n", person.Name, person.Age, person.Country)
	}
}

func parseInt(str string) int {
	var i int
	if _, err := fmt.Sscanf(str, "%d", &i); err != nil {
		return 0
	}
	return i
}

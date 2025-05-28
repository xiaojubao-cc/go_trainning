package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
)

type state struct {
	id               int
	name             string
	abbreviation     string
	censusRegionName string
}

func parseState(columns map[string]int, record []string) (*state, error) {
	id, err := strconv.Atoi(record[columns["id"]])
	name := record[columns["name"]]
	abbreviation := record[columns["abbreviation"]]
	censusRegionName := record[columns["census_region_name"]]
	if err != nil {
		return nil, err
	}
	return &state{
		id:               id,
		name:             name,
		abbreviation:     abbreviation,
		censusRegionName: censusRegionName,
	}, nil

}

func main() {
	file, err := os.Open("D:\\golang projects\\go_training\\code-in-process\\05_review-exercises\\02_csv_state-info\\state_table.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	columns := make(map[string]int)
	for rowcount := 0; ; rowcount++ {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalln("read error:" + err.Error())
		}
		if rowcount == 0 {
			for idx, column := range record {
				columns[column] = idx
			}
		} else {
			//封装数据
			states, err := parseState(columns, record)
			if err != nil {
				log.Fatalln("parse error:" + err.Error())
			}
			log.Println(states)
		}
	}
}

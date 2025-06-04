package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type StateTableOptions func(s *StateTable)

func WithId(id string) StateTableOptions {
	return func(s *StateTable) {
		s.id = id
	}
}
func WithName(name string) StateTableOptions {
	return func(s *StateTable) {
		s.name = name
	}
}
func WithAbbreviation(abbreviation string) StateTableOptions {
	return func(s *StateTable) {
		s.abbreviation = abbreviation
	}
}
func WithCountry(country string) StateTableOptions {
	return func(s *StateTable) {
		s.country = country
	}
}
func WithTypes(types string) StateTableOptions {
	return func(s *StateTable) {
		s.types = types
	}
}
func WithSort(sort string) StateTableOptions {
	return func(s *StateTable) {
		s.sort = sort
	}
}
func WithStatus(status string) StateTableOptions {
	return func(s *StateTable) {
		s.status = status
	}
}
func WithOccupied(occupied string) StateTableOptions {
	return func(s *StateTable) {
		s.occupied = occupied
	}
}
func WithNotes(notes string) StateTableOptions {
	return func(s *StateTable) {
		s.notes = notes
	}
}
func WithFipsState(fipsState string) StateTableOptions {
	return func(s *StateTable) {
		s.fipsState = fipsState
	}
}
func WithAssocPress(assocPress string) StateTableOptions {
	return func(s *StateTable) {
		s.assocPress = assocPress
	}
}
func WithStandardFederalRegion(standardFederalRegion string) StateTableOptions {
	return func(s *StateTable) {
		s.standardFederalRegion = standardFederalRegion
	}
}
func WithCensusRegion(censusRegion string) StateTableOptions {
	return func(s *StateTable) {
		s.censusRegion = censusRegion
	}
}
func WithCensusRegionName(censusRegionName string) StateTableOptions {
	return func(s *StateTable) {
		s.censusRegionName = censusRegionName
	}
}
func WithCensusDivision(censusDivision string) StateTableOptions {
	return func(s *StateTable) {
		s.censusDivision = censusDivision
	}
}
func WithCensusDivisionName(censusDivisionName string) StateTableOptions {
	return func(s *StateTable) {
		s.censusDivisionName = censusDivisionName
	}
}
func WithCircuitCourt(circuitCourt string) StateTableOptions {
	return func(s *StateTable) {
		s.circuitCourt = circuitCourt
	}
}

func NewStateTable(opt ...StateTableOptions) *StateTable {
	stateTable := &StateTable{}
	for _, o := range opt {
		o(stateTable)
	}
	return stateTable
}

type StateTable struct {
	id                    string
	name                  string
	abbreviation          string
	country               string
	types                 string
	sort                  string
	status                string
	occupied              string
	notes                 string
	fipsState             string
	assocPress            string
	standardFederalRegion string
	censusRegion          string
	censusRegionName      string
	censusDivision        string
	censusDivisionName    string
	circuitCourt          string
}

type StateTableSlice []*StateTable

func (s StateTableSlice) Len() int {
	return len(s)
}

//struct实现sort三个函数接口

func (s StateTableSlice) Less(i, j int) bool {
	idi, _ := strconv.Atoi(strings.ReplaceAll(s[i].id, "\"", ""))
	idj, _ := strconv.Atoi(strings.ReplaceAll(s[j].id, "\"", ""))
	return idi < idj
}
func (s StateTableSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func readFile(filePath string) map[string]int {
	content := make(map[string]int)
	row := 0
	openFile, _ := os.OpenFile(filePath, os.O_RDONLY, 0666)
	defer openFile.Close()
	scanner := bufio.NewScanner(openFile)
	for scanner.Scan() {
		content[scanner.Text()] = row
		row++
	}
	return content
}
func main() {
	//读取文件内容
	content := readFile("D:\\golang projects\\go_training\\01_basic_training\\15_os\\state_table.csv")
	tables := make([]*StateTable, 0)
	//映射到实体类
	for key, value := range content {
		if value > 0 {
			split := strings.Split(key, ",")
			stateTable := NewStateTable(
				WithId(split[0]),
				WithName(split[1]),
				WithAbbreviation(split[2]),
				WithCountry(split[3]),
				WithTypes(split[4]),
				WithSort(split[5]),
				WithStatus(split[6]),
				WithOccupied(split[7]),
				WithNotes(split[8]),
				WithFipsState(split[9]),
				WithAssocPress(split[10]),
				WithStandardFederalRegion(split[11]),
				WithCensusRegion(split[12]),
				WithCensusRegionName(split[13]),
				WithCensusDivision(split[14]),
				WithCensusDivisionName(split[15]),
				WithCircuitCourt(split[16]),
			)
			tables = append(tables, stateTable)
		}
	}
	/*这里是强转*/
	sort.Sort(StateTableSlice(tables))
	for _, table := range tables {
		fmt.Printf("%+v\n", table)
	}
}

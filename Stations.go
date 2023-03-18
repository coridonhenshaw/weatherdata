package main

import "database/sql"

type StationsStruct struct {
	SQLdb *sql.DB
}

type StationDataStruct struct {
	Identifier string
	Name       string
	Province   string
	Latitude   string
	Longitude  string
	Elevation  string
	TimeOffset int
	URL        string
}

package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

type StationsUIStruct struct {
	gsQuery string
	gsKML   string
}

func (o *StationsUIStruct) Get() error {

	Query := o.gsQuery
	KMLFile := o.gsKML

	var err error
	var Row *sql.Rows

	SQLdb := Cache.GetSQLConnection()

	if len(Query) == 0 {
		SQL := `SELECT Identifier, Name, Province, Latitude, Longitude FROM StationList ORDER BY Province, Name`
		Row, err = SQLdb.Query(SQL)
		if err != nil {
			log.Panic(err)
		}

	} else {
		SQL := `SELECT Identifier, Name, Province, Latitude, Longitude FROM StationList WHERE Name LIKE ? OR Identifier LIKE ? ORDER BY Province, Name`
		Row, err = SQLdb.Query(SQL, Query, Query)
		if err != nil {
			log.Panic(err)
		}
	}

	defer Row.Close()

	if len(KMLFile) == 0 {

		for Row.Next() {
			var Identifier, Name, Province string
			var Latitude, Longitude float64

			err = Row.Scan(&Identifier, &Name, &Province, &Latitude, &Longitude)
			if err != nil {
				log.Panic(err)
			}

			fmt.Printf("%-48s %2s (%6.03f, % 8.03f) - %s\n", Identifier, Province, Latitude, Longitude, Name)
		}
	} else {
		var Head = `<?xml version="1.0" encoding="UTF-8"?>` + "\n" + `<kml xmlns="http://www.opengis.net/kml/2.2">` + "\n<Document><name>MSC Datamart Weather Stations (via Weatherdata)</name>\n"
		var Placemark = `  <Placemark>
    <name>%s</name>
    <description>%s</description>
    <Point><coordinates>%.04f,%.04f,0</coordinates></Point>
  </Placemark>` + "\n"
		var Foot = `</Document></kml>`

		f, err := os.Create(KMLFile)
		if err != nil {
			log.Panic(err)
		}

		defer f.Close()

		fmt.Fprint(f, Head)

		for Row.Next() {
			var Identifier, Name, Province string
			var Latitude, Longitude float64

			err = Row.Scan(&Identifier, &Name, &Province, &Latitude, &Longitude)
			if err != nil {
				log.Panic(err)
			}

			fmt.Fprintf(f, Placemark, Identifier, Name, Longitude, Latitude)
		}

		fmt.Fprint(f, Foot)

	}

	return nil
}

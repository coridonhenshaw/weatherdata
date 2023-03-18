package main

import (
	"database/sql"
	"fmt"
	"log"
)

type AutotestStruct struct {
}

func (o *AutotestStruct) Run() error {
	var err error
	var SQL string

	SQLdb := Cache.GetSQLConnection()

	SQL = `select "prefix", "name" from providerlist`
	Row, err := SQLdb.Query(SQL)
	if err != nil {
		log.Panic(err)
	}

	var OUI = ObservationUIStruct{ShowIdentifiers: true}

	defer Row.Close()
	for Row.Next() {
		var Prefix, Name string

		err = Row.Scan(&Prefix, &Name)
		if err != nil {
			log.Panic(err)
		}

		if len(Prefix) == 0 {
			Prefix = "-"
		}

		SQL = `SELECT Identifier FROM Stationlist WHERE Identifier LIKE "` + Prefix + `%" ORDER BY RANDOM() LIMIT 1`

		var Identifier string
		err = SQLdb.QueryRow(SQL).Scan(&Identifier)

		if err == sql.ErrNoRows {
			fmt.Println("Prefix", Prefix, "empty")
			continue
		} else if err != nil {
			log.Panic(err)
		}

		OUI.Station = OUI.Station + Identifier + " "
	}

	err = OUI.Get()

	return err
}

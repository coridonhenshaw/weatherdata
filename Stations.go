package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type StationsStruct struct {
	SQLdb *sql.DB

	statement *sql.Stmt
	// ProviderStatement *sql.Stmt
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

func MakeStationsStruct() (o StationsStruct) {

	o.SQLdb = Cache.GetSQLConnection()

	RefreshCache, err := o.CacheIsStale()
	if err != nil {
		log.Panic(err)
	}

	if RefreshCache == false {
		if Verbose {
			fmt.Println("Used cached stations list.")
		}
	} else {
		err = o.Import()
		if err != nil {
			log.Panic(err)
		}
	}

	return
}

func (o *StationsStruct) CacheIsStale() (bool, error) {
	var err error

	if !RefreshStationCache {
		SQL := `CREATE TABLE IF NOT EXISTS KeyValueStore (
		"Key" TEXT PRIMARY KEY,
		"Value" TEXT
		) WITHOUT ROWID;`

		_, err := o.SQLdb.Exec(SQL)
		if err != nil {
			log.Panic(err)
		}

		Now := time.Now().Unix()

		SQL = `SELECT Value FROM KeyValueStore WHERE Key = ?`
		var LastUpdated int64
		err = o.SQLdb.QueryRow(SQL, "StationsLastUpdated").Scan(&LastUpdated)
		if err == nil {
			const Day = 24 * 60 * 60
			if (Now - LastUpdated) <= (Day) {
				return false, err
			}
		}
	}
	return true, err
}

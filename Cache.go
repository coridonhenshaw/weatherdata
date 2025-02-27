package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type CacheStruct struct {
	File           string
	sqliteDatabase *sql.DB
}

func (c *CacheStruct) Open() {

	var err error

	Path, err := os.UserCacheDir()
	c.File = filepath.Join(Path, "weatherdata-cache.sqlite3")

	c.sqliteDatabase, err = sql.Open("sqlite3", c.File)

	SQL := `PRAGMA SYNCHRONOUS = NORMAL`
	_, err = c.sqliteDatabase.Exec(SQL)
	if err != nil {
		log.Panic(err)
	}
	SQL = `PRAGMA journal_mode = WAL2`
	_, err = c.sqliteDatabase.Exec(SQL)
	if err != nil {
		log.Panic(err)
	}
	c.sqliteDatabase.SetMaxOpenConns(4)

	SQL = `CREATE TABLE IF NOT EXISTS Cache (
		"URL" TEXT PRIMARY KEY,
		"XML" TEXT,
		"Timestamp" INTEGER
	  ) WITHOUT ROWID;`

	statement, err := c.sqliteDatabase.Prepare(SQL)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()

}

func (c *CacheStruct) GetSQLConnection() *sql.DB {
	return c.sqliteDatabase
}

func (c *CacheStruct) Get(URL string, MaxAge time.Duration) (string, error) {
	var err error
	var Text string
	var Timestamp int

	SQL := "SELECT Timestamp, Xml FROM Cache WHERE Url = ?"

	err = c.sqliteDatabase.QueryRow(SQL, URL).Scan(&Timestamp, &Text)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}

	if err == sql.ErrNoRows {
		return "", nil
	}

	t := time.Unix(int64(Timestamp), 0)
	if t.Before(time.Now().UTC().Add(-MaxAge)) {
		if Verbose {
			fmt.Println("Expired", URL)
		}

		SQL := "DELETE FROM Cache WHERE Url = ?"
		_, err = c.sqliteDatabase.Exec(SQL, URL)
		if err != nil {
			log.Panic(err)
		}

		return "", nil
	}

	return Text, nil
}

func (c *CacheStruct) Put(URL string, Text string) error {
	SQL := "INSERT OR REPLACE INTO Cache (Xml, Url, Timestamp) VALUES (?, ?, ?)"

	Timestamp := time.Now().Unix()

	_, err := c.sqliteDatabase.Exec(SQL, Text, URL, Timestamp)
	if err != nil {
		return err
	}

	return nil
}

func (c *CacheStruct) Close() {

	SQL := "DELETE FROM Cache WHERE Timestamp < ?"
	Timestamp := time.Now().Add(-30 * 24 * time.Hour).Unix() // 30 days

	_, err := c.sqliteDatabase.Exec(SQL, Timestamp)
	if err != nil {
		log.Panic(err)
	}

	_, err = c.sqliteDatabase.Exec("VACUUM")
	if err != nil {
		log.Panic(err)
	}

	c.sqliteDatabase.Close()
}

package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
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

	SQL := `CREATE TABLE IF NOT EXISTS Cache (
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

func (c *CacheStruct) Get(URL string) (string, error) {
	var err error
	var Text string

	SQL := "SELECT Xml FROM Cache WHERE Url = ?"

	err = c.sqliteDatabase.QueryRow(SQL, URL).Scan(&Text)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}

	return Text, nil
}

func (c *CacheStruct) Put(URL string, Text string) error {
	SQL := "INSERT OR REPLACE INTO Cache (Xml, Url, Timestamp) VALUES (?, ?, ?)"

	Timestamp := time.Now().Unix()

	statement, err := c.sqliteDatabase.Prepare(SQL)
	if err != nil {
		return err
	}

	_, err = statement.Exec(Text, URL, Timestamp)
	if err != nil {
		return err
	}

	return nil
}

func (c *CacheStruct) Close() {

	SQL := "DELETE FROM Cache WHERE Timestamp < ?"
	Timestamp := time.Now().Add(-30 * 24 * time.Hour).Unix() // 30 days

	statement, err := c.sqliteDatabase.Prepare(SQL)
	if err == nil {
		_, err = statement.Exec(Timestamp)
	} else {
		fmt.Println("SQLite error when expiring cache", err)
		os.Exit(1)
	}

	_, err = c.sqliteDatabase.Exec("VACUUM")
	if err != nil {
		fmt.Println("SQLite error when vaccuming", err)
		os.Exit(1)
	}

	c.sqliteDatabase.Close()
}

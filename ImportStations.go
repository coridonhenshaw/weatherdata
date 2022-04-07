package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"strings"
)

type StationsStruct struct {
	SQLdb *sql.DB
}

func (o *StationsStruct) Import() error {

	SQLdb := Cache.GetSQLConnection()
	o.SQLdb = SQLdb

	SQL := `CREATE TEMPORARY TABLE IF NOT EXISTS StationList (
    "Identifier" TEXT PRIMARY KEY,
    "Name" TEXT,
    "Province" TEXT,
    "Latitude" FLOAT,
    "Longitude" FLOAT,
    "Elevation" FLOAT,
    "WMOID" INTEGER
    ) WITHOUT ROWID;`

	statement, err := SQLdb.Prepare(SQL)
	if err != nil {
		fmt.Println("0")
		log.Fatal(err.Error())
	}
	statement.Exec()

	SQL = "INSERT OR REPLACE INTO StationList (Identifier, Name, Province, Latitude, Longitude, Elevation, WMOID) VALUES (?, ?, ?, ?, ?, ?, ?)"

	statement, err = SQLdb.Prepare(SQL)
	if err != nil {
		fmt.Println("1")
		return err
	}

	URL := "https://dd.weather.gc.ca/observations/doc/swob-xml_partner_station_list.csv"

	StationString, err := HTTPSGet(URL)
	if err != nil {
		fmt.Println("HTTP Error", err, "acquiring", URL)
		return err
	}

	reader := csv.NewReader(strings.NewReader(StationString))
	_, _ = reader.Read()

	UnknownProviderMap := make(map[string]bool)

	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			return error
		}

		I := line[0]
		Name := line[2]
		Province := line[3]
		//		AUTOMAN := line[4]
		Latitude := line[5]
		Longitude := line[6]
		Elevation := line[7]
		WMOID := line[9]
		MSCID := line[10]
		DataProvider := line[13]

		var Site = "partners/"

		if DataProvider == "Department of Fisheries and Ocean Canada" && Elevation == "0" {
			//			fmt.Fprintln(os.Stderr, "Skip inactive DFO station", Name)
			continue
		} else if DataProvider == "Government of Canada: Fisheries and Oceans Canada; Canadian Coast Guard" {
			//Site += "dfo-ccg-lighthouse/" + strings.ReplaceAll(strings.ToLower(I), " ", "_") // Non-functional
			//			fmt.Fprintln(os.Stderr, "Skip DFO station", Name)
			continue
		} else if DataProvider == "Government of British Columbia: Ministry of Environment" {
			if strings.HasPrefix(MSCID, "BC_ENV-ASW_") {
				Site += "bc-env-snow/" + strings.ToLower(I)
			} else if strings.HasPrefix(MSCID, "BC_ENV-AQ_") {
				Site += "bc-env-aq/" + strings.ToLower(I)
			} else {
				fmt.Println("Unknown BC MoE subtype for MSC ID", MSCID)
			}
		} else if strings.HasPrefix(MSCID, "BC_WMB_") {
			Site += "bc-forestry/" + I
		} else if DataProvider == "Government of British Columbia: Ministry of Transportation and Infrastructure" {
			Site += "bc-tran/" + I
		} else if strings.HasPrefix(MSCID, "NL-DECCM-WRMD_") {
			Site += "nl-water/" + strings.ToLower(I)
		} else if DataProvider == "Government of Northwest Territories: Department of Environment and Natural Resources; Forest Management Division" {
			Site += "nt-forestry/" + strings.ToLower(I)
		} else if DataProvider == "Government of Northwest Territories: Department of Environment and Natural Resources; Water Resources Division" {
			Site += "nt-water/" + strings.ToLower(I)
		} else if DataProvider == "Government of Saskatchewan: Public Safety Agency - Wildfire Management Branch" {
			Site += "sk-forestry/" + strings.ToLower(I)
		} else if DataProvider == "Government of Yukon" || DataProvider == "Government of Yukon on behalf of The Yukon Avalanche Association" {
			Site += "yt-gov/" + strings.ReplaceAll(strings.ToLower(Name), " ", "_")
		} else if strings.HasPrefix(MSCID, "YT-DE-WRB_") {
			Site += "yt-water/" + strings.ToLower(I)
		} else {
			UnknownProviderMap[DataProvider] = true
		}

		Identifier := Site

		//fmt.Printf("%-9s %2s %-25s\n", Identifier, Province, Name)

		_, err = statement.Exec(Identifier, Name, Province, Latitude, Longitude, Elevation, WMOID)
		if err != nil {
			fmt.Println("2a")
			return err
		}

	}

	var Flag bool
	for k := range UnknownProviderMap {
		fmt.Println("Unknown data provider:", k)
		Flag = true
	}
	if Flag {
		fmt.Println()
	}

	var ProvinceMap = map[string]string{
		"Alberta":                   "AB",
		"British Columbia":          "BC",
		"Newfoundland and Labrador": "NL",
		"Saskatchewan":              "SK",
		"Ontario":                   "ON",
		"Northwest Territories":     "NT",
		"Yukon":                     "YT",
		"Prince Edward Island":      "PE",
		"Manitoba":                  "MB",
		"Nunavut":                   "NU",
		"Quebec":                    "QC",
		"New Brunswick":             "NB",
		"Nova Scotia":               "NS",
	}

	URL = "https://dd.weather.gc.ca/observations/doc/swob-xml_station_list.csv"

	StationString, err = HTTPSGet(URL)
	if err != nil {
		fmt.Println("HTTP Error", err, "acquiring", URL)
		return err
	}

	reader = csv.NewReader(strings.NewReader(StationString))
	_, _ = reader.Read()

	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			return error
		}

		Name := line[1]
		WMOID := line[2]
		Latitude := line[4]
		Longitude := line[5]
		Elevation := line[6]
		Province := ProvinceMap[line[10]]

		Identifier := line[0] + "-" + line[9]

		_, err = statement.Exec(Identifier, Name, Province, Latitude, Longitude, Elevation, WMOID)
		if err != nil {
			fmt.Println("2")
			return err
		}

		//			fmt.Printf("%-9s  %-25s  %s\n", Identifier, Province, Name)
	}
	return nil
}

func (o *StationsStruct) Validate(Name string) bool {

	var Count int

	SQL := `SELECT COUNT() FROM StationList WHERE Identifier = ?`
	err := o.SQLdb.QueryRow(SQL, Name).Scan(&Count)
	if err != nil {
		return false
	}
	if Count > 0 {
		return true
	}
	return false

}

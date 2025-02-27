package main

import (
	"bytes"
	_ "embed"
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"html/template"
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//go:embed stations.xml
var EmbeddedStationsXMLBytes []byte

type StationsXMLStruct struct {
	XMLName xml.Name               `xml:"WeatherData"`
	CSV     []StationsXMLCSVStruct `xml:"CSV"`
}

type StationsXMLCSVStruct struct {
	URL            string `xml:"URL,attr"`
	Root           string `xml:"Root,attr"`
	IdentifierCase string `xml:"IdentifierCase,attr"`

	Columns StationsXMLColumnsStruct `xml:"Columns"`

	Provider []StationsXMLProviderStruct `xml:"Provider"`
}

type StationsXMLColumnsStruct struct {
	Name      int `xml:"Name,attr"`
	Province  int `xml:"Province,attr"`
	Latitude  int `xml:"Latitude,attr"`
	Longitude int `xml:"Longitude,attr"`
	Elevation int `xml:"Elevation,attr"`
	Provider  int `xml:"Provider,attr"`
}

type StationsXMLProviderStruct struct {
	Prefix     string                  `xml:"Prefix,attr"`
	Identifier string                  `xml:"Identifier,attr"`
	URL        string                  `xml:"URL,attr"`
	Ignore     string                  `xml:"Ignore,attr"`
	TimeOffset int                     `xml:"TimeOffset,attr"`
	Detect     StationsXMLDetectStruct `xml:"Detect"`
}

type StationsXMLDetectStruct struct {
	Key    string `xml:"Key,attr"`
	Column int    `xml:"Column,attr"`
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
	"AB":                        "AB",
	"BC":                        "BC",
	"NL":                        "NL",
	"SK":                        "SK",
	"ON":                        "ON",
	"NT":                        "NT",
	"YT":                        "YT",
	"PE":                        "PE",
	"MB":                        "MB",
	"NU":                        "NU",
	"QC":                        "QC",
	"NB":                        "NB",
	"NS":                        "NS",
}

type LoadtimeTemplateStruct struct {
	Date          string
	FullTimestamp string
	DFOTimestamp  string

	Identifier string
	Column     []string
}

func (o *StationsStruct) Import() error {

	var err error
	var SQL string

	Transaction := Cache.GetSQLConnection()

	/* Force cache refresh if update fails */
	SQL = "INSERT OR REPLACE INTO KeyValueStore (Key, Value) VALUES (?, ?)"
	_, err = Transaction.Exec(SQL, "StationsLastUpdated", 0)
	if err != nil {
		log.Panic(err)
	}

	SQL = `DROP TABLE IF EXISTS ProviderList`
	_, err = Transaction.Exec(SQL)
	if err != nil {
		log.Panic(err)
	}

	// SQL = `CREATE TABLE IF NOT EXISTS ProviderList (
	// 	"Prefix" TEXT PRIMARY KEY,
	// 	"Quantity" INT DEFAULT 1,
	// 	"Name" TEXT
	// 	) WITHOUT ROWID;`

	// _, err = Transaction.Exec(SQL)
	// if err != nil {
	// 	log.Panic(err)
	// }

	SQL = `DROP TABLE IF EXISTS StationList`
	_, err = Transaction.Exec(SQL)
	if err != nil {
		log.Panic(err)
	}

	SQL = `CREATE TABLE IF NOT EXISTS StationList (
    "Identifier" TEXT PRIMARY KEY,
    "Name" TEXT,
    "Province" TEXT,
    "Latitude" FLOAT,
    "Longitude" FLOAT,
    "Elevation" FLOAT,
	"TimeOffset" INTEGER,
    "URL" TEXT
    ) WITHOUT ROWID;`

	_, err = Transaction.Exec(SQL)
	if err != nil {
		log.Panic(err)
	}

	SQL = "INSERT OR REPLACE INTO StationList (Identifier, Name, Province, Latitude, Longitude, Elevation, TimeOffset, URL) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	StoreStatement, err := Transaction.Prepare(SQL)
	if err != nil {
		log.Panic(err)
	}

	// SQL = "INSERT OR REPLACE INTO ProviderList (Prefix, Name) VALUES (?, ?) ON CONFLICT DO UPDATE SET Quantity=Quantity+1"
	// o.ProviderStatement, err = Transaction.Prepare(SQL)
	// if err != nil {
	// 	log.Panic(err)
	// }

	var Stations StationsXMLStruct
	err = xml.Unmarshal(EmbeddedStationsXMLBytes, &Stations)
	if err != nil {
		log.Panic(err)
	}

	UnknownProviderMap := make(map[string]int)

	for _, StationListEntry := range Stations.CSV {
		if Verbose {
			fmt.Println("Station list CSV:", StationListEntry.URL)
		}

		RawCSVData, err := HTTPSGet(StationListEntry.URL, time.Hour*12)
		if err != nil {
			fmt.Println("HTTP Error", err, "acquiring", StationListEntry.URL)
			return err
		}

		StoreFunc := func(FinalIdentifier string, Name string, Province string, Latitude string, Longitude string, Elevation string, TimeOffset int, URL string) error {
			_, err := StoreStatement.Exec(FinalIdentifier, Name, Province, Latitude, Longitude, Elevation, TimeOffset, URL)
			return err
		}

		err = o.ParseCSV(StationListEntry, strings.NewReader(RawCSVData), StoreFunc, UnknownProviderMap)
		if err != nil {
			log.Panic(err)
		}

	}

	var Flag bool
	for k, b := range UnknownProviderMap {
		fmt.Println("Unknown data provider:", b, k)
		Flag = true
	}
	if Flag {
		fmt.Println()
	}

	SQL = "INSERT OR REPLACE INTO KeyValueStore (Key, Value) VALUES (?, ?)"
	_, err = Transaction.Exec(SQL, "StationsLastUpdated", time.Now().Unix())
	if err != nil {
		log.Panic(err)
	}

	//	Transaction.Commit()

	return nil
}

func (o *StationsStruct) ParseCSV(v StationsXMLCSVStruct,
	RawCSVData io.Reader,
	Store func(FinalIdentifier string, Name string, Province string, Latitude string, Longitude string, Elevation string, TimeOffset int, URL string) error,
	UnknownProviderMap map[string]int) error {

	var err error

	CSVReader := csv.NewReader(RawCSVData)
	_, err = CSVReader.Read()
	if err != nil {
		log.Panic(err)
	}

	var BasePath = v.Root

	var LowerCaseIdentifiers = false
	if v.IdentifierCase == "lower" {
		LowerCaseIdentifiers = true
	}

	var LineCount int
	for {
		line, err := CSVReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		var p StationsXMLProviderStruct
		var Matched bool
		for _, Provider := range v.Provider {
			p = Provider
			if len(Provider.Detect.Key) == 0 && Provider.Detect.Column == 0 {
				Matched = true
				break
			}

			Match, err := regexp.MatchString(Provider.Detect.Key, line[Provider.Detect.Column])
			if err != nil {
				log.Panic(err)
			}

			if Match {
				Matched = true
				break
			}
		}

		if !Matched {
			_, o := UnknownProviderMap[line[v.Columns.Provider]]
			if !o {
				UnknownProviderMap[line[v.Columns.Provider]] = 0
			}
			UnknownProviderMap[line[v.Columns.Provider]]++

			continue
		}

		if len(p.Ignore) > 0 {
			continue
		}

		var Identifier string

		if len(p.Identifier) > 0 {
			CompositeIdentifier := strings.Split(p.Identifier, ",")
			for _, k := range CompositeIdentifier {
				i, err := strconv.Atoi(k)
				Identifier += line[i] + "-"
				if err != nil {
					log.Panic(err)
				}
			}
			Identifier = Identifier[:len(Identifier)-1]
		} else {
			Identifier = line[0]
		}

		if LowerCaseIdentifiers {
			Identifier = strings.ToLower(Identifier)
		}
		Identifier = strings.ReplaceAll(Identifier, " ", "_")

		Name := line[v.Columns.Name]
		Province := ProvinceMap[line[v.Columns.Province]]
		Latitude := line[v.Columns.Latitude]
		Longitude := line[v.Columns.Longitude]
		Elevation := line[v.Columns.Elevation]
		TimeOffset := p.TimeOffset

		TemplateName := fmt.Sprintf("T:%v", LineCount)

		ut, err := template.New(TemplateName).Funcs(template.FuncMap{
			"SpaceToUnderscore": func(Input string) string {
				return strings.ReplaceAll(Input, " ", "_")
			},
			"SplitWord": func(Input string, SepIn string, SepOut string, Begin int, End int) string {
				Words := strings.Split(Input, SepIn)

				if End == 0 {
					End = len(Words)
				}

				return strings.Join(Words[Begin:End], SepOut)
			},
		}).Parse(p.URL)
		if err != nil {
			fmt.Print("Error in template:\n\n")
			fmt.Println(p.URL)
			fmt.Println()
			log.Panic(err)
		}

		var u LoadtimeTemplateStruct
		u.Date = "{{ .Date }}"
		u.FullTimestamp = "{{ .FullTimestamp }}"
		u.DFOTimestamp = "{{ .DFOTimestamp }}"
		u.Identifier = Identifier
		u.Column = line

		var tpl bytes.Buffer
		err = ut.Execute(&tpl, u)

		if err != nil {
			log.Panic(err)
		}

		URL := tpl.String()

		var FinalIdentifier string
		if len(BasePath) > 0 {
			FinalIdentifier = BasePath + "/"
		}
		if len(p.Prefix) > 0 {
			FinalIdentifier += p.Prefix + "/"
		}

		// _, err = o.ProviderStatement.Exec(FinalIdentifier, "ni")
		// if err != nil {
		// 	log.Panic(err)
		// }

		FinalIdentifier += Identifier

		//fmt.Printf("%-9s %2s %-25s %s\n", FinalIdentifier, Province, Name, URL)

		err = Store(FinalIdentifier, Name, Province, Latitude, Longitude, Elevation, TimeOffset, URL)
		if err != nil {
			log.Panic(err)
		}

	}

	return err
}

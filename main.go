package main

import (
	"fmt"
	"time"

	"github.com/integrii/flaggy"
)

var Cache CacheStruct
var UTCLoc *time.Location
var LocalLoc *time.Location

func main() {

	var err error

	//	flaggy.DefaultParser.ShowHelpOnUnexpected = false

	subcommandGetStations := flaggy.NewSubcommand("getstations")
	subcommandGetTotals := flaggy.NewSubcommand("totalize")
	subcommandGetObservation := flaggy.NewSubcommand("observation")
	flaggy.AttachSubcommand(subcommandGetStations, 1)
	flaggy.AttachSubcommand(subcommandGetTotals, 1)
	flaggy.AttachSubcommand(subcommandGetObservation, 1)

	var gsQuery string
	var gsKML string
	subcommandGetStations.AddPositionalValue(&gsQuery, "Query", 1, false, "Station name or identifier to search for. Use SQL LIKE syntax.")
	subcommandGetStations.String(&gsKML, "k", "kml", "Export results to specified KML file.")

	var gtStation string
	var gtHour int
	var gtStartTime string
	var gtEndTime string
	subcommandGetTotals.AddPositionalValue(&gtStation, "Station", 1, false, "Station name.")
	subcommandGetTotals.Int(&gtHour, "o", "hours", "Relative Hour")
	subcommandGetTotals.String(&gtStartTime, "s", "starttime", "Absolute Time")
	subcommandGetTotals.String(&gtStartTime, "e", "endtime", "Absolute Time")

	var goStation string
	var goHour int
	var goTime string
	subcommandGetObservation.AddPositionalValue(&goStation, "Station", 1, false, "Station name. Use double quotes to specify multiple stations.")
	subcommandGetObservation.Int(&goHour, "o", "hours", "Relative Hour")
	subcommandGetObservation.String(&goTime, "d", "datetime", "Absolute Time")

	flaggy.Parse()

	Cache.Open()
	defer Cache.Close()

	UTCLoc, err = time.LoadLocation("UTC")
	if err != nil {
		fmt.Println(`Failed to load location "UTC"`)
	}

	LocalLoc, err = time.LoadLocation("Local")
	if err != nil {
		fmt.Println(`Failed to load location "Local"`)
	}

	if subcommandGetStations.Used {
		err = GetStations(gsQuery, gsKML)
		//
		//
	} else if subcommandGetTotals.Used {
		EndTime, err := time.Parse("2006-01-02 15 MST", gtEndTime)
		if err != nil {
			EndTime = time.Now()
		}
		EndTime = EndTime.Truncate(1 * time.Hour)

		StartTime, err := time.Parse("2006-01-02 15 MST", gtStartTime)
		if err != nil {
			if gtHour > 0 {
				gtHour = -gtHour
			}

			StartTime = EndTime.Add(time.Hour * time.Duration(gtHour))
		}
		StartTime = StartTime.Truncate(1 * time.Hour)

		err = GetTotals(gtStation, StartTime, EndTime)
		//
		//
	} else if subcommandGetObservation.Used {

		StartTime, err := time.Parse("2006-01-02 15 MST", goTime)
		if err != nil {
			if goHour > 0 {
				goHour = -goHour
			}

			StartTime = time.Now().Add(time.Hour * time.Duration(goHour))
		}
		StartTime = StartTime.Truncate(1 * time.Hour)

		err = PresentObservation(goStation, StartTime)
	} else {

		return
	}

	if err != nil {
		fmt.Println(err)
	}

	return
}

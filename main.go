package main

import (
	"fmt"
	"time"

	"github.com/integrii/flaggy"
)

var Cache CacheStruct
var Stations StationsStruct
var UTCLoc *time.Location
var LocalLoc *time.Location

func main() {

	var err error

	flaggy.DefaultParser.ShowHelpOnUnexpected = true
	flaggy.SetName("Weatherdata")
	flaggy.SetDescription("Extract and display Canadian weather observations from Environment Canada SWOB feeds.")

	subcommandGetStations := flaggy.NewSubcommand("getstations")
	subcommandGetStations.Description = "List or search for SWOB weather stations."
	flaggy.AttachSubcommand(subcommandGetStations, 1)

	subcommandGetTotals := flaggy.NewSubcommand("totalize")
	subcommandGetTotals.Description = "Totalize observations from one weather station over a specified time period."
	flaggy.AttachSubcommand(subcommandGetTotals, 1)

	subcommandGetObservation := flaggy.NewSubcommand("observation")
	subcommandGetObservation.Description = "Get observations from one or more weather stations at a specified time."
	flaggy.AttachSubcommand(subcommandGetObservation, 1)

	var gsQuery string
	var gsKML string

	subcommandGetStations.AddPositionalValue(&gsQuery, "Query", 1, false, "Station name or identifier to search for. Use SQL LIKE syntax.")
	subcommandGetStations.String(&gsKML, "k", "kml", "Export results to specified KML file.")

	var gtStation string
	var gtHour int = 6
	var gtStartTime string
	var gtEndTime string
	subcommandGetTotals.AddPositionalValue(&gtStation, "Station", 1, false, "Station name.")
	subcommandGetTotals.Int(&gtHour, "o", "hours", "Totalize observations over the past N hours.")
	subcommandGetTotals.String(&gtStartTime, "s", "starttime", "Totalize observations from the specified date and hour. Use the format \"YYYY-MM-DD HH TZ\"")
	subcommandGetTotals.String(&gtStartTime, "e", "endtime", "Totalize observations up to the specified date and hour. Not valid without -s. Use the format \"YYYY-MM-DD HH TZ\"")

	var goStation string
	var goHour int
	var goTime string
	subcommandGetObservation.AddPositionalValue(&goStation, "Station", 1, false, "Station name. Use double quotes to specify multiple stations.")
	subcommandGetObservation.Int(&goHour, "o", "hours", "Show observation from N hours ago.")
	subcommandGetObservation.String(&goTime, "d", "datetime", "Show observation from the specified date and hour. Use the format \"YYYY-MM-DD HH TZ\"")

	fmt.Print("\nWeatherdata Release 0 -- https://github.com/coridonhenshaw/weatherdata\n\n")

	flaggy.Parse()

	Cache.Open()
	defer Cache.Close()

	Stations.Import()

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

		err = PresentTotals(gtStation, StartTime, EndTime)
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
		flaggy.ShowHelp("")
	}

	if err != nil {
		fmt.Println(err)
	}

	return
}

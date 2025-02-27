package main

import (
	"fmt"
	"log"
	"time"
	"weatherdata/psychrometrics"

	"github.com/integrii/flaggy"
)

var Cache CacheStruct
var Stations StationsStruct
var UTCLoc *time.Location
var LocalLoc *time.Location

const TimeFmt string = "2006-01-02 15:04 MST"

var RefreshStationCache bool = false
var Verbose bool = false
var UTC bool = false

func main() {
	var err error

	flaggy.DefaultParser.ShowHelpOnUnexpected = true
	flaggy.SetName("Weatherdata")
	flaggy.SetDescription("Extract and display Canadian weather observations from Environment Canada SWOB feeds.")

	flaggy.Bool(&Verbose, "v", "verbose", "Enable verbose (debugging) output")
	flaggy.Bool(&RefreshStationCache, "r", "refreshcache", "Force refresh of station cache")
	flaggy.Bool(&UTC, "u", "utc", "Display all times in UTC")

	var SUI StationsUIStruct
	subcommandStations := flaggy.NewSubcommand("stations")
	subcommandStations.Description = "List or search for SWOB weather stations."
	subcommandStations.AddPositionalValue(&SUI.gsQuery, "Query", 1, false, "Station name or identifier to search for. Use SQL LIKE syntax.")
	subcommandStations.String(&SUI.gsKML, "k", "kml", "Export results to the specified KML file.")
	flaggy.AttachSubcommand(subcommandStations, 1)

	var OUI ObservationUIStruct
	subcommandObservation := flaggy.NewSubcommand("observation")
	subcommandObservation.AddPositionalValue(&OUI.Station, "Station", 1, false, "Station name. Use double quotes to specify multiple stations.")
	subcommandObservation.String(&OUI.AtTime, "t", "time", "Collect observation(s) closest to the time specified. Use the format \"YYYY-MM-DD HH:MM TZ\"")
	subcommandObservation.String(&OUI.StartTime, "s", "starttime", "Totalize observations from the specified date and hour. Use the format \"YYYY-MM-DD HH:MM TZ\"")
	subcommandObservation.String(&OUI.EndTime, "e", "endtime", "Totalize observations up to the specified date and hour. Not valid without -s. Use the format \"YYYY-MM-DD HH:MM TZ\"")
	subcommandObservation.Bool(&OUI.ShowIdentifiers, "i", "identifiers", "Show station identifiers in results table.")
	subcommandObservation.String(&OUI.RHAT, "r", "rhat", "Calculate relative humidity if outdoor air is heated to the specified temperature.")
	flaggy.AttachSubcommand(subcommandObservation, 1)

	var ATI AutotestStruct
	subcommandAutotest := flaggy.NewSubcommand("autotest")
	subcommandAutotest.Description = "Harvest observations from randomly selected weather stations from every network that reports into the SWOB system. Used for testing purposes."
	flaggy.AttachSubcommand(subcommandAutotest, 1)

	fmt.Print("\nWeatherdata Release 2 -- https://github.com/coridonhenshaw/weatherdata\n\n")

	flaggy.Parse()

	Cache.Open()
	defer Cache.Close()

	Stations = MakeStationsStruct()

	UTCLoc, err = time.LoadLocation("UTC")
	if err != nil {
		log.Panic(err)
	}

	LocalLoc, err = time.LoadLocation("Local")
	if err != nil {
		log.Panic(err)
	}

	psychrometrics.SetUnitSystem(psychrometrics.SI)

	if subcommandStations.Used {
		err = SUI.Get()
	} else if subcommandObservation.Used && len(OUI.Station) > 0 {
		err = OUI.Get()
	} else if subcommandAutotest.Used {
		err = ATI.Run()
	} else {
		flaggy.ShowHelp("")
	}

	if err != nil {
		fmt.Println(err)
	}

	return
}

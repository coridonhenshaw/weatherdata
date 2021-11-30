package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

var Cache CacheStruct

func main() {

	var err error

	Cache.Open()
	defer Cache.Close()

	Station := flag.String("s", "", "Station Identifier")
	Hour := flag.Int("hours", 0, "Relative Hour")
	StartTime := flag.String("starttime", "", "Absolute Time")
	EndTime := flag.String("endtime", "", "Absolute Time")
	Totals := flag.Bool("total", false, "Summarize Totals")
	flag.Parse()

	if *Hour > 0 {
		*Hour = -*Hour
	}

	if *Totals == false {
		StartTime, err := time.Parse("2006-01-02 15 MST", *StartTime)
		if err != nil {
			StartTime = time.Now().Add(time.Hour * time.Duration(*Hour))
		}
		StartTime = StartTime.Truncate(1 * time.Hour)

		err = GetOneObservation(*Station, StartTime)
	} else {
		EndTime, err := time.Parse("2006-01-02 15 MST", *EndTime)
		if err != nil {
			EndTime = time.Now()
		}
		EndTime = EndTime.Truncate(1 * time.Hour)

		StartTime, err := time.Parse("2006-01-02 15 MST", *StartTime)
		if err != nil {
			StartTime = EndTime.Add(time.Hour * time.Duration(*Hour))
		}
		StartTime = StartTime.Truncate(1 * time.Hour)

		err = GetTotals(*Station, StartTime, EndTime)
	}

	if err != nil {
		fmt.Println(err)
	}

	os.Exit(0)
}

package main

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

type ObservationUIStruct struct {
	Station         string
	AtTime          string
	StartTime       string
	EndTime         string
	RHAT            string
	ShowIdentifiers bool
}

type HarvestStruct struct {
	Station string
	ObsList []ObservationStruct
}

func (o *ObservationUIStruct) Get() error {
	var err error

	Now := time.Now().Truncate(1 * time.Minute)

	var Totalize bool

	var StartTime time.Time
	var EndTime time.Time
	var AtTime time.Time

	if len(o.StartTime) > 0 {
		Totalize = true

		EndTime, err = time.Parse(TimeFmt, o.EndTime)
		if err != nil {
			EndTime = Now
		}

		if strings.HasPrefix(o.StartTime, "-") {
			var Offset int
			_, err = fmt.Sscanf(o.StartTime, "-%d", &Offset)
			if err != nil {
				log.Panic(err)
			}

			StartTime = time.Now().UTC().Add(-time.Hour * time.Duration(Offset))

		} else {
			StartTime, err = time.Parse(TimeFmt, o.StartTime)
			if err != nil {
				log.Panic(err)
			}
		}

		StartTime = StartTime.In(UTCLoc)
		EndTime = EndTime.In(UTCLoc)

	} else {
		AtTime, err = time.Parse(TimeFmt, o.AtTime)
		if err != nil {
			AtTime = Now
			err = nil
		}
		AtTime = AtTime.In(UTCLoc)
	}

	if !Totalize {
		var Time string
		if UTC {
			Time = AtTime.Format(TimeFmt)
		} else {
			Time = AtTime.In(LocalLoc).Format(TimeFmt)
		}

		fmt.Printf("Reports from %s at %s:\n\n", o.Station, Time)
	} else {
		var ST string
		var ET string
		if UTC {
			ST = StartTime.Format(TimeFmt)
			ET = EndTime.Format(TimeFmt)
		} else {
			ST = StartTime.In(LocalLoc).Format(TimeFmt)
			ET = EndTime.In(LocalLoc).Format(TimeFmt)
		}
		fmt.Printf("Reports from %s over %s to %s:\n\n", o.Station, ST, ET)
	}

	var ObsList []ObservationStruct

	StationArray := strings.Fields(o.Station)

	var List []*HarvestStruct
	var Queue chan *HarvestStruct = make(chan *HarvestStruct, len(StationArray))

	for _, Station := range StationArray {
		var E HarvestStruct
		E.Station = Station
		List = append(List, &E)
		Queue <- &E
	}

	var wg sync.WaitGroup

	for i := 0; i < 8; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			Harvest(Queue, Totalize, StartTime, EndTime, AtTime)
		}()
	}

	wg.Wait()

	for _, i := range List {
		for _, j := range i.ObsList {
			ObsList = append(ObsList, j)
		}
	}

	OutTableText(ObsList, o.ShowIdentifiers, o.RHAT)

	return err
}

func Harvest(Queue chan *HarvestStruct,
	Totalize bool,
	StartTime time.Time,
	EndTime time.Time,
	AtTime time.Time) {

	var err error
	var Piece *HarvestStruct

Outer:
	for {
		select {
		case Piece = <-Queue:
		default:
			break Outer
		}

		if !Totalize {

			var obs ObservationStruct

			// fmt.Println(Piece.Station, AtTime, false)

			obs, err = Stations.GetObservation(Piece.Station, AtTime, false)
			if err != nil {
				obs.err = err
			}
			Piece.ObsList = append(Piece.ObsList, obs)
		} else {
			var obs []ObservationStruct
			obs = Stations.GetObservationSpan(Piece.Station, StartTime, EndTime)

			for _, v := range obs {
				Piece.ObsList = append(Piece.ObsList, v)
			}
		}
	}
}

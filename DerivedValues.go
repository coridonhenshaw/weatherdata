package main

import (
	"math"
	"strconv"
	"time"
	"weatherdata/psychrometrics"
)

type DerivedValuesStruct struct {
	Obs *ObservationStruct

	Timestamp     string
	DewPointDelta ValueStruct
	WindDirection ValueStruct
	Remarks       string
	RHAT          ValueStruct
}

func GenerateDerivedValues(Obs *ObservationStruct, RHAT string) (rc DerivedValuesStruct) {
	var err error
	rc.Obs = Obs

	if Obs.err != nil {
		rc.Remarks = Obs.err.Error()
	}

	rc.CalcDewPointDifference()
	rc.RewriteTimestamp()
	rc.RewriteWindDirection()

	RHATf, err := strconv.ParseFloat(RHAT, 64)

	if len(RHAT) > 0 && err == nil {
		rc.CalcRHAT(RHATf)
	}

	return rc
}

func (o *DerivedValuesStruct) CalcDewPointDifference() {

	if !o.Obs.Temperature.Valid || !o.Obs.DewPoint.Valid {
		return
	}

	o.DewPointDelta.LoadF(math.Abs((o.Obs.Temperature.Float - o.Obs.DewPoint.Float)))

	return
}

func (o *DerivedValuesStruct) RewriteTimestamp() {
	t, err := time.Parse("2006-01-02T15:04:05.000Z", o.Obs.Timestamp)
	if err != nil {
		o.Timestamp = "------ ----"
		return
	}

	if UTC {
		o.Timestamp = t.Format("060102 1504")
	} else {
		o.Timestamp = t.In(LocalLoc).Format("060102 1504")
	}
}

func (o *DerivedValuesStruct) RewriteWindDirection() {

	if !o.Obs.AverageWindDirection.Valid {
		return
	}

	Directions := []string{"  N", "NNE", " NE", "ENE", "  E", "ESE", " SE", "SSE", "  S", "SSW", " SW", "WSW", "  W", "WNW", " NW", "NNW", "  N"}

	var Index int = (int)(math.Round((float64)(o.Obs.AverageWindDirection.Float) / (float64)(22.5)))

	o.WindDirection.String = o.Obs.AverageWindDirection.String + " " + Directions[Index]
}

func (o *DerivedValuesStruct) CalcRHAT(TDryBulbOut float64) {
	var TDryBulb float64
	var RelHum float64
	var Pressure float64

	if !o.Obs.Temperature.Valid || !o.Obs.Pressure.Valid || !o.Obs.Humidity.Valid {
		return
	}

	TDryBulb = o.Obs.Temperature.Float
	RelHum = o.Obs.Humidity.Float / 100
	Pressure = o.Obs.Pressure.Float * 100

	HumRatio := psychrometrics.GetHumRatioFromRelHum(TDryBulb, RelHum, Pressure)

	o.RHAT.LoadF(
		psychrometrics.GetRelHumFromHumRatio(TDryBulbOut, HumRatio, Pressure) * 100)
}

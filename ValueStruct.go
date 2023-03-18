package main

import "strconv"

type ValueStruct struct {
	Float  float64
	String string
	Valid  bool
}

func (o *ValueStruct) Load(Val string) {
	var err error
	o.Float, err = strconv.ParseFloat(Val, 64)
	if err == nil {
		o.Valid = true
		o.String = RoundFloatTwoPlaces(o.Float)
	} else {
		o.String = Val
	}
}

func (o *ValueStruct) LoadF(Val float64) {
	o.Float = Val
	o.String = RoundFloatTwoPlaces(Val)
	o.Valid = true
}

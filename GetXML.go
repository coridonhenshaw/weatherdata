package main

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func MakeBaseURL(Station string) (string, error) {
	var OutString string

	reECStation := regexp.MustCompile(`C...-((AUTO)|(MAN))`)
	rePartners := regexp.MustCompile(`partners/(.*)/(.*)`)

	if reECStation.MatchString(Station) == true {
		OutString = "https://dd.weather.gc.ca/observations/swob-ml/{{ .Date }}/" + Station[0:4] + "/{{ .FullTimestamp }}-" + Station + "-swob.xml"
	} else if rePartners.MatchString(Station) == true {
		Parts := rePartners.FindStringSubmatch(Station)
		OutString = "https://dd.weather.gc.ca/observations/swob-ml/partners/" + Parts[1] + "/{{ .Date }}/" + Parts[2] + "/{{ .FullTimestamp }}-" + Parts[1] + "-" + Parts[2] + "-AUTO-swob.xml"
	} else {
		return "", errors.New("Invalid station name")
	}

	return OutString, nil
}

type URLDateTime struct {
	Date          string
	FullTimestamp string
}

func InjectURLDateTime(BaseURL string, DateTime time.Time) (string, error) {

	u := URLDateTime{DateTime.Format("20060102"), DateTime.Format("2006-01-02-1500")}

	ut, err := template.New("foo").Parse(BaseURL)

	if err != nil {
		panic(err)
	}

	var tpl bytes.Buffer
	err = ut.Execute(&tpl, u)

	if err != nil {
		panic(err)
	}

	return tpl.String(), nil
}

func HTTPSGet(URL string) (string, error) {

	BodyText, err := Cache.Get(URL)

	if err != nil {
		fmt.Println("HTTPSGet cache failure", err)
		os.Exit(1)
	}

	if len(BodyText) > 0 {
		return BodyText, nil
	}

	//	return "", nil

	resp, err := http.Get(URL)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if !(resp.StatusCode >= 200 && resp.StatusCode <= 299) {
		return "", errors.New(strconv.Itoa(resp.StatusCode))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	BodyText = string(body)

	err = Cache.Put(URL, BodyText)
	if err != nil {
		fmt.Println("Cache put failure", err)
		os.Exit(1)
	}

	return BodyText, nil
}

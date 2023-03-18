package main

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"
)

func HTTPSGet(URL string, MaxAge time.Duration) (string, error) {

	BodyText, err := Cache.Get(URL, MaxAge)
	if err != nil {
		return "", err
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	BodyText = string(body)

	err = Cache.Put(URL, BodyText)
	if err != nil {
		return "", err
	}

	return BodyText, nil
}

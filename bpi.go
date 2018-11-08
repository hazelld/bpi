package bpi

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

const URL string = "http://data.nba.net/10s%s"

func MakeRequest(endpoint string) (string, error) {
	response, err := http.Get(fmt.Sprintf(URL, endpoint))

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	json_resp, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return string(json_resp), nil
}

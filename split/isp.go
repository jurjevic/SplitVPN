package split

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"time"
)

type Isp struct {
	Query         string  `json:"query"`
	Status        string  `json:"status"`
	Continent     string  `json:"continent"`
	ContinentCode string  `json:"continentCode"`
	Country       string  `json:"country"`
	CountryCode   string  `json:"countryCode"`
	Region        string  `json:"region"`
	RegionName    string  `json:"regionName"`
	City          string  `json:"city"`
	District      string  `json:"district"`
	Zip           string  `json:"zip"`
	Lat           float64 `json:"lat"`
	Lon           float64 `json:"lon"`
	Timezone      string  `json:"timezone"`
	Offset        int     `json:"offset"`
	Currency      string  `json:"currency"`
	Isp           string  `json:"isp"`
	Org           string  `json:"org"`
	As            string  `json:"as"`
	Asname        string  `json:"asname"`
	Mobile        bool    `json:"mobile"`
	Proxy         bool    `json:"proxy"`
	Hosting       bool    `json:"hosting"`
}

func fetchISP() Isp {
	cmd := "nscurl http://ip-api.com/json?fields=66846719"
	out, err := exec.Command("/bin/sh", "-c", cmd).Output()
	Debugln(string(out))
	if err != nil {
		log.Printf("Failed to fetch ISP. %s", err.Error())
		return Isp{}
	}

	isp := Isp{}
	jsonErr := json.Unmarshal(out, &isp)
	if jsonErr != nil {
		log.Println(jsonErr.Error())
		return Isp{}
	}

	return isp
}

func fetchNoProxyISP() Isp {
	url := "http://ip-api.com/json?fields=66846719"

	ispClient := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println(err.Error())
		return Isp{}
	}

	res, getErr := ispClient.Do(req)
	if getErr != nil {
		log.Println(getErr.Error())
		return Isp{}
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Println(readErr.Error())
		return Isp{}
	}

	Debugln(string(body))
	isp := Isp{}
	jsonErr := json.Unmarshal(body, &isp)
	if jsonErr != nil {
		log.Println(jsonErr.Error())
		return Isp{}
	}

	return isp
}

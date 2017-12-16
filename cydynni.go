package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const apiURL = "https://cydynni.org.uk/"

func login(password string, email string) *http.Cookie {

	login := "login?"

	fullurl := apiURL + login + "email=" + email + "&password=" + url.QueryEscape(password)
	print("fullurl:" + fullurl + "\n")
	u, _ := url.ParseRequestURI(fullurl)

	//u.Path = login + "email=" + email + "&password=" + password
	//urlStr := fmt.Sprintf("%v", u) // "https://api.com/user/"
	client := &http.Client{}
	r, err := http.NewRequest("GET", u.String(), bytes.NewBufferString("")) // <-- URL-encoded payload
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Do(r)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(resp.Status)
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	newStr := buf.String()

	//io.Copy(os.Stdout, resp.Body)

	var res map[string]interface{}
	json.Unmarshal([]byte(newStr), &res)
	fmt.Println(newStr)
	//res["apikey"](string)
	//if (strings.Compare("a","") >0)
	if s, ok := res["apikey"].(string); ok {
		if strings.Compare(s, "") > 0 {
		}
		return resp.Cookies()[1]
	}
	log.Fatal(err)
	panic(res)

}

func getString(path string, password string, email string) string {

	uhome, _ := url.ParseRequestURI(apiURL + path)
	rhome, err := http.NewRequest("GET", uhome.String(), bytes.NewBufferString("")) // <-- URL-encoded payload
	logincookie := login(password, email)
	rhome.Header.Set("Cookie", logincookie.Name+"="+logincookie.Value)

	//fmt.Println(resp.Cookies()[0].Name + "=" + resp.Cookies()[0].Value)
	//rhome.Header.Add("Cookie", resp.Cookies()[0].Name+"="+resp.Cookies()[0].Value)
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{}

	resphome, err := client.Do(rhome)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(resp.Status)
	bufhome := new(bytes.Buffer)
	bufhome.ReadFrom(resphome.Body)
	newStrHome := bufhome.String()
	return newStrHome

}

func summary(summaryJSON string) (Summary, error) {
	s := Summary{}
	err := json.Unmarshal([]byte(summaryJSON), &s)
	return s, err
}

type DataPoints []DataPoint
type DataPoint struct {
	Time  time.Time
	Value float64
}

func (dp *DataPoint) UnmarshalJSON(b []byte) error {
	var raw []float64
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	if len(raw) != 2 {
		return fmt.Errorf("Invalid datapoint data; expected 2 items in list but had %d", len(raw))
	}

	dp.Time = time.Unix(int64(raw[0])/1000, 0)
	dp.Value = raw[1]
	return nil
}

type powerMap map[time.Time]float64

func timePowerMap(jsonData string) (powerMap, error) {
	var dps DataPoints
	pm := make(powerMap)

	if err := json.Unmarshal([]byte(jsonData), &dps); err != nil {
		return pm, err
	}

	for _, dp := range dps {
		pm[dp.Time] = dp.Value
	}
	return pm, nil
}

func main() {
	email := ""
	password := ""
	flag.StringVar(&email, "email", "email", "email address used to log in")
	flag.StringVar(&password, "password", "password", "password needed to log in")
	flag.Parse()

	//communitySummary := "community/data"
	homeHalfHourData := "data"
	comunityHalfHourData := "community/halfhourlydata"
	//homeSummary, _ := summary(getString("household/data"))
	//communitySummary, _ := summary(getString("community/data"))
	hydroHalfHour := "hydro"

	home, err := timePowerMap(getString(homeHalfHourData, password, email))
	if err != nil {
		panic(err)
	}
	fmt.Println(home)

	hydroHalf, err := timePowerMap(getString(hydroHalfHour, password, email))
	if err != nil {
		panic(err)
	}
	fmt.Println(hydroHalf)

	communityHalf, err := timePowerMap(getString(comunityHalfHourData, password, email))
	if err != nil {
		panic(err)
	}
	fmt.Println(communityHalf)
	//fmt.Println(getasJSON(getString(communitySummary)))
	//fmt.Println(getString(comunityHalfHourData))

}

type Power struct {
	Morning   float32 `json:"Morning,omitempty"`
	Midday    float32 `json:"Midday,omitempty"`
	Evening   float32 `json:"Evening,omitempty"`
	Overnight float32 `json:"Overnight,omitempty"`
	Hydro     float32 `json:"Hydro,omitempty"`
	Total     float32 `json:"Total,omitempty"`
}

type Summary struct {
	KWH       *Power `json:"kwh,omitempty"`
	Cost      *Power `json:"cost,omitempty"`
	Month     string `json:"month,omitempty"`
	Day       string `json:"day,omitempty"`
	DayOffset int    `json:"dayoffset,omitempty"`
}

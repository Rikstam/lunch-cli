package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type sodexoLunchCli struct {
	apiBase string
}

type Course struct {
	Title     string `json:"title_fi"`
	DietCodes string `json:"dietcodes"`
}

type LunchInfo struct {
	Meta struct {
		RestaurantName string `json:"ref_title"`
	} `json:"meta"`
	Courses map[int]Course `json:"courses"`
}

func (s sodexoLunchCli) getLunch(restaurantUrl string) error {
	resp, err := http.Get(s.apiBase + restaurantUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var data struct {
		LunchInfo
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Fatal(err)
		return err
	}

	fmt.Println(data.Meta.RestaurantName)
	for k, v := range data.Courses {
		fmt.Printf("%v -> %s\n", k, v)
	}
	return nil
}

func main() {
	now := time.Now()
	todayString := now.Format("2006-01-02")
	fmt.Println(todayString)
	lc := sodexoLunchCli{apiBase: "https://www.sodexo.fi/ruokalistat/output/daily_json/"}
	lc.getLunch("107/" + todayString)
	lc.getLunch("110/" + todayString)
}

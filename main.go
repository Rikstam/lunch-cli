package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

var client = http.Client{}

const apiBase = "https://www.sodexo.fi/ruokalistat/output/daily_json/"

type Restaurant struct {
	Id   string
	Name string
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

func getLunch(ctx context.Context, label, restaurantUrl string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, restaurantUrl, nil)

	if err != nil {
		fmt.Println(label, "request err:", err)
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(label, "response err:", err)
		return err
	}

	err = printLunch(resp.Body)

	if err != nil {
		fmt.Println(label, "decoding err for:", label, err)
		return err
	}

	return nil
}

func printLunch(body io.ReadCloser) error {
	var data struct {
		LunchInfo
	}
	if err := json.NewDecoder(body).Decode(&data); err != nil {
		log.Fatal("error decoding response: ", err)
		return err
	}

	fmt.Println(data.Meta.RestaurantName)
	for k, v := range data.Courses {
		fmt.Printf("%v -> %s\n", k, v)
	}
	return nil
}

func callAll(ctx context.Context, restaurants []Restaurant, timeNow string) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	var wg sync.WaitGroup
	for _, r := range restaurants {
		wg.Add(1)
		restaurant := r
		go func() {
			defer wg.Done()
			err := getLunch(ctx, restaurant.Name, apiBase+restaurant.Id+"/"+timeNow)
			if err != nil {
				cancel()
			}
		}()
	}
	wg.Wait()
	fmt.Println("done with both")
}

func main() {
	now := time.Now()
	todayString := now.Format("2006-01-02")
	restaurants := []Restaurant{
		{Id: "107", Name: "sodexo5"},
		{Id: "110", Name: "sodexo6"},
	}
	fmt.Println("Sodexo lunch for date: ", todayString)

	ctx := context.Background()
	callAll(ctx, restaurants, todayString)
}

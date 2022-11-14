package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

var client = http.Client{}

const apiBase = "https://www.sodexo.fi/ruokalistat/output/daily_json/"

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

	var data struct {
		LunchInfo
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		fmt.Println(label, "decoding err:", err)
		log.Fatal(err)
		return err
	}

	fmt.Println(data.Meta.RestaurantName)
	for k, v := range data.Courses {
		fmt.Printf("%v -> %s\n", k, v)
	}
	return nil
}

func callAll(ctx context.Context, r1, r2 string) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		err := getLunch(ctx, "sodexo5", r1)
		if err != nil {
			cancel()
		}
	}()

	go func() {
		defer wg.Done()
		err := getLunch(ctx, "sodexo6", r2)
		if err != nil {
			cancel()
		}
	}()

	wg.Wait()
	fmt.Println("done with both")
}

func main() {
	now := time.Now()
	todayString := now.Format("2006-01-02")
	fmt.Println(todayString)

	ctx := context.Background()
	s5Url := apiBase + "107/" + todayString
	s6Url := apiBase + "110/" + todayString

	callAll(ctx, s5Url, s6Url)
}

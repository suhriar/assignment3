package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"text/template"
	"time"
)

type StatusData struct {
	Status struct {
		Water int `json:"water"`
		Wind  int `json:"wind"`
	} `json:"status"`
}

func main() {
	go AutoReloadJSON()
	http.HandleFunc("/", AutoReloadWeb)
	http.Handle("/asset/", http.StripPrefix("/asset/", http.FileServer(http.Dir("asset"))))
	fmt.Println("listening on PORT:", ":8080")
	http.ListenAndServe(":8080", nil)
}

func AutoReloadJSON() {
	for {
		min := 1
		max := 25
		wind := rand.Intn(max-min) + min
		water := rand.Intn(max-min) + min

		data := StatusData{}
		data.Status.Wind = wind
		data.Status.Water = water

		jsonData, err := json.Marshal(data)

		if err != nil {
			log.Fatal("error occured while marshalling status data:", err.Error())
		}
		err = ioutil.WriteFile("data.json", jsonData, 0644)

		if err != nil {
			log.Fatal("error occured while writing data to data.json file", err.Error())
		}
		time.Sleep(15 * time.Second)
	}
}

func AutoReloadWeb(w http.ResponseWriter, r *http.Request) {
	fileData, err := ioutil.ReadFile("data.json")

	if err != nil {
		log.Fatal("error occured while reading data from data.json file", err.Error())
	}

	var statusData StatusData

	err = json.Unmarshal(fileData, &statusData)
	if err != nil {
		log.Fatal("error occured while unMarshalling from data.json file", err.Error())
	}

	waterVal := statusData.Status.Water
	windVal := statusData.Status.Wind

	var (
		waterStatus string
		windStatus  string
	)

	waterValue := strconv.Itoa(waterVal)
	windValue := strconv.Itoa(windVal)

	switch {
	case waterVal <= 5:
		waterStatus = "Aman"
	case waterVal >= 6 && waterVal <= 8:
		waterStatus = "Siaga"
	case waterVal > 8:
		waterStatus = "Bahaya"
	default:
		waterStatus = "Water Value not defined"
	}

	switch {
	case windVal <= 6:
		windStatus = "Aman"
	case windVal >= 7 && windVal <= 15:
		windStatus = "Siaga"
	case windVal > 15:
		windStatus = "Bahaya"
	default:
		windStatus = "Wind Value not defined"
	}

	data := map[string]string{
		"waterStatus": waterStatus,
		"windStatus":  windStatus,
		"waterValue":  waterValue,
		"windValue":   windValue,
	}

	tpl, err := template.ParseFiles("index.html")

	if err != nil {
		log.Fatal("error parsing html:", err.Error())
	}

	tpl.Execute(w, data)

}

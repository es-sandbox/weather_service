package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/boltdb/bolt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func dataHandler(resp http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		allRecords := getAllRecordsWithErrorSuppressing()

		jsonText, err := json.Marshal(allRecords)
		if err != nil {
			fmt.Println(err)
			return
		}

		resp.Header().Set("Access-Control-Allow-Origin", "*")

		if _, err := resp.Write(jsonText); err != nil {
			fmt.Println(err)
		}

	case "POST":
		data, err := ioutil.ReadAll(req.Body)
		if err != nil {
			fmt.Printf("can't read data: %v\n", err)
			return
		}

		fmt.Printf("RECEIVED LEN: %v\n", len(data))
		fmt.Printf("RECEIVED: %v\n", string(data))

		weatherInfo := weatherInfo{}
		if err := json.Unmarshal(data, &weatherInfo); err != nil {
			fmt.Printf("can't parse data: %v\n", err)
			return
		}

		err = db.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists(bucketName)
			if err != nil {
				return err
			}

			id, _ := b.NextSequence()
			weatherInfo.ID = id

			weatherInfo.TimeStamp = time.Now().UnixNano()

			// TODO(evg): check if TimeStamp exists
			/*
			if weatherInfo.TimeStamp == 0 {
				weatherInfo.TimeStamp = time.Now().UnixNano()
			} else {
				weatherInfo.TimeStamp *= 1e9
			}
			*/

			data, err := weatherInfo.Serialize()
			if err != nil {
				return err
			}

			return b.Put(itob(id), data)
		})
		if err != nil {
			fmt.Printf("can't write data: %v\n", err)
			return
		}
	}
}

func dataLastHandler(resp http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		weatherInfo, err := getLastRecord()
		if err != nil {
			fmt.Println(err)
			return
		}

		jsonText, err := json.Marshal(weatherInfo)
		if err != nil {
			fmt.Println(err)
			return
		}

		resp.Header().Set("Access-Control-Allow-Origin", "*")

		if _, err := resp.Write(jsonText); err != nil {
			fmt.Println(err)
		}
	}
}

func dataLastHourAvgHandler(resp http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		avg := avg(getLastHourRecords())

		jsonText, err := json.Marshal(avg)
		if err != nil {
			fmt.Println(err)
			return
		}

		if _, err := resp.Write(jsonText); err != nil {
			fmt.Println(err)
		}
	}
}

func dataLastHourHandler(resp http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		lastHourRecords := getLastHourRecords()

		jsonText, err := json.Marshal(lastHourRecords)
		if err != nil {
			fmt.Println(err)
			return
		}

		resp.Header().Set("Access-Control-Allow-Origin", "*")

		if _, err := resp.Write(jsonText); err != nil {
			fmt.Println(err)
		}
	}
}

func dataLastDayAvgHandler(resp http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		avgLastDay := make([]*weatherInfo, 24)
		for i := 0; i < 24; i++ {
			left := time.Now().Add(-time.Hour * time.Duration(i+1))
			right := time.Now().Add(-time.Hour * time.Duration(i))

			hourlyRecords := getBoundedInTimeRecords(left, right)
			avgLastDay[i] = avg(hourlyRecords)
			if avgLastDay[i] != nil {
				avgLastDay[i].ID = uint64(i)
			}
		}

		jsonText, err := json.Marshal(avgLastDay)
		if err != nil {
			fmt.Println(err)
			return
		}

		if _, err := resp.Write(jsonText); err != nil {
			fmt.Println(err)
		}
	}
}

func dataLastDayHandler(resp http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		lastDayRecords := getLastDayRecords()

		jsonText, err := json.Marshal(lastDayRecords)
		if err != nil {
			fmt.Println(err)
			return
		}

		resp.Header().Set("Access-Control-Allow-Origin", "*")

		if _, err := resp.Write(jsonText); err != nil {
			fmt.Println(err)
		}
	}
}

func dataLastMinuteHandler(resp http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		lastMinuteRecords := getLastMinuteRecords()

		jsonText, err := json.Marshal(lastMinuteRecords)
		if err != nil {
			fmt.Println(err)
			return
		}

		resp.Header().Set("Access-Control-Allow-Origin", "*")

		if _, err := resp.Write(jsonText); err != nil {
			fmt.Println(err)
		}
	}
}

func launchCloudTest() {
	weatherInfo := &weatherInfo{
		TempOUT:       42,
		Humidity:      43,
		TempIN:        44,
		Pressure:      45,
		WindSpeed:     46,
		WindDirection: 47,
		Rainfall:      48,
		Battery:       49,
		Thunder:       50,
		Light:         51,
		Charging:      52,
	}

	if err := weatherInfo.sendToCloud(); err != nil {
		fmt.Println(err)
		return
	}

	weatherInfo, err := newWeatherInfoFromCloud()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(weatherInfo)
}

func getWeatherInfoFromCloud() (*weatherInfo, error) {
	return newWeatherInfoFromCloud()
}

func main() {
	listenAddr := flag.String("listen", "0.0.0.0:9000", "address of http server, format: host:port")
	daPath := flag.String("dbpath", "data.db", "absolute path to database, example: /tmp/my.db")
	cloudTest := flag.Bool("cloud", false, "launch cloud testing and exit")
	cloudInfo := flag.Bool("cloud_info", false, "show current cloud info")
	flag.Parse()

	if *cloudTest {
		launchCloudTest()
		return
	}

	if *cloudInfo {
		weatherInfo, err := getWeatherInfoFromCloud()
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(weatherInfo)
		return
	}

	var err error
	db, err = bolt.Open(*daPath, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/data", dataHandler)
	http.HandleFunc("/data/last", dataLastHandler)
	http.HandleFunc("/data/last_hour", dataLastHourHandler)
	http.HandleFunc("/data/last_hour/avg", dataLastHourAvgHandler)
	http.HandleFunc("/data/last_day", dataLastDayHandler)
	http.HandleFunc("/data/last_day/avg", dataLastDayAvgHandler)
	http.HandleFunc("/data/last_minute", dataLastMinuteHandler)
	fmt.Printf("listen on: %v\n", *listenAddr)
	http.ListenAndServe(*listenAddr, nil)
}

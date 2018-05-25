package main

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"time"
)

var (
	bucketDoesNotExistError = errors.New("bucket does not exist")
	keyDoesNotExist         = errors.New("key does not exist")

	bucketName = []byte("bucket")
	keyName    = []byte("key")
)

var db *bolt.DB

type weatherInfo struct {
	ID        uint64
	TimeStamp int64 // Unix TimeStamp in nsec

	TempOUT       int
	Humidity      int
	TempIN        float64
	Pressure      float64
	WindSpeed     float64
	WindDirection int
	Rainfall      int
	Battery       int
	Thunder       int
	Light         float64
	Charging      int

	Latitude  float64
	Longitude float64

	Fire int
	Snow int
	SOS  int
}

func (w *weatherInfo) GetTime() time.Time {
	sec := w.TimeStamp / 1e9
	nsec := w.TimeStamp % 1e9
	return time.Unix(sec, nsec)
}

func (w *weatherInfo) String() string {
	tmpl := `
	ID            %v
	TimeStamp     %v
	TempOUT       %v
	Humidity      %v
	TempIN        %v
	Pressure      %v
	WindSpeed     %v
	WindDirection %v
	Rainfall      %v
	Battery       %v
	Thunder       %v
	Light         %v
	Charging      %v
	`
	return fmt.Sprintf(tmpl, w.ID, w.TimeStamp, w.TempOUT, w.Humidity, w.TempIN, w.Pressure, w.WindSpeed, w.WindDirection,
		w.Rainfall, w.Battery, w.Thunder, w.Light, w.Charging)
}

func (w *weatherInfo) Serialize() ([]byte, error) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(w); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (w *weatherInfo) Deserialize(data []byte) error {
	buf := bytes.NewBuffer(data)
	return gob.NewDecoder(buf).Decode(w)
}

func getAllRecords() ([]*weatherInfo, error) {
	dataSlice := make([][]byte, 0)

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		if b == nil {
			return bucketDoesNotExistError
		}

		return b.ForEach(func(key, value []byte) error {
			data := make([]byte, len(value))
			copy(data, value)
			dataSlice = append(dataSlice, data)
			return nil
		})
	})
	if err != nil {
		return nil, fmt.Errorf("db's view error: %v\n", err)
	}

	weatherInfoSlice := make([]*weatherInfo, len(dataSlice))

	for i, data := range dataSlice {
		weatherInfo := &weatherInfo{}
		if err := weatherInfo.Deserialize(data); err != nil {
			return nil, err
		}

		weatherInfoSlice[i] = weatherInfo
	}
	return weatherInfoSlice, nil
}

func getAllRecordsWithErrorSuppressing() []*weatherInfo {
	allRecords, err := getAllRecords()
	if err != nil {
		panic(err)
	}
	return allRecords
}

func getBoundedInTimeRecords(left, right time.Time) []*weatherInfo {
	allRecords := getAllRecordsWithErrorSuppressing()

	boundedInTimeRecords := make([]*weatherInfo, 0)
	for _, record := range allRecords {
		if record.GetTime().After(left) && record.GetTime().Before(right) {
			boundedInTimeRecords = append(boundedInTimeRecords, record)
		}
	}
	return boundedInTimeRecords
}

func getLastRecord() (*weatherInfo, error) {
	allRecords := getAllRecordsWithErrorSuppressing()

	if len(allRecords) == 0 {
		return nil, errors.New("element does not exist")
	}

	return allRecords[len(allRecords)-1], nil
}

func getLastHourRecords() []*weatherInfo {
	left := time.Now().Add(-time.Hour)
	right := time.Now()
	return getBoundedInTimeRecords(left, right)
}

func getLastDayRecords() []*weatherInfo {
	left := time.Now().Add(-time.Hour * 24)
	right := time.Now()
	return getBoundedInTimeRecords(left, right)
}

func getLastMinuteRecords() []*weatherInfo {
	left := time.Now().Add(-time.Minute)
	right := time.Now()
	return getBoundedInTimeRecords(left, right)
}

func avg(recordSlice []*weatherInfo) *weatherInfo {
	avg := weatherInfo{}

	null := 0

	for _, record := range recordSlice {
		if record == nil {
			null++
			continue
		}

		avg.TempOUT += record.TempOUT
		avg.Humidity += record.Humidity
		avg.TempIN += record.TempIN
		avg.Pressure += record.Pressure
		avg.WindSpeed += record.WindSpeed
		avg.WindDirection += record.WindDirection
		avg.Rainfall += record.Rainfall
	}

	div := len(recordSlice) - null

	if div == 0 {
		return nil
	}

	avg.TempOUT /= div
	avg.Humidity /= div
	avg.TempIN /= float64(div)
	avg.Pressure /= float64(div)
	avg.WindSpeed /= float64(div)
	avg.WindDirection /= div
	avg.Rainfall /= div

	return &avg
}

// itob returns an 8-byte big endian representation of v.
func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

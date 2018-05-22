package main

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/boltdb/bolt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var (
	bucketDoesNotExistError = errors.New("bucket does not exist")
	keyDoesNotExist         = errors.New("key does not exist")

	bucketName = []byte("bucket")
	keyName    = []byte("key")
)

var db *bolt.DB

// TODO(evg): review it
type weatherInfo struct {
	ID        uint64
	TimeStamp int64 // Unix TimeStamp

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
}

// TODO(evg): remove it?
func (w *weatherInfo) String() string {
	return fmt.Sprintf("Temp: %v", w.TempOUT)
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

func dataHandler(resp http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		// var jsonText []byte
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
			fmt.Printf("db's view error: %v\n", err)
			return
		}

		weatherInfoSlice := make([]*weatherInfo, len(dataSlice))

		for i, data := range dataSlice {
			weatherInfo := &weatherInfo{}
			if err := weatherInfo.Deserialize(data); err != nil {
				fmt.Println(err)
				return
			}

			weatherInfoSlice[i] = weatherInfo
		}

		jsonText, err := json.Marshal(weatherInfoSlice)
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

// itob returns an 8-byte big endian representation of v.
func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func main() {
	listenAddr := flag.String("listen", "0.0.0.0:9000", "address of http server, format: host:port")
	daPath := flag.String("dbpath", "data.db", "absolute path to database, example: /tmp/my.db")
	flag.Parse()

	var err error
	db, err = bolt.Open(*daPath, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/data", dataHandler)
	http.ListenAndServe(*listenAddr, nil)
}

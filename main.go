package main

import (
	"net/http"
	"io/ioutil"
	"fmt"
	"encoding/json"

	"github.com/boltdb/bolt"
	"log"
	"encoding/gob"
	"bytes"
	"errors"
	"flag"
)

var (
	bucketDoesNotExistError = errors.New("bucket does not exist")
	keyDoesNotExist = errors.New("key does not exist")

	bucketName = []byte("bucket")
	keyName = []byte("key")
)

var db *bolt.DB

type weatherInfo struct {
	Temp int
}

func (w *weatherInfo) String() string {
	return fmt.Sprintf("Temp: %v", w.Temp)
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
		var jsonText []byte

		err := db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket(bucketName)
			if b == nil {
				return bucketDoesNotExistError
			}

			data := b.Get(keyName)
			if data == nil {
				return keyDoesNotExist
			}

			weatherInfo := weatherInfo{}
			if err := weatherInfo.Deserialize(data); err != nil {
				return err
			}

			jsonTextLocal, err := json.Marshal(weatherInfo)
			if err != nil {
				return err
			}

			jsonText = make([]byte, len(jsonTextLocal))
			copy(jsonText, jsonTextLocal)
			return nil
		})
		if err != nil {
			fmt.Printf("db's view error: %v\n", err)
			return
		}

		resp.Write(jsonText)

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

			data, err := weatherInfo.Serialize()
			if err != nil {
				return err
			}

			return b.Put(keyName, data)
		})
		if err != nil {
			fmt.Printf("can't write data: %v\n", err)
			return
		}
	}
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
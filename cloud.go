package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

const (
	personalCloudKey = "w1co9s82"
)

func (w *weatherInfo) sendToCloud() error {
	cloudSender := newCloudSender()
	cloudSender.send("TempOUT", strconv.Itoa(w.TempOUT))
	cloudSender.send("Humidity", strconv.Itoa(w.Humidity))
	cloudSender.send("TempIN", strconv.Itoa(int(w.TempIN)))
	cloudSender.send("Pressure", strconv.Itoa(int(w.Pressure)))
	cloudSender.send("WindSpeed", strconv.Itoa(int(w.WindSpeed)))
	cloudSender.send("WindDirection", strconv.Itoa(w.WindDirection))
	cloudSender.send("Rainfall", strconv.Itoa(w.Rainfall))
	cloudSender.send("Battery", strconv.Itoa(w.Battery))
	cloudSender.send("Thunder", strconv.Itoa(w.Thunder))
	cloudSender.send("Light", strconv.Itoa(int(w.Light)))
	cloudSender.send("Charging", strconv.Itoa(w.Charging))
	return cloudSender.error()
}

func newWeatherInfoFromCloud() (*weatherInfo, error) {
	cloudGetter := newCloudGetter()
	temperatureOutside := cloudGetter.getAndConv("TempOUT")
	humidity := cloudGetter.getAndConv("Humidity")
	temperatureInside := cloudGetter.getAndConv("TempIN")
	pressure := cloudGetter.getAndConv("Pressure")
	windSpeed := cloudGetter.getAndConv("WindSpeed")
	windDirection := cloudGetter.getAndConv("WindDirection")
	rainfall := cloudGetter.getAndConv("Rainfall")
	battery := cloudGetter.getAndConv("Battery")
	thunder := cloudGetter.getAndConv("Thunder")
	light := cloudGetter.getAndConv("Light")
	charging := cloudGetter.getAndConv("Charging")

	if err := cloudGetter.error(); err != nil {
		return nil, err
	}

	return &weatherInfo{
		TempOUT:       temperatureOutside,
		Humidity:      humidity,
		TempIN:        float64(temperatureInside),
		Pressure:      float64(pressure),
		WindSpeed:     float64(windSpeed),
		WindDirection: windDirection,
		Rainfall:      rainfall,
		Battery:       battery,
		Thunder:       thunder,
		Light:         float64(light),
		Charging:      charging,
	}, nil
}

type cloudSender struct {
	err error
}

func newCloudSender() *cloudSender {
	return &cloudSender{}
}

func (cloudSender *cloudSender) send(key, value string) {
	if cloudSender.err != nil {
		return
	}
	cloudSender.err = sendValueToCloud(key, value)
}

func (cloudSender *cloudSender) error() error {
	return cloudSender.err
}

type cloudGetter struct {
	err error
}

func newCloudGetter() *cloudGetter {
	return &cloudGetter{}
}

func (cloudGetter *cloudGetter) get(key string) string {
	if cloudGetter.err != nil {
		return ""
	}
	value, err := getValueFromCloud(key)
	if err != nil {
		cloudGetter.err = err
		return ""
	}
	return value
}

func (cloudGetter *cloudGetter) getAndConv(key string) int {
	value := cloudGetter.get(key)
	value = value[1:len(value)-1]
	intValue, err := strconv.Atoi(value)
	if err != nil {
		cloudGetter.err = err
		return 0
	}
	return intValue
}

func (cloudGetter *cloudGetter) error() error {
	return cloudGetter.err
}

func sendValueToCloud(key, value string) error {
	url := fmt.Sprintf("https://keyvalue.immanuel.co/api/KeyVal/UpdateValue/%v/%v/%v",
		personalCloudKey, key, value)
	resp, err := http.Post(url, "", nil)
	if err != nil {
		return err
	}
	return resp.Body.Close()
}

func getValueFromCloud(key string) (string, error) {
	url := fmt.Sprintf("https://keyvalue.immanuel.co/api/KeyVal/GetValue/%v/%v",
		personalCloudKey, key)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	value, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(value), nil
}

package cbr_api

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"test-service-a/internal/entity"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

func NatsGetWithDelay(logger *zap.Logger, delay time.Duration, cbr_url string) {
	startTime := time.Now()
	for i := 1; ; i++ {

		// Получение структуры курса валют
		valCurs, err := getValCurse(cbr_url)
		if err != nil {
			logger.Error("Error NatsGetWithDelay getValCurse", zap.Error(err))
			time.Sleep(delay*time.Duration(i) - time.Since(startTime))
			continue
		}

		// Фильтрация валют
		var valute []entity.Valute
		for _, v := range valCurs.Valute {
			if v.CharCode == "EUR" || v.CharCode == "USD" {
				valute = append(valute, v)
			}
		}
		valute = append(valute, entity.Valute{
			NumCode:  643,
			CharCode: "RUB",
			Nominal:  1,
			Name:     "Российский рубль",
			Value:    1,
		})
		valCurs.Valute = valute

		// Отправка в шину
		jsonValCurs, _ := json.Marshal(valCurs)
		nc, _ := nats.Connect(viper.GetString("NATS_CONNECT"))
		nc.Publish("ValCurs", jsonValCurs)
		nc.Close()

		logger.Info("NatsGetWithDelay", zap.String(("message"), time.Now().Format("02 Jan 06 15:04:05.999")), zap.Duration("delay", delay))

		time.Sleep(delay*time.Duration(i) - time.Since(startTime))
	}
}

func NatsGetEveryday(logger *zap.Logger, requestTime time.Time, cbr_url string) {
	startTime := time.Now()
	requestTime = requestTime.AddDate(startTime.Year(), int(startTime.Month()-1), startTime.Day()-1)

	// Если в этот день уже поздно - переносим на следующий
	if startTime.After(requestTime) {
		requestTime = requestTime.Add(time.Second * 5)
	}
	time.Sleep(time.Until(requestTime))

	startTime = time.Now()
	var timeOut int // Длительность паузы при ошибке
	for {

		// Получение структуры курса валют
		valCurs, err := getValCurse(cbr_url)
		if err != nil {
			logger.Error("Error NatsGetEveryday getValCurse", zap.Error(err))
			time.Sleep(time.Second * time.Duration(timeOut))
			timeOut++
			continue
		}

		// Фильтрация структуры
		var valute []entity.Valute
		for _, v := range valCurs.Valute {
			if v.CharCode == "EUR" || v.CharCode == "USD" {
				valute = append(valute, v)
			}
		}
		valute = append(valute, entity.Valute{
			NumCode:  643,
			CharCode: "RUB",
			Nominal:  1,
			Name:     "Российский рубль",
			Value:    1,
		})
		valCurs.Valute = valute

		// Отправка в шину
		jsonValCurs, _ := json.Marshal(valCurs)
		nc, _ := nats.Connect(viper.GetString("NATS_CONNECT"))
		nc.Publish("ValCurs", jsonValCurs)
		nc.Close()

		logger.Info("NatsGetEveryday", zap.String("message", time.Now().Format("02 Jan 06 15:04:05.999")))

		time.Sleep(time.Minute*time.Duration(time.Since(startTime).Minutes()+1) - time.Since(startTime))
		timeOut = 0
	}
}

func getValCurse(cbr_url string) (entity.ValCurs, error) {
	
	// Запрос к CBR_URL
	resp, err := http.Get(cbr_url)
	if err != nil {
		return entity.ValCurs{}, err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return entity.ValCurs{}, fmt.Errorf("error resp.StatusCode %d", resp.StatusCode)
	}

	// Конвертация Windows-1251 в UTF-8
	tr := transform.NewReader(resp.Body, charmap.Windows1251.NewDecoder())
	xmlText, err := ioutil.ReadAll(tr)
	if err != nil {
		return entity.ValCurs{}, err
	}
	xmlText = []byte(strings.Replace(string(xmlText), "windows-1251", "UTF-8", 1))
	xmlText = []byte(strings.Replace(string(xmlText), `,`, `.`, -1))

	// Парсинг XML в Структуру
	valCurs := entity.ValCurs{}
	err = xml.Unmarshal(xmlText, &valCurs)
	if err != nil {
		return entity.ValCurs{}, err
	}

	return valCurs, nil
}

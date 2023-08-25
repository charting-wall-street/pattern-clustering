package db

import (
	"encoding/json"
	"github.com/northberg/candlestick"
	"io"
	"log"
	"net/http"
)

type listAlgorithmsPayload struct {
	Algorithms []AlgorithmDefinition `json:"algorithms"`
}

func Algorithms() []string {

	resp, err := http.Get("http://localhost:9706/algorithms")
	if err != nil {
		log.Fatal(err)
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	algos := new(listAlgorithmsPayload)
	err = json.Unmarshal(body, algos)
	if err != nil {
		log.Fatal(err)
	}

	allAlgos := make([]string, 0)
	for _, algo := range algos.Algorithms {
		allAlgos = append(allAlgos, algo.Name)
	}

	return allAlgos
}

type AlgorithmDefinition struct {
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Entry       string `json:"entry"`
}

func Symbols() []string {

	resp, err := http.Get("http://192.168.1.7:9702/market/info")
	if err != nil {
		log.Fatal(err)
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	exchangeList := new(candlestick.ExchangeList)
	err = json.Unmarshal(body, exchangeList)
	if err != nil {
		log.Fatal(err)
	}

	allSymbols := make([]string, 0)
	for _, exchange := range exchangeList.Exchanges {
		if exchange.BrokerId != "UNICORN" {
			continue
		}
		for s := range exchange.Symbols {
			allSymbols = append(allSymbols, s)
		}
	}

	return allSymbols
}

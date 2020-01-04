package main

import (
	"encoding/json"
	"fmt"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const CurrencyUrl = "http://localhost:5000"

type Currencies []string

type CurrencyOutput struct {
	From            string    `json:"from"`
	To              string    `json:"to"`
	ExchangerName   string    `json:"exchangerName"`
	ExchangeValue   float64   `json:"exchangeValue"`
	OriginalAmount  int       `json:"originalAmount"`
	ConvertedAmount float64   `json:"convertedAmount"`
	ConvertedText   string    `json:"convertedText"`
	RateDateTime    time.Time `json:"rateDateTime"`
	RateFromCache   bool      `json:"rateFromCache"`
}

type OrderBook struct {
	Bids [][]string `json:"bids"`
	Asks [][]string `json:"asks"`
}

func main() {
	currencyFile, err := ioutil.ReadFile("currencies.json")
	if err != nil {
		log.Printf("currencies err: %v ", err)
		os.Exit(1)
	}
	log.Println("currencies.json:", string(currencyFile))
	var currencies Currencies
	err = json.Unmarshal(currencyFile, &currencies)
	if err != nil {
		log.Printf("currency unmarshal err: %v ", err)
		os.Exit(1)
	}
	for _, v := range currencies {
		log.Println("processing:", v)
		orderBook := getOrderBook(v)
		log.Println(orderBook)
		makePlot(v, orderBook)
		time.Sleep(1)
	}
}

func makePlot(currency string, orderBook OrderBook) {
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = currency + " bids and asks"
	p.X.Label.Text = "Value In USD"
	p.Y.Label.Text = "Limits"

	ptsAsk := make(plotter.XYs, len(orderBook.Asks))
	for i := range ptsAsk {
		tmpX, _ := strconv.ParseFloat(orderBook.Asks[i][0], 64)
		tmpY, _ := strconv.ParseFloat(orderBook.Asks[i][1], 64)
		tmpX = convertCurrency(currency, int(tmpX))
		ptsAsk[i].X = tmpX
		ptsAsk[i].Y = tmpY
	}

	ptsBid := make(plotter.XYs, len(orderBook.Bids))
	for i := range ptsBid {
		tmpX, _ := strconv.ParseFloat(orderBook.Bids[i][0], 64)
		tmpY, _ := strconv.ParseFloat(orderBook.Bids[i][1], 64)
		tmpX = convertCurrency(currency, int(tmpX))
		ptsBid[i].X = tmpX
		ptsBid[i].Y = tmpY
	}

	err = plotutil.AddLinePoints(p,
		"ask", ptsAsk, "bid", ptsBid)
	if err != nil {
		panic(err)
	}

	if err := p.Save(4*vg.Inch, 4*vg.Inch, "images/"+currency+".png"); err != nil {
		panic(err)
	}
}

func convertCurrency(currency string, value int) float64 {
	var output CurrencyOutput
	client := &http.Client{
		Timeout: time.Second * 30,
	}
	url := fmt.Sprintf("%s/convert?from=%s&to=USD&amount=%d&exchanger=yahoo", CurrencyUrl, currency, value)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
		return 0
	}
	result, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return 0
	}
	defer result.Body.Close()
	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		log.Println(err)
		return 0
	}
	log.Println(string(body))
	err = json.Unmarshal(body, &output)
	if err != nil {
		log.Println(err)
	}
	return output.ConvertedAmount
}

func getOrderBook(currency string) OrderBook {
	var output OrderBook
	client := &http.Client{
		Timeout: time.Second * 30,
	}
	url := "https://localbitcoins.com/bitcoincharts/" + currency + "/orderbook.json"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
		return output
	}
	result, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return output
	}
	defer result.Body.Close()
	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		log.Println(err)
		return output
	}
	err = json.Unmarshal(body, &output)
	if err != nil {
		log.Println(err)
		return output
	}
	return output
}

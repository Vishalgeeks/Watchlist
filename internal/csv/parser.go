package csv

import (
	"encoding/csv"
	"net/http"
	"strconv"
	"strings"
	"watchlist-backend/pkg/models"
)

func ParseCSV(url string) ([]models.Stock, error) {
	// URL se CSV fetch karo
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	reader := csv.NewReader(resp.Body)
	reader.Comma = '\t'
	reader.LazyQuotes = true

	// Header skip karo
	_, err = reader.Read()
	if err != nil {
		return nil, err
	}

	var stocks []models.Stock
	for {
		row, err := reader.Read()
		if err != nil {
			break
		}

		if len(row) < 32 {
			continue
		}

		stock := models.Stock{
			ExchangeInstrumentID:  strings.TrimSpace(row[0]),
			Segment:               strings.TrimSpace(row[1]),
			InstrumentType:        strings.TrimSpace(row[2]),
			Symbol:                strings.TrimSpace(row[3]),
			DisplayName:           strings.TrimSpace(row[4]),
			CompanyName:           strings.TrimSpace(row[5]),
			ISIN:                  strings.TrimSpace(row[6]),
			Series:                strings.TrimSpace(row[7]),
			Exchange:              strings.TrimSpace(row[8]),
			ContractExpiration:    strings.TrimSpace(row[9]),
			OptionType:            strings.TrimSpace(row[11]),
			UnderlyingSymbolID:    strings.TrimSpace(row[12]),
			UnderlyingSymbol:      strings.TrimSpace(row[13]),
			Description:           strings.TrimSpace(row[19]),
			CautionaryMessageInfo: strings.TrimSpace(row[31]),
		}

		stock.Strike = parseFloat(row[10])
		stock.TickSize = parseFloat(row[15])
		stock.UpperCircuit = parseFloat(row[16])
		stock.LowerCircuit = parseFloat(row[17])
		stock.LTP = parseFloat(row[20])
		stock.Open = parseFloat(row[21])
		stock.High = parseFloat(row[22])
		stock.Low = parseFloat(row[23])
		stock.Close = parseFloat(row[24])
		stock.Bid = parseFloat(row[27])
		stock.Ask = parseFloat(row[28])
		stock.LotSize = parseInt(row[14])
		stock.FreezeQty = parseInt(row[18])
		stock.BidQty = parseInt(row[29])
		stock.AskQty = parseInt(row[30])
		stock.Vol = parseInt64(row[25])
		stock.OI = parseInt64(row[26])

		stocks = append(stocks, stock)
	}

	return stocks, nil
}

func parseFloat(s string) float64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	val, _ := strconv.ParseFloat(s, 64)
	return val
}

func parseInt(s string) int {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	val, _ := strconv.Atoi(s)
	return val
}

func parseInt64(s string) int64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	val, _ := strconv.ParseInt(s, 10, 64)
	return val
}

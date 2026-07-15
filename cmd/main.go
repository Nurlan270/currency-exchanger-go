package main

import (
	"context"
	"currency-exchanger/internal/http/handlers"
	"currency-exchanger/internal/services"
	"currency-exchanger/internal/stores"
	"database/sql"
	"fmt"
	_ "github.com/ncruces/go-sqlite3/driver"
	"log"
	"net/http"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	db, err := sql.Open("sqlite3", "./internal/db/sqlite.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		log.Fatal(err)
	}

	// Dependencies
	currencyHandler, exchangeHandler := registerHandlers(db)

	// HTTP Handlers
	http.HandleFunc("/currencies", currencyHandler.Currencies)
	http.HandleFunc("/currency/", currencyHandler.GetCurrency)
	http.HandleFunc("/exchangeRates", exchangeHandler.Exchanges)
	http.HandleFunc("/exchangeRate/", exchangeHandler.Exchange)
	http.HandleFunc("/exchange", exchangeHandler.CalcExchangeRate)

	fmt.Println("Listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func registerHandlers(db *sql.DB) (*handlers.CurrencyHandler, *handlers.ExchangeHandler) {
	// Currency
	currencyStore := stores.NewCurrencyStore(db)
	currencyService := services.NewCurrencyService(currencyStore)
	currencyHandler := handlers.NewCurrencyHandler(currencyService)

	// ExchangeRates
	exchangeStore := stores.NewExchangeStore(db)
	exchangeService := services.NewExchangeService(exchangeStore)
	exchangeHandler := handlers.NewExchangeHandler(exchangeService)

	return currencyHandler, exchangeHandler
}

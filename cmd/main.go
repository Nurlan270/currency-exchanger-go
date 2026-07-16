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

	mux := http.NewServeMux()

	// API handlers
	mux.HandleFunc("/currencies", currencyHandler.Currencies)
	mux.HandleFunc("/currency/", currencyHandler.GetCurrency)
	mux.HandleFunc("/exchangeRates", exchangeHandler.Exchanges)
	mux.HandleFunc("/exchangeRate/", exchangeHandler.Exchange)
	mux.HandleFunc("/exchange", exchangeHandler.CalcExchangeRate)

	// Static assets
	mux.HandleFunc("/", serveIndex)
	mux.HandleFunc("/index.html", serveIndex)
	mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./css"))))
	mux.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./js"))))

	fmt.Println("Listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", mux))
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

func serveIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" && r.URL.Path != "/index.html" {
		http.NotFound(w, r)
		return
	}

	http.ServeFile(w, r, "./index.html")
}

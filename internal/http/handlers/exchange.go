package handlers

import (
	"currency-exchanger/internal/db"
	"currency-exchanger/internal/helpers"
	"currency-exchanger/internal/http/responses"
	"currency-exchanger/internal/services"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

type ExchangeHandler struct {
	service *services.ExchangeService
}

func NewExchangeHandler(service *services.ExchangeService) *ExchangeHandler {
	return &ExchangeHandler{service: service}
}

func (h *ExchangeHandler) Exchanges(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.getExchanges(w, r)
	} else if r.Method == http.MethodPost {
		h.createNewExchange(w, r)
	}
}

func (h *ExchangeHandler) Exchange(w http.ResponseWriter, r *http.Request) {
	code := strings.ToUpper(strings.TrimPrefix(r.URL.Path, "/exchangeRate/"))
	if code == "" {
		helpers.SendError(w, "Коды валют пары отсутствуют в адресе", http.StatusBadRequest)
		return
	}

	if len(code) != 6 {
		helpers.SendError(w, "Коды валют пары должны содержать 6 символов, пример: USDRUB", http.StatusBadRequest)
		return
	}

	baseCurrencyCode, targetCurrencyCode := code[:3], code[3:]

	if r.Method == http.MethodGet {
		h.getExchange(w, r, baseCurrencyCode, targetCurrencyCode)
	} else if r.Method == http.MethodPatch {
		h.updateExchange(w, r, baseCurrencyCode, targetCurrencyCode)
	}
}

func (h *ExchangeHandler) getExchange(
	w http.ResponseWriter, r *http.Request,
	baseCurrencyCode string, targetCurrencyCode string,
) {
	exchange, err := h.service.GetByCodes(baseCurrencyCode, targetCurrencyCode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			helpers.SendError(w, "Обменный курс для пары не найден", http.StatusNotFound)
		} else {
			helpers.SendError(w, "Что-то пошло не так", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(exchange)
}

func (h *ExchangeHandler) getExchanges(w http.ResponseWriter, r *http.Request) {
	exchanges, err := h.service.GetAll()
	if err != nil {
		helpers.SendError(w, "Что-то пошло не так", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(exchanges)
}

func (h *ExchangeHandler) createNewExchange(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		helpers.SendError(w, "Что-то пошло не так", http.StatusInternalServerError)
		return
	}

	baseCurrencyCode := strings.ToUpper(r.PostForm.Get("baseCurrencyCode"))
	targetCurrencyCode := strings.ToUpper(r.PostForm.Get("targetCurrencyCode"))
	rate := r.PostForm.Get("rate")
	if baseCurrencyCode == "" || targetCurrencyCode == "" || rate == "" {
		helpers.SendError(w, "Отсутствует нужное поле формы", http.StatusBadRequest)
		return
	}

	if len(targetCurrencyCode) != 3 || len(baseCurrencyCode) != 3 {
		helpers.SendError(w, "Коды валют должен содержать 3 символа", http.StatusBadRequest)
		return
	}

	f64Rate, err := strconv.ParseFloat(rate, 64)
	if err != nil {
		helpers.SendError(w, "Некорректное значение курса", http.StatusBadRequest)
		return
	}

	newExchange, err := h.service.Create(baseCurrencyCode, targetCurrencyCode, f64Rate)
	if err != nil {
		if errors.Is(err, db.ErrRowAlreadyExists) {
			helpers.SendError(w, "Валютная пара с таким кодом уже существует", http.StatusConflict)
		} else if errors.Is(err, db.ErrNotFound) {
			helpers.SendError(w, "Одна (или обе) валюта из валютной пары не существует в БД", http.StatusNotFound)
		} else {
			helpers.SendError(w, "Что-то пошло не так", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newExchange)
}

func (h *ExchangeHandler) updateExchange(
	w http.ResponseWriter, r *http.Request,
	baseCurrencyCode string, targetCurrencyCode string,
) {
	if err := r.ParseForm(); err != nil {
		helpers.SendError(w, "Что-то пошло не так", http.StatusInternalServerError)
		return
	}

	rate := r.PostForm.Get("rate")
	if rate == "" {
		helpers.SendError(w, "Отсутствует нужное поле формы", http.StatusBadRequest)
		return
	}

	f64Rate, err := strconv.ParseFloat(rate, 64)
	if err != nil {
		helpers.SendError(w, "Некорректное значение курса", http.StatusBadRequest)
		return
	}

	exchange, err := h.service.Update(baseCurrencyCode, targetCurrencyCode, f64Rate)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			helpers.SendError(w, "Валютная пара отсутствует в базе данных", http.StatusNotFound)
		} else {
			helpers.SendError(w, "Что-то пошло не так", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(exchange)
}

func (h *ExchangeHandler) CalcExchangeRate(w http.ResponseWriter, r *http.Request) {
	baseCode := strings.ToUpper(r.URL.Query().Get("from"))
	targetCode := strings.ToUpper(r.URL.Query().Get("to"))
	amountString := r.URL.Query().Get("amount")

	if baseCode == "" || targetCode == "" || amountString == "" {
		helpers.SendError(w, "Отсутствует нужное поле", http.StatusBadRequest)
		return
	}

	// Scenario 1 and 2
	exchange, err := h.service.GetByCodesWithRevert(baseCode, targetCode)

	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			helpers.SendError(w, "Что-то пошло не так", http.StatusInternalServerError)
			return
		}

		// Scenario 3
		exchange, err = h.service.GetByForeignRate(baseCode, targetCode)
		if err != nil {
			helpers.SendError(w, "Обменный курс для пары не найден", http.StatusNotFound)
			return
		}
	}

	amount, err := strconv.ParseFloat(amountString, 64)
	if err != nil {
		helpers.SendError(w, "Некорректное значение суммы", http.StatusBadRequest)
		return
	}

	convertedAmount := h.service.ConvertAmount(amount, exchange.Rate)

	res := responses.ExchangeResponse{
		BaseCurrency:    exchange.BaseCurrency,
		TargetCurrency:  exchange.TargetCurrency,
		Rate:            exchange.Rate,
		Amount:          amount,
		ConvertedAmount: convertedAmount,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

package handlers

import (
	"currency-exchanger/internal/db"
	"currency-exchanger/internal/helpers"
	"currency-exchanger/internal/services"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type CurrencyHandler struct {
	service *services.CurrencyService
}

func NewCurrencyHandler(service *services.CurrencyService) *CurrencyHandler {
	return &CurrencyHandler{service: service}
}

func (h *CurrencyHandler) Currencies(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.getCurrencies(w, r)
	} else if r.Method == http.MethodPost {
		h.createNewCurrency(w, r)
	}
}

func (h *CurrencyHandler) GetCurrency(w http.ResponseWriter, r *http.Request) {
	code := strings.ToUpper(strings.TrimPrefix(r.URL.Path, "/currency/"))
	if code == "" {
		helpers.SendError(w, "Код валюты отсутствует в адресе", http.StatusBadRequest)
		return
	}

	currency, err := h.service.GetByCode(code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			helpers.SendError(w, "Валюта не найдена", http.StatusNotFound)
		} else {
			helpers.SendError(w, "Что-то пошло не так", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(currency)
}

func (h *CurrencyHandler) getCurrencies(w http.ResponseWriter, r *http.Request) {
	currencies, err := h.service.GetAll()
	if err != nil {
		helpers.SendError(w, "Что-то пошло не так", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(currencies)
}

func (h *CurrencyHandler) createNewCurrency(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		helpers.SendError(w, "Что-то пошло не так", http.StatusInternalServerError)
		return
	}

	name, code, sign := r.PostForm.Get("name"), r.PostForm.Get("code"), r.PostForm.Get("sign")
	if err := helpers.ValidateEmpty(name, code, sign); err != nil {
		helpers.SendError(w, "Отсутствует нужное поле формы", http.StatusBadRequest)
		return
	}

	if len(code) != 3 {
		helpers.SendError(w, "Код валюты должен содержать 3 символа", http.StatusBadRequest)
		return
	}

	if len(sign) > 3 {
		helpers.SendError(w, "Символ валюты должен содержать не больше 3 символов", http.StatusBadRequest)
		return
	}

	newCurrency, err := h.service.Create(name, code, sign)
	if err != nil {
		if errors.Is(err, db.ErrRowAlreadyExists) {
			helpers.SendError(w, "Валюта с таким кодом уже существует", http.StatusConflict)
		} else {
			helpers.SendError(w, "Что-то пошло не так", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newCurrency)
}

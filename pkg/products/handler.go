package products

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/keregocen/go-product-crud-api/pkg/mockexchange"
	"github.com/keregocen/go-product-crud-api/pkg/models"
	"github.com/keregocen/go-product-crud-api/pkg/storage"
)

type Handler struct {
	StorageAPI           *storage.Store
	CurrencyConverterAPI mockexchange.CurrencyConverter
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		h.renderError(w, "failed to decode request body", http.StatusBadRequest)
		return
	}

	// check if product exists
	if h.StorageAPI.Exist(product.Name) {
		h.renderError(w, fmt.Sprintf("product %v already exists", product.Name), http.StatusBadRequest)
		return
	}

	if err := h.StorageAPI.Save(product.Name, product); err != nil {
		h.renderError(w, fmt.Sprintf("failed to create product %v", product.Name), http.StatusInternalServerError)
		return
	}

	responseErr := h.response(w, http.StatusCreated, product)
	if responseErr != nil {
		fmt.Printf("failed to write http response %v", responseErr)
	}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	products, err := h.StorageAPI.LoadAll()
	if err != nil {
		h.renderError(w, "failed to list all products", http.StatusInternalServerError)
		return
	}

	responseErr := h.response(w, http.StatusOK, products)
	if responseErr != nil {
		fmt.Printf("failed to write http response %v", responseErr)
	}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	var getProductRequest GetProductRequest
	if err := json.NewDecoder(r.Body).Decode(&getProductRequest); err != nil {
		h.renderError(w, "failed to decode request body", http.StatusBadRequest)
		return
	}

	if !h.StorageAPI.Exist(getProductRequest.Name) {
		h.renderError(w, fmt.Sprintf("missing product with name %v", getProductRequest.Name), http.StatusNotFound)
		return
	}

	product, ok := h.StorageAPI.Load(getProductRequest.Name)
	if !ok {
		h.renderError(w, fmt.Sprintf("failed to get product %v", getProductRequest.Name), http.StatusBadRequest)
		return
	}

	p, ok := product.(models.Product)
	if !ok {
		h.renderError(w, fmt.Sprintf("unexpected object returned for %v", getProductRequest.Name),
			http.StatusInternalServerError)
		return
	}

	// currency will be converted if the request has a Currency field and it's not the default GBP
	if len(getProductRequest.Currency) > 0 && getProductRequest.Currency != "GBP" {
		rate, err := h.CurrencyConverterAPI.ConvertExchangeRate("GBP", getProductRequest.Currency)
		if err != nil {
			h.renderError(w, fmt.Sprintf("failed to get exchange rate information for %v",
				getProductRequest.Currency), http.StatusNotFound)
			return
		}
		p.Price *= rate
	}

	responseErr := h.response(w, http.StatusOK, p)
	if responseErr != nil {
		fmt.Printf("failed to write http response %v", responseErr)
	}
}

func (h *Handler) response(w http.ResponseWriter, statusCode int, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal response JSON: %w", err)
	}

	w.WriteHeader(statusCode)

	if _, writeErr := w.Write(jsonData); writeErr != nil {
		return fmt.Errorf("failed to write response JSON: %w", writeErr)
	}

	return nil
}

func (h *Handler) renderError(writer http.ResponseWriter, errorMessage string, statusCode int) {
	// todo use logger
	errResponse := &ErrorResponse{
		Message: errorMessage,
	}
	errJSON, err := json.Marshal(errResponse)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Error(writer, string(errJSON), statusCode)
}

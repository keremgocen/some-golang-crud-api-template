package products

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/keregocen/go-product-crud-api/pkg/mockexchange"
	"github.com/keregocen/go-product-crud-api/pkg/models"
	"github.com/keregocen/go-product-crud-api/pkg/storage"
)

type Service struct {
	StorageAPI           *storage.Store
	CurrencyConverterAPI mockexchange.CurrencyConverter
}

func NewService(storage *storage.Store, currencyConverter mockexchange.CurrencyConverter) *Service {
	return &Service{
		StorageAPI:           storage,
		CurrencyConverterAPI: currencyConverter,
	}
}

func (s *Service) Create(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		s.renderError(w, "failed to decode request body", http.StatusBadRequest)
		return
	}

	// check if product exists
	if s.StorageAPI.Exist(product.Name) {
		s.renderError(w, fmt.Sprintf("product %v already exists", product.Name), http.StatusBadRequest)
		return
	}

	if err := s.StorageAPI.Save(product.Name, product); err != nil {
		s.renderError(w, fmt.Sprintf("failed to create product %v", product.Name), http.StatusInternalServerError)
		return
	}

	responseErr := s.response(w, http.StatusCreated, product)
	if responseErr != nil {
		log.Printf("failed to write http response %v", responseErr)
	}
}

func (s *Service) List(w http.ResponseWriter, r *http.Request) {
	products, err := s.StorageAPI.LoadAll()
	if err != nil {
		s.renderError(w, "failed to list all products", http.StatusInternalServerError)
		return
	}

	responseErr := s.response(w, http.StatusOK, products)
	if responseErr != nil {
		log.Printf("failed to write http response %v", responseErr)
	}
}

func (s *Service) Get(w http.ResponseWriter, r *http.Request) {
	var getProductRequest GetProductRequest
	if err := json.NewDecoder(r.Body).Decode(&getProductRequest); err != nil {
		s.renderError(w, "failed to decode request body", http.StatusBadRequest)
		return
	}

	if !s.StorageAPI.Exist(getProductRequest.Name) {
		s.renderError(w, fmt.Sprintf("missing product with name %v", getProductRequest.Name), http.StatusNotFound)
		return
	}

	product, ok := s.StorageAPI.Load(getProductRequest.Name)
	if !ok {
		s.renderError(w, fmt.Sprintf("failed to get product %v", getProductRequest.Name), http.StatusBadRequest)
		return
	}

	p, ok := product.(models.Product)
	if !ok {
		s.renderError(w, fmt.Sprintf("unexpected object returned for %v", getProductRequest.Name),
			http.StatusInternalServerError)
		return
	}

	// currency will be converted if the request has a Currency field and it's not the default GBP
	if len(getProductRequest.Currency) > 0 && getProductRequest.Currency != "GBP" {
		rate, err := s.CurrencyConverterAPI.ConvertExchangeRate("GBP", getProductRequest.Currency)
		if err != nil {
			s.renderError(w, fmt.Sprintf("failed to get exchange rate information for %v",
				getProductRequest.Currency), http.StatusNotFound)
			return
		}
		p.Price *= rate
	}

	responseErr := s.response(w, http.StatusOK, p)
	if responseErr != nil {
		log.Printf("failed to write http response %v", responseErr)
	}
}

func (s *Service) response(w http.ResponseWriter, statusCode int, data interface{}) error {
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

func (s *Service) renderError(writer http.ResponseWriter, errorMessage string, statusCode int) {
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

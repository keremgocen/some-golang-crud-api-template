package products

type ErrorResponse struct {
	Message string `json:"message"`
}

type GetProductRequest struct {
	Name     string `json:"name" validate:"required"`
	Currency string `json:"currency"`
}

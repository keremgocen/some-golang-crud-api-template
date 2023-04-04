// Keyvalue package provides the service layer.
package keyvalue

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Service is the keyvalue service responsible for the business logic, implementing the HTTP handlers.
// The storage interface is implemented by the storage package.
type Service struct {
	StorageAPI Storage
}

// NewService returns a new keyvalue service. Accepts a storage interface.
func NewService(storage Storage) *Service {
	return &Service{
		StorageAPI: storage,
	}
}

// GetKeyValue returns a keyvalue pair from the storage.
func (s *Service) GetKeyValue(c *gin.Context) {
	var req EntryResponse
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "field validation failed" + err.Error()})
		return
	}

	key := req.Name

	if cachedValue, ok := s.StorageAPI.Load(key); ok {
		c.JSON(http.StatusOK, cachedValue)
		return
	}

	c.AbortWithStatusJSON(http.StatusNotFound, fmt.Errorf("key %s not found", key))
}

// PostKeyValue saves a keyvalue pair to the storage.
func (s *Service) PostKeyValue(c *gin.Context) {
	var req EntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "field validation failed" + err.Error()})
		return
	}

	// ids –:= c.QueryMap("ids")
	// names := c.PostFormMap("names")

	// todo input validation error handling

	if err := s.StorageAPI.Save(req.Name, req.Value); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

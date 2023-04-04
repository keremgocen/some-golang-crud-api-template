package keyvalue_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/keregocen/go-product-crud-api/pkg/keyvalue"
	"github.com/keregocen/go-product-crud-api/pkg/storage"
	"github.com/stretchr/testify/assert"
)

// TestGetKeyvalue tests the Get method of the keyvalue service.
func TestGetKeyvalue(t *testing.T) {
	tests := map[string]struct {
		name        string
		wantStatus  int
		wantErr     error
		getEntryErr error
	}{
		"get keyvalue returns 200 with expected value": {
			name:        "foo",
			wantStatus:  http.StatusOK,
			wantErr:     nil,
			getEntryErr: nil,
		},
		// "get keyvalue returns 404 when when the key is missing": {
		// 	name:        "missing",
		// 	wantStatus:  http.StatusNotFound,
		// 	wantErr:     nil,
		// 	getEntryErr: gin.Error{Err: fmt.Errorf("not found")},
		// },
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			storageAPI := storage.NewStore()
			err := storageAPI.Save(tc.name, "bar")
			assert.Nil(t, err)

			// ctrl := gomock.NewController(t)
			// mockStorageAPI := mock.NewMockStorage(ctrl)
			// mockStorageAPI.EXPECT().Load(tc.name).Return("bar")

			service := keyvalue.NewService(storageAPI)

			router := gin.Default()
			router.GET("/keyvalue/:name", service.GetKeyValue)

			req, err := http.NewRequest(http.MethodGet, "/keyvalue/"+tc.name, nil)
			assert.Equal(t, err, tc.wantErr)

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)
			assert.Equal(t, tc.wantStatus, rr.Code)

			var entry keyvalue.EntryResponse
			err = json.Unmarshal(rr.Body.Bytes(), &entry)
			assert.Nil(t, err)
		})
	}
}

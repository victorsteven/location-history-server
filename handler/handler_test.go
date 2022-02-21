package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var h = NewService()

func TestService_Create(t *testing.T) {
	tests := []struct {
		name       string
		inputJSON  string
		statusCode int
		errMessage string
	}{
		{
			name:       "Creation success",
			inputJSON:  `{"lat": 12.4, "lng": 41.3}`,
			statusCode: 200,
			errMessage: "",
		},
		{
			name:       "Creation failure",
			inputJSON:  `{"lat": "12.4", "lng": "41.3"}`,
			statusCode: 400,
			errMessage: `{"error":"json: cannot unmarshal string into Go struct field Location.lat of type float64"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			orderID := "abc123"
			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/location/%s/now", orderID), bytes.NewBufferString(test.inputJSON))
			assert.NoError(t, err)

			r := gin.Default()
			r.POST("/location/:order_id/now", h.Create)
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			assert.Equal(t, rr.Code, test.statusCode)
			if rr.Code == 200 {
				var response []Payload
				err = json.Unmarshal(rr.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.EqualValues(t, 200, rr.Code)
				assert.EqualValues(t, 1, len(response))
			} else {
				assert.EqualValues(t, test.errMessage, rr.Body.String())
			}
		})
	}
}

func TestService_Get(t *testing.T) {
	seedStorage()

	tests := []struct {
		name       string
		max        string
		statusCode int
		errMessage string
	}{
		{
			name:       "Success",
			max:        "1",
			statusCode: 200,
			errMessage: "",
		},
		{
			name:       "Failure - invalid max value",
			max:        "invalid", // a positive number is expected
			statusCode: 400,
			errMessage: `{"error":"kindly provide a valid positive number"}`,
		},
		{
			name:       "Failure - greater max value provided",
			max:        "2", // only one location is available
			statusCode: 400,
			errMessage: `{"error":"the max input cannot be greater than the available locations"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			orderID := "abc123"
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/location/%s?max=%s", orderID, test.max), nil)
			assert.NoError(t, err)

			r := gin.Default()
			r.GET("/location/:order_id", h.Get)
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			assert.Equal(t, rr.Code, test.statusCode)
			if rr.Code == 200 {
				var response Payload
				err = json.Unmarshal(rr.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.EqualValues(t, 200, rr.Code)
				assert.EqualValues(t, orderID, response.OrderID)
				assert.EqualValues(t, 1, len(response.History))
			} else {
				assert.EqualValues(t, test.errMessage, rr.Body.String())
			}
		})
	}
}

func TestService_Delete(t *testing.T) {
	seedStorage()

	orderID := "abc123"
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/location/%s", orderID), nil)
	assert.NoError(t, err)

	r := gin.Default()
	r.DELETE("/location/:order_id", h.Delete)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, 200)
	assert.Equal(t, len(storage), 0)
}

// seedStorage is used to seed test data
func seedStorage() {
	storage = []Payload{
		{
			OrderID: "abc123",
			History: []Location{
				{
					Latitude:  12.4,
					Longitude: 13.6,
				},
			},
		},
	}
}

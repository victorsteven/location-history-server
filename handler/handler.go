package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Payload struct {
	OrderID string     `json:"order_id"`
	History []Location `json:"history"`
}

type Location struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
}

var (
	storage         []Payload
	orderPresent    = make(map[string]bool)
	locationPresent = make(map[Location]bool)
)

type service struct{}

func NewService() *service {
	return &service{}
}

func (s *service) Create(c *gin.Context) {
	orderID := c.Param("order_id")

	var location Location

	if err := c.ShouldBindJSON(&location); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// first entry in the storage
	if len(storage) == 0 {
		s.create(orderID, location)
		orderPresent[orderID] = true
		locationPresent[location] = true
		c.JSON(http.StatusOK, storage)
		return
	}

	for i, store := range storage {
		// create new entry for a new order
		if _, present := orderPresent[orderID]; !present {
			orderPresent[orderID] = true
			locationPresent[location] = true
			s.create(orderID, location)
		} else {
			if store.OrderID == orderID {
				// avoid duplicate location insertion
				if _, ok := locationPresent[location]; !ok {
					locationPresent[location] = true
					storage[i].History = append(storage[i].History, location)
				}
			}
		}
	}
	c.JSON(http.StatusOK, storage)
}

func (s *service) Get(c *gin.Context) {
	var (
		max int
		err error
	)
	orderID := c.Param("order_id")
	maxQuery := c.Query("max")
	if maxQuery != "" {
		max, err = strconv.Atoi(maxQuery)
		if err != nil || max <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "kindly provide a valid positive number",
			})
			return
		}
	}
	for _, store := range storage {
		if store.OrderID == orderID {
			if max > len(store.History) {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "the max input cannot be greater than the available locations",
				})
				return
			}
			var history []Location
			if max > 0 {
				if max == 1 {
					history = []Location{store.History[len(store.History)-1]}
				} else {
					history = store.History[len(store.History)-max:]
				}
			} else {
				history = store.History[:]
			}
			payload := Payload{
				OrderID: orderID,
				History: history,
			}
			c.JSON(http.StatusOK, payload)
			return
		}
	}
	c.JSON(http.StatusOK, nil)
}

func (s *service) Delete(c *gin.Context) {
	orderID := c.Param("order_id")
	for i, store := range storage {
		if store.OrderID == orderID {
			storage = append(storage[:i], storage[i+1:]...)
		}
	}
	c.JSON(http.StatusOK, "OK")
}

func (s *service) create(orderID string, location Location) {
	payload := Payload{
		OrderID: orderID,
		History: []Location{location},
	}
	storage = append(storage, payload)
}

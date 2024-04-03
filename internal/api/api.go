package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rnkjnk/decks-api/internal/services"
)

type handlers struct {
	service services.DecksServicer
}

func NewHandlers(service services.DecksServicer) *handlers {
	return &handlers{
		service: service,
	}
}

func (h *handlers) SetupRoutes(router *gin.Engine) {
	// Define routes and attach handler functions
	router.POST("/deck", h.createDeck)
	router.GET("/deck/:id/open", h.openDeck)
	router.POST("/deck/:id/draw-cards", h.drawCards)
}

// Creates a deck
func (h *handlers) createDeck(c *gin.Context) {
	shuffle := stringToBoolDefault(c.DefaultQuery("shuffle", "true"), true)

	cards, err := stringToCardCodeSlice(c.Query("cards"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	deck, err := h.service.CreateDeck(shuffle, cards)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, deck)
}

// Opens a deck
func (h *handlers) openDeck(c *gin.Context) {
	id := c.Param("id")

	deck, err := h.service.OpenDeck(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, deck)
}

// Draws cards
func (h *handlers) drawCards(c *gin.Context) {
	id := c.Param("id")
	draw, err := strconv.ParseUint(c.Query("draw"), 10, 8)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	cards, err := h.service.DrawCards(id, uint8(draw))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, cards)
}

func stringToBoolDefault(s string, def bool) bool {
	var r bool
	s = strings.ToLower(s)
	if s == "false" {
		r = false
	} else if s == "true" {
		r = true
	} else {
		r = def
	}
	return r
}

func stringToCardCodeSlice(cards string) ([][2]rune, error) {
	if cards == "" {
		return make([][2]rune, 0), nil
	}
	strings := strings.Split(cards, ",")
	output := make([][2]rune, len(strings))
	for i, s := range strings {
		if len(s) != 2 {
			return nil, fmt.Errorf("invalid card code: %s", s)
		}
		runes := []rune(s)
		output[i][0] = runes[0]
		output[i][1] = runes[1]
	}
	return output, nil
}

package handlers

import (
	"github.com/gofiber/fiber/v2"
	"strconv"
	"techical/internal/domain"
	"techical/internal/service"
)

type Handler struct {
	service *service.CurrencyService
}

func NewHandler(service *service.CurrencyService) *Handler {
	return &Handler{
		service: service,
	}
}

// GetRate godoc
// @Summary      Get Cryptocurrency Rate
// @Description  Get Cryptocurrency Rate
// @Tags         cryptocurrency
// @Accept       json
// @Produce      json
// @Param        from   path      string  true    "From currency"
// @Param        to     path      string  true    "To currency"
// @Param        from   path      number  true   "Amount to convert"
// @Success      200  {object}  domain.CurrencyResponse
// @Failure      400  {object}  domain.ErrorResponse
// @Router       /rates [get]
func (h *Handler) GetRate(c *fiber.Ctx) error {
	from := c.Query("from")
	to := c.Query("to")
	amountStr := c.Query("amount")

	if from == "" || to == "" || amountStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "missing required query parameters 'from', 'to' or 'amount'",
		})
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid amount parameter",
		})
	}

	result, err := h.service.Convert(from, to, amount)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(
		domain.CurrencyResponse{Result: result})
}

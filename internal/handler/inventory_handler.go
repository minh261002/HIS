package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/minhtran/his/internal/dto"
	"github.com/minhtran/his/internal/pkg/response"
	"github.com/minhtran/his/internal/service"
)

// InventoryHandler handles inventory HTTP requests
type InventoryHandler struct {
	inventoryService *service.InventoryService
}

// NewInventoryHandler creates a new inventory handler
func NewInventoryHandler(inventoryService *service.InventoryService) *InventoryHandler {
	return &InventoryHandler{inventoryService: inventoryService}
}

// AddStock handles adding inventory
func (h *InventoryHandler) AddStock(c *gin.Context) {
	var req dto.CreateInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, map[string]interface{}{"error": err.Error()})
		return
	}

	inventory, err := h.inventoryService.AddStock(&req)
	if err != nil {
		if err == service.ErrExpiredStock {
			response.BadRequest(c, "Expiry date must be in the future", nil)
			return
		}
		response.InternalServerError(c, "Failed to add stock")
		return
	}

	response.Created(c, "Stock added successfully", inventory)
}

// GetMedicationStock handles getting medication stock
func (h *InventoryHandler) GetMedicationStock(c *gin.Context) {
	medicationID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid medication ID", nil)
		return
	}

	stock, err := h.inventoryService.GetMedicationStock(uint(medicationID))
	if err != nil {
		response.InternalServerError(c, "Failed to get medication stock")
		return
	}

	response.Success(c, "Medication stock retrieved successfully", stock)
}

// GetLowStockAlerts handles getting low stock alerts
func (h *InventoryHandler) GetLowStockAlerts(c *gin.Context) {
	threshold, _ := strconv.Atoi(c.DefaultQuery("threshold", "10"))

	alerts, err := h.inventoryService.GetLowStockAlerts(threshold)
	if err != nil {
		response.InternalServerError(c, "Failed to get low stock alerts")
		return
	}

	response.Success(c, "Low stock alerts retrieved successfully", alerts)
}

// GetExpiringSoonAlerts handles getting expiring soon alerts
func (h *InventoryHandler) GetExpiringSoonAlerts(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))

	alerts, err := h.inventoryService.GetExpiringSoonAlerts(days)
	if err != nil {
		response.InternalServerError(c, "Failed to get expiring soon alerts")
		return
	}

	response.Success(c, "Expiring soon alerts retrieved successfully", alerts)
}

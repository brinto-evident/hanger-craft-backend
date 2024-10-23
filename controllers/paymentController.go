package controllers

import (
	"backend/config"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreatePayment handles creating a new payment record
func CreatePayment(c *gin.Context) {
	var payment *models.Payment

	// Bind the JSON request to the Payment struct
	if err := c.BindJSON(&payment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set the payment date to current time if it's not provided
	// if payment.PaymentDate.IsZero() {
	// 	payment.PaymentDate = time.Now()
	// }

	// Create the payment in the database
	if err := config.DB.Create(&payment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment"})
		return
	}

	// Return the created payment
	c.JSON(http.StatusOK, gin.H{"payment": payment})
}

// GetPaymentsByOrder retrieves all payments for a specific order
func GetPaymentsByOrder(c *gin.Context) {
	orderID := c.Param("order_id")
	var payments []*models.Payment

	// Find payments by the associated order ID
	if err := config.DB.Where("order_id = ?", orderID).Find(&payments).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "No payments found for this order"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Return the list of payments
	c.JSON(http.StatusOK, payments)
}

// UpdatePaymentStatus updates the payment status for a specific payment
func UpdatePaymentStatus(c *gin.Context) {
	paymentID := c.Param("id")
	var payment *models.Payment

	// Find the payment by ID
	if err := config.DB.First(&payment, paymentID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Bind the JSON request to the Payment struct (for status update)
	if err := c.BindJSON(&payment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the payment status
	if err := config.DB.Model(&payment).Update("payment_status", payment.PaymentStatus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update payment status"})
		return
	}

	// Return the updated payment
	c.JSON(http.StatusOK, gin.H{"payment": payment})
}
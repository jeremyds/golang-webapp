package api

import (
	"errors"
	"fmt"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"test_jump/database"
)

type AddInvoiceRequest struct {
	UserID uint    `json:"user_id" binding:"required"`
	Amount float64 `json:"amount" binding:"required"`
	Label  string  `json:"label" binding:"required"`
}

type TransactionRequest struct {
	InvoiceID uint    `json:"invoice_id" binding:"required"`
	Amount    float64 `json:"amount" binding:"required"`
	Reference string  `json:"reference,omitempty"`
}

func roundFloat2Dec(value float64) float64 {
	return math.Round(value*100) / 100
}

func AddInvoice(ctx *gin.Context) {
	var data AddInvoiceRequest
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user database.User
	if err := database.DB.First(&user, data.UserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "user does not exist"})
		} else {
			fmt.Printf("Failed to look for user: %v", err)
			ctx.Status(http.StatusInternalServerError)
		}
		return
	}

	if err := database.DB.Create(
		&database.Invoice{UserID: data.UserID, Label: data.Label, Amount: roundFloat2Dec(data.Amount)},
	).Error; err != nil {
		fmt.Printf("Failed to create invoice: %v", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func CreateTransaction(ctx *gin.Context) {
	var data TransactionRequest
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var invoice database.Invoice
	if err := database.DB.Joins("User").First(&invoice, data.InvoiceID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "invoice does not exist"})
		} else {
			fmt.Printf("Failed to look for invoice: %v", err)
			ctx.Status(http.StatusInternalServerError)
		}
		return
	}
	if invoice.Amount != roundFloat2Dec(data.Amount) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "unexpected amount"})
		return
	}
	if invoice.Paid {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invoice already paid"})
		return
	}

	invoice.Paid = true
	tx := database.DB.Begin()

	ret := tx.Model(&invoice).Update("Paid", true)
	if ret.Error != nil {
		tx.Rollback()
		fmt.Printf("Failed to look for invoice: %v", ret.Error)
		ctx.Status(http.StatusInternalServerError)
		return
	}
	if ret.RowsAffected == 0 {
		// The invoice status has been set to Paid since we fetched it
		tx.Rollback()
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invoice already paid"})
		return
	}
	if err := tx.Model(&invoice.User).Update("Balance", gorm.Expr("Balance + ?", invoice.Amount)).Error; err != nil {
		tx.Rollback()
		fmt.Printf("Failed to update user balance: %v", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	tx.Commit()

	ctx.Status(http.StatusNoContent)
}

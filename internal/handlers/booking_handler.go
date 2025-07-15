package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/Itaro18/bt_be/internal/db"
	"github.com/Itaro18/bt_be/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetBookings(w http.ResponseWriter, r *http.Request) {
	var bookings []models.Booking

	today := time.Now().Truncate(24 * time.Hour) // Remove time part

	if err := db.DB.
		Preload("Customer").
		Preload("Property").
		Where("check_in_date >= ? OR check_out_date >= ?", today, today).
		Find(&bookings).Error; err != nil {
		http.Error(w, "Failed to fetch bookings", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(bookings); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func CreateBooking(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name         string    `json:"name"`
		PhoneNo      string    `json:"phNo"`
		CheckInDate  time.Time `json:"checkIn"`
		CheckOutDate time.Time `json:"checkOut"`
		PropertyName string    `json:"property"`
		City         string    `json:"city"`
		Handler      string    `json:"handler"`
		Through      string    `json:"through"`
		AdvancePaid  float64   `json:"advance"`
		TotalAmount  float64   `json:"total"`
		Status       string    `json:"status"`
		Floor        string    `json:"floor"`
		Remarks      string    `json:"remarks"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Start Transaction
	tx := db.DB.Begin()
	if tx.Error != nil {
		http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
		return
	}

	// 1. Check if customer exists
	var customer models.Customer
	if err := tx.Where("phone = ?", input.PhoneNo).First(&customer).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			customer = models.Customer{
				CustomerID: uuid.NewString(),
				Name:       input.Name,
				Phone:      input.PhoneNo,
			}
			if err := tx.Create(&customer).Error; err != nil {
				tx.Rollback()
				http.Error(w, "Failed to create customer", http.StatusInternalServerError)
				return
			}
		} else {
			tx.Rollback()
			http.Error(w, "Database error while checking customer", http.StatusInternalServerError)
			return
		}
	}

	// 2. Get property by name
	var property models.Property
	if err := tx.Where("name = ?", input.PropertyName).First(&property).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Property not found", http.StatusBadRequest)
		return
	}

	// 3. Create the booking
	booking := models.Booking{
		CustomerID:   customer.CustomerID,
		PhoneNo:      input.PhoneNo,
		PropertyID:   property.PropertyID,
		Handler:      input.Handler,
		CheckInDate:  input.CheckInDate,
		CheckOutDate: input.CheckOutDate,
		TotalAmount:  input.TotalAmount,
		AdvancePaid:  input.AdvancePaid,
		Status:       input.Status,
		Through:      input.Through,
		Floor:        input.Floor,
		City:         input.City,
		Remarks:      input.Remarks,
	}

	if err := tx.Create(&booking).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Failed to create booking", http.StatusInternalServerError)
		return
	}

	// Commit if everything succeeded
	if err := tx.Commit().Error; err != nil {
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(booking)
}

func GetBookingsByPropertyID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var bookings []models.Booking
	if err := db.DB.Preload("Customer").Preload("Property").Where("property_id = ?", id).Find(&bookings).Error; err != nil {
		http.Error(w, "Failed to fetch bookings", http.StatusInternalServerError)
		return
	}

	// Return bookings as JSON
	if err := json.NewEncoder(w).Encode(bookings); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func GetPropertiesTodayStatus(w http.ResponseWriter, r *http.Request) {
	var properties []models.Property
	err := db.DB.Find(&properties).Error
	if err != nil {
		http.Error(w, "Failed to fetch properties", http.StatusInternalServerError)
		return
	}

	type PropertyStatus struct {
		PropertyID string `json:"PropertyID"`
		Name       string `json:"Name"`
		City       string `json:"City"`
		Status     string `json:"Status"`
	}

	var response []PropertyStatus

	loc, _ := time.LoadLocation("Asia/Kolkata")
	today := time.Now().In(loc).Format("2006-01-02") // formatted as date-only string

	for _, prop := range properties {
		var booking models.Booking
		err := db.DB.Model(&models.Booking{}).
			Where("property_id = ? AND DATE(check_in_date) <= ? AND DATE(check_out_date) > ?", prop.PropertyID, today, today).
			First(&booking).Error

		status := "Available"
		if err == nil {
			status = booking.Status // actual booking status like Confirmed, Pending, etc.
		}

		response = append(response, PropertyStatus{
			PropertyID: prop.PropertyID,
			Name:       prop.Name,
			City:       prop.City,
			Status:     status,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetBookingsByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var booking models.Booking
	if err := db.DB.Preload("Customer").Preload("Property").Where("booking_id = ?", id).Find(&booking).Error; err != nil {
		http.Error(w, "Failed to fetch bookings", http.StatusInternalServerError)
		return
	}

	// Return bookings as JSON
	if err := json.NewEncoder(w).Encode(booking); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func UpdateBooking(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req struct {
		Name     string    `json:"name"`
		PhNo     string    `json:"phNo"`
		CheckIn  time.Time `json:"checkIn"`
		CheckOut time.Time `json:"checkOut"`
		City     string    `json:"city"`
		Property string    `json:"property"`
		Handler  string    `json:"handler"`
		Through  string    `json:"through"`
		Advance  float64   `json:"advance"`
		Total    float64   `json:"total"`
		Status   string    `json:"status"`
		Remarks  string    `json:"remarks"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Fetch booking with customer and property
	var booking models.Booking
	if err := db.DB.Preload("Customer").
		Where("booking_id = ?", id).
		First(&booking).Error; err != nil {
		http.Error(w, "Booking not found", http.StatusNotFound)
		return
	}

	// Validate city & property, then fetch the new property ID
	var newProperty models.Property
	if err := db.DB.Where("city = ? AND name = ?", req.City, req.Property).
		First(&newProperty).Error; err != nil {
		http.Error(w, "New property not found for the given city and name", http.StatusBadRequest)
		return
	}

	// Update booking fields
	booking.CheckInDate = req.CheckIn
	booking.CheckOutDate = req.CheckOut
	booking.City = req.City
	booking.Handler = req.Handler
	booking.Through = req.Through
	booking.AdvancePaid = req.Advance
	booking.TotalAmount = req.Total
	booking.Status = req.Status
	booking.Remarks = req.Remarks
	booking.PropertyID = newProperty.PropertyID // assign new property

	// Update customer details
	if booking.Customer != nil {
		booking.Customer.Name = req.Name
		booking.PhoneNo = req.PhNo
		db.DB.Save(booking.Customer)
	}

	if err := db.DB.Save(&booking).Error; err != nil {
		http.Error(w, "Failed to update booking", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Booking updated successfully"})
}

func DeleteBooking(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	log.Println("Deleting booking:", id)

	var booking models.Booking

	// Check if booking exists
	if err := db.DB.Where("booking_id = ?", id).First(&booking).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Booking not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to find booking", http.StatusInternalServerError)
		}
		return
	}

	// Delete the booking
	if err := db.DB.Delete(&booking).Error; err != nil {
		http.Error(w, "Failed to delete booking", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Booking deleted successfully"})
}

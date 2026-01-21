package handler

import (
  "database/sql"
  "net/http"

  "github.com/gin-gonic/gin"

  "project_ap3/internal/service"
)

type paymentCreateRequest struct {
  BookingID int64  `json:"booking_id"`
  Method    string `json:"method"`
}

type paymentUpdateRequest struct {
  Status string `json:"status"`
}

func (h *Handler) CreatePaymentAPI(c *gin.Context) {
  userID, _ := getUserID(c)

  var req paymentCreateRequest
  if err := c.ShouldBindJSON(&req); err != nil {
    respondJSON(c, http.StatusBadRequest, nil, "invalid_request")
    return
  }
  if req.BookingID == 0 || req.Method == "" {
    respondJSON(c, http.StatusBadRequest, nil, "invalid_input")
    return
  }

  id, err := h.PaymentService.Create(c.Request.Context(), service.PaymentCreateInput{
    BookingID: req.BookingID,
    UserID:    userID,
    Method:    req.Method,
  })
  if err != nil {
    switch err {
    case service.ErrNotFound:
      respondJSON(c, http.StatusNotFound, nil, "not_found")
      return
    case service.ErrForbidden:
      respondJSON(c, http.StatusForbidden, nil, "forbidden")
      return
    case service.ErrInvalid:
      respondJSON(c, http.StatusBadRequest, nil, "invalid_booking")
      return
    case service.ErrConflict:
      respondJSON(c, http.StatusConflict, nil, "conflict")
      return
    default:
      respondJSON(c, http.StatusInternalServerError, nil, "server_error")
      return
    }
  }

  respondJSON(c, http.StatusCreated, gin.H{"id": id}, "")
}

func (h *Handler) UpdatePaymentAPI(c *gin.Context) {
  userID, _ := getUserID(c)
  paymentID, err := parseID(c.Param("id"))
  if err != nil {
    respondJSON(c, http.StatusBadRequest, nil, "invalid_id")
    return
  }

  var req paymentUpdateRequest
  if err := c.ShouldBindJSON(&req); err != nil {
    respondJSON(c, http.StatusBadRequest, nil, "invalid_request")
    return
  }
  if req.Status != "PAID" {
    respondJSON(c, http.StatusBadRequest, nil, "invalid_status")
    return
  }

  if err := h.PaymentService.MarkPaid(c.Request.Context(), paymentID, userID); err != nil {
    if err == sql.ErrNoRows {
      respondJSON(c, http.StatusNotFound, nil, "not_found")
      return
    }
    respondJSON(c, http.StatusInternalServerError, nil, "server_error")
    return
  }

  respondJSON(c, http.StatusOK, gin.H{"status": "PAID"}, "")
}

func (h *Handler) MyPaymentsAPI(c *gin.Context) {
  userID, _ := getUserID(c)

  rows, err := h.Payments.ListByUser(c.Request.Context(), userID)
  if err != nil {
    respondJSON(c, http.StatusInternalServerError, nil, "server_error")
    return
  }
  respondJSON(c, http.StatusOK, rows, "")
}


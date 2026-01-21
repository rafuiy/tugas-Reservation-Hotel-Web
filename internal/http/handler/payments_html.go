package handler

import (
  "database/sql"
  "net/http"

  "github.com/gin-gonic/gin"

  "project_ap3/internal/service"
)

func (h *Handler) MyPaymentsPage(c *gin.Context) {
  userID, _ := getUserID(c)
  payments, err := h.Payments.ListByUser(c.Request.Context(), userID)
  if err != nil {
    h.setDebugErrorHeader(c, err)
    h.render(c, "My Payments", "my_payments", gin.H{
      "Payments": []interface{}{},
      "Bookings": []interface{}{},
      "Error":    "server_error",
    })
    return
  }

  bookings, _ := h.Bookings.ListByUser(c.Request.Context(), userID)
  h.render(c, "My Payments", "my_payments", gin.H{
    "Payments": payments,
    "Bookings": bookings,
  })
}

func (h *Handler) CreatePayment(c *gin.Context) {
  userID, _ := getUserID(c)
  bookingID, err := parseID(c.PostForm("booking_id"))
  if err != nil {
    c.Redirect(http.StatusFound, "/payments/my?err=invalid_booking")
    return
  }
  method := c.PostForm("method")
  if method == "" {
    c.Redirect(http.StatusFound, "/payments/my?err=invalid_input")
    return
  }

  _, err = h.PaymentService.Create(c.Request.Context(), service.PaymentCreateInput{
    BookingID: bookingID,
    UserID:    userID,
    Method:    method,
  })
  if err != nil {
    h.setDebugErrorHeader(c, err)
    switch err {
    case service.ErrNotFound:
      c.Redirect(http.StatusFound, "/payments/my?err=not_found")
      return
    case service.ErrForbidden:
      c.Redirect(http.StatusFound, "/payments/my?err=forbidden")
      return
    case service.ErrInvalid:
      c.Redirect(http.StatusFound, "/payments/my?err=invalid_booking")
      return
    case service.ErrConflict:
      c.Redirect(http.StatusFound, "/payments/my?err=already_paid")
      return
    default:
      c.Redirect(http.StatusFound, "/payments/my?err=server_error")
      return
    }
  }

  c.Redirect(http.StatusFound, "/payments/my?msg=payment_created")
}

func (h *Handler) PayPayment(c *gin.Context) {
  userID, _ := getUserID(c)
  paymentID, err := parseID(c.Param("id"))
  if err != nil {
    c.Redirect(http.StatusFound, "/payments/my?err=invalid_payment")
    return
  }

  err = h.PaymentService.MarkPaid(c.Request.Context(), paymentID, userID)
  if err != nil {
    h.setDebugErrorHeader(c, err)
    if err == sql.ErrNoRows {
      c.Redirect(http.StatusFound, "/payments/my?err=not_found")
      return
    }
    c.Redirect(http.StatusFound, "/payments/my?err=server_error")
    return
  }

  c.Redirect(http.StatusFound, "/payments/my?msg=paid")
}


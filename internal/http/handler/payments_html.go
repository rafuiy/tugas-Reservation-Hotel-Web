package handler

import (
  "database/sql"
  "net/http"
  "strconv"

  "github.com/gin-gonic/gin"

  "project_ap3/internal/model"
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
      "PaymentByBooking": map[int64]model.PaymentView{},
      "PaymentStatusByBooking": map[int64]string{},
      "Error":    "server_error",
    })
    return
  }

  bookings, _ := h.Bookings.ListByUser(c.Request.Context(), userID)
  paymentByBooking := map[int64]model.PaymentView{}
  paymentStatusByBooking := map[int64]string{}
  for _, payment := range payments {
    paymentByBooking[payment.BookingID] = payment
    paymentStatusByBooking[payment.BookingID] = payment.Status
  }
  h.render(c, "My Payments", "my_payments", gin.H{
    "Payments": payments,
    "Bookings": bookings,
    "PaymentByBooking": paymentByBooking,
    "PaymentStatusByBooking": paymentStatusByBooking,
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
  var amount int64
  if amountStr := c.PostForm("amount"); amountStr != "" {
    if parsed, err := strconv.ParseInt(amountStr, 10, 64); err == nil {
      amount = parsed
    }
  }

  _, err = h.PaymentService.Create(c.Request.Context(), service.PaymentCreateInput{
    BookingID: bookingID,
    UserID:    userID,
    Method:    method,
    Amount:    amount,
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
    case service.ErrAmountMismatch:
      c.Redirect(http.StatusFound, "/payments/my?err=amount_mismatch")
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


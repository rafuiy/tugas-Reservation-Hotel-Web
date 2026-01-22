package handler

import (
  "html/template"
  "log"
  "net/http"
  "os"
  "strconv"
  "time"

  "github.com/gin-gonic/gin"

  "project_ap3/internal/config"
  "project_ap3/internal/http/middleware"
  "project_ap3/internal/repo"
  "project_ap3/internal/service"
  "project_ap3/web"
)

type Handler struct {
  Cfg            config.Config
  Users          *repo.UserRepo
  Rooms          *repo.RoomRepo
  Availability   *repo.AvailabilityRepo
  Bookings       *repo.BookingRepo
  Payments       *repo.PaymentRepo
  Invoice        *repo.InvoiceRepo
  BookingService *service.BookingService
  PaymentService *service.PaymentService
}

func New(cfg config.Config, users *repo.UserRepo, rooms *repo.RoomRepo, availability *repo.AvailabilityRepo,
  bookings *repo.BookingRepo, payments *repo.PaymentRepo, invoice *repo.InvoiceRepo,
  bookingSvc *service.BookingService, paymentSvc *service.PaymentService) *Handler {
  return &Handler{
    Cfg:            cfg,
    Users:          users,
    Rooms:          rooms,
    Availability:   availability,
    Bookings:       bookings,
    Payments:       payments,
    Invoice:        invoice,
    BookingService: bookingSvc,
    PaymentService: paymentSvc,
  }
}

func (h *Handler) render(c *gin.Context, title, contentTemplate string, data gin.H) {
  if data == nil {
    data = gin.H{}
  }
  data["Title"] = title
  if _, ok := data["Flash"]; !ok {
    data["Flash"] = c.Query("msg")
  }
  if _, ok := data["Error"]; !ok {
    data["Error"] = c.Query("err")
  }
  if role, ok := c.Get(middleware.CtxRole); ok {
    data["Role"] = role
  }
  if errMsg, ok := data["Error"].(string); ok {
    switch errMsg {
    case "conflict":
      data["Error"] = "Jadwal tersebut sudah ter booked."
    }
  }

  tmpl, err := template.New("").ParseFS(
    web.FS,
    "templates/layout.html",
    "templates/"+contentTemplate+".html",
  )
  if err != nil {
    c.String(http.StatusInternalServerError, "template_error")
    return
  }

  c.Header("Content-Type", "text/html; charset=utf-8")
  c.Status(http.StatusOK)
  if err := tmpl.ExecuteTemplate(c.Writer, "layout", data); err != nil {
    c.String(http.StatusInternalServerError, "template_error")
  }
}

func respondJSON(c *gin.Context, status int, data interface{}, errMsg string) {
  if errMsg != "" {
    c.JSON(status, gin.H{"success": false, "error": errMsg})
    return
  }
  c.JSON(status, gin.H{"success": true, "data": data})
}

func (h *Handler) setDebugErrorHeader(c *gin.Context, err error) {
  if err == nil {
    return
  }
  log.Printf("handler error: %v", err)
  if os.Getenv("DEBUG_ERRORS") == "1" {
    c.Header("X-Error-Detail", err.Error())
  }
}

func parseID(param string) (int64, error) {
  return strconv.ParseInt(param, 10, 64)
}

func getUserID(c *gin.Context) (int64, bool) {
  val, ok := c.Get(middleware.CtxUserID)
  if !ok || val == nil {
    return 0, false
  }
  return val.(int64), true
}

func tokenTTL() time.Duration {
  return 24 * time.Hour
}


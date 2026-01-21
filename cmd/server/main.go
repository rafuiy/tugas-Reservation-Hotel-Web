package main

import (
  "io/fs"
  "log"
  "net/http"

  "github.com/gin-gonic/gin"

  "project_ap3/internal/config"
  "project_ap3/internal/db"
  "project_ap3/internal/http/handler"
  "project_ap3/internal/http/middleware"
  "project_ap3/internal/repo"
  "project_ap3/internal/service"
  "project_ap3/web"
)

func main() {
  cfg, err := config.Load()
  if err != nil {
    log.Fatal(err)
  }

  database, err := db.New(cfg.DatabaseURL)
  if err != nil {
    log.Fatal(err)
  }
  defer database.Close()

  users := repo.NewUserRepo(database)
  rooms := repo.NewRoomRepo(database)
  availability := repo.NewAvailabilityRepo(database)
  bookings := repo.NewBookingRepo(database)
  payments := repo.NewPaymentRepo(database)
  invoice := repo.NewInvoiceRepo(database)

  bookingSvc := service.NewBookingService(database)
  paymentSvc := service.NewPaymentService(database)

  h := handler.New(cfg, users, rooms, availability, bookings, payments, invoice, bookingSvc, paymentSvc)

  r := gin.Default()

  staticFS, err := fs.Sub(web.FS, "static")
  if err != nil {
    log.Fatal(err)
  }
  r.StaticFS("/static", http.FS(staticFS))

  r.GET("/", middleware.TryAuth(cfg.JWTSecret), h.Root)
  r.GET("/login", h.ShowLogin)
  r.POST("/login", h.PostLogin)
  r.GET("/register", h.ShowRegister)
  r.POST("/register", h.PostRegister)
  r.GET("/logout", h.Logout)
  r.POST("/logout", h.Logout)

  user := r.Group("/")
  user.Use(middleware.RequireAuth(cfg.JWTSecret, false), middleware.RequireRole("USER", false))
  user.GET("/rooms", h.RoomsPage)
  user.GET("/availability", h.AvailabilityPage)
  user.GET("/bookings/my", h.MyBookingsPage)
  user.POST("/bookings", h.CreateBooking)
  user.GET("/payments/my", h.MyPaymentsPage)
  user.POST("/payments", h.CreatePayment)
  user.POST("/payments/:id/pay", h.PayPayment)
  user.GET("/invoice/:booking_id", h.InvoicePage)

  admin := r.Group("/admin")
  admin.Use(middleware.RequireAuth(cfg.JWTSecret, false), middleware.RequireRole("ADMIN", false))
  admin.GET("", h.AdminDashboard)
  admin.GET("/rooms", h.AdminRoomsPage)
  admin.POST("/rooms", h.AdminCreateRoom)
  admin.POST("/rooms/:id/update", h.AdminUpdateRoom)
  admin.POST("/rooms/:id/delete", h.AdminDeleteRoom)

  admin.GET("/availability", h.AdminAvailabilityPage)
  admin.POST("/availability", h.AdminCreateAvailability)
  admin.POST("/availability/:id/update", h.AdminUpdateAvailability)
  admin.POST("/availability/:id/delete", h.AdminDeleteAvailability)

  admin.GET("/bookings", h.AdminBookingsPage)
  admin.POST("/bookings/:id/approve", h.AdminApproveBooking)
  admin.POST("/bookings/:id/reject", h.AdminRejectBooking)

  admin.GET("/payments", h.AdminPaymentsPage)
  admin.GET("/users", h.AdminUsersPage)

  api := r.Group("/api")
  api.POST("/auth/register", h.RegisterAPI)
  api.POST("/auth/login", h.LoginAPI)

  apiAuth := api.Group("/")
  apiAuth.Use(middleware.RequireAuth(cfg.JWTSecret, true))
  apiAuth.GET("/auth/me", h.MeAPI)

  apiUser := apiAuth.Group("/")
  apiUser.Use(middleware.RequireRole("USER", true))
  apiUser.GET("/rooms", h.RoomsAPI)
  apiUser.GET("/rooms/:id", h.RoomByIDAPI)
  apiUser.GET("/availability", h.AvailabilityAPI)
  apiUser.POST("/bookings", h.CreateBookingAPI)
  apiUser.GET("/bookings/my", h.MyBookingsAPI)
  apiUser.POST("/payments", h.CreatePaymentAPI)
  apiUser.PATCH("/payments/:id", h.UpdatePaymentAPI)
  apiUser.GET("/payments/my", h.MyPaymentsAPI)
  apiUser.GET("/bookings/:id/invoice", h.InvoiceAPI)

  apiAdmin := apiAuth.Group("/admin")
  apiAdmin.Use(middleware.RequireRole("ADMIN", true))
  apiAdmin.GET("/rooms", h.AdminRoomsAPI)
  apiAdmin.POST("/rooms", h.AdminCreateRoomAPI)
  apiAdmin.PATCH("/rooms/:id", h.AdminUpdateRoomAPI)
  apiAdmin.DELETE("/rooms/:id", h.AdminDeleteRoomAPI)

  apiAdmin.GET("/availability", h.AdminAvailabilityAPI)
  apiAdmin.POST("/availability", h.AdminCreateAvailabilityAPI)
  apiAdmin.PATCH("/availability/:id", h.AdminUpdateAvailabilityAPI)
  apiAdmin.DELETE("/availability/:id", h.AdminDeleteAvailabilityAPI)

  apiAdmin.GET("/bookings", h.AdminBookingsAPI)
  apiAdmin.PATCH("/bookings/:id", h.AdminUpdateBookingAPI)

  apiAdmin.GET("/payments", h.AdminPaymentsAPI)
  apiAdmin.PATCH("/payments/:id", h.AdminUpdatePaymentAPI)

  apiAdmin.GET("/users", h.AdminUsersAPI)

  log.Fatal(r.Run(":" + cfg.Port))
}

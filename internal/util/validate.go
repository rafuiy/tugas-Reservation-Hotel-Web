package util

import (
  "regexp"
  "strings"
  "time"
)

var emailRegex = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)

func ValidateEmail(email string) bool {
  return emailRegex.MatchString(strings.TrimSpace(email))
}

func ValidatePassword(password string) bool {
  return len(strings.TrimSpace(password)) >= 6
}

func ValidateTimeRange(start, end string) bool {
  tStart, err := time.Parse("15:04", start)
  if err != nil {
    return false
  }
  tEnd, err := time.Parse("15:04", end)
  if err != nil {
    return false
  }
  return tStart.Before(tEnd)
}

func ValidateRole(role string) bool {
  return role == "ADMIN" || role == "USER"
}

func ValidateRoomStatus(status string) bool {
  return status == "ACTIVE" || status == "INACTIVE"
}

func NormalizeRoomType(roomType string) string {
  return strings.ToUpper(strings.TrimSpace(roomType))
}

func ValidateRoomType(roomType string) bool {
  _, ok := roomTypeMultiplier()[NormalizeRoomType(roomType)]
  return ok
}

func ApplyRoomTypeMultiplier(basePrice int64, roomType string) (int64, bool) {
  if basePrice < 0 {
    return 0, false
  }
  percent, ok := roomTypeMultiplier()[NormalizeRoomType(roomType)]
  if !ok {
    return 0, false
  }
  adjusted := basePrice + (basePrice*percent)/100
  return adjusted, true
}

func roomTypeMultiplier() map[string]int64 {
  return map[string]int64{
    "STANDARD":      0,
    "SUPERIOR":      10,
    "DELUXE":        20,
    "SUITE":         35,
    "EXECUTIVE":     45,
    "PRESIDENTIAL":  70,
  }
}

func ValidateBookingStatus(status string) bool {
  switch status {
  case "PENDING", "APPROVED", "REJECTED", "CANCELLED":
    return true
  default:
    return false
  }
}

func ValidatePaymentStatus(status string) bool {
  switch status {
  case "UNPAID", "PAID", "FAILED":
    return true
  default:
    return false
  }
}

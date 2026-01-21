package util

import "testing"

func TestHashPassword(t *testing.T) {
  hash, err := HashPassword("secret123")
  if err != nil {
    t.Fatalf("hash error: %v", err)
  }
  if !CheckPassword(hash, "secret123") {
    t.Fatal("expected password to match")
  }
  if CheckPassword(hash, "wrongpass") {
    t.Fatal("expected password mismatch")
  }
}

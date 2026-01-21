# ERD Sistem Reservasi Hotel

```
users (1) --------< bookings >-------- (1) rooms
   |                                  |
   |                                  |
   |                                  v
   |                               availability
   |                                  |
   |                                  v
   +-----------------------------< payments
```

## Relasi
- `users` 1..* `bookings` (satu user bisa banyak booking)
- `rooms` 1..* `availability` (satu room punya banyak slot availability)
- `rooms` 1..* `bookings` (booking mengacu ke room)
- `availability` 1..1 `bookings` (dibatasi `UNIQUE(availability_id)` untuk cegah double booking)
- `bookings` 1..1 `payments` (dibatasi `UNIQUE(booking_id)`)

## Constraint Utama
- `users.email` UNIQUE
- `rooms.room_no` UNIQUE
- `availability` UNIQUE(`room_id`, `date`, `time_start`, `time_end`)
- `availability` `time_start < time_end`
- `bookings` UNIQUE(`availability_id`)
- `payments` UNIQUE(`booking_id`)
- `role` dan status memakai CHECK constraint:
  - `users.role`: `ADMIN` / `USER`
  - `rooms.status`: `ACTIVE` / `INACTIVE`
  - `bookings.status`: `PENDING` / `APPROVED` / `REJECTED` / `CANCELLED`
  - `payments.status`: `UNPAID` / `PAID` / `FAILED`

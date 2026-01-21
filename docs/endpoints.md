# Daftar Endpoint

Format response API JSON:
- Sukses: `{ "success": true, "data": ... }`
- Gagal: `{ "success": false, "error": "..." }`

## HTML Pages
Public:
- GET  `/login`
- POST `/login`
- GET  `/register`
- POST `/register`
- POST `/logout` (GET juga disediakan untuk kemudahan)
- GET  `/` (redirect ke `/rooms` atau `/login`)

User (auth USER):
- GET  `/rooms`
- GET  `/availability` (query: `date`, `room_id` optional)
- GET  `/bookings/my`
- POST `/bookings`
- GET  `/payments/my`
- POST `/payments`
- POST `/payments/:id/pay`
- GET  `/invoice/:booking_id`

Admin (auth ADMIN):
- GET  `/admin`
- GET  `/admin/rooms`
- POST `/admin/rooms`
- POST `/admin/rooms/:id/update`
- POST `/admin/rooms/:id/delete`

- GET  `/admin/availability`
- POST `/admin/availability`
- POST `/admin/availability/:id/update`
- POST `/admin/availability/:id/delete`

- GET  `/admin/bookings`
- POST `/admin/bookings/:id/approve`
- POST `/admin/bookings/:id/reject`

- GET  `/admin/payments`
- GET  `/admin/users`

## API JSON
Auth:
- POST `/api/auth/register`
  - Request: `{ "name":"User A", "email":"a@b.com", "password":"secret123", "phone":"0812" }`
  - Response: `{ "success": true, "data": { "id": 1 } }`
- POST `/api/auth/login`
  - Request: `{ "email":"a@b.com", "password":"secret123" }`
  - Response: `{ "success": true, "data": { "token": "<jwt>" } }`
- GET  `/api/auth/me`

User:
- GET  `/api/rooms`
- GET  `/api/rooms/:id`
- GET  `/api/availability` (query: `date`, `room_id` optional)
- POST `/api/bookings`
  - Request: `{ "availability_id": 1, "guest_name":"Budi", "guest_phone":"0812", "notes":"..." }`
- GET  `/api/bookings/my`
- POST `/api/payments`
  - Request: `{ "booking_id": 10, "method":"TRANSFER" }`
- PATCH `/api/payments/:id`
  - Request: `{ "status":"PAID" }`
- GET  `/api/payments/my`
- GET  `/api/bookings/:id/invoice`

Admin (`/api/admin`):
- Rooms: GET/POST/PATCH/DELETE `/api/admin/rooms`
- Availability: GET/POST/PATCH/DELETE `/api/admin/availability`
- Bookings: GET `/api/admin/bookings`, PATCH `/api/admin/bookings/:id` (approve/reject)
- Payments: GET `/api/admin/payments`, PATCH `/api/admin/payments/:id`
- Users: GET `/api/admin/users`

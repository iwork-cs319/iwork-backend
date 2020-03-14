# I.Work Go-API
This is the backend REST API for booking/viewing workspaces

## Users
### POST /users
- Bulk create users. You have to send a `multipart/form-data` with `users=<users-csv>`. 
- The CSV should have the following format `email, name, id, department, isAdmin` 

## Workspaces
### GET /workspaces
- Get All workspaces objects

### GET /workspaces/:id
- Get workspace object with `id`

### POST /workspaces
- Create new workspace object

### PATCH /workspaces/:id
- Update workspace object with `id`

### DELETE /workspaces/:id
- Delete workspace object with `id`

### GET /workspaces/available?start={start_timestamp}&end={end_timestamp}&floor={floor_id}
- Get ids for all workspaces available to book between `start_time` and `end_time`, where `start_time` and `end_time` are unix timestamps.

### POST /assignments
- Bulk create assignments. You have to send a `multipart/form-data` with `assignments=<assignments-csv>`. 
- The CSV should have the following format `WorkspaceName, FloorName, UserId` 

## Bookings / Offerings (same syntax)
Notes: Any endpoint can have `?start={start_timestamp}&end={end_timestamp}` added to search Bookings/Offerings based on date range.
### GET /bookings
- Get All booking objects

### GET /bookings/:id
- Get booking object with `id`

### GET /bookings/workspaces/:id
- Get booking object with `workspace_id`

### GET /bookings/users/:id
- Get booking object with `user_id`

### POST /bookings
- Create new booking object

### PATCH /bookings/:id
- Update booking object with `id`

### DELETE /booking/:id
- Delete booking object with `id`

## Floor
### GET /floors
- Get All floors objects

### GET /floors/:id
- Get floors object with `id`

### POST /floors
- Create a floor object. You have to send a `multipart/form-data` with `image=<image-data>` and `name=<floor-name>`

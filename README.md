# I.Work Go-API
This is the backend REST API for booking/viewing workspaces

## Users
### GET /users  
- Get All users objects  

### GET /users:id  
- Get user object with `id`  

### POST /users
- Bulk create users. You have to send a `multipart/form-data` with `users=<users-csv>`.
- The CSV should have the following format `email, name, id, department, isAdmin`

### GET /users/assigned?start={start}&end={end}
- Get the users that were assigned to a particular workspace for at least the duration `start_time` to `end_time`

### GET /users/assigned?now={now}
- Get the users that were assigned to a particular workspace at the given time 'now'

## Workspaces
### GET /workspaces
- Get All workspaces objects

### GET /workspaces/:id
- Get workspace object with `id`

### GET /workspaces/bulk/available
- Get all workspaces available, count of available workspaces, and total workspaces for the floor for all floors by `start_time`, `end_time`

### GET /workspaces/bulk/countavailable
- Get count of available workspaces, and total workspaces for the floor for all floors by `start_time`, `end_time`

### POST /workspaces
- Create new workspace object

### POST /bulk/workspaces
- Create workspaces by CSV

### PATCH /workspaces/:id
- Update workspace object with `id`

### PATCH /workspaces/{id}/props
- Update workspace props with `id`

### DELETE /workspaces/:id
- Delete workspace object with `id`

### GET /workspaces/available?start={start_timestamp}&end={end_timestamp}&floor={floor_id}
- Get ids for all workspaces available to book between `start_time` and `end_time`, where `start_time` and `end_time` are unix timestamps.

### GET /workspaces/available?start={start}&?end={end}
- Get ids for all workspaces available to book between `start_time` and `end_time` for all floors, where `start_time` and `end_time` are unix timestamps.

### POST /assignments
- Bulk create assignments. You have to send a `multipart/form-data` with `assignments=<assignments-csv>`.
- The CSV should have the following format `WorkspaceName, FloorName, UserId`

## Bookings / Offerings (same syntax)
Notes: Any endpoint can have `?start={start_timestamp}&end={end_timestamp}` added to search Bookings/Offerings based on date range. Also, any GET endpoint can use `?expand=true` to return additional fields (workspace_name, user_name, floor_id, floor_name).  
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

package connectors

type User struct {
	ID          int    `json:"id"`
	Password    string `json:"password"`
	LastLogin   *int64 `json:"last_login"`
	IsSuperuser int    `json:"is_superuser"`
	Username    string `json:"username"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	IsStaff     uint   `json:"is_staff"`
	IsActive    uint   `json:"is_active"`
	DateJoined  int64  `json:"date_joined"`
}

package http

type CreatePuskesmasAccountRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Contact  string `json:"contact"`
	Password    string `json:"password" validate:"required,min=6"`
	Role        string `json:"role" validate:"required"` // ROLE_PETUGAS_PUSKESMAS or ROLE_KEPALA_PUSKESMAS
	PuskesmasID string `json:"puskesmas_id" validate:"required"`
}

type UpdatePuskesmasAccountRequest struct {
	Name         string `json:"name" validate:"required"`
	Email        string `json:"email" validate:"required,email"`
	Contact      string `json:"contact"`
	ActiveStatus bool   `json:"active_status"`
	Role         string `json:"role"`
	PuskesmasID  string `json:"puskesmas_id"`
	Password     string `json:"password"`
}

type CreateDinkesAccountRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Contact  string `json:"contact"`
	Password string `json:"password" validate:"required,min=6"`
}

type UpdateDinkesAccountRequest struct {
	Name         string `json:"name" validate:"required"`
	Email        string `json:"email" validate:"required,email"`
	Contact      string `json:"contact"`
	ActiveStatus bool   `json:"active_status"`
	Password     string `json:"password"`
}

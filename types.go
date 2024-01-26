// This file contains types that are used in the repository layer.
package repository

type (
	UserRequest struct {
		ID          string
		PhoneNumber string
		FullName    string
	}

	UserResponse struct {
		ID          string
		PhoneNumber string
		FullName    string
		Password    string
	}

	PasswordRequest struct {
		ID       string
		UserID   string
		Password string
	}

	PasswordResponse struct {
		ID       string
		Password string
		UserID   string
	}

	LoginRequest struct {
		ID     string
		UserID string
	}
)

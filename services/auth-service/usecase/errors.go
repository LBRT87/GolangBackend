package usecase

import "errors"

var (
	ErrEmailTaken            = errors.New("email already registered")
	ErrUsernameTaken         = errors.New("username already taken")
	ErrInvalidOTP            = errors.New("OTP code invalid or expired")
	ErrInvalidCreds          = errors.New("invalid email or password")
	ErrInvalidRefresh        = errors.New("refresh token invalid or expired")
	ErrNotVerified           = errors.New("account not verified yet")
	ErrUserNotFound          = errors.New("user not found")
	ErrGoogleLecturerBlocked = errors.New("lecturer account cannot use Google login")
	ErrWrongPassword         = errors.New("old password is incorrect")
)

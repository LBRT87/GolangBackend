package usecase

import (
	"context"
	"crypto/rand"
	"strings"

	"github.com/LBRT87/GolangBackend/services/auth-service/entity"
	"github.com/LBRT87/GolangBackend/services/auth-service/pkg/email"
	"github.com/LBRT87/GolangBackend/services/auth-service/pkg/hash"
	"github.com/LBRT87/GolangBackend/services/auth-service/pkg/jwt"
	"github.com/LBRT87/GolangBackend/services/auth-service/usecase/dto"
)

type AuthUsecase interface {
	Register(ctx context.Context, req dto.RegisterRequest) error
	Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
	VerifyOTP(ctx context.Context, req dto.VerifyOTPRequest) (*dto.LoginResponse, error)
	ResendOTP(ctx context.Context, req dto.ResendOTPRequest) error
	RefreshToken(ctx context.Context, req dto.RefreshTokenRequest) (*dto.LoginResponse, error)
	Logout(ctx context.Context, userID uint) error
	ChangePassword(ctx context.Context, userID uint, req dto.ChangePasswordRequest) error
	UpdateUsername(ctx context.Context, userID uint, req dto.UpdateUsernameRequest) error
	ForgotPassword(ctx context.Context, req dto.ForgotPasswordRequest) error
	ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) error
	LoginWithGoogle(ctx context.Context, info dto.GoogleUserInfo) (*dto.LoginResponse, error)
}

type authUsecase struct {
	userRepo  entity.UserRepository
	cacheRepo entity.CacheRepository
	mailer    *email.Mailer
	jwtMgr    *jwt.Manager
}

func NewAuthUseCase(userrepo entity.UserRepository, cacherepo entity.CacheRepository, mailer *email.Mailer, jwtMgr *jwt.Manager) AuthUsecase {
	return &authUsecase{
		userRepo:  userrepo,
		cacheRepo: cacherepo,
		mailer:    mailer,
		jwtMgr:    jwtMgr,
	}
}

func generateOTP() (string, error) {
	const digits = "0123456789"
	buf := make([]byte, 6)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	code := make([]byte, 6)
	for i, b := range buf {
		code[i] = digits[int(b)%len(digits)]
	}
	return string(code), nil
}

func (u *authUsecase) sendOTP(ctx context.Context, emailAddr string) error {
	code, err := generateOTP()
	if err != nil {
		return err
	}
	if err := u.cacheRepo.SetOTP(ctx, emailAddr, code); err != nil {
		return err
	}
	return u.mailer.SendOTP(emailAddr, code)
}

func (u *authUsecase) Register(ctx context.Context, req dto.RegisterRequest) error {
	if existing, _ := u.userRepo.GetByEmail(ctx, req.Email); existing != nil {
		return ErrEmailTaken
	}
	if existing, _ := u.userRepo.GetByUsername(ctx, req.Username); existing != nil {
		return ErrUsernameTaken
	}

	hashed, err := hash.HashPassword(req.Password)
	if err != nil {
		return err
	}

	user := &entity.User{
		Username:       req.Username,
		FullName:       req.FullName,
		Email:          req.Email,
		Password:       hashed,
		Bio:            "-",
		DOB:            req.DateofBirth,
		ProfilePicture: "-",
	}

	if err := u.userRepo.Create(ctx, user); err != nil {
		return err
	}
	return u.sendOTP(ctx, req.Email)
}

func (u *authUsecase) VerifyOTP(ctx context.Context, req dto.VerifyOTPRequest) (*dto.LoginResponse, error) {
	storedCode, err := u.cacheRepo.GetOTP(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if storedCode == "" || storedCode != req.Code {
		return nil, ErrInvalidOTP
	}

	user, err := u.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	user.IsVerified = true
	if err := u.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}
	if err := u.cacheRepo.DeleteOTP(ctx, req.Email); err != nil {
		return nil, err
	}

	return u.issueTokens(ctx, user)
}

func (u *authUsecase) ResendOTP(ctx context.Context, req dto.ResendOTPRequest) error {
	user, err := u.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return err
	}
	if user == nil {
		return nil
	}
	return u.sendOTP(ctx, req.Email)
}

func (u *authUsecase) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := u.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCreds
	}

	if !hash.CheckPassword(req.Password, user.Password) {
		return nil, ErrInvalidCreds
	}

	if !user.IsVerified {
		return nil, ErrNotVerified
	}

	return u.issueTokens(ctx, user)
}

func (u *authUsecase) issueTokens(ctx context.Context, user *entity.User) (*dto.LoginResponse, error) {
	accessToken, err := u.jwtMgr.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	refreshToken, err := u.jwtMgr.GenerateRefreshToken(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	if err := u.cacheRepo.SetRefreshToken(ctx, user.ID, refreshToken); err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Role:         user.Role,
	}, nil
}

func (u *authUsecase) RefreshToken(ctx context.Context, req dto.RefreshTokenRequest) (*dto.LoginResponse, error) {
	claims, err := u.jwtMgr.Verify(req.RefreshToken)
	if err != nil {
		return nil, ErrInvalidRefresh
	}

	storedToken, err := u.cacheRepo.GetRefreshToken(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}
	if storedToken == "" || storedToken != req.RefreshToken {
		return nil, ErrInvalidRefresh
	}

	user, err := u.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	return u.issueTokens(ctx, user)
}

func (u *authUsecase) Logout(ctx context.Context, userID uint) error {
	return u.cacheRepo.DeleteRefreshToken(ctx, userID)
}

func (u *authUsecase) ChangePassword(ctx context.Context, userID uint, req dto.ChangePasswordRequest) error {
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}
	if !hash.CheckPassword(req.OldPassword, user.Password) {
		return ErrWrongPassword
	}

	hashedPassword, err := hash.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	if err := u.userRepo.Update(ctx, user); err != nil {
		return err
	}

	return u.cacheRepo.DeleteRefreshToken(ctx, user.ID)
}

func (u *authUsecase) UpdateUsername(ctx context.Context, userID uint, req dto.UpdateUsernameRequest) error {
	if existing, _ := u.userRepo.GetByUsername(ctx, req.Username); existing != nil && existing.ID != userID {
		return ErrUsernameTaken
	}

	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	user.Username = req.Username
	return u.userRepo.Update(ctx, user)
}

func (u *authUsecase) ForgotPassword(ctx context.Context, req dto.ForgotPasswordRequest) error {
	user, err := u.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return err
	}
	if user == nil {
		return nil
	}
	return u.sendOTP(ctx, req.Email)
}

func (u *authUsecase) ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) error {
	storedCode, err := u.cacheRepo.GetOTP(ctx, req.Email)
	if err != nil {
		return err
	}
	if storedCode == "" || storedCode != req.Code {
		return ErrInvalidOTP
	}

	user, err := u.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	hashedPassword, err := hash.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	if err := u.userRepo.Update(ctx, user); err != nil {
		return err
	}

	if err := u.cacheRepo.DeleteOTP(ctx, req.Email); err != nil {
		return err
	}
	return u.cacheRepo.DeleteRefreshToken(ctx, user.ID)
}

func (u *authUsecase) LoginWithGoogle(ctx context.Context, info dto.GoogleUserInfo) (*dto.LoginResponse, error) {
	user, err := u.userRepo.GetByGoogle(ctx, info.Sub)
	if err != nil {
		return nil, err
	}

	if user == nil {
		user, err = u.userRepo.GetByEmail(ctx, info.Email)
		if err != nil {
			return nil, err
		}

		if user == nil {
			username := strings.Split(info.Email, "@")[0]
			googleID := info.Sub
			user = &entity.User{
				Username:       username,
				Email:          info.Email,
				FullName:       info.Name,
				ProfilePicture: info.Picture,
				GoogleID:       &googleID,
				Role:           entity.RoleStudent,
				IsVerified:     true,
			}
			if err := u.userRepo.Create(ctx, user); err != nil {
				return nil, err
			}
		} else {
			googleID := info.Sub
			user.GoogleID = &googleID
			if info.Picture != "" && user.ProfilePicture == "" {
				user.ProfilePicture = info.Picture
			}
			if err := u.userRepo.Update(ctx, user); err != nil {
				return nil, err
			}
		}
	}

	if user.Role == entity.RoleLecturer {
		return nil, ErrGoogleLecturerBlocked
	}

	return u.issueTokens(ctx, user)
}

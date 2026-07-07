package handler

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/LBRT87/GolangBackend/services/auth-service/pkg/jwt"
	"github.com/LBRT87/GolangBackend/services/auth-service/pkg/response"
	"github.com/LBRT87/GolangBackend/services/auth-service/usecase"
	"github.com/LBRT87/GolangBackend/services/auth-service/usecase/dto"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/oauth2"
)

const oauthStateCookie = "oauth_state"

type AuthHandler struct {
	usecase      usecase.AuthUsecase
	oauthCfg     *oauth2.Config
	cookieSecure bool
	frontendURL  string
}

func NewAuthHandler(uc usecase.AuthUsecase, oauthCfg *oauth2.Config, cookieSecure bool, frontendURL string) *AuthHandler {
	return &AuthHandler{
		usecase:      uc,
		oauthCfg:     oauthCfg,
		cookieSecure: cookieSecure,
		frontendURL:  frontendURL,
	}
}

func (h *AuthHandler) setAuthCookies(c *gin.Context, accessToken, refreshToken string) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("access_token", accessToken, int(jwt.AccessTokenTTL.Seconds()), "/", "", h.cookieSecure, true)
	c.SetCookie("refresh_token", refreshToken, int(jwt.RefreshTokenTTL.Seconds()), "/", "", h.cookieSecure, true)
}

func (h *AuthHandler) clearAuthCookies(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("access_token", "", -1, "/", "", h.cookieSecure, true)
	c.SetCookie("refresh_token", "", -1, "/", "", h.cookieSecure, true)
}

func generateState() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

func friendlyError(err error) string {
	var ve validator.ValidationErrors
	if !errors.As(err, &ve) {
		return err.Error()
	}
	messages := map[string]map[string]string{
		"Username": {
			"required": "Username must be filled",
			"min":      "Username minimal 3 characters",
			"max":      "Username maximal 50 characters",
		},
		"Email": {
			"required": "Email must be filled",
			"email":    "Format email not valid",
		},
		"Password": {
			"required": "Password must be filled",
			"min":      "Password minimal 6 characters",
		},
		"NewPassword": {
			"required": "New Password must be filled",
			"min":      "New Password minimal 6 characters",
		},
		"FullName": {
			"required": "Full name must be filled",
		},
		"Code": {
			"required": "OTP Code must be filled",
			"len":      "OTP Code must 6 digit"},
		"Identity": {
			"required": "Email must be filled",
		},
		"RefreshToken": {
			"required": "Refresh token must be filled"},
	}
	for _, fe := range ve {
		if fieldMap, ok := messages[fe.Field()]; ok {
			if msg, ok := fieldMap[fe.Tag()]; ok {
				return msg
			}
		}
		return fmt.Sprintf("Field '%s' not valid (%s)", fe.Field(), fe.Tag())
	}
	return err.Error()
}

func statusFor(err error) int {
	switch {
	case errors.Is(err, usecase.ErrEmailTaken),
		errors.Is(err, usecase.ErrUsernameTaken):
		return http.StatusConflict
	case errors.Is(err, usecase.ErrInvalidOTP),
		errors.Is(err, usecase.ErrInvalidCreds),
		errors.Is(err, usecase.ErrInvalidRefresh),
		errors.Is(err, usecase.ErrNotVerified):
		return http.StatusUnauthorized
	case errors.Is(err, usecase.ErrUserNotFound):
		return http.StatusNotFound
	case errors.Is(err, usecase.ErrGoogleLecturerBlocked):
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, friendlyError(err))
		return
	}
	if err := h.usecase.Register(c.Request.Context(), req); err != nil {
		response.Error(c, statusFor(err), err.Error())
		return
	}
	response.Success(c, http.StatusCreated, "registration success, check email", nil)
}

func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var req dto.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	result, err := h.usecase.VerifyOTP(c.Request.Context(), req)
	if err != nil {
		response.Error(c, statusFor(err), err.Error())
		return
	}
	h.setAuthCookies(c, result.AccessToken, result.RefreshToken)
	response.Success(c, http.StatusOK, "account success verified", result)
}

func (h *AuthHandler) ResendOTP(c *gin.Context) {
	var req dto.ResendOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.usecase.ResendOTP(c.Request.Context(), req); err != nil {
		response.Error(c, statusFor(err), err.Error())
		return
	}
	response.Success(c, http.StatusOK, "OTP code has been sent", nil)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	result, err := h.usecase.Login(c.Request.Context(), req)
	if err != nil {
		response.Error(c, statusFor(err), err.Error())
		return
	}
	h.setAuthCookies(c, result.AccessToken, result.RefreshToken)
	response.Success(c, http.StatusOK, "login success", result)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil || refreshToken == "" {
		response.Error(c, http.StatusUnauthorized, "refresh token not found")
		return
	}

	result, err := h.usecase.RefreshToken(c.Request.Context(), dto.RefreshTokenRequest{RefreshToken: refreshToken})
	if err != nil {
		response.Error(c, statusFor(err), err.Error())
		return
	}
	h.setAuthCookies(c, result.AccessToken, result.RefreshToken)
	response.Success(c, http.StatusOK, "token successfully updated", result)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	userID := c.GetUint("user_id")
	if err := h.usecase.Logout(c.Request.Context(), userID); err != nil {
		response.Error(c, statusFor(err), err.Error())
		return
	}
	h.clearAuthCookies(c)
	response.Success(c, http.StatusOK, "Log out success", nil)
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, friendlyError(err))
		return
	}
	if err := h.usecase.ChangePassword(c.Request.Context(), userID, req); err != nil {
		response.Error(c, statusFor(err), err.Error())
		return
	}

	h.clearAuthCookies(c)
	response.Success(c, http.StatusOK, "Password success changed", nil)
}

func (h *AuthHandler) UpdateUsername(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req dto.UpdateUsernameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, friendlyError(err))
		return
	}
	if err := h.usecase.UpdateUsername(c.Request.Context(), userID, req); err != nil {
		response.Error(c, statusFor(err), err.Error())
		return
	}
	response.Success(c, http.StatusOK, "username updated", nil)
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req dto.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.usecase.ForgotPassword(c.Request.Context(), req); err != nil {
		response.Error(c, statusFor(err), err.Error())
		return
	}
	response.Success(c, http.StatusOK, "if email registered, OTP code will be send", nil)
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.usecase.ResetPassword(c.Request.Context(), req); err != nil {
		response.Error(c, statusFor(err), err.Error())
		return
	}
	response.Success(c, http.StatusOK, "password success reset", nil)
}

func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	if h.oauthCfg == nil {
		response.Error(c, http.StatusNotImplemented, "Google OAuth notyet config")
		return
	}

	state, err := generateState()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to prepare Google login")
		return
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(oauthStateCookie, state, 5*60, "/", "", h.cookieSecure, true)

	url := h.oauthCfg.AuthCodeURL(state, oauth2.AccessTypeOnline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	if h.oauthCfg == nil {
		response.Error(c, http.StatusNotImplemented, "Google OAuth not yet configured")
		return
	}

	expectedState, err := c.Cookie(oauthStateCookie)
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(oauthStateCookie, "", -1, "/", "", h.cookieSecure, true)
	if err != nil || expectedState == "" || c.Query("state") != expectedState {
		response.Error(c, http.StatusBadRequest, "invalid or expired OAuth state")
		return
	}

	code := c.Query("code")
	if code == "" {
		response.Error(c, http.StatusBadRequest, "code empty from Google")
		return
	}

	token, err := h.oauthCfg.Exchange(context.Background(), code)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed change code with token Google: "+err.Error())
		return
	}

	httpClient := h.oauthCfg.Client(context.Background(), token)
	resp, err := httpClient.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed get info user from Google")
		return
	}
	defer resp.Body.Close()

	var userInfo dto.GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		response.Error(c, http.StatusInternalServerError, "failed parse info user Google")
		return
	}

	result, err := h.usecase.LoginWithGoogle(c.Request.Context(), userInfo)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/login?error=%s", h.frontendURL, url.QueryEscape(err.Error())))
		return
	}
	h.setAuthCookies(c, result.AccessToken, result.RefreshToken)
	c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/oauth/callback?role=%s", h.frontendURL, result.Role))

}

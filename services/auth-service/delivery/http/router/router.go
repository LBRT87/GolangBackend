package router

import (
	"github.com/gin-gonic/gin"
	"github.com/LBRT87/GolangBackend/services/auth-service/delivery/http/handler"
	"github.com/LBRT87/GolangBackend/services/auth-service/delivery/http/middleware"
	"github.com/LBRT87/GolangBackend/services/auth-service/pkg/jwt"
)

func NewRouter(h *handler.AuthHandler, jwtMgr *jwt.Manager) *gin.Engine {
	r := gin.Default()
	r.GET("/health", func(ctx *gin.Context) {
		ctx.String(200, "auth-service oke")
	})
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/verify-otp", h.VerifyOTP)
		auth.POST("/resend-otp", h.ResendOTP)
		auth.POST("/login", h.Login)
		auth.POST("/refresh-token", h.RefreshToken)
		auth.POST("/forgot-password", h.ForgotPassword)
		auth.POST("/reset-password", h.ResetPassword)
		auth.POST("/logout", middleware.JWTAuth(jwtMgr), h.Logout)
		auth.POST("/change-password", middleware.JWTAuth(jwtMgr), h.ChangePassword)
		auth.PUT("/username", middleware.JWTAuth(jwtMgr), h.UpdateUsername)

		auth.GET("/google", h.GoogleLogin)
		auth.GET("/google/callback", h.GoogleCallback)
	}

	return r
}

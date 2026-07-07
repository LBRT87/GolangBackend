package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	grpcdelivery "github.com/LBRT87/GolangBackend/services/auth-service/delivery/grpc"
	"github.com/LBRT87/GolangBackend/services/auth-service/delivery/http/handler"
	"github.com/LBRT87/GolangBackend/services/auth-service/delivery/http/router"
	"github.com/LBRT87/GolangBackend/services/auth-service/entity"
	gen "github.com/LBRT87/GolangBackend/contracts/auth-service/gen"
	"github.com/LBRT87/GolangBackend/services/auth-service/pkg/email"
	"github.com/LBRT87/GolangBackend/services/auth-service/pkg/jwt"
	postgresRepo "github.com/LBRT87/GolangBackend/services/auth-service/repository/postgres"
	redisrepo "github.com/LBRT87/GolangBackend/services/auth-service/repository/redis"
	"github.com/LBRT87/GolangBackend/services/auth-service/usecase"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	gormpkg "gorm.io/gorm"
)

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func mustGetEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("environment variable %s wajib diisi, cek file .env", key)
	}
	return v
}

func newRedisClient() redis.UniversalClient {
	password := getEnv("REDIS_PASSWORD", "")
	if addrs := os.Getenv("REDIS_ADDRS"); addrs != "" {
		return redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    strings.Split(addrs, ","),
			Password: password,
		})
	}
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", getEnv("REDIS_HOST", "localhost"), getEnv("REDIS_PORT", "6379")),
		Password: password,
	})
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("tidak ada file .env, lanjut pakai environment variable system")
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		mustGetEnv("DB_USER"),
		mustGetEnv("DB_PASSWORD"),
		getEnv("DB_NAME","onsite_auth_db"),
	)

	db, err := gormpkg.Open(postgres.Open(dsn), &gormpkg.Config{})
	if err != nil {
		log.Fatalf("gagal konek database: %v", err)
	}
	if err := db.AutoMigrate(&entity.User{}); err != nil {
		log.Fatalf("gagal migrasi database: %v", err)
	}

	redisClient := newRedisClient()

	mailer, err := email.NewMailer(email.Config{
		Host:     getEnv("SMTP_HOST", "smtp.gmail.com"),
		Port:     587,
		Username: getEnv("SMTP_USERNAME", ""),
		Password: getEnv("SMTP_PASSWORD", ""),
		From:     getEnv("SMTP_FROM", ""),
	})
	if err != nil {
		log.Fatalf("gagal setup mailer: %v", err)
	}

	jwtMgr := jwt.NewManager(mustGetEnv("JWT_SECRET"))

	userRepo := postgresRepo.NewUserRepository(db)
	cacheRepo := redisrepo.NewCacheRepository(redisClient)

	authUsecase := usecase.NewAuthUseCase(userRepo, cacheRepo, mailer, jwtMgr)

	var oauthCfg *oauth2.Config
	if clientID := getEnv("GOOGLE_CLIENT_ID", ""); clientID != "" {
		oauthCfg = &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
			RedirectURL:  getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8001/api/auth/google/callback"),
			Scopes:       []string{"openid", "email", "profile"},
			Endpoint:     google.Endpoint,
		}
		log.Println("Google OAuth aktif")
	} else {
		log.Println("gugel client id kosong, endpoint /api/auth/google tidak aktif")
	}

	cookieSecure := getEnv("COOKIE_SECURE", "false") == "true"
	frontendURL := getEnv("FRONTEND_URL", "http://localhost:5173")

	authHandler := handler.NewAuthHandler(authUsecase, oauthCfg, cookieSecure, frontendURL)
	r := router.NewRouter(authHandler, jwtMgr)

	go func() {
		grpcPort := getEnv("GRPC_PORT", "9001")
		lis, err := net.Listen("tcp", ":"+grpcPort)
		if err != nil {
			log.Fatalf("gagal listen gRPC port %s: %v", grpcPort, err)
		}
		grpcSrv := grpc.NewServer()
		gen.RegisterAuthServiceServer(grpcSrv, grpcdelivery.NewAuthGRPCServer(jwtMgr, userRepo))
		log.Printf("auth-service gRPC jalan di port %s", grpcPort)
		if err := grpcSrv.Serve(lis); err != nil {
			log.Fatalf("gRPC server error: %v", err)
		}
	}()

	port := getEnv("APP_PORT", "8001")
	log.Printf("auth-service HTTP jalan di port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("gagal start server: %v", err)
	}
}

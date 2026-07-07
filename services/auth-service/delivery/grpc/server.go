package grpc

import (
	"context"

	gen "github.com/LBRT87/GolangBackend/contracts/auth-service/gen"
	"github.com/LBRT87/GolangBackend/services/auth-service/entity"
	"github.com/LBRT87/GolangBackend/services/auth-service/pkg/jwt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthGRPCServer struct {
	gen.UnimplementedAuthServiceServer
	jwtMgr   *jwt.Manager
	userRepo entity.UserRepository
}

func NewAuthGRPCServer(jwtMgr *jwt.Manager, userRepo entity.UserRepository) *AuthGRPCServer {
	return &AuthGRPCServer{
		jwtMgr:   jwtMgr,
		userRepo: userRepo,
	}
}

func (s *AuthGRPCServer) VerifyToken(ctx context.Context, req *gen.VerifyTokenRequest) (*gen.VerifyTokenResponse, error) {
	claims, err := s.jwtMgr.Verify(req.GetToken())
	if err != nil {
		return &gen.VerifyTokenResponse{IsVerified: false}, nil
	}

	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil || user == nil {
		return &gen.VerifyTokenResponse{IsVerified: false}, nil
	}

	return &gen.VerifyTokenResponse{
		UserId:     uint64(user.ID),
		Email:      user.Email,
		Username:   user.Username,
		Role:       user.Role,
		IsVerified: true,
	}, nil
}

func (s *AuthGRPCServer) GetUser(ctx context.Context, req *gen.GetUserRequest) (*gen.GetUserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, uint(req.GetUserId()))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if user == nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return &gen.GetUserResponse{
		Id:       uint64(user.ID),
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}, nil
}

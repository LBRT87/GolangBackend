package usecase

import (
	"github.com/LBRT87/GolangBackend/services/payment-service/config"
	"github.com/LBRT87/GolangBackend/services/payment-service/internal/domain"
)

const (
	TaskQueue = "payment-task-queue"
	CourseCheckoutFlowName = "CourseCheckoutWorkflow"
)

type PaymentUsecase struct {
	orderRepo domain.OrderRepository
	orderItemRepo domain.OrderItemRepository
	lectureBalanceRepo domain.LecturerBalanceRepository
	withdrawalRepo domain.WithdrawalRepository
	cfg *config.Config
}

func NewPaymentUsecase(orderRepo domain.OrderRepository,
	orderItemRepo domain.OrderItemRepository,
	lectureBalanceRepo domain.LecturerBalanceRepository,
	withdrawalRepo domain.WithdrawalRepository,
	cfg *config.Config) *PaymentUsecase{
		return &PaymentUsecase{orderRepo: orderRepo, orderItemRepo: orderItemRepo,lectureBalanceRepo: lectureBalanceRepo,
		withdrawalRepo: withdrawalRepo, cfg: cfg,}
	}

func (u *PaymentUsecase) Checkout(userId string, courseId []string) (string, error) {
	
}
package repository

import (
	"errors"

	"github.com/LBRT87/GolangBackend/services/payment-service/internal/domain"
	"gorm.io/gorm"
)

type PostgresBalanceRepository struct {
	db *gorm.DB
}

func NewPostgresBalanceRepository(db *gorm.DB) *PostgresBalanceRepository{
	return &PostgresBalanceRepository{db: db}
}

func (r PostgresBalanceRepository) FindByLecturer(lecturerId string) (*domain.LecturerBalance, error){
	var LecturerBalance domain.LecturerBalance
	err := r.db.Where("lecturer_id = ?", lecturerId).First(&LecturerBalance).Error; 
	if errors.Is(err, gorm.ErrRecordNotFound){
		return nil, gorm.ErrRecordNotFound
	}
	if err != nil{
		return nil, err
	}
	return &LecturerBalance, nil
}

func (r *PostgresBalanceRepository) CreditEarnings(lecturerId string, amount int64) error{
	lecturerBalance, err := r.FindByLecturer(lecturerId)
	if err != nil{
		return err
	}
	var amountCasted float64 = float64(amount)
	var curr float64 = lecturerBalance.TotalEarnings + amountCasted
	if err := r.db.Model(&domain.LecturerBalance{}).Where("lecturer_id = ?", lecturerId).Update("total_earnings", curr).Error; err != nil{
		return err
	}
	return nil
}

type PostgresWithdrawalRepository struct{
	db *gorm.DB
}

func NewPostgresWithdrawalRepository(db *gorm.DB) *PostgresWithdrawalRepository{
	return &PostgresWithdrawalRepository{db: db}
}

func (r *PostgresWithdrawalRepository) CreateWithdrawal(withdrawal *domain.Withdrawal) error {
	return r.db.Create(withdrawal).Error
}

func (r *PostgresWithdrawalRepository)ListByLecturer(lecturerId string) (*[]domain.Withdrawal, error){
	var Withdrawal []domain.Withdrawal
	if err := r.db.Where("lecturer_id = ?", lecturerId).Find(&Withdrawal).Error; err != nil{
		return nil, err
	}
	return &Withdrawal, nil
}
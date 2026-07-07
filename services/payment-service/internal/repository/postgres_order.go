package repository

import (
	"errors"

	"github.com/LBRT87/GolangBackend/services/payment-service/internal/domain"
	"gorm.io/gorm"
)

type PostgresOrderRepository struct {
	db *gorm.DB
}

func NewPostgresOrderRepository(db *gorm.DB) *PostgresOrderRepository{
	return &PostgresOrderRepository{db: db}
}

func (r *PostgresOrderRepository) Create(order *domain.Order) error{
	if err:= r.db.Create(order); err != nil {return err.Error}
	return nil
}

func (r *PostgresOrderRepository) FindByGatewayOrderId(GatewayOrderId string) (*domain.Order, error){
	var order domain.Order
	err := r.db.Where("gateway_order_id = ?", GatewayOrderId).First(&order); 
	if errors.Is(err.Error, gorm.ErrRecordNotFound){
		return nil, gorm.ErrRecordNotFound
	}
	if err != nil{
		return nil, err.Error
	}
	return &order, nil
}

func (r *PostgresOrderRepository) FindByOrderId(orderId string) (*domain.Order, error){
	var order domain.Order
	err := r.db.Where("order_id = ?", orderId).First(&order); 
	if errors.Is(err.Error, gorm.ErrRecordNotFound){
		return nil, gorm.ErrRecordNotFound
	}
	if err != nil{
		return nil, err.Error
	}
	return &order, nil
}

func (r *PostgresOrderRepository) ListByUser(userId string) (*[]domain.Order, error){
	var order []domain.Order
	err := r.db.Where("user_id = ?", userId).Find(&order); 
	if errors.Is(err.Error, gorm.ErrRecordNotFound){
		return nil, gorm.ErrRecordNotFound
	}
	if err != nil{
		return nil, err.Error
	}
	return &order, nil
}


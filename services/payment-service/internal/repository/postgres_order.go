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

func (r *PostgresOrderRepository) UpdateStatus(orderId string, status string) error {
	if err := r.db.Model(&domain.Order{}).Where("order_id = ?", orderId).Update("status", status).Error; err != nil {
		return err
	}
	return nil
}

func (r *PostgresOrderRepository) SetGatewayInfo(orderId string, GatewayOrderId string, paymentUrl string) error{
	if err := r.db.Model(&domain.Order{}).Where("order_id = ?", orderId).Updates(map[string]interface{}{
		"gateway_order_id": GatewayOrderId,
		"payment_url" : paymentUrl,
	}); err != nil{
		return err.Error
	}
	return nil
}

type PostgresOrderItemRepository struct{
	db *gorm.DB
}

func NewPostgresOrderItemRepository(db *gorm.DB) *PostgresOrderItemRepository{
	return &PostgresOrderItemRepository{db: db}
}

func (r *PostgresOrderItemRepository) Create(orderItem *domain.OrderItem) error {
	return r.db.Create(orderItem).Error
}

func (r *PostgresOrderItemRepository) FindOrderItemByOrderId(orderId string) ([]domain.OrderItem, error){
	var orderItems []domain.OrderItem
	err := r.db.Where("order_id = ?", orderId).Find(&orderItems).Error
	if err != nil{
		return nil, err
	}
	return orderItems, err
}

func (r *PostgresOrderItemRepository) CreateBatch(orderItems []domain.OrderItem) error{
	return r.db.Create(orderItems).Error
}

func (r *PostgresOrderItemRepository) ListByLecturer(lecturerId string) ([]domain.OrderItem, error){
	var orderItems []domain.OrderItem
	err := r.db.Where("lecturer_id= ?", lecturerId).Find(&orderItems).Error
	if err != nil{
		return nil, err
	}
	return orderItems, err
}
package domain

type Order struct{
	ID string `gorm:"primaryKey;type:uuid"`
	UserId string
	GatewayOrderId string
	PaymentUrl string
	Status string `gorm:"default:PENDING"`
	Total float64
}

type OrderItem struct{
	Id string `gorm:"primaryKey;type:uuid"`
	OrderId string
	CourseId string
	UserId string
	LecturerId string
}

type OrderRepository interface{
	Create(order *Order) error
	FindByGatewayOrderId(GatewayOrderId string) (*Order, error)
	FindById(orderId string) (*Order, error)
	setGatewayInfo(id, GatewayOrderId, PaymentUrl string) error
	UpdateStatus(id, status string) error
	ListByUser(UserId string) ([]Order, error)
}

type OrderItemRepository interface{
	Create(orderItem OrderItem) error
	FindById(orderItemId string) (OrderItem, error)
	FindByOrderId(orderItemId string) ([]OrderItem, error)
	CreateBatch(orderItems []OrderItem) error
	ListByLecturer(LecturerId string) ([]OrderItem, error)
}
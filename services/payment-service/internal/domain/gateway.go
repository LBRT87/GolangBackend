package domain

type CustomerInfo struct{
	name string
	email string
}

type GatewayRepository interface{
	CreateTransaction(orderId string, amount int64, cust CustomerInfo) (paymentUrl string, GatewayOrderId string, err error)
	VerifyWebhokSignature(payload []byte) bool
	ParseWebhook(payload []byte)(orderId string, status string, err error)
}
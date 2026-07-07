package domain

import "time"

type LecturerBalance struct {
	Lecturerid     string `gorm:"primaryKey;type:uuid"`
	TotalEarnings  float64
	TotalWithdrawn float64
	UpdatedAt      time.Time
}

type Withdrawal struct{
	Id string `gorm:"primaryKey;type:uuid"`
	LecturerId string
	amount int64
	WithdrawnAt time.Time
}

type LecturerBalanceRepository interface{
	FindByLecturer(lecturerId string) (LecturerBalance, error)
	CreditEarnings(lecturerId string, amount int64) error
	DebitFromWithdrawal(lecturerId string, amount int64) error
}

type WithdrawalRepository interface{
	Create(withdraw *Withdrawal) error
	ListByLecturer(lecturerId string) ([]Withdrawal, error)
}


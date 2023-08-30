package database

type User struct {
	ID        uint      `gorm:"primaryKey"`
	FirstName string    `gorm:"not null"`
	LastName  string    `gorm:"not null"`
	Balance   float64   `gorm:"not null"`
	Invoices  []Invoice `gorm:"constraint:OnDelete:CASCADE;"`
}

type Invoice struct {
	ID     uint `gorm:"primaryKey"`
	UserID uint `gorm:"not null"`
	User   User
	Paid   bool    `gorm:"default:false"`
	Label  string  `gorm:"not null"`
	Amount float64 `gorm:"not null"`
}

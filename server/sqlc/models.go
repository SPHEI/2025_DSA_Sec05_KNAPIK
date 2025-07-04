// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package sqlc

import (
	"time"
	"server/types"
)

type Apartment struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	Street         string `json:"street"`
	BuildingNumber string `json:"building_number"`
	BuildingName   string `json:"building_name"`
	FlatNumber     string `json:"flat_number"`
	OwnerID        int64  `json:"owner_id"`
}

type Expense struct {
	ID          int64       `json:"id"`
	Amount      float64     `json:"amount"`
	ExpenseDate time.Time   `json:"expense_date"`
	Description string      `json:"description"`
	CategoryID  int64       `json:"category_id"`
	RepairID    interface{} `json:"repair_id"`
}

type ExpenseCategory struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type FaultReport struct {
	ID           int64     `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	DateReported time.Time `json:"date_reported"`
	StatusID     int64     `json:"status_id"`
	ApartmentID  int64     `json:"apartment_id"`
	UserID       int64     `json:"user_id"`
}

type FaultStatus struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type FinancialRecord struct {
	Type             string       `json:"type"`
	SourceID         int64        `json:"source_id"`
	Amount           float64      `json:"amount"`
	RecordDate       types.JSONNullTime `json:"record_date"`
	Description      string       `json:"description"`
	RelatedPaymentID int64        `json:"related_payment_id"`
	RelatedExpenseID interface{}  `json:"related_expense_id"`
	UserName         string       `json:"user_name"`
	ApartmentName    string       `json:"apartment_name"`
}

type Payment struct {
	ID                   int64          `json:"id"`
	Amount               float64        `json:"amount"`
	PaymentDate          types.JSONNullTime   `json:"payment_date"`
	DueDate              time.Time      `json:"due_date"`
	StatusID             int64          `json:"status_id"`
	RentingID            int64          `json:"renting_id"`
	TransactionReference types.JSONNullString `json:"transaction_reference"`
}

type PaymentStatus struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type PricingHistory struct {
	ID          int64         `json:"id"`
	ApartmentID int64         `json:"apartment_id"`
	Date        time.Time     `json:"date"`
	Price       float64       `json:"price"`
	IsCurrent   types.JSONNullInt64 `json:"is_current"`
}

type RentingHistory struct {
	ID          int64        `json:"id"`
	ApartmentID int64        `json:"apartment_id"`
	UserID      int64        `json:"user_id"`
	StartDate   time.Time    `json:"start_date"`
	EndDate     types.JSONNullTime `json:"end_date"`
	IsCurrent   int64        `json:"is_current"`
}

type Repair struct {
	ID              int64         `json:"id"`
	Title           string        `json:"title"`
	FaultReportID   int64         `json:"fault_report_id"`
	DateAssigned    time.Time     `json:"date_assigned"`
	DateCompleted   types.JSONNullTime  `json:"date_completed"`
	StatusID        int64         `json:"status_id"`
	SubcontractorID types.JSONNullInt64 `json:"subcontractor_id"`
}

type RepairStatus struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Role struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Speciality struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Subcontractor struct {
	ID           int64  `json:"id"`
	UserID       int64  `json:"user_id"`
	Address      string `json:"address"`
	Nip          string `json:"nip"`
	SpecialityID int64  `json:"speciality_id"`
}

type User struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	RoleID   int64  `json:"role_id"`
	Password string `json:"password"`
}

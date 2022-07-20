package models

type CreditLine struct {
	UserId             string
	ProductId          string
	SinglePaymentLimit string
	InstallmentsLimit  string
	CashAdvanceLimits  string
	OfferStartDate     string
	OfferEndDate       string
	DueDateDay         string
	UserScoring        string
}

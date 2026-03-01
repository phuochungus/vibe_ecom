package dto

type CreatePaymentRequest struct {
	Provider  string `json:"provider"`
	ReturnURL string `json:"return_url"`
	CancelURL string `json:"cancel_url"`
}

type CreatePaymentResponseDTO struct {
	PaymentID   string `json:"payment_id"`
	Status      string `json:"status"`
	CheckoutURL string `json:"checkout_url"`
}

type PaymentItemDTO struct {
	ID              string `json:"id"`
	TxnType         string `json:"txn_type"`
	Provider        string `json:"provider"`
	ProviderTxnCode string `json:"provider_txn_code,omitempty"`
	Status          string `json:"status"`
	Amount          string `json:"amount"`
	CurrencyCode    string `json:"currency_code"`
	CreatedAt       string `json:"created_at"`
}
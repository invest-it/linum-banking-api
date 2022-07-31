package nordigen

type tokenInfo struct {
	Access         string `json:"access"`
	AccessExpires  int    `json:"access_expires"`
	Refresh        string `json:"refresh"`
	RefreshExpires string `json:"refresh_expires"`
}

type institution struct {
	Id                   string   `json:"id"`
	Name                 string   `json:"name"`
	Bic                  string   `json:"bic,omitempty"`
	TransactionTotalDays string   `json:"transaction_total_days,omitempty"`
	Countries            []string `json:"countries"`
	Logo                 string   `json:"logo"`
}

type requisitionRequest struct {
	Redirect      string `json:"redirect"`
	InstitutionId string `json:"institution_id"`
	Reference     string `json:"reference"`
	Agreement     string `json:"agreement,omitempty"`
	UserLanguage  string `json:"user_language"`
}

type requisition struct {
	Id                string `json:"id"`      // TODO: Decide if it makes sense to parse as UUID
	Created           string `json:"created"` // TODO: Decide if it makes sense to parse as DateTime
	Redirect          string `json:"redirect"`
	Status            string `json:"status"`
	InstitutionId     string `json:"institution_id"`
	Agreement         string `json:"agreement"`
	UserLanguage      string `json:"user_language"`
	Link              string `json:"link"`
	Ssn               string `json:"ssn"`
	AccountSelection  bool   `json:"account_selection"`
	RedirectImmediate bool   `json:"redirect_immediate"`
}

type accountTransactions struct {
	Booked  []transaction `json:"booked"`
	Pending []transaction `json:"pending"`
}

type transaction struct {
	TransactionId                     string            `json:"transactionId"`
	DebtorName                        string            `json:"debtorName"`
	DebtorAccount                     debtorAccount     `json:"debtorAccount"`
	TransactionAmount                 transactionAmount `json:"transactionAmount"`
	BankTransactionCode               string            `json:"bankTransactionCode"`
	BookingDate                       string            `json:"bookingDate"`
	ValueDate                         string            `json:"valueDate"`
	RemittanceInformationUnstructured string            `json:"remittanceInformationUnstructured"`
}

type transactionAmount struct {
	Currency string `json:"currency"`
	Amount   string `json:"amount"`
}

type debtorAccount struct {
	Iban string `json:"iban"`
}

package nordigen

import "fmt"

type TokenInfo struct {
	Access         string `json:"access"`
	AccessExpires  int    `json:"access_expires"`
	Refresh        string `json:"refresh"`
	RefreshExpires int    `json:"refresh_expires"`
}

type Institution struct {
	Id                   string   `json:"id"`
	Name                 string   `json:"name"`
	Bic                  string   `json:"bic,omitempty"`
	TransactionTotalDays string   `json:"transaction_total_days,omitempty"`
	Countries            []string `json:"countries"`
	Logo                 string   `json:"logo"`
}

type RequisitionRequest struct {
	Redirect          string `json:"redirect"`
	InstitutionId     string `json:"institution_id"`
	Reference         string `json:"reference"`
	Agreement         string `json:"agreement,omitempty"`
	UserLanguage      string `json:"user_language"`
	Ssn               string `json:"ssn,omitempty"`
	AccountSelection  bool   `json:"account_selection,omitempty"`
	RedirectImmediate bool   `json:"redirect_immediate,omitempty"`
}

type Requisition struct {
	Id                string   `json:"id"`      // TODO: Decide if it makes sense to parse as UUID
	Created           string   `json:"created"` // TODO: Decide if it makes sense to parse as DateTime
	Redirect          string   `json:"redirect"`
	Status            string   `json:"status"`
	InstitutionId     string   `json:"institution_id"`
	Agreement         string   `json:"agreement"`
	Reference         string   `json:"reference"`
	Accounts          []string `json:"accounts"`
	UserLanguage      string   `json:"user_language"`
	Link              string   `json:"link"`
	Ssn               string   `json:"ssn,omitempty"`
	AccountSelection  bool     `json:"account_selection,omitempty"`
	RedirectImmediate bool     `json:"redirect_immediate,omitempty"`
}

type AccountTransactions struct {
	Booked  []Transaction `json:"booked"`
	Pending []Transaction `json:"pending"`
}

type TransactionsResponse struct {
	Transactions AccountTransactions `json:"transactions"`
}

type Transaction struct {
	TransactionId                     string            `json:"transactionId"`
	DebtorName                        string            `json:"debtorName"`
	DebtorAccount                     DebtorAccount     `json:"debtorAccount"`
	DebtorAgent                       string            `json:"debtorAgent"`
	CreditorAgent                     string            `json:"creditorAgent"`
	TransactionAmount                 TransactionAmount `json:"transactionAmount"`
	BankTransactionCode               string            `json:"bankTransactionCode"`
	ProprietaryBankTransactionCode    string            `json:"ProprietaryBankTransactionCode"`
	EntryReference                    string            `json:"entryReference"`
	BookingDate                       string            `json:"bookingDate"`
	ValueDate                         string            `json:"valueDate"`
	RemittanceInformationStructured   string            `json:"remittanceInformationStructured"`
	RemittanceInformationUnstructured string            `json:"remittanceInformationUnstructured"`
	AdditionalInformation             string            `json:"additionalInformation"`
}

type TransactionAmount struct {
	Currency string `json:"currency"`
	Amount   string `json:"amount"`
}

type DebtorAccount struct {
	Iban string `json:"iban"`
}

type ApiError struct {
	Summary    string `json:"summary"`
	Detail     string `json:"detail"`
	StatusCode int    `json:"status_code"`
	Raw        string `json:"raw"`
	Decodable  bool   `json:"decodeable"`
}

func (e *ApiError) Error() string {
	return fmt.Sprintf("api returned status code: %d", e.StatusCode)
}

func HandleApiError(err error) { // TODO: Move to helpers
	apiError, ok := err.(*ApiError)
	if ok {
		if apiError.Decodable {
			fmt.Println("Summary: ", apiError.Summary)
			fmt.Println("Detail: ", apiError.Detail)
			fmt.Println("StatusCode: ", apiError.StatusCode)
		} else {
			fmt.Println("Raw: ", apiError.Raw)
		}
	} else {
		fmt.Println(err)
	}
}

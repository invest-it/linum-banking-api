package main

type CreateRequisitionRequest struct {
	InstitutionId string `json:"institution_id"`
	RedirectUrl   string `json:"redirect_url"`
	UserLanguage  string `json:"user_language"`
}

type CreateRequisitionResponse struct {
	Link string `json:"link"`
}

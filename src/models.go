package main

type createRequisitionRequest struct {
	InstitutionId string `json:"institution_id"`
	RedirectUrl   string `json:"redirect_url"`
	UserLanguage  string `json:"user_language"`
}

type createRequisitionResponse struct {
	Link string `json:"link"`
}

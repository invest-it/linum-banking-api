package nordigen

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"linum-banking-api/nordigen/endpoints"
	"net/http"
)

const (
	redirect = "https://www.example.com"
)

func getToken() (*tokenInfo, error) {
	url := endpoints.UseEndpoint(endpoints.Token)

	data := map[string]string{"secret_id": "SECRET", "secret_key": "SECRET"}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		// TODO: Log server response
		return nil, fmt.Errorf("server returned status code %d", resp.StatusCode)
	}

	res := tokenInfo{}
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}

func getInstitutionsForCountry(countryCode string, token *tokenInfo) (*[]institution, error) {
	url := endpoints.UseEndpoint(endpoints.Institutions)
	endpoints.AddQuery(&url, "country", countryCode)

	resp, err := getWithAuthorization(countryCode, token)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned status code %d", resp.StatusCode)
	}

	var res []institution
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}

func createRequisition(inst *institution, userLanguage string, token *tokenInfo) (*requisition, error) {
	reference := uuid.New()
	requisitionReq := requisitionRequest{
		Redirect:      redirect,
		InstitutionId: inst.Id,
		Reference:     reference.String(),
		UserLanguage:  userLanguage,
	}

	jsonData, err := json.Marshal(requisitionReq)
	if err != nil {
		return nil, err
	}

	url := endpoints.UseEndpoint(endpoints.Requisitions)
	resp, err := postWithAuthorization(url, token, bytes.NewBuffer(jsonData), "application/json")
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("server returned status code %d", resp.StatusCode)
	}

	var res = requisition{}
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}

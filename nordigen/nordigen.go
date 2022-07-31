package nordigen

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"linum-banking-api/nordigen/endpoints"
	"net/http"
	"time"
)

var token *tokenInfo
var tokenExpiration int64

func GetToken() (*tokenInfo, error) {
	if token == nil {
		var err error
		token, err = getNewToken()
		if err != nil {
			return nil, err
		}
		tokenExpiration = int64(token.AccessExpires) + time.Now().Unix()
	} else if tokenExpiration <= time.Now().Unix() {
		var err error
		token, err = getNewToken()
		if err != nil {
			token = nil
			tokenExpiration = 0
			return nil, err
		}
		tokenExpiration = int64(token.AccessExpires) + time.Now().Unix()
	}
	return token, nil
}

func getNewToken() (*tokenInfo, error) {
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

	return mapResponseToStruct[tokenInfo](resp)
}

func GetInstitutionsForCountry(countryCode string, token *tokenInfo) (*[]institution, error) {
	url := endpoints.UseEndpoint(endpoints.Institutions)
	endpoints.AddQuery(&url, "country", countryCode)

	return getAndMapWithAuthorization[[]institution](url, token)
}

func CreateRequisition(institutionId string, userLanguage string, redirectUrl string, token *tokenInfo) (*requisition, error) {
	reference := uuid.New()
	requisitionReq := requisitionRequest{
		Redirect:      redirectUrl,
		InstitutionId: institutionId,
		Reference:     reference.String(),
		UserLanguage:  userLanguage,
	}

	jsonData, err := json.Marshal(requisitionReq)
	if err != nil {
		return nil, err
	}

	url := endpoints.UseEndpoint(endpoints.Requisitions)
	return postAndMapWithAuthorization[requisition](url, token, bytes.NewBuffer(jsonData), "application/json")
}

func GetRequisitionById(id string, token *tokenInfo) (*requisition, error) {
	url := endpoints.UseEndpoint(endpoints.Requisitions) + id
	return getAndMapWithAuthorization[requisition](url, token)
}

func GetTransactionsForAccountId(id string, token *tokenInfo) (*accountTransactions, error) {
	url := endpoints.UseEndpoint(endpoints.Accounts) + id + "/transactions/"
	return getAndMapWithAuthorization[accountTransactions](url, token)
}

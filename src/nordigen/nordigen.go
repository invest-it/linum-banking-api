package nordigen

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"linum-banking-api/src/nordigen/endpoints"
	"net/http"
	"os"
	"time"
)

var token *TokenInfo
var tokenExpiration int64

func GetToken() (*TokenInfo, error) {
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

func getNewToken() (*TokenInfo, error) {
	url := endpoints.UseEndpoint(endpoints.Token)

	data := map[string]string{"secret_id": os.Getenv("NORDIGEN_SECRET_ID"), "secret_key": os.Getenv("NORDIGEN_SECRET_KEY")}
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

	return mapResponseToStruct[TokenInfo](resp)
}

func GetInstitutionsForCountry(countryCode string, token *TokenInfo) (*[]Institution, error) {
	url := endpoints.UseEndpoint(endpoints.Institutions)
	endpoints.AddQuery(&url, "country", countryCode)

	return getAndMapWithAuthorization[[]Institution](url, token)
}

func CreateRequisition(institutionId string, userLanguage string, redirectUrl string, token *TokenInfo) (*Requisition, error) {
	reference := uuid.New()
	requisitionReq := RequisitionRequest{
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
	return postAndMapWithAuthorization[Requisition](url, token, bytes.NewBuffer(jsonData), "application/json")
}

func GetRequisitionById(id string, token *TokenInfo) (*Requisition, error) {
	url := endpoints.UseEndpoint(endpoints.Requisitions) + id
	return getAndMapWithAuthorization[Requisition](url, token)
}

func GetTransactionsForAccountId(id string, token *TokenInfo) (*TransactionsResponse, error) {
	url := endpoints.UseEndpoint(endpoints.Accounts) + id + "/transactions/"
	return getAndMapWithAuthorization[TransactionsResponse](url, token)
}

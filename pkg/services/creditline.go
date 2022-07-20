package services

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"job-throtling-go-poc/pkg/models"
	"net/http"
)

type CreditLineService struct{}

func NewCreditLineService() *CreditLineService {
	return &CreditLineService{}
}

func (a *CreditLineService) Post(address models.CreditLine) (*models.CreditLine, error) {
	body, err := json.Marshal(address)

	if err != nil {
		return nil, err
	}

	resp, err := http.Post("http://localhost:8080", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := &models.CreditLine{}
	err = json.Unmarshal(respBody, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

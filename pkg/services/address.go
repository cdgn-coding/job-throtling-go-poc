package services

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"job-throtling-go-poc/pkg/models"
	"net/http"
)

type AddressService struct{}

func NewAddressService() *AddressService {
	return &AddressService{}
}

func (a *AddressService) Post(address models.Address) (*models.Address, error) {
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

	result := &models.Address{}
	err = json.Unmarshal(respBody, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

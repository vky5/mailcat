package client

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/vky5/mailcat/internal/db/models"
)

type AccountClient struct {
	base *Client
}

// handle call to apis and return model pointer
func (a *AccountClient) CreateAcc(acc models.Account) (*models.Account, error) {
	data, _ := json.Marshal(acc)
	resp, err := a.base.HTTP.Post(
		a.base.BaseURL+"/account", "application/json", bytes.NewBuffer(data),
	)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("server error :%s", resp.Status)
	}

	var out models.Account
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil

}


// get the list of the accounts 
func (a *AccountClient) AccList() ([]models.Account, error) {
	resp, err := a.base.HTTP.Get(a.base.BaseURL + "/account")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var accs []models.Account
	if err := json.NewDecoder(resp.Body).Decode(&accs); err != nil {
		return nil, err
	}
	return accs, nil
}

package gapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

type User struct {
	Id      int64
	Email   string
	Name    string
	Login   string
	IsAdmin bool
}
type SwitchUserContextResponse struct {
	Message string
}

func (c *Client) Users() ([]User, error) {
	users := make([]User, 0)
	req, err := c.newRequest("GET", "/api/users", nil)
	if err != nil {
		return users, err
	}
	resp, err := c.Do(req)
	if err != nil {
		return users, err
	}
	if resp.StatusCode != 200 {
		return users, errors.New(resp.Status)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return users, err
	}
	err = json.Unmarshal(data, &users)
	if err != nil {
		return users, err
	}
	return users, err
}

func (c *Client) SwitchUserContext(Id int64) (SwitchUserContextResponse, error) {
	var message SwitchUserContextResponse
	path := fmt.Sprintf("/api/user/using/%d", Id)
	req, err := c.newRequest("POST", path, nil)
	if err != nil {
		return message, err
	}
	resp, err := c.Do(req)
	if err != nil {
		return message, err
	}
	switch resp.StatusCode {
	case 200:
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return message, err
		}
		err = json.Unmarshal(data, &message)
		if err != nil {
			return message, err
		}
	case 401:
		origErr := errors.New(resp.Status)
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return message, origErr
		}
		err = json.Unmarshal(data, &message)
		if err != nil {
			return message, origErr
		}
	default:
		return message, errors.New(resp.Status)
	}

	return message, err
}

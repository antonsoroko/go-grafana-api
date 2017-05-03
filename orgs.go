package gapi

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type Org struct {
	Id   int64
	Name string
}

type OrgResponse struct {
	Id      int64  `json:"orgId"`
	Message string `json:"message"`
}

func (c *Client) Orgs() ([]Org, error) {
	orgs := make([]Org, 0)

	req, err := c.newRequest("GET", "/api/orgs/", nil)
	if err != nil {
		return orgs, err
	}
	data, err := c.DoRead(req)
	if err != nil {
		return orgs, err
	}
	err = json.Unmarshal(data, &orgs)
	return orgs, err
}

func (c *Client) NewOrg(name string) (OrgResponse, error) {
	settings := map[string]string{
		"name": name,
	}
	result := OrgResponse{}
	postData, err := json.Marshal(settings)
	req, err := c.newRequest("POST", "/api/orgs/", bytes.NewBuffer(postData))
	if err != nil {
		return result, err
	}
	data, err := c.DoRead(req)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(data, &result)
	return result, err
}

func (c *Client) UpdateOrg(Id int64, name string) (OrgResponse, error) {
	settings := map[string]string{
		"name": name,
	}
	result := OrgResponse{}
	postData, err := json.Marshal(settings)
	path := fmt.Sprintf("/api/orgs/%d", Id)
	req, err := c.newRequest("PUT", path, bytes.NewBuffer(postData))
	if err != nil {
		return result, err
	}
	data, err := c.DoRead(req)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(data, &result)
	return result, err
}

func (c *Client) GetOrgByName(name string) (Org, error) {
	result := Org{}
	path := fmt.Sprintf("/api/orgs/name/%s", name)
	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return result, err
	}
	data, err := c.DoRead(req)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(data, &result)
	return result, err
}

func (c *Client) DeleteOrg(id int64) (OrgResponse, error) {
	result := OrgResponse{}
	req, err := c.newRequest("DELETE", fmt.Sprintf("/api/orgs/%d", id), nil)
	if err != nil {
		return result, err
	}
	data, err := c.DoRead(req)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(data, &result)
	return result, err
}

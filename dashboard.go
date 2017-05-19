package gapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path"

	log "github.com/Sirupsen/logrus"
)

type DashboardMeta struct {
	IsStarred bool   `json:"isStarred"`
	Slug      string `json:"slug"`
}

type DashboardSaveResponse struct {
	Slug    string `json:"slug"`
	Status  string `json:"status"`
	Version int64  `json:"version"`
	Message string `json:"message"`
}

type DashboardImportResponse struct {
	Message     string `json:"message"`
	Title       string `json:"title"`
	ImportedUri string `json:"importedUri"`
	Slug        string `json:"slug"`
}

type Dashboard struct {
	Meta  DashboardMeta          `json:"meta"`
	Model map[string]interface{} `json:"dashboard"`
}

type DashboardList []DashboardEntry

type DashboardEntry struct {
	Id        int      `json:"id"`
	Title     string   `json:"title"`
	URI       string   `json:"uri"`
	Type      string   `json:"type"`
	Tags      []string `json:"tags"`
	IsStarred bool     `json:"isStarred"`
}

func (c *Client) ListDashboards() (*DashboardList, error) {
	req, err := c.newRequest("GET", "/api/search", nil)
	if err != nil {
		return nil, err
	}
	d, err := c.DoRead(req)
	if err != nil {
		return nil, err
	}

	var dl DashboardList
	if err = json.Unmarshal(d, &dl); err != nil {
		return nil, err
	}
	return &dl, nil
}

func (c *Client) SaveDashboard(model map[string]interface{}, overwrite bool) (*DashboardSaveResponse, error) {
	wrapper := map[string]interface{}{
		"dashboard": model,
		"overwrite": overwrite,
	}
	data, err := json.Marshal(wrapper)
	if err != nil {
		return nil, err
	}
	req, err := c.newRequest("POST", "/api/dashboards/db", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := &DashboardSaveResponse{}
	if err = json.Unmarshal(data, &result); err != nil {
		return result, err
	}

	switch resp.StatusCode {
	case 200:
	case 400, 412:
		log.Error(resp.Status)
		return result, errors.New(result.Message)
	default:
		log.Error(resp.Status)
		return result, errors.New(resp.Status)
	}

	return result, err
}

func (c *Client) ImportDashboard(model map[string]interface{}, overwrite bool, inputs []interface{}) (*DashboardImportResponse, error) {
	wrapper := map[string]interface{}{
		"dashboard": model,
		"overwrite": overwrite,
		"inputs":    inputs,
	}
	data, err := json.Marshal(wrapper)
	if err != nil {
		return nil, err
	}
	req, err := c.newRequest("POST", "/api/dashboards/import", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := &DashboardImportResponse{}
	if err = json.Unmarshal(data, &result); err != nil {
		return result, err
	}

	switch resp.StatusCode {
	case 200:
	case 500:
		log.Error(resp.Status)
		return result, errors.New(result.Message)
	default:
		log.Error(resp.Status)
		return result, errors.New(resp.Status)
	}

	return result, err
}

// Dashboard fetches the dashboard with the given uri,
// and unmarshals it into a Dashboard structure.
// If uri does not contain a path then "db/" is prepended.
func (c *Client) Dashboard(uri string) (*Dashboard, error) {
	if path.Dir(uri) == "." {
		uri = path.Join("db", uri)
	}
	uri = path.Join("/api/dashboards", uri)

	req, err := c.newRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := &Dashboard{}
	err = json.Unmarshal(data, &result)
	return result, err
}

func (c *Client) DeleteDashboard(slug string) error {
	path := fmt.Sprintf("/api/dashboards/db/%s", slug)
	req, err := c.newRequest("DELETE", path, nil)
	if err != nil {
		return err
	}

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}

	return nil
}

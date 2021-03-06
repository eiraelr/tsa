package client

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/juliengk/stack/client"
	"github.com/juliengk/stack/jsonapi"
	"github.com/kassisol/tsa/api/types"
)

// Get CA public Key
func (c *Config) GetCACertificate() ([]byte, error) {
	cc := &client.Config{
		Scheme: c.URL.Scheme,
		Host:   c.URL.Host,
		Port:   c.URL.Port,
		Path:   c.Directory.CAInfo,
	}

	req, err := client.New(cc)
	if err != nil {
		return nil, err
	}

	req.HeaderAdd("Accept", "application/json")

	result := req.Get()
	if result.Error != nil {
		return nil, result.Error
	}

	var response jsonapi.Response
	if err := json.Unmarshal(result.Body, &response); err != nil {
		return nil, err
	}

	if result.Response.StatusCode != 200 {
		if response.Errors == (jsonapi.ResponseMessage{}) {
			return nil, fmt.Errorf("Could not fetch CA public key")
		}

		return nil, fmt.Errorf(response.Errors.Message)
	}

	info := GetReflectStringValue(response.Data)

	return []byte(info), nil
}

// Get Certificate
func (c *Config) GetCertificate(token string, certType string, csr []byte, duration int) ([]byte, error) {
	cc := &client.Config{
		Scheme: c.URL.Scheme,
		Host:   c.URL.Host,
		Port:   c.URL.Port,
		Path:   c.Directory.NewApp,
	}

	req, err := client.New(cc)
	if err != nil {
		return nil, err
	}

	req.HeaderAdd("Accept", "application/json")
	req.HeaderAdd("Content-Type", "application/json")
	req.HeaderAdd("Authorization", fmt.Sprintf("Bearer %s", token))

	newcert := types.NewCert{
		Type:     certType,
		CSR:      csr,
		Duration: duration,
	}

	data, err := json.Marshal(newcert)
	if err != nil {
		return nil, err
	}

	result := req.Post(bytes.NewBuffer(data))
	if result.Error != nil {
		return nil, result.Error
	}

	var response jsonapi.Response
	if err := json.Unmarshal(result.Body, &response); err != nil {
		return nil, err
	}

	if result.Response.StatusCode != 200 {
		return nil, fmt.Errorf(response.Errors.Message)
	}

	rc := GetReflectStringValue(response.Data)

	return []byte(rc), nil
}

// Revoke Certificate
func (c *Config) RevokeCertificate(token string, serialNumber int) error {
	revokecert := types.RevokeCert{
		SerialNumber: serialNumber,
	}

	data, err := json.Marshal(revokecert)
	if err != nil {
		return err
	}

	cc := &client.Config{
		Scheme: c.URL.Scheme,
		Host:   c.URL.Host,
		Port:   c.URL.Port,
		Path:   c.Directory.RevokeCert,
	}

	req, err := client.New(cc)
	if err != nil {
		return err
	}

	req.HeaderAdd("Accept", "application/json")
	req.HeaderAdd("Content-Type", "application/json")
	req.HeaderAdd("Authorization", fmt.Sprintf("Bearer %s", token))

	result := req.Post(bytes.NewBuffer(data))
	if result.Error != nil {
		return result.Error
	}

	var response jsonapi.Response
	if err := json.Unmarshal(result.Body, &response); err != nil {
		return err
	}

	if result.Response.StatusCode != 200 {
		return fmt.Errorf(response.Errors.Message)
	}

	return nil
}

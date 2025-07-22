package services

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

// EnviarXMLAlPAC envía el XML firmado al PAC y regresa el XML timbrado
func EnviarXMLAlPAC(xmlFirmado []byte, usuarioPAC, clavePAC, endpoint string) ([]byte, error) {
	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(xmlFirmado))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(usuarioPAC, clavePAC)
	req.Header.Set("Content-Type", "application/xml")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

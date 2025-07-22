package pac

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

// TimbrarConPAC timbra el XML usando el API de Solución Factible
func TimbrarConPAC(xmlCFDI, rfc, pacClave, endpoint string) ([]byte, error) {
	payload := map[string]interface{}{
		"usuario":    rfc,
		"password":   pacClave,
		"produccion": "SI",
		"cfdi":       base64.StdEncoding.EncodeToString([]byte(xmlCFDI)),
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	request, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Error PAC: " + string(body))
	}

	var res map[string]interface{}
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}
	cfdiTimbrado64, ok := res["cfdi"].(string)
	if !ok || cfdiTimbrado64 == "" {
		return nil, errors.New("No se recibió XML timbrado del PAC")
	}
	cfdiTimbrado, err := base64.StdEncoding.DecodeString(cfdiTimbrado64)
	if err != nil {
		return nil, err
	}
	return cfdiTimbrado, nil
}

package monoliticComunication

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/xubiosueldos/conexionBD/Autenticacion/structAutenticacion"
)

func reqMonolitico(w http.ResponseWriter, r *http.Request, tokenAutenticacion *structAutenticacion.Security, codigo string, id string, options string, url string) string {
	var strReqMonolitico strRequestMonolitico
	token := *tokenAutenticacion
	strReqMonolitico.Options = options
	strReqMonolitico.Tenant = token.Tenant
	strReqMonolitico.Token = token.Token
	strReqMonolitico.Username = token.Username
	strReqMonolitico.Id = id

	pagesJson, err := json.Marshal(strReqMonolitico)
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	fmt.Println("URL:>", url)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(pagesJson))

	if err != nil {
		fmt.Println("Error: ", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error: ", err)
	}

	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("Error: ", err)
	}

	str := string(body)
	fmt.Println("BYTES RECIBIDOS :", len(str))

	return str
}

func Obtenercentrodecosto(w http.ResponseWriter, r *http.Request, tokenAutenticacion *structAutenticacion.Security, codigo string, id string, options string, url string) string {
	str := reqMonolitico(w, r, tokenAutenticacion, codigo, id, options, url)
	return str
}

func ChequeoAuthenticationMonolitico(tokenEncode string, url string, r *http.Request) bool {

	infoUserValida := false
	var prueba []byte = []byte("xubiosueldosimplementadocongo")
	tokenSecurity := base64.StdEncoding.EncodeToString(prueba)

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println("Error: ", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	req.Header.Add("Authorization", tokenEncode)
	req.Header.Add("SecurityToken", tokenSecurity)

	client := &http.Client{}

	res, err := client.Do(req)

	if err != nil {
		fmt.Println("Error: ", err)
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusAccepted {
		infoUserValida = true
	}

	return infoUserValida
}

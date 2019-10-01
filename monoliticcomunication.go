package monoliticComunication

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/xubiosueldos/conexionBD/Autenticacion/structAutenticacion"
	"github.com/xubiosueldos/conexionBD/Helper/structHelper"
	"github.com/xubiosueldos/conexionBD/Legajo/structLegajo"
	"github.com/xubiosueldos/framework"
	"github.com/xubiosueldos/framework/configuracion"
)

type requestMono struct {
	Value interface{}
	Error error
}

func reqMonolitico(w http.ResponseWriter, r *http.Request, tokenAutenticacion *structAutenticacion.Security, codigo string, id string, options string) string {
	var strReqMonolitico strRequestMonolitico
	token := *tokenAutenticacion
	strReqMonolitico.Options = options
	strReqMonolitico.Tenant = token.Tenant
	strReqMonolitico.Token = token.Token
	strReqMonolitico.Username = token.Username
	strReqMonolitico.Id = id

	pagesJson, err := json.Marshal(strReqMonolitico)
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	url := configuracion.GetUrlMonolitico() + codigo + "GoServlet"
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

func Obtenercentrodecosto(w http.ResponseWriter, r *http.Request, tokenAutenticacion *structAutenticacion.Security, id string) structLegajo.Centrodecosto {
	var centroDeCosto structLegajo.Centrodecosto
	str := reqMonolitico(w, r, tokenAutenticacion, "centrodecosto", id, "CANQUERY")
	json.Unmarshal([]byte(str), &centroDeCosto)
	return centroDeCosto
}

func Checkexistecuenta(w http.ResponseWriter, r *http.Request, tokenAutenticacion *structAutenticacion.Security, id string) *requestMono {
	var s *requestMono

	str := reqMonolitico(w, r, tokenAutenticacion, "cuenta", id, "CANQUERY")
	if str == "0" {
		framework.RespondError(w, http.StatusNotFound, "Cuenta Inexistente")
		s.Error = errors.New("Cuenta Inexistente")
	}
	return s
}

func Checkexistebanco(w http.ResponseWriter, r *http.Request, tokenAutenticacion *structAutenticacion.Security, id string) *requestMono {
	var s *requestMono

	str := reqMonolitico(w, r, tokenAutenticacion, "banco", id, "CANQUERY")
	if str == "0" {
		framework.RespondError(w, http.StatusNotFound, "Banco Inexistente")
		s.Error = errors.New("Banco Inexistente")
	}
	return s
}

func Obtenercodigohelper(w http.ResponseWriter, r *http.Request, tokenAutenticacion *structAutenticacion.Security, codigo string, id string) *requestMono {

	var s *requestMono
	str := reqMonolitico(w, r, tokenAutenticacion, codigo, id, "HLP")

	fixUtf := func(r rune) rune {
		if r == utf8.RuneError {
			return -1
		}
		return r
	}

	var dataStruct []structHelper.Helper
	json.Unmarshal([]byte(strings.Map(fixUtf, str)), &dataStruct)

	framework.RespondJSON(w, http.StatusOK, dataStruct)

	return s
}

func CheckAuthenticationMonolitico(tokenEncode string, r *http.Request) bool {

	infoUserValida := false
	var prueba []byte = []byte("xubiosueldosimplementadocongo")
	tokenSecurity := base64.StdEncoding.EncodeToString(prueba)

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	url := configuracion.GetUrlMonolitico() + "SecurityAuthenticationGo"
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

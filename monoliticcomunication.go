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

	"github.com/xubiosueldos/conexionBD/Autenticacion/structAutenticacion"
	"github.com/xubiosueldos/conexionBD/Helper/structHelper"
	"github.com/xubiosueldos/conexionBD/Legajo/structLegajo"
	"github.com/xubiosueldos/conexionBD/Liquidacion/structLiquidacion"
	"github.com/xubiosueldos/framework"
	"github.com/xubiosueldos/framework/configuracion"
)

type requestMono struct {
	Value interface{}
	Error error
}

type strCuentaImporte struct {
	Cuentaid      int     `json:"cuentaid"`
	Importecuenta float32 `json:"importecuenta"`
}

type strContabilizarDescontabilizar struct {
	requestMonolitico          strRequestMonolitico
	Descripcion                string             `json:"descripcion"`
	Cuentasimportes            []strCuentaImporte `json:"cuentasimportes"`
	Asientomanualtransaccionid int                `json:"asientomanualtransaccionid"`
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

func Obtenercentrodecosto(w http.ResponseWriter, r *http.Request, tokenAutenticacion *structAutenticacion.Security, id string) *structLegajo.Centrodecosto {
	var centroDeCosto structLegajo.Centrodecosto
	str := reqMonolitico(w, r, tokenAutenticacion, "centrodecosto", id, "CANQUERY")
	json.Unmarshal([]byte(str), &centroDeCosto)
	return &centroDeCosto
}

func Checkexistecentrodecosto(w http.ResponseWriter, r *http.Request, tokenAutenticacion *structAutenticacion.Security, id string) *requestMono {
	var s requestMono
	centroDeCosto := Obtenercentrodecosto(w, r, tokenAutenticacion, id)
	if centroDeCosto.ID == 0 {
		s.Error = errors.New("El Centro de costo con ID: " + id + " no existe")
	}

	return &s
}

func Checkexistecuenta(w http.ResponseWriter, r *http.Request, tokenAutenticacion *structAutenticacion.Security, id string) *requestMono {
	var s requestMono

	str := reqMonolitico(w, r, tokenAutenticacion, "cuenta", id, "CANQUERY")
	if str == "0" {
		s.Error = errors.New("La Cuenta con ID: " + id + " no existe")
	}
	return &s
}

func Obtenerbanco(w http.ResponseWriter, r *http.Request, tokenAutenticacion *structAutenticacion.Security, id string) *structLiquidacion.Banco {
	var banco structLiquidacion.Banco
	str := reqMonolitico(w, r, tokenAutenticacion, "banco", id, "CANQUERY")
	json.Unmarshal([]byte(str), &banco)
	return &banco
}

func Checkexistebanco(w http.ResponseWriter, r *http.Request, tokenAutenticacion *structAutenticacion.Security, id string) *requestMono {
	var s requestMono
	banco := Obtenerbanco(w, r, tokenAutenticacion, id)
	if banco.ID == 0 {
		s.Error = errors.New("El Banco con ID: " + id + " no existe")
	}

	return &s
}

func Gethelpers(w http.ResponseWriter, r *http.Request, tokenAutenticacion *structAutenticacion.Security, codigo string, id string) *requestMono {

	var s requestMono
	str := reqMonolitico(w, r, tokenAutenticacion, codigo, id, "HLP")

	var dataHelper []structHelper.Helper
	json.Unmarshal([]byte(str), &dataHelper)
	framework.RespondJSON(w, http.StatusOK, dataHelper)

	return &s
}

func Obtenerdatosempresa(w http.ResponseWriter, r *http.Request, tokenAutenticacion *structAutenticacion.Security, codigo string, id string) *requestMono {
	var emp requestMono
	str := reqMonolitico(w, r, tokenAutenticacion, codigo, id, "CANQUERY")

	var dataEmpresa structHelper.Empresa
	json.Unmarshal([]byte(str), &dataEmpresa)

	framework.RespondJSON(w, http.StatusOK, dataEmpresa)

	return &emp

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

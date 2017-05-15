package cgp

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"net/http"
)

type envelope struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	Body    body     `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`
}

type body struct {
	Response response `xml:"response"`
}

type response struct {
	Object interface{} `xml:"object"`
}

// SOAPNotFoundError is returned when we cannot find a soap resource.
type SOAPNotFoundError struct{}

func (err SOAPNotFoundError) Error() string {
	return "The requested resource could not be found"
}

func (cgp CGP) request(req, res interface{}) error {
	// Build request
	var b bytes.Buffer
	b.WriteString("<SOAP-ENV:Envelope xmlns:SOAP-ENV=\"http://schemas.xmlsoap.org/soap/envelope/\"><SOAP-ENV:Body>")
	err := xml.NewEncoder(&b).Encode(req)
	if err != nil {
		return err
	}
	b.WriteString("</SOAP-ENV:Body></SOAP-ENV:Envelope>")
	httpReq, err := http.NewRequest("POST", fmt.Sprintf("https://%s/CLI/", cgp.url), &b)
	httpReq.SetBasicAuth(cgp.user, cgp.pass)
	httpReq.Header.Set("Content-Type", "application/soap+xml")

	// Send request
	httpClient := http.Client{}
	httpRes, err := httpClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer httpRes.Body.Close()

	// Unmarshal response
	if httpRes.StatusCode == http.StatusInternalServerError {
		return SOAPNotFoundError{}
	}
	envel := envelope{Body: body{Response: response{Object: res}}}
	err = xml.NewDecoder(httpRes.Body).Decode(&envel)
	if err != nil {
		return err
	}
	return nil
}

type valueList struct {
	SubValues []string `xml:"subValue"`
}

// compact removes empty strings from the value list
func (vl valueList) compact() []string {
	return compact(vl.SubValues)
}

func compact(src []string) []string {
	res := make([]string, 0, len(src))
	for _, s := range src {
		if s == "" {
			continue
		}
		res = append(res, s)
	}
	return res
}

type keyValuePair struct {
	Key   string `xml:"key,attr"`
	Value string `xml:",innerxml"`
}

type dictionary struct {
	KeyValuePairs []keyValuePair `xml:"subKey"`
}

type dictionaryList struct {
	SubValues []dictionary `xml:"subValue"`
}

func (d dictionary) toMap() map[string]string {
	m := make(map[string]string, len(d.KeyValuePairs))
	for _, kvp := range d.KeyValuePairs {
		m[kvp.Key] = kvp.Value
	}
	return m
}

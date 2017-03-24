package cgp

import "encoding/xml"

// Domain represents a single domain
type Domain struct {
	cgp  *CGP
	Name string
}

type getDomainSettings struct {
	XMLName xml.Name `xml:"getDomainSettings"`
	Domain  string   `xml:"param"`
}

// Exists returns true if the domain
func (dom Domain) Exists() (bool, error) {
	var d dictionary
	err := dom.cgp.request(getDomainSettings{Domain: dom.Name}, &d)
	if _, ok := err.(SOAPNotFoundError); ok {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

type listDomains struct {
	XMLName xml.Name `xml:"listDomains"`
}

// Domain creates a domain type with the given name
func (cgp *CGP) Domain(name string) *Domain {
	return &Domain{cgp: cgp, Name: name}
}

// Domains lists the domains on the server
func (cgp *CGP) Domains() ([]*Domain, error) {
	var vl valueList
	err := cgp.request(listDomains{}, &vl)
	if err != nil {
		return []*Domain{}, err
	}
	vals := vl.SubValues
	ds := make([]*Domain, len(vals))
	for i, d := range vals {
		ds[i] = cgp.Domain(d)
	}
	return ds, nil
}

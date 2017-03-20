package cgp

import "encoding/xml"

// Domain represents a single domain
type Domain struct {
	cgp  *CGP
	Name string
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

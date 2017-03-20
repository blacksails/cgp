package cgp

import (
	"encoding/xml"
	"fmt"
)

// Forwarder represents a forwarder
type Forwarder struct {
	Domain *Domain
	Name   string
	To     string
}

// Forwarder returns a forwarder type with the given from and dest
func (dom *Domain) Forwarder(name, to string) *Forwarder {
	return &Forwarder{Domain: dom, Name: name, To: to}
}

type listForwarders struct {
	XMLName xml.Name `xml:"listForwarders"`
	Param   string   `xml:"param"`
}

// Forwarders lists the forwarders of a domain
func (dom *Domain) Forwarders() ([]*Forwarder, error) {
	var vl valueList
	err := dom.cgp.request(listForwarders{Param: dom.Name}, &vl)
	if err != nil {
		return []*Forwarder{}, err
	}
	vals := vl.compact()
	fs := make([]*Forwarder, len(vals))
	for i, v := range vals {
		f, err := dom.GetForwarder(v)
		if err != nil {
			return fs, err
		}
		fs[i] = f
	}
	return fs, err
}

type getForwarder struct {
	XMLName xml.Name `xml:"getForwarder"`
	Param   string   `xml:"param"`
}

// GetForwarder retreives a forwarder with the given name
func (dom *Domain) GetForwarder(name string) (*Forwarder, error) {
	var f string
	err := dom.cgp.request(getForwarder{Param: fmt.Sprintf("%s@%s", name, dom.Name)}, &f)
	if err != nil {
		return &Forwarder{}, err
	}
	return &Forwarder{Domain: dom, Name: name, To: f}, nil
}

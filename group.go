package cgp

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
)

// Group represents a group of a domain
type Group struct {
	Domain  *Domain
	Name    string
	Members []*Account
}

// Group creates a Group type from a domain, with the given name and members
func (dom *Domain) Group(name string, members []*Account) *Group {
	return &Group{Domain: dom, Name: name, Members: members}
}

type listGroups struct {
	XMLName xml.Name `xml:"listGroups"`
	Domain  string   `xml:"param"`
}

// Groups lists the groups of a domain
func (dom *Domain) Groups() ([]*Group, error) {
	var vl valueList
	err := dom.cgp.request(listGroups{Domain: dom.Name}, &vl)
	if err != nil {
		return []*Group{}, err
	}
	vals := vl.compact()
	grps := make([]*Group, len(vals))
	for i, v := range vals {
		g, err := dom.GetGroup(v)
		if err != nil {
			return grps, err
		}
		grps[i] = g
	}
	return grps, nil
}

type getGroup struct {
	XMLName xml.Name `xml:"getGroup"`
	Name    string   `xml:"param"`
}

// GetGroup retrieves a group from a domain with the given group name
func (dom *Domain) GetGroup(name string) (*Group, error) {
	var d dictionary
	err := dom.cgp.request(getGroup{Name: fmt.Sprintf("%s@%s", name, dom.Name)}, &d)
	if err != nil {
		return &Group{}, err
	}
	memStr := d.toMap()["Members"]
	var mems []*Account
	dec := xml.NewDecoder(bytes.NewBufferString(memStr))
	for {
		var a string
		err := dec.Decode(&a)
		if err == io.EOF {
			break
		}
		if err != nil {
			return dom.Group(name, mems), err
		}
		if a == "" {
			continue
		}
		mems = append(mems, dom.Account(a))
	}
	return dom.Group(name, mems), nil
}

package cgp

import "encoding/xml"

// MailingList represents a malinglist of a domain
type MailingList struct {
	Domain *Domain
	Name   string
}

// MailingList creates a MailingList type from a domain, with the given name
func (dom *Domain) MailingList(name string) *MailingList {
	return &MailingList{Domain: dom, Name: name}
}

type listLists struct {
	XMLName xml.Name `xml:"listLists"`
	Domain  string   `xml:"param"`
}

// MailingLists lists the mailing lists of a domain
func (dom *Domain) MailingLists() ([]*MailingList, error) {
	var vl valueList
	err := dom.cgp.request(listLists{Domain: dom.Name}, &vl)
	if err != nil {
		return []*MailingList{}, err
	}
	vals := vl.compact()
	mls := make([]*MailingList, len(vals))
	for i, v := range vals {
		mls[i] = dom.MailingList(v)
	}
	return mls, nil
}

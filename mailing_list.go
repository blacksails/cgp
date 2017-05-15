package cgp

import (
	"encoding/xml"
	"fmt"
)

// MailingList represents a malinglist of a domain
type MailingList struct {
	Domain *Domain
	Name   string
}

// MailingList creates a MailingList type from a domain, with the given name
func (dom *Domain) MailingList(name string) *MailingList {
	return &MailingList{Domain: dom, Name: name}
}

// Subscriber represents a subscription to a mailinglist
type Subscriber struct {
	MailingList *MailingList
	Email       string
	RealName    string
}

// Subscriber create a Subscriber type from a MalingList, with the given email
// and name
func (ml *MailingList) Subscriber(email, name string) *Subscriber {
	return &Subscriber{MailingList: ml, Email: email, RealName: name}
}

type readSubscribers struct {
	XMLName xml.Name `xml:"readSubscribers"`
	Name    string   `xml:"param"`
}

type readSubscribersResponse struct {
	SubValues []dictionaryList `xml:"subValue"`
}

// Subscribers returns a list of subscriber of a mailing list.
func (ml *MailingList) Subscribers() ([]*Subscriber, error) {
	var res readSubscribersResponse
	err := ml.Domain.cgp.request(readSubscribers{Name: fmt.Sprintf("%s@%s", ml.Name, ml.Domain.Name)}, &res)
	if err != nil {
		return []*Subscriber{}, err
	}
	ds := res.SubValues[1].SubValues
	subs := make([]*Subscriber, len(ds))
	for i, d := range ds {
		m := d.toMap()
		subs[i] = ml.Subscriber(m["Sub"], m["RealName"])
	}
	return subs, nil
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

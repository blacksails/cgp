package cgp

import (
	"encoding/xml"
	"fmt"
)

// Account represents an account under a domain
type Account struct {
	Domain *Domain
	Name   string
}

type getAccountSettings struct {
	XMLName xml.Name `xml:"getAccountSettings"`
	Account string   `xml:"param"`
}

// RealName return the real name of the account as registered
func (a Account) RealName() (string, error) {
	var d dictionary
	err := a.Domain.cgp.request(getAccountSettings{Account: a.Email()}, &d)
	if err != nil {
		return "", err
	}
	return d.toMap()["RealName"], nil
}

// Email returns the primary email of the account
func (a Account) Email() string {
	return fmt.Sprintf("%s@%s", a.Name, a.Domain.Name)
}

// Account returns an account type with the given name
func (dom *Domain) Account(name string) *Account {
	return &Account{Domain: dom, Name: name}
}

type listAccounts struct {
	XMLName xml.Name `xml:"listAccounts"`
	Domain  string   `xml:"param"`
}

type accountList struct {
	SubKeys []accountKey `xml:"subKey"`
}

type accountKey struct {
	Name string `xml:"key,attr"`
}

// Accounts lists the acounts of a domain
func (dom *Domain) Accounts() ([]*Account, error) {
	var al accountList
	err := dom.cgp.request(listAccounts{Domain: dom.Name}, &al)
	if err != nil {
		return []*Account{}, err
	}
	keys := al.SubKeys
	as := make([]*Account, len(keys))
	for i, k := range keys {
		as[i] = dom.Account(k.Name)
	}
	return as, nil
}

package cgp

import (
	"encoding/xml"
	"fmt"
)

// Alias represents an alias of an account
type Alias struct {
	account *Account
	Name    string
}

// Alias creates an Alias type from an account
func (acc *Account) Alias(name string) *Alias {
	return &Alias{account: acc, Name: name}
}

type listAliases struct {
	XMLName xml.Name `xml:"getAccountAliases"`
	Param   string   `xml:"param"`
}

// Aliases lists the aliases of an account
func (acc *Account) Aliases() ([]*Alias, error) {
	var vl valueList
	err := acc.Domain.cgp.request(listAliases{Param: fmt.Sprintf("%s@%s", acc.Name, acc.Domain.Name)}, &vl)
	if err != nil {
		return []*Alias{}, err
	}
	vals := vl.compact()
	as := make([]*Alias, len(vals))
	for i, v := range vals {
		as[i] = acc.Alias(v)
	}
	return as, nil
}

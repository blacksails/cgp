package cgp

// CGP represents a configured go client
type CGP struct {
	url  string
	user string
	pass string
}

// New returns a new Communigate Pro client
func New(url, user, pass string) *CGP {
	return &CGP{url: url, user: user, pass: pass}
}

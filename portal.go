package jcdportal

const (
	baseUrl = "https://api.jcdecaux.com/vls/v3/"
)

type APIRequester string

func (r APIRequester) Refresh(d jcdData) {}
func (r APIRequester) Execute(d jcdData) {}

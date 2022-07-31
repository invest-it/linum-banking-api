package endpoints

const (
	webAPI       = "https://ob.nordigen.com/api/v2"
	Token        = "/token/new/"
	Institutions = "/institutions/"
	Requisitions = "/requisitions/"
)

type UrlStr string

func UseEndpoint(endpoint string) string {
	return webAPI + endpoint
}

func AddQuery(url *string, param string, value string) {
	lastChar := (*url)[len(*url)-1]
	if lastChar == '/' {
		*url += "?"
	} else {
		*url += "&"
	}
	*url += param + "=" + value
}

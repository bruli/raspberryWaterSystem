package ws

import "net/url"

func New(serverURL url.URL, cl HTTPClient, token string) Handlers {
	cli := client{
		cl:        cl,
		serverURL: serverURL,
		token:     token,
	}
	return Handlers{
		GetStatus:   GetStatus(cli),
		Weather:     nil,
		Logs:        nil,
		ExecuteZone: nil,
	}
}

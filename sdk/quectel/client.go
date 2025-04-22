package quectel

const (
	api = "http://api.quectel.com/openapi/router"
)

type Client struct {
	appKey    string
	appSecret string
}

func NewClient(key, secret string) *Client {
	return &Client{appKey: key, appSecret: secret}
}

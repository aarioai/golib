package quectel

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Response struct {
	ResultCode   int64  `json:"resultCode"`
	ErrorMessage string `json:"errorMessage"`
}

func (c *Client) getUrl(param url.Values) string {
	return fmt.Sprintf("%s?%s", api, param.Encode())
}

func (c *Client) post(param url.Values, result interface{}) error {
	if param == nil {
		param = url.Values{}
	}
	param.Add("appKey", c.appKey)
	param.Add("t", fmt.Sprintf("%d", time.Now().Unix()))
	param.Add("sign", c.sign(param))

	postUrl := c.getUrl(param)
	resp, err := http.Post(postUrl, "application/json", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	//default 5MB
	body, err := io.ReadAll(io.LimitReader(resp.Body, int64(5<<20)))
	if err != nil {
		return err
	}

	return json.Unmarshal(body, result)
}

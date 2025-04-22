package weixinpay

import (
	"crypto/rsa"
	"fmt"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
	"path"
)

func loadMchCertClientKey(certDir string) (*rsa.PrivateKey, error) {
	pemFile := path.Join(certDir, "apiclient_key.pem")
	pem, err := utils.LoadPrivateKeyWithPath(pemFile)
	if err != nil {
		return nil, fmt.Errorf("load private key %s failed %v ", pemFile, err.Error())
	}
	return pem, nil
}

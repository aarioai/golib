package weixinpay

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"github.com/aarioai/airis/aa"
	"github.com/aarioai/airis/pkg/afmt"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"sync"
	"time"
)

// Config 微信支付配置
//
//	https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_1_5.shtml
//
// https://github.com/wechatpay-apiv3/wechatpay-go?tab=readme-ov-file
// 通过基础下单接口中的请求参数“notify_url”设置，要求必须为https地址。请确保回调URL是外部可正常访问的，且不能携带后缀参数
type Config struct {
	Mchid string `json:"mchid"` // 商户号 e.g. 190000****
	// 商户APIv3密钥，调用APIv3的下载平台证书接口、处理回调通知中报文时，要通过该密钥来解密信息
	// API安全 -> 解密回调 -> APIv3密钥   自定义32字符
	MchApiV3Key string `json:"mch_api_v3_key"` // 商户APIv3密钥

	// 商户API证书 用于证实商户身份。APIv3 支付所有接口都需要使用
	// https://kf.qq.com/faq/161222NneAJf161222U7fARv.html
	// API安全 -> 验证商户身份 -> 商户API证书
	// 	- 商户API证书序列号：3775B6A45ACD588826D15E583A95F5DD********
	// 	- 商户API私钥：商户申请商户 API 证书时，会生成商户私钥，并保存在本地证书文件夹的文件 apiclient_key.pem
	// 按上面流程会生成一个zip文件，解压后文件夹内有：apiclient_cert.p12, apiclient_cert.pem, apiclient_key.pem 三个密钥
	MchCertSerial    string `json:"mch_cert_serial"` // 商户证书序列号
	MchCertDir       string `json:"mch_cert_dir"`    // 商户API私钥
	mchCertClientKey *rsa.PrivateKey

	// API安全 -> 验证微信支付身份 -> 平台证书 商户调用业务API后，微信支付回调会使用平台证书的私钥生成签名，商户需要使用平台证书的公钥验签。平台证书5年过期一次
	// 平台证书 --> 通过SDK就可以直接下载了。不用手动保存
	// 证书序列号：1EF6365B6FC23B5D9D31C1053BDEEABB8C******

	NotifyUrl string `json:"notify_url"` // 由prepay的时候，提供给微信。自己配置即可，不用在微信支付后台配置

}

func (c *Config) PrivateKey() *rsa.PrivateKey {
	return c.mchCertClientKey
}

type ConfigKey struct {
	MchIdKeyName         string
	MchCertSerialKeyName string
	MchApiV3KeyKeyName   string
	MchCertDirKeyName    string
	NotifyUrlKeyName     string
}
type Service struct {
	app *aa.App
	loc *time.Location

	Config           Config
	rsaNotifyHandler *notify.Handler
}

var (
	mtx sync.RWMutex
)

func validateMchConfig(c Config) error {
	if c.Mchid == "" {
		return errors.New("config Mchid is empty")
	}
	if c.MchCertSerial == "" {
		return errors.New("config MchCertSerial is empty")
	}
	if c.MchApiV3Key == "" {
		return errors.New("config MchApiV3Key is empty")
	}
	if c.MchCertDir == "" {
		return errors.New("config MchCertDir is empty")
	}
	if c.NotifyUrl == "" {
		return errors.New("config NotifyUrl is empty")
	}
	return nil
}

func ParseConfig(app *aa.App, key ConfigKey) Config {
	mchId, _ := app.Config.MustGetString(key.MchIdKeyName)
	mchCertSerial, _ := app.Config.MustGetString(key.MchCertSerialKeyName)
	mchApiKey, _ := app.Config.MustGetString(key.MchApiV3KeyKeyName)
	mchPrivatePem, _ := app.Config.MustGetString(key.MchCertDirKeyName)
	notifyUrl, _ := app.Config.MustGetString(key.NotifyUrlKeyName)
	return Config{
		Mchid:         mchId,
		MchCertSerial: mchCertSerial,
		MchApiV3Key:   mchApiKey,
		MchCertDir:    mchPrivatePem,
		NotifyUrl:     notifyUrl,
	}
}

var (
	services sync.Map
)

// New
// 一个支付mch账号，可以对应多个公众号、小程序、APP等
func New(app *aa.App, config Config) (*Service, error) {
	var s *Service
	sv, ok := services.Load(config.Mchid)
	if ok {
		if s, ok = sv.(*Service); ok && s != nil {
			return s, nil
		}
		services.Delete(config.Mchid)
	}
	err := validateMchConfig(config)
	if err != nil {
		return nil, err
	}

	config.mchCertClientKey, err = loadMchCertClientKey(config.MchCertDir)
	if err != nil {
		return nil, err
	}
	s = &Service{
		app:    app,
		loc:    app.Config.TimeLocation,
		Config: config,
	}
	services.LoadOrStore(config.Mchid, s)
	return s, nil
}

func (s *Service) NewError(msg string, a ...any) error {
	msg = afmt.Sprintf(msg, a...)
	return fmt.Errorf("weixin pay (mchid:%s): %s", s.Config.Mchid, msg)
}

func (s *Service) NewClient(ctx context.Context) (*core.Client, error) {
	c := s.Config
	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(c.Mchid, c.MchCertSerial, c.mchCertClientKey, c.MchApiV3Key),
	}
	return core.NewClient(ctx, opts...)
}

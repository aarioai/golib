package metrics

import "github.com/kataras/iris/v12"

type Party struct {
	p iris.Party
}

func New(p iris.Party) *Party {
	return &Party{p: p}
}

func (p *Party) WithAll() *Party {
	return p.WithPing().WithPrometheusMetrics()
}

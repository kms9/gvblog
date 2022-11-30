package base

import (
	"github.com/kms9/gvblog/conf"
	"github.com/kms9/gvblog/model"
)

type R map[string]interface{}

func NewR() R {
	mp := make(R, 0)
	mp["appjs"] = conf.AppJsUrl
	mp["appcss"] = conf.AppCssUrl
	mp["global"] = model.Gcfg()
	return mp
}

func (r R) Add(m map[string]interface{}) map[string]interface{} {
	for k, v := range m {
		r[k] = v
	}
	return r
}

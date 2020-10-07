package logging

import (
	"sync"

	"go.uber.org/zap"
)

var prodZapSingleton = new(sync.Once)
var prodZap *zap.Logger

func NewZapProd() (*zap.Logger, func(), error) {
	var err error
	prodZapSingleton.Do(func() {
		if prodZap == nil {
			l, e := zap.NewProduction()
			if e != nil {
				err = e
			}
			prodZap = l
		}
	})
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		_ = prodZap.Sync()
	}
	return prodZap, cleanup, nil
}

package scripts

import (
	"github.com/gembackend/gembackendlog"
	"sync"
)

var (
	log = gembackendlog.Log
	wg  sync.WaitGroup
)

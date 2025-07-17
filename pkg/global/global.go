package global

import (
	"github.com/valyala/fasthttp"

	"main/pkg/types"
)

var AccountsList []types.AccountData
var Clients []*fasthttp.Client
var Config *types.Settings
var Module string
var TargetProgress int64
var CurrentProgress int64 = 0

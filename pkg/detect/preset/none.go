package preset

import (
	"github.com/crissyfield/troll-a/pkg/detect"
)

// None is an emply list of secrets, access and refresh token detection rules.
var None = []detect.GitleaksRuleFunction{}

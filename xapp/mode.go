package xapp

const (
	ReleaseMode = "release"
	DebugMode   = "debug"
	TestMode    = "test"
)

var appMode string

func SetMode(mode string) {
	if !ValidMode(mode) {
		panic(`invalid application mode "` + mode + `"`)
	}
	appMode = mode
}

func Mode() string {
	return appMode
}

func ValidMode(mode string) bool {
	switch mode {
	case ReleaseMode, DebugMode, TestMode:
		return true
	default:
		return false
	}
}

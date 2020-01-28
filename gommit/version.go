package gommit

var appVersion = ""

// GetVersion return app version
func GetVersion() string {
	return "v" + appVersion
}

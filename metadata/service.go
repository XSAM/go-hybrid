package metadata

// AppInfo return program's info
func AppInfo() Info {
	return appInfo
}

// SetAppName set program's name
func SetAppName(name string) {
	appInfo.AppName = name
}

// AppName return program's name
func AppName() string {
	return appInfo.AppName
}

// RuntimeID return program's runtime id (uuid)
func RuntimeID() string {
	return appInfo.RuntimeID
}

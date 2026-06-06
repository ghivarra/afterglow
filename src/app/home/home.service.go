package home

import "ghivarra/afterglow/environment"

func loadAppData() AppData {
	return AppData{
		Environment: environment.APP_ENV,
		Name:        environment.APP_NAME,
		Version:     environment.APP_VERSION,
		Port:        environment.SERVER_PORT,
	}
}

package home

type AppData struct {
	Environment string `json:"environment"`
	Name        string `json:"name"`
	Version     string `json:"version"`
	Port        int    `json:"port"`
}

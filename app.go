package cli

const APP_KEY = "app"

type AppContext struct {
	Active   *Project
	Projects []*Project
	Services *ServicesView
	Console  *ConsoleView
	ConBuff  *ConsoleBuff
}

type ConsoleBuff struct {
	Project *Project
	Data    []string
}

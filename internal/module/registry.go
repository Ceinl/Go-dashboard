package module

var availableModules = []string{
	"linksaver",
	"placeholder",
	"kanban",
	"twitter",
	// "profile",
}

func GetAvailableModules() []string {
	return availableModules
}

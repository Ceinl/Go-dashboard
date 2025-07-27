package module

var availableModules = []string{
	"linksaver",
	"placeholder",
	"kanban",
	// "profile",
}

func GetAvailableModules() []string {
	return availableModules
}

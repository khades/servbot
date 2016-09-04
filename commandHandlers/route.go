package commandHandlers

// RouterStruct is struct for handling command to handler
type RouterStruct struct {
	routes map[string]CommandHandler
}

// Go returns command for route
func (router RouterStruct) Go(command string) CommandHandler {
	handler, found := router.routes[command]
	if found {
		return handler
	}
	return Custom
}

// Router aa
var Router = RouterStruct{routes: map[string]CommandHandler{
	"new":   New,
	"alias": Alias}}

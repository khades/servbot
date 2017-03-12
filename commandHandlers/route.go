package commandHandlers

// RouterStruct is struct for handling command to handler
type RouterStruct struct {
	routes map[string]CommandHandler
}

// Go returns from router to work with
func (router RouterStruct) Go(command string) CommandHandler {
	handler, found := router.routes[command]
	if found {
		return handler
	}
	return Custom
}

// Router explicitly maps input chat command to a handler
var Router = RouterStruct{routes: map[string]CommandHandler{}}

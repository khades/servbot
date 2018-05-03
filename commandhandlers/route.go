package commandhandlers

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
	return custom
}

// Router explicitly maps input chat command to a handler
var Router = RouterStruct{routes: map[string]CommandHandler{
	"new":   newCommand,
	"alias": aliasCommand,
	"subdaynew": subdayNew}}
//	"subdayEnd": subdayEnd

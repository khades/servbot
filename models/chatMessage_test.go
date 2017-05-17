package models

// func TestEmptyMessage(t *testing.T) {
// 	v := ChatMessage{"a", "a", "", true, true, time.Now()}
// 	isCommand, _ := v.isCommand()
// 	if isCommand == true {
// 		t.Error("It is not command")
// 	}
// }

// func TestEmptyCommandBody(t *testing.T) {
// 	v := ChatMessage{"a", "a", "!halp", true, true, time.Now()}
// 	isCommand, chatCommand := v.isCommand()
// 	if isCommand == false {
// 		t.Error("It is actually a command")
// 	}
// 	if chatCommand.command != "halp" {
// 		t.Error("Invalid unpacking of string")
// 	}
// 	if chatCommand.body != "" {
// 		t.Error("chatCommand  body should be empty")
// 	}
// }

// func TestParsing(t *testing.T) {
// 	v := ChatMessage{"a", "a", "!halp  mother fucker", true, true, time.Now()}
// 	isCommand, chatCommand := v.isCommand()
// 	if isCommand == false {
// 		t.Error("It is actually a command")
// 	}
// 	if chatCommand.command != "halp" {
// 		t.Error("Invalid unpacking of string")
// 	}
// 	if chatCommand.body != " mother fucker" {
// 		t.Error("chatCommand  body should be spacemother fucker")
// 	}
// }

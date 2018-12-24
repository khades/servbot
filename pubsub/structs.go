package pubsub

type twitchWSOutgoingMessage struct {
	Type  string          `json:"type"`
	Nonce string          `json:"nonce"`
	Data  authMessageData `json:"data"`
}
type authMessageData struct {
	AuthToken string   `json:"auth_token"`
	Topics    []string `json:"topics"`
}

type wsMessage struct {
	Type string `json:"type"`
	Data wsData `json:"data"`
}
type wsData struct {
	Topic   string `json:"topic"`
	Message string `json:"message"`
}

type moderationActionMessage struct {
	Data moderationActionData `json:"data"`
}

type moderationActionData struct {
	Type            string   `json:"type"`
	ModeratorAction string   `json:"moderation_action"`
	Args            []string `json:"args"`
	User            string   `json:"created_by"`
	UserID          string   `json:"created_by_user_id"`
	RecipientID     string   `json:"target_user_id"`
}

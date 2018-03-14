package l10n

// L10n stores all
type L10nItem struct {
	AliasCreationSuccess           string
	CantMuteModerator              string
	CommandCyclingIsForbidden      string
	EmptyCommandNameIsForbidden    string
	ReservedCommandNameIsForbidden string
	CommandCreationSuccess         string
	InvalidCommandTemplate         string
	SubdayEndSuccess               string
	SubdayEndNothingToClose        string
	SubdayCreationSuccess          string
	SubdayCreationGeneralError     string
	SubdayCreationAlreadyExists    string
	SubdayCreationPrefix           string
	SubdayVoteNoActiveSubday       string
	SubdayVoteYouReNotSub          string
	GameHistoryActivePrefix        string
}

var l10n = map[string]L10nItem{
	"en": L10nItem{
		AliasCreationSuccess:           "Alias creation: Success",
		CantMuteModerator:              "Can't mute moderator",
		CommandCyclingIsForbidden:      "Command creation: Command cycling is now allowed",
		EmptyCommandNameIsForbidden:    "Command creation: Empty command name is now allowed",
		ReservedCommandNameIsForbidden: "Command creation: Reserved command name usage is now allowed",
		CommandCreationSuccess:         "Command creation: Success",
		InvalidCommandTemplate:         "Command creation: Invalid Template",
		SubdayEndSuccess:               "Subday end: Success",
		SubdayEndNothingToClose:        "Subday end: Nothing to close",
		SubdayCreationSuccess:          "Subday creation: Success",
		SubdayCreationGeneralError:     "Subday creation: General error",
		SubdayCreationAlreadyExists:    "Subday creation: Active subday already exists",
		SubdayCreationPrefix:           "Subday, created at ",
		SubdayVoteNoActiveSubday:       "Subday vote: No active subdays",
		SubdayVoteYouReNotSub:          "Subday vote: You're not subscriber",
		GameHistoryActivePrefix:        "NOW"},
	"ru": L10nItem{
		AliasCreationSuccess:           "Создание алиaса: Успешно",
		CantMuteModerator:              "Модератора нельзя затаймаутить",
		CommandCyclingIsForbidden:      "Создание команды: Запрещено зацикливать команды",
		EmptyCommandNameIsForbidden:    "Создание команды: Запрещено создавать пустые команды",
		ReservedCommandNameIsForbidden: "Создание команды: Запрещено создавать команды для зарезервированных слов",
		CommandCreationSuccess:         "Создание команды: Успешно",
		InvalidCommandTemplate:         "Создание команды: Невалидный шаблон для команды",
		SubdayEndSuccess:               "Закрытие сабдея: Успешно",
		SubdayEndNothingToClose:        "Закрытие сабдея: Нечего закрывать",
		SubdayCreationSuccess:          "Создание сабдея: Успешно",
		SubdayCreationGeneralError:     "Создание сабдея: Неизвестная ошибка",
		SubdayCreationAlreadyExists:    "Создание сабдея: Существует незакрытый сабдей",
		SubdayCreationPrefix:           "Сабдей, созданый ",
		SubdayVoteNoActiveSubday:       "Голосование на сабдее: Нет открытых сабдеев",
		SubdayVoteYouReNotSub:          "Голосование на сабдее: Вы не подписчик",
		GameHistoryActivePrefix:        "CЕЙЧАС"}}

func GetL10n(lang string) *L10nItem {
	item, found := l10n[lang]
	if found == true {
		return &item
	}
	result, _ := l10n["en"]
	return &result

}

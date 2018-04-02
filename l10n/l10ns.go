package l10n

import (
	"fmt"
	"log"
	"math"
	"strings"
	"time"
)

// Item  stores all l10n strings used by bot
type Item struct {
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
	NotAFollower                   string
}

var l10n = map[string]Item{
	"en": Item{
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
		GameHistoryActivePrefix:        "NOW",
		NotAFollower:                   "Not a follower"},
	"ru": Item{
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
		GameHistoryActivePrefix:        "CЕЙЧАС",
		NotAFollower:                   "Не фолловер"}}

// GetL10n returns l10n object for specified lang
func GetL10n(lang string) *Item {
	item, found := l10n[lang]
	if found == true {
		return &item
	}
	result, _ := l10n["en"]
	return &result

}

//HumanizeDuration converts time.Duration to human-readable string
func HumanizeDuration(duration time.Duration, lang string) string {

	result := ""

	years := math.Floor(duration.Hours() / (24 * 365))
	log.Println(math.Floor(duration.Hours() / (24 * 365)))
	log.Println(years)
	log.Println(years - math.Floor(years/10)*10)
	days := math.Floor(duration.Hours()/24) - years*365
	hours := math.Floor(duration.Hours()) - days*24 - years*365*24
	minutes := float64(int(duration.Minutes() - math.Floor(duration.Minutes()/60)*60))

	if int64(years) > 0 {
		if lang == "ru" {
			switch int64(years - math.Floor(years/10)*10) {
			case 1:
				result = fmt.Sprintf("%d год", int64(years))
				break
			case 2, 3, 4:

				result = fmt.Sprintf("%d года", int64(years))
				break
			default:
				result = fmt.Sprintf("%d лет", int64(years))
			}
		} else {
			if int64(years) == 1 {
				result = "1 year"
			} else {
				result = fmt.Sprintf("%d years", int64(years))
			}
		}
	}

	if int(days) > 0 {
		if lang == "ru" {
			switch int(days - math.Floor(days/10)*10) {
			case 1:
				result = result + fmt.Sprintf(" %d день", int(days))
				break
			case 2, 3, 4:
				result = result + fmt.Sprintf(" %d дня", int(days))
				break
			default:
				result = result + fmt.Sprintf(" %d дней", int(days))
			}
		} else {
			if int(days) == 1 {
				result = result + " 1 day"
			} else {
				result = result + fmt.Sprintf(" %s days", int(days))
			}
		}
	}

	if int(hours) > 0 {
		if lang == "ru" {
			switch int(hours - math.Floor(hours/10)*10) {
			case 1:
				result = result + fmt.Sprintf(" %d час", int(hours))
				break
			case 2, 3, 4:
				result = result + fmt.Sprintf(" %d часа", int(hours))
				break
			default:
				result = result + fmt.Sprintf(" %d часов", int(hours))
			}
		} else {
			if int(days) == 1 {
				result = result + " 1 hour"
			} else {
				result = result + fmt.Sprintf(" %s hours", int(hours))
			}
		}
	}

	if int(minutes) > 0 {
		if lang == "ru" {
			switch int(minutes - math.Floor(minutes/10)*10) {
			case 1:
				result = result + fmt.Sprintf(" %d минуту", int(minutes))
				break
			case 2, 3, 4:
				result = result + fmt.Sprintf(" %d минуты", int(minutes))
				break
			default:
				result = result + fmt.Sprintf(" %d минут", int(minutes))
			}
		} else {
			if int(days) == 1 {
				result = result + " 1 minute"
			} else {
				result = result + fmt.Sprintf(" %s minutes", int(minutes))
			}
		}
	}
	return strings.TrimSpace(result)
}

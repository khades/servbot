package l10n

import (
	"fmt"
	"math"
	"strings"
	"time"
)

// Item  stores all l10n strings used by bot
type Item struct {
	AliasCreationSuccess               string
	CantMuteModerator                  string
	CommandCyclingIsForbidden          string
	EmptyCommandNameIsForbidden        string
	ReservedCommandNameIsForbidden     string
	CommandCreationSuccess             string
	InvalidCommandTemplate             string
	SubdayEndSuccess                   string
	SubdayEndNothingToClose            string
	SubdayCreationSuccess              string
	SubdayCreationGeneralError         string
	SubdayCreationAlreadyExists        string
	SubdayCreationPrefix               string
	SubdayVoteNoActiveSubday           string
	SubdayVoteYouReNotSub              string
	GameHistoryActivePrefix            string
	NotAFollower                       string
	InvalidValue                       string
	VolumeChangeInvalidValue           string
	VolumeChangeSuccess                string
	SongRequestYoutubeRestricted       string
	SongRequestTwitchRestricted        string
	SongRequestChannelRestricted       string
	SongRequestTagRestricted           string
	SongRequestOffline                 string
	SongRequestPlaylistIsFull          string
	SongRequestAlreadyInPlaylist       string
	SongRequestTooManyRequests         string
	SongRequestInvalidLink             string
	SongRequestNothingFound            string
	SongRequestInternalError           string
	SongRequestTooLong                 string
	SongRequestTooLittleViews          string
	SongRequestMoreDislikes            string
	SongRequestSuccess                 string
	SongRequestPulledYoutubeRestricted string
	SongRequestPulledTwitchRestricted  string
	SongRequestPulledChannelRestricted string
	SongRequestPulledTagRestricted     string
	SongRequestPulled                  string
	SongRequestNoRequests              string
	SongRequestNotAModerator           string
}

var l10n = map[string]Item{
	"en": Item{
		AliasCreationSuccess:               "Alias creation: Success",
		CantMuteModerator:                  "Can't mute moderator",
		CommandCyclingIsForbidden:          "Command creation: Command cycling is now allowed",
		EmptyCommandNameIsForbidden:        "Command creation: Empty command name is now allowed",
		ReservedCommandNameIsForbidden:     "Command creation: Reserved command name usage is now allowed",
		CommandCreationSuccess:             "Command creation: Success",
		InvalidCommandTemplate:             "Command creation: Invalid Template",
		SubdayEndSuccess:                   "Subday end: Success",
		SubdayEndNothingToClose:            "Subday end: Nothing to close",
		SubdayCreationSuccess:              "Subday creation: Success",
		SubdayCreationGeneralError:         "Subday creation: General error",
		SubdayCreationAlreadyExists:        "Subday creation: Active subday already exists",
		SubdayCreationPrefix:               "Subday, created at ",
		SubdayVoteNoActiveSubday:           "Subday vote: No active subdays",
		SubdayVoteYouReNotSub:              "Subday vote: You're not subscriber",
		GameHistoryActivePrefix:            "NOW",
		NotAFollower:                       "Not a follower",
		InvalidValue:                       "Invalid Value",
		VolumeChangeInvalidValue:           "Song Request Volume Change: Invalid volume",
		VolumeChangeSuccess:                "Song Request Volume Change: Success",
		SongRequestYoutubeRestricted:       "Song Request: Playback of track \"%s\" is restricted by YouTube",
		SongRequestTwitchRestricted:        "Song Request: Playback of track \"%s\" is restricted by TOS of Twitch",
		SongRequestChannelRestricted:       "Song Request: Playback of track \"%s\" is restricted on that channel",
		SongRequestTagRestricted:           "Song Request: Playback of track \"%s\" is restricted on that channel due to tag \"%s\"",
		SongRequestOffline:                 "Song Request: Stream is offline",
		SongRequestPlaylistIsFull:          "Song Request: Playlist is full",
		SongRequestAlreadyInPlaylist:       "Song Request: Track \"%s\" is already in the playlist",
		SongRequestTooManyRequests:         "Song Request: You have too many requests",
		SongRequestInvalidLink:             "Song Request: Invalid song request link",
		SongRequestNothingFound:            "Song Request: Nothing found by that link",
		SongRequestInternalError:           "Song Request: Internal error, try again later or contact Bot administrator",
		SongRequestTooLong:                 "Song Request: Track \"%s\" is too long",
		SongRequestTooLittleViews:          "Song Request: Track \"%s\" has not enough views",
		SongRequestMoreDislikes:            "Song Request: Track \"%s\" has more dislikes than likes",
		SongRequestSuccess:                 "Song Request: Track \"%s\", duration %s has been added to playlist",
		SongRequestPulledTwitchRestricted:  "Song Request: Track \"%s\" was pulled from playlist due to Twitch TOS restrictions",
		SongRequestPulledYoutubeRestricted: "Song Request: Track \"%s\" was pulled from playlist due to YouTube restrictions",
		SongRequestPulledChannelRestricted: "Song Request: Track \"%s\" was pulled from playlist and restrictedon that channel",
		SongRequestPulledTagRestricted:     "Song Request: Track \"%s\" was pulled from playlist due to adding restricted tag \"%s\"",
		SongRequestPulled:                  "Song Request: Track \"%s\" was pulled from playlist",
		SongRequestNoRequests:              "Song Request: Nothing to remove",
		SongRequestNotAModerator:           "Song request: You're not a moderator"},
	"ru": Item{
		AliasCreationSuccess:               "Создание алиaса: Успешно",
		CantMuteModerator:                  "Модератора нельзя затаймаутить",
		CommandCyclingIsForbidden:          "Создание команды: Запрещено зацикливать команды",
		EmptyCommandNameIsForbidden:        "Создание команды: Запрещено создавать пустые команды",
		ReservedCommandNameIsForbidden:     "Создание команды: Запрещено создавать команды для зарезервированных слов",
		CommandCreationSuccess:             "Создание команды: Успешно",
		InvalidCommandTemplate:             "Создание команды: Невалидный шаблон для команды",
		SubdayEndSuccess:                   "Закрытие сабдея: Успешно",
		SubdayEndNothingToClose:            "Закрытие сабдея: Нечего закрывать",
		SubdayCreationSuccess:              "Создание сабдея: Успешно",
		SubdayCreationGeneralError:         "Создание сабдея: Неизвестная ошибка",
		SubdayCreationAlreadyExists:        "Создание сабдея: Существует незакрытый сабдей",
		SubdayCreationPrefix:               "Сабдей, созданый ",
		SubdayVoteNoActiveSubday:           "Голосование на сабдее: Нет открытых сабдеев",
		SubdayVoteYouReNotSub:              "Голосование на сабдее: Вы не подписчик",

		GameHistoryActivePrefix:            "CЕЙЧАС",
		NotAFollower:                       "Не фолловер",
		InvalidValue:                       "Неверное значение",
		VolumeChangeInvalidValue:           "Смена громкости заказа песен: Неверное значение",
		VolumeChangeSuccess:                "Смена громкости заказа песен: Успешно",
		SongRequestYoutubeRestricted:       "Заказ песен: Воспроизведение трека \"%s\" запрещено сервисом YouTube",
		SongRequestTwitchRestricted:        "Заказ песен: Воспроизведение трека \"%s\" запрещено сервисом Twitch",
		SongRequestChannelRestricted:       "Заказ песен: Воспроизведение трека  \"%s\" запрещено на этом канале",
		SongRequestTagRestricted:           "Заказ песен: Воспроизведение трека  \"%s\" запрещено на этом канале изза тега \"%s\"",
		SongRequestOffline:                 "Заказ песен: Трансляция не запущена",
		SongRequestPlaylistIsFull:          "Заказ песен: Список заказов полон",
		SongRequestAlreadyInPlaylist:       "Заказ песен: Трек \"%s\" уже в списке заказов",
		SongRequestTooManyRequests:         "Заказ песен: Вы делаете слишком много заказов",
		SongRequestInvalidLink:             "Заказ песен: Неверная ссылка",
		SongRequestNothingFound:            "Заказ песен: По ссылке ничего не найдено",
		SongRequestInternalError:           "Заказ песен: Внутренняя ошибка, попробуйте позже или свяжитесь с Администратором Бота",
		SongRequestTooLong:                 "Заказ песен: Трек \"%s\" слишком длинный",
		SongRequestTooLittleViews:          "Заказ песен: Трек \"%s\" не имеет достаточного количества просмотров",
		SongRequestMoreDislikes:            "Заказ песен: Трек \"%s\" имеет больше дизлайков чем лайков",
		SongRequestSuccess:                 "Заказ песен: Трек \"%s\" длительностью %s добавлен в список заказов",
		SongRequestPulledTwitchRestricted:  "Заказ песен: Трек \"%s\" удалён из списка заказов изза ограничений сервиса Twitch",
		SongRequestPulledYoutubeRestricted: "Заказ песен: Трек \"%s\" удалён из списка заказов изза ограничений сервиса YouTube",
		SongRequestPulledChannelRestricted: "Заказ песен: Трек \"%s\" удалён из списка заказов и запрещён на этом канале",
		SongRequestPulledTagRestricted:     "Заказ песен: Трек \"%s\" удалён из списка заказов изза добавления запрещённого тега \"%s\"",
		SongRequestPulled:                  "Заказ песен: Трек \"%s\" удалён из списка заказов",
		SongRequestNoRequests:              "Заказ песен: Нечего удалять",
		SongRequestNotAModerator:           "Заказ песен: Вы не модератор"}}

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

	days := math.Floor(duration.Hours()/24) - years*365
	hours := math.Floor(duration.Hours()) - days*24 - years*365*24
	minutes := float64(int(duration.Minutes() - math.Floor(duration.Minutes()/60)*60))
	seconds := math.Floor(duration.Seconds() - math.Floor(duration.Minutes())*60)
	if int64(years) > 0 {

		if lang == "ru" {
			if int64(years) > 10 && int(years) < 20 {
				result = result + fmt.Sprintf(" %d лет", int64(years))
			} else {
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
			if int(days) > 10 && int(days) < 20 {
				result = result + fmt.Sprintf(" %d дней", int(days))
			} else {
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
			if int(hours) > 10 && int(hours) < 20 {
				result = result + fmt.Sprintf(" %d часов", int(hours))
			} else {
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
			}
		} else {
			if int(hours) == 1 {
				result = result + " 1 hour"
			} else {
				result = result + fmt.Sprintf(" %s hours", int(hours))
			}
		}
	}

	if int(minutes) > 0 {
		if lang == "ru" {
			if int(minutes) > 10 && int(minutes) < 20 {
				result = result + fmt.Sprintf(" %d минут", int(minutes))
			} else {
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
			}
		} else {
			if int(minutes) == 1 {
				result = result + " 1 minute"
			} else {
				result = result + fmt.Sprintf(" %s minutes", int(minutes))
			}
		}
	}

	if int(seconds) > 0 {
		if lang == "ru" {
			if int(seconds) > 10 && int(seconds) < 20 {
				result = result + fmt.Sprintf(" %d секунд", int(seconds))
			} else {
				switch int(seconds - math.Floor(seconds/10)*10) {
				case 1:
					result = result + fmt.Sprintf(" %d секунду", int(seconds))
					break
				case 2, 3, 4:
					result = result + fmt.Sprintf(" %d секунды", int(seconds))
					break
				default:
					result = result + fmt.Sprintf(" %d секунд", int(seconds))
				}
			}
		} else {
			if int(seconds) == 1 {
				result = result + " 1 second"
			} else {
				result = result + fmt.Sprintf(" %s seconds", int(seconds))
			}
		}
	}

	return strings.TrimSpace(result)
}

package songRequest

import (
	"strings"
	"unicode/utf8"

	"github.com/khades/servbot/utils"
)

func parseYoutubeLink(input string) (string, bool) {

	if utf8.RuneCountInString(input) == 11 {
		return input, true
	}
	if strings.Contains(input, "youtube.com/watch?") {
		result := ""
		params := strings.Split(input, "youtube.com/watch?")[1]
		paramsSplit := strings.Split(params, "&")
		for _, param := range paramsSplit {
			paramSplit := strings.Split(param, "=")
			if paramSplit[0] == "v" {
				result = utils.Short(paramSplit[1], 11)
				break
			}
		}
		return result, true
	}
	if strings.Contains(input, "youtube.com/v/") {
		return utils.Short(strings.Split(input, "youtube.com/v/")[1], 11), true

	}

	if strings.Contains(input, "youtu.be/") {
		return utils.Short(strings.Split(input, "youtu.be/")[1], 11), true
	}

	return input, false
}

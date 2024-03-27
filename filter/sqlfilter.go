package filter

import (
	"log"
	"regexp"
	"strings"
)

var reStr = `(?:')|(?:\%)|(?:\\)|(?:--)|(/\\*(?:.|[\\n\\r])*?\\*/)|(\b(update|and|or|delete|insert|trancate|char|chr|into|substr|ascii|declare|exec|count|master|into|drop|execute)\b)`

func SqlFilter(matchStr string, exactly bool) string {
	re, err := regexp.Compile(reStr)
	if err != nil {
		log.Println(err)
		return ``
	}
	return re.ReplaceAllStringFunc(matchStr, func(dest string) string {
		if dest == `\` && !exactly {
			dest = strings.ReplaceAll(dest, dest, `\\\`+dest)
		} else {
			dest = strings.ReplaceAll(dest, dest, `\`+dest)
		}
		return dest
	})
}

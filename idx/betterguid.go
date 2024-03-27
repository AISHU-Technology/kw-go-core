package idx

import (
	"github.com/kjk/betterguid"
	"strings"
)

func BetterGuId() string {
	return strings.TrimLeft(betterguid.New(), "-")
}

package withdraw

import (
	"os"
	"strconv"
)

var (
	// os env
	startCursor = os.Getenv("START_CURSOR")

	START_CURSOR = 0
)

func ParseEnv() {
	var err error

	if startCursor != "" {
		START_CURSOR, err = strconv.Atoi(startCursor)
		if err != nil {
			panic(err)
		}
	}
}

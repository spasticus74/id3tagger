package parse

import (
	"errors"
	"regexp"
	"strconv"
)

func ParseFilename(filename string) (string, string, error) {
	re := regexp.MustCompile(`(\d*) (.*).mp3`)
	matches := re.FindAllStringSubmatch(filename, -1)
	if len(matches) == 0 {
		return "-1", "", errors.New("unable to parse filename")
	}

	n, err := strconv.Atoi(matches[0][1])
	if err != nil {
		return "-1", "", errors.New("unable to parse filename")
	}

	return strconv.Itoa(n), matches[0][2], nil
}

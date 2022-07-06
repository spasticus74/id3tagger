package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func GenrePrompt() string {
	var s string
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprintf(os.Stderr, "Enter a genre, or press enter to ignore: ")
		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}

	ret := strings.TrimSpace(s)
	return ret
}

func YearPrompt() (string, error) {
	var s string
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprintf(os.Stderr, "Enter a year, or press enter to ignore: ")
		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}

	ret := strings.TrimSpace(s)
	if strings.ToUpper(ret) == "" {
		return "-1", nil
	} else {
		_, e := strconv.Atoi(ret)
		if e != nil {
			return "-1", e
		}
		return ret, nil
	}
}

func NamePrompt(label, value string) string {
	var s string
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprintf(os.Stderr, "%s %s. Enter new value or press enter to accept: ", label, value)
		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}

	ret := strings.TrimSpace(s)
	if strings.ToUpper(ret) == "" {
		return value
	} else {
		return ret
	}
}

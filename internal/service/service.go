package service

import (
	"errors"
	"log"
	"strings"

	morse "github.com/Yandex-Practicum/go1fl-sprint6-final/pkg/morse"
)

func IsItMorse(line string) (result bool) {
	if len(line) == 0 {
		err := errors.New("line is empty")
		log.Println(err)
		return false
	}

	f := func(r rune) bool {
		return r == '.' || r == '-'
	}
	result = strings.ContainsFunc(line, f)
	return result
}

func LineConverter(line string) string {
	var result string
	if IsItMorse(line) {
		result = morse.ToText(line)
	} else {
		result = morse.ToMorse(line)
	}
	return result

}

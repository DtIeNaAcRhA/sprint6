package service

import (
	"strings"
	"unicode/utf8"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/pkg/morse"
)

func Translate(s string) (string, error) {
	var strResult string
	if strings.Count(s, ".")+strings.Count(s, "-")+strings.Count(s, " ") == utf8.RuneCountInString(s) {
		strResult = morse.ToText(s)
	} else {
		strResult = morse.ToMorse(s)
	}

	return strResult, nil
}

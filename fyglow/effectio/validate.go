package effectio

import (
	"fmt"
	"gglow/text"
	"strings"
	"unicode"
)

func ValidateFolderName(title string) error {
	title = strings.TrimSpace(title)

	if len(title) < 1 {
		return fmt.Errorf(text.MsgRequired.String())
	}

	if title == "NULL" {
		return fmt.Errorf(text.MsgRequired.String())
	}

	for _, c := range title {
		if !(c == '_' || unicode.IsLetter(c) || unicode.IsDigit(c)) {
			return fmt.Errorf(text.MsgAlphaNumeric.String())
		}
	}
	return nil
}

func ValidateEffectName(title string) error {
	title = strings.TrimSpace(title)

	if len(title) < 1 {
		return fmt.Errorf(text.MsgRequired.String())
	}

	if title == "NULL" {
		return fmt.Errorf(text.MsgRequired.String())
	}

	for i, c := range title {
		if i == 0 && !unicode.IsUpper(c) {
			return fmt.Errorf(text.MsgFirstUpper.String())
		}
		if !(c == ' ' || unicode.IsLetter(c) || unicode.IsDigit(c)) {
			return fmt.Errorf(text.MsgAlphaNumeric.String())
		}
	}
	return nil
}

package store

import (
	"fmt"
	"glow-gui/resources"
	"strings"
	"unicode"

	"fyne.io/fyne/v2"
)

func ValidateFolderName(title string) error {
	title = strings.TrimSpace(title)

	if len(title) < 1 {
		return fmt.Errorf(resources.MsgRequired.String())
	}

	if title == "NULL" {
		return fmt.Errorf(resources.MsgRequired.String())
	}

	for _, c := range title {
		if !(c == '_' || unicode.IsLetter(c) || unicode.IsDigit(c)) {
			return fmt.Errorf(resources.MsgAlphaNumeric.String())
		}
	}
	return nil
}

func (st *Store) ValidateNewFolderName(title string) error {
	err := ValidateFolderName(title)
	if err != nil {
		return err
	}

	err = st.IsDuplicate(title)
	return err
}

func (st *Store) ValidateNewEffectName(title string) error {
	err := ValidateEffectName(title)
	if err != nil {
		return err
	}
	err = st.IsDuplicate(title)
	return err
}

func ValidateEffectName(title string) error {
	title = strings.TrimSpace(title)

	if len(title) < 1 {
		return fmt.Errorf(resources.MsgRequired.String())
	}

	if title == "NULL" {
		return fmt.Errorf(resources.MsgRequired.String())
	}

	for i, c := range title {
		if i == 0 && !unicode.IsUpper(c) {
			return fmt.Errorf(resources.MsgFirstUpper.String())
		}
		if !(c == ' ' || unicode.IsLetter(c) || unicode.IsDigit(c)) {
			return fmt.Errorf(resources.MsgAlphaNumeric.String())
		}
	}
	return nil
}

func MakeTitle(uri fyne.URI) (s string) {
	s = uri.Name()
	i := strings.Index(s, uri.Extension())
	if i > 0 {
		s = s[:i]
	}
	s = strings.ReplaceAll(s, "_", " ")
	return
}

func MakeFileName(title string) string {
	s := strings.ReplaceAll(title, " ", "_")
	s += ".yaml"
	return s
}

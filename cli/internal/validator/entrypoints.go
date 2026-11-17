package validator

import (
	"fmt"

	"siguma0013/reskk-dictionary/internal/model"
)

func ValidateCharacters(characters model.Characters, filePath string) []error {
	return validateEntries(characters.Entries, filePath, validateCharacterReadings)
}

func validateCharacterReadings(entry model.CharacterEntry) []error {
	var errs []error

	// 音読みのバリデーション
	if err := validateYomi(entry.GetOnYomi(), allowOnYomiString); err != nil {
		errs = append(errs, fmt.Errorf("音読みの%w", err))
	}

	// 訓読みのバリデーション
	if err := validateYomi(entry.GetKunYomi(), allowYomiString); err != nil {
		errs = append(errs, fmt.Errorf("訓読みの%w", err))
	}

	return errs
}

func ValidateWords(words model.Words, filePath string) []error {
	return validateEntries(words.Entries, filePath, validateWordReadings)
}

func validateWordReadings(entry model.WordEntry) []error {
	// 読みのバリデーション
	if err := validateYomi(entry.GetYomi(), allowYomiString); err != nil {
		return []error{fmt.Errorf("読みの%w", err)}
	}

	return nil
}

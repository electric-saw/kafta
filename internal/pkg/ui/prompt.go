package ui

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

func GetTextOrDef(question string, def string) string {
	prompt := promptui.Prompt{
		Label:   question,
		Default: def,
	}

	result, err := prompt.Run()
	if err != nil || result == "" {
		return def
	}

	return result
}

func GetText(question string, required bool) (string, error) {
	prompt := promptui.Prompt{
		Label: question,
		Validate: func(result string) error {
			if len(result) == 0 && required {
				return fmt.Errorf("Value is required!")
			}
			return nil
		},
	}

	return prompt.Run()
}

func GetPassword(question string, required bool) (string, error) {
	prompt := promptui.Prompt{
		Label: question,
		Mask:  '*',
		Validate: func(result string) error {
			if len(result) == 0 && required {
				return fmt.Errorf("Value is required!")
			}
			return nil
		},
	}
	return prompt.Run()

}

func GetConfirmation(question string, def bool) (bool, error) {
	prompt := promptui.Prompt{
		Label:     question,
		IsConfirm: true,
	}

	result, err := prompt.Run()
	if err != nil {
		return def, nil
	}

	return result == "y", err

}

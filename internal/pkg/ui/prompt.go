package ui

import (
	"github.com/manifoldco/promptui"
	"github.com/pterm/pterm"
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
	for {
		output, err := pterm.DefaultInteractiveTextInput.
			WithMultiLine(false).
			Show(question)

		if err != nil {
			return "", err
		}

		if len(output) == 0 && required {
			println("Value is required!")
			continue
		}

		return output, nil
	}
}

func GetPassword(question string, required bool) (string, error) {
	for {
		output, err := pterm.DefaultInteractiveTextInput.
			WithHide().
			Show(question)

		if err != nil {
			return "", err
		}

		if len(output) == 0 && required {
			println("Value is required!")
			continue
		}

		return output, nil
	}
}

func GetConfirmation(question string, def bool) (bool, error) {
	return pterm.DefaultInteractiveConfirm.
		WithDefaultValue(def).
		Show(question)

}

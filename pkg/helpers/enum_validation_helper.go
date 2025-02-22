package helpers

import (
	"fmt"

	"github.com/sugiiianaa/remember-my-story/internal/models/enums"
)

func ValidateMood(mood string) error {
	if enums.MoodFromString(mood) == enums.Mood.Unknown {
		return fmt.Errorf("invalid mood: %s", mood)
	}
	return nil
}

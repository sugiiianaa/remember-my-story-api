package enums

import (
	"encoding/json"
	"strings"
)

type MoodType int

// Mood "namespace" struct
var Mood = struct {
	Unknown   MoodType
	Happy     MoodType
	Sad       MoodType
	Energized MoodType
	Calm      MoodType
	Anxious   MoodType
}{
	Unknown:   0,
	Happy:     1,
	Sad:       2,
	Energized: 3,
	Calm:      4,
	Anxious:   5,
}

func (m MoodType) String() string {
	switch m {
	case Mood.Happy:
		return "Happy"
	case Mood.Sad:
		return "Sad"
	case Mood.Energized:
		return "Energized"
	case Mood.Calm:
		return "Calm"
	case Mood.Anxious:
		return "Anxious"
	default:
		return "Unknown"
	}
}

// MarshalJSON implements json.Marshaler
func (m MoodType) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.String())
}

// UnmarshalJSON implements json.Unmarshaler
func (m *MoodType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	*m = moodFromString(s)
	return nil
}

func moodFromString(s string) MoodType {
	switch strings.ToLower(s) {
	case "happy":
		return Mood.Happy
	case "sad":
		return Mood.Sad
	case "energized":
		return Mood.Energized
	case "calm":
		return Mood.Calm
	case "anxious":
		return Mood.Anxious
	default:
		return Mood.Unknown
	}
}

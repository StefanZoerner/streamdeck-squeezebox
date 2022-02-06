package plugin

import (
	"encoding/json"
	"github.com/samwho/streamdeck"
)

type PlayerSettings struct {
	PlayerId   string `json:"player_id"`
	PlayerName string `json:"player_name"`
}

func getPlayerSettingsFromKeyDownEvent(event streamdeck.Event) (PlayerSettings, error) {
	settings := PlayerSettings{}

	payload := streamdeck.KeyDownPayload{}
	err := json.Unmarshal(event.Payload, &payload)
	if err != nil {
		return settings, err
	}

	err = json.Unmarshal(payload.Settings, &settings)

	return settings, err
}

func getPlayerSettingsFromWillAppearEvent(event streamdeck.Event) (PlayerSettings, error) {
	settings := PlayerSettings{}

	payload := streamdeck.WillAppearPayload{}
	err := json.Unmarshal(event.Payload, &payload)
	if err != nil {
		return settings, err
	}

	err = json.Unmarshal(payload.Settings, &settings)
	return settings, err
}

func getPlayerSettingsFromWillDisappearEvent(event streamdeck.Event) (PlayerSettings, error) {
	settings := PlayerSettings{}

	payload := streamdeck.WillDisappearPayload{}
	err := json.Unmarshal(event.Payload, &payload)
	if err != nil {
		return settings, err
	}

	err = json.Unmarshal(payload.Settings, &settings)
	return settings, err
}

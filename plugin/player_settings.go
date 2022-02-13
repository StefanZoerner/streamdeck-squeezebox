package plugin

import (
	"context"
	"encoding/json"
	"github.com/samwho/streamdeck"
)

type PlayerSettings struct {
	PlayerId   string `json:"player_id"`
	PlayerName string `json:"player_name"`
}

type PlayerSelection struct {
	Players []PlayerSettings `json:"players"`
}

type DataFromPlayerSelectionPI struct {
	Command string `json:"command"`
	Value   string `json:"value"`
}

func selectPlayerHandlerWillAppear(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
	LogEvent(client, event)

	payload := streamdeck.WillAppearPayload{}
	err := json.Unmarshal(event.Payload, &payload)
	if err != nil {
		LogErrorWithEvent(client, event, err)
		return err
	}

	settings := PlayerSettings{}
	err = json.Unmarshal(payload.Settings, &settings)
	if err != nil {
		LogErrorWithEvent(client, event, err)
		return err
	}

	if settings.PlayerId == "" {
		settings.PlayerName = "(None)"
		err = client.SetSettings(ctx, settings)
		if err != nil {
			LogErrorWithEvent(client, event, err)
			return err
		}
	}

	return nil
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

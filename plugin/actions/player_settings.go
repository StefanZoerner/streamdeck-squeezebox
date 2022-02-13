package actions

import (
	"context"
	"encoding/json"
	"github.com/StefanZoerner/streamdeck-squeezebox/plugin/general"
	"github.com/StefanZoerner/streamdeck-squeezebox/squeezebox"
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

func getPlayerSelection() (PlayerSelection, error) {

	selection := PlayerSelection{}

	globalSettings := general.GetPluginGlobalSettings()
	conProps := globalSettings.ConnectionProps()

	players, err := squeezebox.GetPlayerInfos(conProps)
	if err == nil {
		playerSettings := []PlayerSettings{}
		for _, p := range players {
			np := PlayerSettings{
				PlayerId:   p.ID,
				PlayerName: p.Name,
			}
			playerSettings = append(playerSettings, np)
		}
		selection.Players = playerSettings
	}

	return selection, err
}

func SelectPlayerHandlerWillAppear(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
	general.LogEvent(client, event)

	payload := streamdeck.WillAppearPayload{}
	err := json.Unmarshal(event.Payload, &payload)
	if err != nil {
		general.LogErrorWithEvent(client, event, err)
		return err
	}

	settings := PlayerSettings{}
	err = json.Unmarshal(payload.Settings, &settings)
	if err != nil {
		general.LogErrorWithEvent(client, event, err)
		return err
	}

	if settings.PlayerId == "" {
		settings.PlayerName = "(None)"
		err = client.SetSettings(ctx, settings)
		if err != nil {
			general.LogErrorWithEvent(client, event, err)
			return err
		}
	}

	return nil
}

func GetPlayerSettingsFromKeyDownEvent(event streamdeck.Event) (PlayerSettings, error) {
	settings := PlayerSettings{}

	payload := streamdeck.KeyDownPayload{}
	err := json.Unmarshal(event.Payload, &payload)
	if err != nil {
		return settings, err
	}

	err = json.Unmarshal(payload.Settings, &settings)

	return settings, err
}

func GetPlayerSettingsFromWillAppearEvent(event streamdeck.Event) (PlayerSettings, error) {
	settings := PlayerSettings{}

	payload := streamdeck.WillAppearPayload{}
	err := json.Unmarshal(event.Payload, &payload)
	if err != nil {
		return settings, err
	}

	err = json.Unmarshal(payload.Settings, &settings)
	return settings, err
}

func GetPlayerSettingsFromWillDisappearEvent(event streamdeck.Event) (PlayerSettings, error) {
	settings := PlayerSettings{}

	payload := streamdeck.WillDisappearPayload{}
	err := json.Unmarshal(event.Payload, &payload)
	if err != nil {
		return settings, err
	}

	err = json.Unmarshal(payload.Settings, &settings)
	return settings, err
}

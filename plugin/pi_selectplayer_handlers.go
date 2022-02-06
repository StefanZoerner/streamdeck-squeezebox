package plugin

import (
	"context"
	"encoding/json"
	"github.com/StefanZoerner/streamdeck-squeezebox/squeezebox"
	"github.com/samwho/streamdeck"
)

type PlayerSelection struct {
	Players []PlayerSettings `json:"players"`
}

type DataFromPlayerSelectionPI struct {
	Command string `json:"command"`
	Value   string `json:"value"`
}

func selectPlayerHandlerWillAppear(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
	logEvent(client, event)

	payload := streamdeck.WillAppearPayload{}
	err := json.Unmarshal(event.Payload, &payload)
	if err != nil {
		logError(client, event, err)
		return err
	}

	settings := PlayerSettings{}
	err = json.Unmarshal(payload.Settings, &settings)
	if err != nil {
		logError(client, event, err)
		return err
	}

	if settings.PlayerId == "" {
		settings.PlayerName = "(None)"
		err = client.SetSettings(ctx, settings)
		if err != nil {
			logError(client, event, err)
			return err
		}
	}

	return nil
}

func selectPlayerHandlerSendToPlugin(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
	logEvent(client, event)

	fromPI := DataFromPlayerSelectionPI{}
	err := json.Unmarshal(event.Payload, &fromPI)
	if err != nil {
		logError(client, event, err)
		return err
	}

	globalSettings := GetPluginGlobalSettings()

	if fromPI.Command == "getPlayerSelectionOptions" {

		players, err := squeezebox.GetPlayers(globalSettings.Hostname, globalSettings.CLIPort)
		if err != nil {
			logError(client, event, err)
			return err
		}

		playerSettings := []PlayerSettings{}
		for _, p := range players {
			np := PlayerSettings{
				PlayerId:   p.ID,
				PlayerName: p.Name,
			}
			playerSettings = append(playerSettings, np)
		}

		payload := PlayerSelection{
			Players: playerSettings,
		}

		err = client.SendToPropertyInspector(ctx, &payload)
		if err != nil {
			logError(client, event, err)
			return err
		}
	} else if fromPI.Command == "setSelectedPlayer" {

		player_id := fromPI.Value
		pinfo, err := squeezebox.GetPlayerInfo(globalSettings.Hostname, globalSettings.CLIPort, player_id)
		if err != nil {
			logError(client, event, err)
			return err
		}

		np := PlayerSettings{
			PlayerId:   player_id,
			PlayerName: pinfo.Name,
		}

		err = client.SetSettings(ctx, np)
		if err != nil {
			logError(client, event, err)
			return err
		}
	}

	return nil
}

package plugin

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/StefanZoerner/streamdeck-squeezebox/plugin/keyimages"
	"github.com/StefanZoerner/streamdeck-squeezebox/squeezebox"
	"github.com/samwho/streamdeck"
)

const (
	TRACK_PREV = "prev"
	TRACK_NEXT = "next"
)

type TrackActionSettings struct {
	PlayerSettings
	Direction string `json:"track_direction"`
}

type TrackFromPI struct {
	Command  string              `json:"command"`
	Settings TrackActionSettings `json:"settings"`
}

func setupTrackActions(client *streamdeck.Client) {
	trackAction := client.Action("de.szoerner.streamdeck.squeezebox.actions.track")
	trackAction.RegisterHandler(streamdeck.WillAppear, WillAppearRequestGlobalSettingsHandler)

	trackAction.RegisterHandler(streamdeck.KeyDown, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		logEvent(client, event)

		settings := TrackActionSettings{}
		err := getSettingsFromKeydownEvent(event, &settings)
		if err == nil {
			if settings.PlayerId == "" {
				_ = client.ShowAlert(ctx)
				err = errors.New("no player configured")
			} else {

				globalSettings := GetPluginGlobalSettings()
				delta := 0
				if settings.Direction == TRACK_PREV {
					delta = -1
				} else if settings.Direction == TRACK_NEXT {
					delta = +1
				}
				if delta != 0 {
					_, _, err := squeezebox.ChangePlayerTrack(globalSettings.Hostname, globalSettings.CLIPort, settings.PlayerId, delta)
					if err != nil {
						_ = client.ShowAlert(ctx)
					}
				}
			}
		}

		return err
	})
	trackAction.RegisterHandler(streamdeck.WillAppear, trackActionWillAppear)
	trackAction.RegisterHandler(streamdeck.SendToPlugin, trackSendToPlugin)
}

func getSettingsFromKeydownEvent(event streamdeck.Event, settings interface{}) error {
	var err error

	payload := streamdeck.KeyDownPayload{}
	err = json.Unmarshal(event.Payload, &payload)
	if err == nil {
		err = json.Unmarshal(payload.Settings, &settings)
	}

	return err
}

func trackActionWillAppear(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
	logEvent(client, event)

	payload := streamdeck.WillAppearPayload{}
	err := json.Unmarshal(event.Payload, &payload)
	if err != nil {
		logErrorWithEvent(client, event, err)
		return err
	}

	settings := TrackActionSettings{}
	err = json.Unmarshal(payload.Settings, &settings)
	if err != nil {
		logErrorWithEvent(client, event, err)
		return err
	}

	if settings.PlayerId == "" {
		settings.PlayerName = "(None)"
		err = client.SetSettings(ctx, settings)
		if err != nil {
			logErrorWithEvent(client, event, err)
			return err
		}
	}

	if settings.Direction == "" {
		settings.Direction = TRACK_NEXT
		err = client.SetSettings(ctx, settings)
		if err != nil {
			logErrorWithEvent(client, event, err)
			return err
		}
	}

	err = trackSetKeyImage(ctx, client, settings.Direction)
	if err != nil {
		logErrorWithEvent(client, event, err)
		return err
	}

	return nil
}

func trackSendToPlugin(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
	logEvent(client, event)

	fromPI := TrackFromPI{}
	err := json.Unmarshal(event.Payload, &fromPI)
	if err != nil {
		logErrorWithEvent(client, event, err)
		return err
	}

	globalSettings := GetPluginGlobalSettings()

	if fromPI.Command == "getPlayerSelectionOptions" {

		conProps := globalSettings.connectionProps()

		players, err := squeezebox.GetPlayerInfos(conProps)
		if err != nil {
			logErrorWithEvent(client, event, err)
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
			logErrorWithEvent(client, event, err)
			return err
		}
	} else if fromPI.Command == "sendFormData" {

		err = client.SetSettings(ctx, fromPI.Settings)
		if err != nil {
			logErrorWithEvent(client, event, err)
			return err
		}

		err = trackSetKeyImage(ctx, client, fromPI.Settings.Direction)
		if err != nil {
			logErrorWithEvent(client, event, err)
			return err
		}

	}

	return nil
}

func trackSetKeyImage(ctx context.Context, client *streamdeck.Client, direction string) error {
	var err error

	switch direction {
	case TRACK_NEXT:
		image, err := keyimages.GetStreamDeckImageForIcon("track_next")
		if err == nil {
			err = client.SetImage(ctx, image, streamdeck.HardwareAndSoftware)
		}
	case TRACK_PREV:
		image, err := keyimages.GetStreamDeckImageForIcon("track_prev")
		if err == nil {
			err = client.SetImage(ctx, image, streamdeck.HardwareAndSoftware)
		}
	}

	return err
}

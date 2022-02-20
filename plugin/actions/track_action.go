package actions

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/StefanZoerner/streamdeck-squeezebox/plugin/general"
	"github.com/StefanZoerner/streamdeck-squeezebox/plugin/keyimages"
	"github.com/StefanZoerner/streamdeck-squeezebox/squeezebox"
	"github.com/samwho/streamdeck"
)

const (
	TrackPrev = "prev"
	TrackNext = "next"
)

type trackActionSettings struct {
	PlayerSettings
	Direction string `json:"track_direction"`
}

type trackFromPI struct {
	Command  string              `json:"command"`
	Settings trackActionSettings `json:"Settings"`
}

func SetupTrackAction(client *streamdeck.Client) {
	trackAction := client.Action("de.szoerner.streamdeck.squeezebox.actions.track")
	trackAction.RegisterHandler(streamdeck.WillAppear, general.WillAppearRequestGlobalSettingsHandler)

	trackAction.RegisterHandler(streamdeck.WillAppear, trackHandlerWillAppear)
	trackAction.RegisterHandler(streamdeck.KeyDown, trackHandlerKeyDown)
	trackAction.RegisterHandler(streamdeck.SendToPlugin, trackHandlerSendToPlugin)
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

func trackHandlerKeyDown(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
	general.LogEvent(client, event)

	var err error

	settings := trackActionSettings{}
	err = getSettingsFromKeydownEvent(event, &settings)
	if err == nil {
		if settings.PlayerId == "" {
			_ = client.ShowAlert(ctx)
			err = errors.New("no player configured")
		} else {

			globalSettings := general.GetPluginGlobalSettings()
			delta := 0
			if settings.Direction == TrackPrev {
				delta = -1
			} else if settings.Direction == TrackNext {
				delta = +1
			}
			if delta != 0 {
				_, _, err = squeezebox.ChangePlayerTrack(globalSettings.Hostname, globalSettings.CLIPort, settings.PlayerId, delta)
				if err != nil {
					_ = client.ShowAlert(ctx)
				}
			}
		}
	}

	return err
}

func trackHandlerWillAppear(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
	general.LogEvent(client, event)

	var err error

	payload := streamdeck.WillAppearPayload{}
	err = json.Unmarshal(event.Payload, &payload)
	if err != nil {
		general.LogErrorWithEvent(client, event, err)
		return err
	}

	settings := trackActionSettings{}
	err = json.Unmarshal(payload.Settings, &settings)
	if err == nil {
		var modified bool
		if settings.PlayerId == "" {
			settings.PlayerName = "(None)"
			modified = true
		}
		if settings.Direction == "" {
			settings.Direction = TrackNext
			modified = true
		}
		if modified {
			err = client.SetSettings(ctx, settings)
		}
	}

	if err == nil {
		err = setTrackKeyImage(ctx, client, settings.Direction)
	}

	if err != nil {
		general.LogErrorWithEvent(client, event, err)
	}

	return err
}

func trackHandlerSendToPlugin(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
	general.LogEvent(client, event)

	var err error

	fromPI := trackFromPI{}
	err = json.Unmarshal(event.Payload, &fromPI)
	if err == nil {
		if fromPI.Command == "getPlayerSelectionOptions" {
			var payload PlayerSelection
			payload, err = getPlayerSelection()
			if err == nil {
				err = client.SendToPropertyInspector(ctx, &payload)
			}
		} else if fromPI.Command == "sendFormData" {
			err = client.SetSettings(ctx, fromPI.Settings)
			if err == nil {
				err = setTrackKeyImage(ctx, client, fromPI.Settings.Direction)
			}
		}
	}

	if err != nil {
		general.LogErrorWithEvent(client, event, err)
	}

	return err
}

func setTrackKeyImage(ctx context.Context, client *streamdeck.Client, direction string) error {
	var err error
	var iconName string

	switch direction {
	case TrackNext:
		iconName = "track_next"
	case TrackPrev:
		iconName = "track_prev"
	}

	var image string
	image, err = keyimages.GetStreamDeckImageForIcon(iconName)
	if err == nil {
		err = client.SetImage(ctx, image, streamdeck.HardwareAndSoftware)
	}

	return err
}

package actions

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/StefanZoerner/streamdeck-squeezebox/plugin/general"
	"github.com/StefanZoerner/streamdeck-squeezebox/plugin/keyimages"
	"github.com/StefanZoerner/streamdeck-squeezebox/squeezebox"
	"github.com/samwho/streamdeck"
	"time"
)

const (
	VolumeUp   = "up"
	VolumeDown = "down"
)

type VolumeActionSettings struct {
	PlayerSettings
	Kind string `json:"volume_kind"`
}

type VolumeFromPI struct {
	Command  string               `json:"command"`
	Settings VolumeActionSettings `json:"settings"`
}

func SetupVolumeAction(client *streamdeck.Client) {

	volumeAction := client.Action("de.szoerner.streamdeck.squeezebox.actions.volume")
	volumeAction.RegisterHandler(streamdeck.WillAppear, general.WillAppearRequestGlobalSettingsHandler)
	volumeAction.RegisterHandler(streamdeck.KeyDown, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		general.LogEvent(client, event)

		settings := VolumeActionSettings{}
		err := getSettingsFromKeydownEvent(event, &settings)
		if err == nil {

			playerID := settings.PlayerId
			globalSettings := general.GetPluginGlobalSettings()
			if playerID == "" {
				playerID = globalSettings.DefaultPlayerID
			}

			if playerID == "" {
				_ = client.ShowAlert(ctx)
				return errors.New("No player configured")
			}

			delta := 0
			if settings.Kind == VolumeDown {
				delta = -10
			} else if settings.Kind == VolumeUp {
				delta = +10
			}
			if delta != 0 {
				cp := globalSettings.ConnectionProps()
				volume, err := squeezebox.ChangePlayerVolume(cp, playerID, delta)
				if err != nil {
					_ = client.ShowAlert(ctx)
				} else {
					go displayNumberInKey(ctx, client, volume, settings.Kind)
				}
			}
		}

		return err
	})
	volumeAction.RegisterHandler(streamdeck.WillAppear, volumeHandlerWillAppear)
	volumeAction.RegisterHandler(streamdeck.SendToPlugin, volumeHandlerSendToPlugin)
}

func displayNumberInKey(ctx context.Context, client *streamdeck.Client, n int, volumeKind string) {

	// Display Number in Key
	img := keyimages.CreateKeyImageWithNumber(n)
	s, _ := streamdeck.Image(img)
	err := client.SetImage(ctx, s, streamdeck.HardwareAndSoftware)
	if err != nil {
		_ = client.LogMessage("Error: " + err.Error())
	}

	// Wait 2 seconds, see https://gobyexample.com/timers
	timer := time.NewTimer(2 * time.Second)
	<-timer.C

	// Display "old" Image
	err = volumeSetKeyImage(ctx, client, volumeKind)
	if err != nil {
		_ = client.LogMessage("Error: " + err.Error())
	}
}

func volumeHandlerWillAppear(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
	general.LogEvent(client, event)

	var err error

	payload := streamdeck.WillAppearPayload{}
	err = json.Unmarshal(event.Payload, &payload)
	if err == nil {

		settings := VolumeActionSettings{}
		err = json.Unmarshal(payload.Settings, &settings)
		if err == nil {

			var modified bool
			if settings.PlayerId == "" && settings.PlayerName != "(Default)" {
				settings.PlayerName = "(Default)"
				modified = true
			}
			if settings.Kind == "" {
				settings.Kind = VolumeUp
				modified = true
			}
			if modified {
				err = client.SetSettings(ctx, settings)
			}

			if err == nil {
				err = volumeSetKeyImage(ctx, client, settings.Kind)
			}
		}
	}

	if err != nil {
		general.LogErrorWithEvent(client, event, err)
	}
	return err
}

func volumeHandlerSendToPlugin(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
	general.LogEvent(client, event)

	fromPI := VolumeFromPI{}
	err := json.Unmarshal(event.Payload, &fromPI)
	if err != nil {
		general.LogErrorWithEvent(client, event, err)
		return err
	}

	if fromPI.Command == "getPlayerSelectionOptions" {

		payload, err := getPlayerSelection()
		if err == nil {
			err = client.SendToPropertyInspector(ctx, &payload)
		}

		if err != nil {
			general.LogErrorWithEvent(client, event, err)
			return err
		}

	} else if fromPI.Command == "sendFormData" {

		err = client.SetSettings(ctx, fromPI.Settings)
		if err != nil {
			general.LogErrorWithEvent(client, event, err)
			return err
		}

		err = volumeSetKeyImage(ctx, client, fromPI.Settings.Kind)
		if err != nil {
			general.LogErrorWithEvent(client, event, err)
			return err
		}

	}

	return nil
}

func volumeSetKeyImage(ctx context.Context, client *streamdeck.Client, kind string) error {
	var err error
	var iconName string

	switch kind {
	case VolumeUp:
		iconName = "volume_up"
	case VolumeDown:
		iconName = "volume_down"
	}

	var image string
	image, err = keyimages.GetStreamDeckImageForIcon(iconName)
	if err == nil {
		err = client.SetImage(ctx, image, streamdeck.HardwareAndSoftware)
	}

	return err
}

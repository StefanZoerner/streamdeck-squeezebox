package plugin

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/StefanZoerner/streamdeck-squeezebox/plugin/keyimages"
	"github.com/StefanZoerner/streamdeck-squeezebox/squeezebox"
	"github.com/samwho/streamdeck"
	"time"
)

const (
	VOLUME_UP   = "up"
	VOLUME_DOWN = "down"
)

type VolumeActionSettings struct {
	PlayerSettings
	Kind string `json:"volume_kind"`
}

type VolumeFromPI struct {
	Command  string               `json:"command"`
	Settings VolumeActionSettings `json:"settings"`
}

func setupVolumeActions(client *streamdeck.Client) {

	volumeAction := client.Action("de.szoerner.streamdeck.squeezebox.actions.volume")
	volumeAction.RegisterHandler(streamdeck.WillAppear, WillAppearRequestGlobalSettingsHandler)
	volumeAction.RegisterHandler(streamdeck.KeyDown, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		logEvent(client, event)

		settings := VolumeActionSettings{}
		err := getSettingsFromKeydownEvent(event, &settings)
		if err == nil {
			if settings.PlayerId == "" {
				_ = client.ShowAlert(ctx)
				err = errors.New("no player configured")
			} else {

				globalSettings := GetPluginGlobalSettings()
				delta := 0
				if settings.Kind == VOLUME_DOWN {
					delta = -10
				} else if settings.Kind == VOLUME_UP {
					delta = +10
				}
				if delta != 0 {
					volume, err := squeezebox.ChangePlayerVolume(globalSettings.Hostname, globalSettings.CLIPort, settings.PlayerId, delta)
					if err != nil {
						_ = client.ShowAlert(ctx)
					} else {
						go displayNumberInKey(ctx, client, volume, settings.Kind)
					}
				}
			}
		}

		return err
	})
	volumeAction.RegisterHandler(streamdeck.WillAppear, volumeActionWillAppear)
	volumeAction.RegisterHandler(streamdeck.SendToPlugin, volumeSendToPlugin)
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

func volumeActionWillAppear(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
	logEvent(client, event)

	payload := streamdeck.WillAppearPayload{}
	err := json.Unmarshal(event.Payload, &payload)
	if err != nil {
		logError(client, event, err)
		return err
	}

	settings := VolumeActionSettings{}
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

	if settings.Kind == "" {
		settings.Kind = VOLUME_UP
		err = client.SetSettings(ctx, settings)
		if err != nil {
			logError(client, event, err)
			return err
		}
	}

	err = volumeSetKeyImage(ctx, client, settings.Kind)
	if err != nil {
		logError(client, event, err)
		return err
	}

	return nil
}

func volumeSendToPlugin(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
	logEvent(client, event)

	fromPI := VolumeFromPI{}
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

		var playerSettings []PlayerSettings
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

	} else if fromPI.Command == "sendFormData" {

		err = client.SetSettings(ctx, fromPI.Settings)
		if err != nil {
			logError(client, event, err)
			return err
		}

		err = volumeSetKeyImage(ctx, client, fromPI.Settings.Kind)
		if err != nil {
			logError(client, event, err)
			return err
		}

	}

	return nil
}

func volumeSetKeyImage(ctx context.Context, client *streamdeck.Client, kind string) error {
	var err error

	switch kind {
	case VOLUME_UP:
		image, err := keyimages.GetStreamDeckImageForIcon("volume_up")
		if err == nil {
			err = client.SetImage(ctx, image, streamdeck.HardwareAndSoftware)
		}
		break
	case VOLUME_DOWN:
		image, err := keyimages.GetStreamDeckImageForIcon("volume_down")
		if err == nil {
			err = client.SetImage(ctx, image, streamdeck.HardwareAndSoftware)
		}
	}

	return err
}

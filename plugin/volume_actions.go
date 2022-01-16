package plugin

import (
	"context"
	"errors"
	"fmt"
	"github.com/StefanZoerner/streamdeck-squeezebox/squeezebox"
	"github.com/samwho/streamdeck"
	"time"
)

func setupVolumeActions(client *streamdeck.Client) {

	// Volume Up
	//
	volumeUpAction := client.Action("de.szoerner.streamdeck.squeezebox.actions.volumeup")
	volumeUpAction.RegisterHandler(streamdeck.WillAppear, WillAppearRequestGlobalSettingsHandler)
	volumeUpAction.RegisterHandler(streamdeck.KeyDown, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		logEvent(client, event)

		settings, err := getPlayerSettingsFromKeyDownEvent(event)
		if (err == nil) {
			if settings.PlayerId == "" {
				client.ShowAlert(ctx)
				err = errors.New("No player configured")
			} else {

				globalSettings := GetPluginGlobalSettings()

				volume, err := squeezebox.ChangePlayerVolume(globalSettings.Hostname, globalSettings.CliPort, settings.PlayerId, +10)
				if err != nil {
					client.ShowAlert(ctx)
				} else {
					go displayTextAsTitleForTwoSeconds(ctx, client, fmt.Sprintf("%d", volume))
				}
			}
		}

		return err
	})
	volumeUpAction.RegisterHandler(streamdeck.WillAppear, selectPlayerHandlerWillAppear)
	volumeUpAction.RegisterHandler(streamdeck.SendToPlugin, selectPlayerHandlerSendToPlugin)

	// Volume Down
	//
	volumeDownAction := client.Action("de.szoerner.streamdeck.squeezebox.actions.volumedown")
	volumeUpAction.RegisterHandler(streamdeck.WillAppear, WillAppearRequestGlobalSettingsHandler)
	volumeDownAction.RegisterHandler(streamdeck.KeyDown, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		logEvent(client, event)

		settings, err := getPlayerSettingsFromKeyDownEvent(event)
		if (err == nil) {
			if settings.PlayerId == "" {
				client.ShowAlert(ctx)
				err = errors.New("No player configured")
			} else {

				globalSettings := GetPluginGlobalSettings()

				volume, err := squeezebox.ChangePlayerVolume(globalSettings.Hostname, globalSettings.CliPort, settings.PlayerId, -10)
				if err != nil {
					client.ShowAlert(ctx)
				} else {
					go displayTextAsTitleForTwoSeconds(ctx, client, fmt.Sprintf("%d", volume))
				}
			}
		}

		return err
	})
	volumeDownAction.RegisterHandler(streamdeck.WillAppear, selectPlayerHandlerWillAppear)
	volumeDownAction.RegisterHandler(streamdeck.SendToPlugin, selectPlayerHandlerSendToPlugin)
}


func displayTextAsTitleForTwoSeconds (ctx context.Context, client *streamdeck.Client, text string) {
	client.SetTitle(ctx, text, streamdeck.HardwareAndSoftware)
	time.Sleep(2 * time.Second)
	client.SetTitle(ctx, "", streamdeck.HardwareAndSoftware)
}
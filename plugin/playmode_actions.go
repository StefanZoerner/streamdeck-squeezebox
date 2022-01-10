package plugin

import (
	"context"
	"github.com/StefanZoerner/streamdeck-squeezebox/squeezebox"
	"github.com/samwho/streamdeck"
)

func setupPlaymodeActions(client *streamdeck.Client) {

	// Play Toggle
	//
	playtoggleaction := client.Action("de.szoerner.streamdeck.squeezebox.actions.playtoggle")
	playtoggleaction.RegisterHandler(streamdeck.KeyDown, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		logEvent(client, event)

		status, err := squeezebox.TogglePlayerMode(player)
		if err != nil {
			client.ShowAlert(ctx)
			logError(client, "playtoggle", err)
			return err
		}

		client.LogMessage("New status: "+status)
		if status == "play" {
			image, err := getImageByFilename("./images/PauseKey.png")
			if (err != nil) {
				logError(client, "pause", err)
			} else {
				client.SetImage(ctx, image, streamdeck.HardwareAndSoftware);
			}
		} else if status == "stop" || status == "pause" {
			image, err := getImageByFilename("./images/PlayKey.png")
			if (err != nil) {
				logError(client, "pause", err)
			} else {
				client.SetImage(ctx, image, streamdeck.HardwareAndSoftware);
			}
		}

		return nil
	})

	playtoggleaction.RegisterHandler(streamdeck.WillAppear, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		logEvent(client, event)

		status, err := squeezebox.GetPlayerMode(player)
		if (err != nil) {
			logError(client, "pause", err)
		} else {
			if status == "play" {
				image, err := getImageByFilename("./images/PauseKey.png")
				if (err != nil) {
					logError(client, "pause", err)
				} else {
					client.SetImage(ctx, image, streamdeck.HardwareAndSoftware);
				}
			} else if status == "stop" || status == "pause" {
				image, err := getImageByFilename("./images/PlayKey.png")
				if (err != nil) {
					logError(client, "pause", err)
				} else {
					client.SetImage(ctx, image, streamdeck.HardwareAndSoftware);
				}
			}
		}

		return nil
	})

	playtoggleaction.RegisterHandler(streamdeck.WillAppear, selectPlayerHandlerWillAppear)
	playtoggleaction.RegisterHandler(streamdeck.SendToPlugin, selectPlayerHandlerSendToPlugin)
}
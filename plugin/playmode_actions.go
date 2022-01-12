package plugin

import (
	"context"
	"errors"
	"github.com/StefanZoerner/streamdeck-squeezebox/squeezebox"
	"github.com/samwho/streamdeck"
)

func setupPlaymodeActions(client *streamdeck.Client) {

	// Play Toggle
	//
	playtoggleaction := client.Action("de.szoerner.streamdeck.squeezebox.actions.playtoggle")
	playtoggleaction.RegisterHandler(streamdeck.KeyDown, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		logEvent(client, event)

		settings, err := getPlayerSettingsFromKeyDownEvent(event)
		if (err != nil) {
			logError(client, event, err)
			return err
		}

		if settings.PlayerId == "" {
			client.ShowAlert(ctx)
			err = errors.New("No player configured")
			logError(client, event, err)
			return err
		}

		mode, err := squeezebox.TogglePlayerMode(settings.PlayerId)
		if err != nil {
			client.ShowAlert(ctx)
		} else {
			err = setImageForPlayMode(ctx, client, mode)
		}

		if err != nil {
			logError(client, event, err)
		}
		return err
	})

	playtoggleaction.RegisterHandler(streamdeck.WillAppear, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		logEvent(client, event)

		settings, err := getPlayerSettingsFromWillAppearEvent(event)
		if err != nil {
			logError(client, event, err)
			return err
		}

		if settings.PlayerId == "" {
			logError(client, event, errors.New("No player configured"))
			return err
		}

		status, err := squeezebox.GetPlayerMode(settings.PlayerId)
		if err != nil {
			logError(client, event, err)
		} else {
			err = setImageForPlayMode(ctx, client, status)
			if err != nil {
				logError(client, event, err)
			}
		}

		return nil
	})

	playtoggleaction.RegisterHandler(streamdeck.WillAppear, selectPlayerHandlerWillAppear)
	playtoggleaction.RegisterHandler(streamdeck.SendToPlugin, selectPlayerHandlerSendToPlugin)

	// Play
	//
	playaction := client.Action("de.szoerner.streamdeck.squeezebox.actions.play")
	playaction.RegisterHandler(streamdeck.KeyDown, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		logEvent(client, event)

		settings, err := getPlayerSettingsFromKeyDownEvent(event)
		if (err == nil) {
			if settings.PlayerId == "" {
				client.ShowAlert(ctx)
				err = errors.New("No player configured")
			} else {
				_, err = squeezebox.SetPlayerMode(settings.PlayerId, "play")
			}
		}

		if err != nil {
			logError(client, event, err)
		}

		return err
	})
	playaction.RegisterHandler(streamdeck.WillAppear, selectPlayerHandlerWillAppear)
	playaction.RegisterHandler(streamdeck.SendToPlugin, selectPlayerHandlerSendToPlugin)

    // Pause
    //
	pauseaction := client.Action("de.szoerner.streamdeck.squeezebox.actions.pause")
	pauseaction.RegisterHandler(streamdeck.KeyDown, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		logEvent(client, event)

		settings, err := getPlayerSettingsFromKeyDownEvent(event)
		if (err == nil) {
			if settings.PlayerId == "" {
				client.ShowAlert(ctx)
				err = errors.New("No player configured")
			} else {
				_, err = squeezebox.SetPlayerMode(settings.PlayerId, "pause")
			}
		}

		if err != nil {
			logError(client, event, err)
		}

		return nil
	})

	pauseaction.RegisterHandler(streamdeck.WillAppear, selectPlayerHandlerWillAppear)
	pauseaction.RegisterHandler(streamdeck.SendToPlugin, selectPlayerHandlerSendToPlugin)
}

func setImageForPlayMode(ctx context.Context, client *streamdeck.Client, mode string) error {
	var err error = nil

	if mode == "play" {
		image, err := getImageByFilename("./images/PauseKey.png")
		if err == nil {
			err = client.SetImage(ctx, image, streamdeck.HardwareAndSoftware)
		}
	} else if mode == "stop" || mode == "pause" {
		image, err := getImageByFilename("./images/PlayKey.png")
		if err == nil {
			err = client.SetImage(ctx, image, streamdeck.HardwareAndSoftware)
		}
	}
	return err
}
package plugin

import (
	"context"
	"errors"
	"fmt"
	"github.com/StefanZoerner/streamdeck-squeezebox/plugin/keyimages"
	"github.com/StefanZoerner/streamdeck-squeezebox/squeezebox"
	"github.com/samwho/streamdeck"
)

func setupPlaymodeActions(client *streamdeck.Client) {

	// Play Toggle
	//
	playtoggleaction := client.Action("de.szoerner.streamdeck.squeezebox.actions.playtoggle")
	playtoggleaction.RegisterHandler(streamdeck.WillAppear, WillAppearRequestGlobalSettingsHandler)
	playtoggleaction.RegisterHandler(streamdeck.KeyDown, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		logEvent(client, event)

		settings, err := getPlayerSettingsFromKeyDownEvent(event)
		if err != nil {
			logError(client, event, err)
			return err
		}

		if settings.PlayerId == "" {
			client.ShowAlert(ctx)
			err = errors.New("No player configured")
			logError(client, event, err)
			return err
		}

		gs := GetPluginGlobalSettings()
		mode, err := squeezebox.TogglePlayerMode(gs.Hostname, gs.CliPort, settings.PlayerId)
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

		gs := GetPluginGlobalSettings()
		status, err := squeezebox.GetPlayerMode(gs.Hostname, gs.CliPort, settings.PlayerId)
		if err != nil {
			logError(client, event, err)
		} else {
			err = setImageForPlayMode(ctx, client, status)
			if err != nil {
				logError(client, event, err)
			}
		}

		// go updatePlayToggle(ctx, client, settings.PlayerId)

		return nil
	})

	playtoggleaction.RegisterHandler(streamdeck.WillAppear, selectPlayerHandlerWillAppear)
	playtoggleaction.RegisterHandler(streamdeck.SendToPlugin, selectPlayerHandlerSendToPlugin)

	// Play
	//
	playaction := client.Action("de.szoerner.streamdeck.squeezebox.actions.play")
	playaction.RegisterHandler(streamdeck.WillAppear, WillAppearRequestGlobalSettingsHandler)
	playaction.RegisterHandler(streamdeck.KeyDown, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		logEvent(client, event)

		settings, err := getPlayerSettingsFromKeyDownEvent(event)
		if err == nil {
			if settings.PlayerId == "" {
				client.ShowAlert(ctx)
				err = errors.New("No player configured")
			} else {
				gs := GetPluginGlobalSettings()
				_, err = squeezebox.SetPlayerMode(gs.Hostname, gs.CliPort, settings.PlayerId, "play")
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
	pauseaction.RegisterHandler(streamdeck.WillAppear, WillAppearRequestGlobalSettingsHandler)
	pauseaction.RegisterHandler(streamdeck.KeyDown, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		logEvent(client, event)

		settings, err := getPlayerSettingsFromKeyDownEvent(event)
		if err == nil {
			if settings.PlayerId == "" {
				client.ShowAlert(ctx)
				err = errors.New("No player configured")
			} else {
				globalSettings := GetPluginGlobalSettings()
				_, err = squeezebox.SetPlayerMode(globalSettings.Hostname, globalSettings.CliPort, settings.PlayerId, "pause")
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

	icon := ""
	switch mode {
	case "play":
		icon = "pause"
	case "pause", "stop":
		icon = "play"
	default:
		err = errors.New(fmt.Sprintf("Unknown play mode: %s", mode))
	}

	if err == nil {
		image, err := keyimages.GetStreamDeckImageForIcon(icon)
		if err == nil {
			err = client.SetImage(ctx, image, streamdeck.HardwareAndSoftware)
		}
	}

	return err
}

/*

func updatePlayToggle(ctx context.Context, client *streamdeck.Client, player_id string) {
	time.Sleep(5 * time.Second)
	for {
		gs := GetPluginGlobalSettings()
		mode, err := squeezebox.GetPlayerMode(gs.Hostname, gs.CLIPort, player_id)
		if err == nil {
		setImageForPlayMode(ctx, client, mode)
		} else {
			client.LogMessage(err.Error())
		}
	}
}

*/

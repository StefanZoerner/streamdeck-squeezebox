package plugin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/StefanZoerner/streamdeck-squeezebox/plugin/keyimages"
	"github.com/StefanZoerner/streamdeck-squeezebox/squeezebox"
	"github.com/samwho/streamdeck"
	sdcontext "github.com/samwho/streamdeck/context"
)

type PlayModeObserver struct {
	client *streamdeck.Client
	ctx    context.Context
}

func (pmo PlayModeObserver) playmodeChanged(s string) {
	err := setImageForPlayMode(pmo.ctx, pmo.client, s)
	if err != nil {
		pmo.client.LogMessage(err.Error())
	}
}

func (pmo PlayModeObserver) albumArtChanged(_ string) {
}

func (pmo PlayModeObserver) getID() string {
	return sdcontext.Context(pmo.ctx)
}

func (pmo PlayModeObserver) String() string {
	return "PlayModeObserver " + pmo.getID()[:5] + "..."
}

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
			err = errors.New("no player configured")
			logError(client, event, err)
			return err
		}

		gs := GetPluginGlobalSettings()
		mode, err := squeezebox.TogglePlayerMode(gs.Hostname, gs.CLIPort, settings.PlayerId)
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
			logError(client, event, errors.New("no player configured"))
			return err
		}

		gs := GetPluginGlobalSettings()
		status, err := squeezebox.GetPlayerMode(gs.Hostname, gs.CLIPort, settings.PlayerId)
		if err != nil {
			logError(client, event, err)
		} else {
			err = setImageForPlayMode(ctx, client, status)
			if err != nil {
				logError(client, event, err)
			}
		}

		pmo := PlayModeObserver{
			client: client,
			ctx:    ctx,
		}
		count := addOberserverForPlayer(settings.PlayerId, pmo)
		client.LogMessage(fmt.Sprintf("added %s for player %s, now total %d", pmo, settings.PlayerId, count))

		return nil
	})

	playtoggleaction.RegisterHandler(streamdeck.WillDisappear, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		logEvent(client, event)

		settings, err := getPlayerSettingsFromWillDisappearEvent(event)
		if err != nil {
			logError(client, event, err)
			return err
		}

		pmo := PlayModeObserver{
			client: client,
			ctx:    ctx,
		}
		count := removeOberserverForPlayer(settings.PlayerId, pmo)
		client.LogMessage(fmt.Sprintf("remove %s for player %s, now total %d", pmo, settings.PlayerId, count))

		return nil
	})

	playtoggleaction.RegisterHandler(streamdeck.WillAppear, selectPlayerHandlerWillAppear)
	playtoggleaction.RegisterHandler(streamdeck.SendToPlugin, playToggleHandlerSendToPlugin)

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
				err = errors.New("no player configured")
			} else {
				gs := GetPluginGlobalSettings()
				_, err = squeezebox.SetPlayerMode(gs.Hostname, gs.CLIPort, settings.PlayerId, "play")
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
				err = errors.New("no player configured")
			} else {
				globalSettings := GetPluginGlobalSettings()
				_, err = squeezebox.SetPlayerMode(globalSettings.Hostname, globalSettings.CLIPort, settings.PlayerId, "pause")
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

func playToggleHandlerSendToPlugin(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
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

		pmo := PlayModeObserver{
			client: client,
			ctx:    ctx,
		}
		removeOberserverForAllPlayers(pmo)
		client.LogMessage(fmt.Sprintf("remove observer for all players"))

		playerID := fromPI.Value
		count := addOberserverForPlayer(playerID, pmo)
		client.LogMessage(fmt.Sprintf("add observer for player %s, now %d", playerID, count))

		pinfo, err := squeezebox.GetPlayerInfo(globalSettings.Hostname, globalSettings.CLIPort, playerID)
		if err != nil {
			logError(client, event, err)
			return err
		}

		np := PlayerSettings{
			PlayerId:   playerID,
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

func setImageForPlayMode(ctx context.Context, client *streamdeck.Client, mode string) error {
	var err error = nil

	icon := ""
	switch mode {
	case "play":
		icon = "pause"
	case "pause", "stop":
		icon = "play"
	default:
		err = fmt.Errorf("unknown play mode: %s", mode)
	}

	if err == nil {
		image, err := keyimages.GetStreamDeckImageForIcon(icon)
		if err == nil {
			err = client.SetImage(ctx, image, streamdeck.HardwareAndSoftware)
		}
	}

	return err
}

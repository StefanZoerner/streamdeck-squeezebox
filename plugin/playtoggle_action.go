package plugin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/StefanZoerner/streamdeck-squeezebox/plugin/general"
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

func setupPlaytoggleAction(client *streamdeck.Client) {

	// Play Toggle
	//
	playtoggleaction := client.Action("de.szoerner.streamdeck.squeezebox.actions.playtoggle")

	playtoggleaction.RegisterHandler(streamdeck.WillAppear, WillAppearRequestGlobalSettingsHandler)
	playtoggleaction.RegisterHandler(streamdeck.KeyDown, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		general.LogEvent(client, event)

		settings, err := getPlayerSettingsFromKeyDownEvent(event)
		if err != nil {
			general.LogErrorWithEvent(client, event, err)
			return err
		}

		if settings.PlayerId == "" {
			client.ShowAlert(ctx)
			err = errors.New("no player configured")
			general.LogErrorWithEvent(client, event, err)
			return err
		}

		gs := GetPluginGlobalSettings()
		mode, err := squeezebox.TogglePlayerMode(gs.connectionProps(), settings.PlayerId)
		if err != nil {
			client.ShowAlert(ctx)
		} else {
			err = setImageForPlayMode(ctx, client, mode)
		}

		if err != nil {
			general.LogErrorWithEvent(client, event, err)
		}
		return err
	})

	playtoggleaction.RegisterHandler(streamdeck.WillAppear, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		general.LogEvent(client, event)

		settings, err := getPlayerSettingsFromWillAppearEvent(event)
		if err != nil {
			general.LogErrorWithEvent(client, event, err)
			return err
		}

		if settings.PlayerId == "" {
			general.LogErrorWithEvent(client, event, errors.New("no player configured"))
			return err
		}

		gs := GetPluginGlobalSettings()
		status, err := squeezebox.GetPlayerMode(gs.connectionProps(), settings.PlayerId)
		if err != nil {
			general.LogErrorWithEvent(client, event, err)
		} else {
			err = setImageForPlayMode(ctx, client, status)
			if err != nil {
				general.LogErrorWithEvent(client, event, err)
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
		general.LogEvent(client, event)

		settings, err := getPlayerSettingsFromWillDisappearEvent(event)
		if err != nil {
			general.LogErrorWithEvent(client, event, err)
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
}

func playToggleHandlerSendToPlugin(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
	general.LogEvent(client, event)

	fromPI := DataFromPlayerSelectionPI{}
	err := json.Unmarshal(event.Payload, &fromPI)
	if err != nil {
		general.LogErrorWithEvent(client, event, err)
		return err
	}

	globalSettings := GetPluginGlobalSettings()

	if fromPI.Command == "getPlayerSelectionOptions" {

		conProps := globalSettings.connectionProps()

		players, err := squeezebox.GetPlayerInfos(conProps)
		if err != nil {
			general.LogErrorWithEvent(client, event, err)
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
			general.LogErrorWithEvent(client, event, err)
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

		conProps := globalSettings.connectionProps()

		pinfo, err := squeezebox.GetPlayerInfo(conProps, playerID)
		if err != nil {
			general.LogErrorWithEvent(client, event, err)
			return err
		}

		np := PlayerSettings{
			PlayerId:   playerID,
			PlayerName: pinfo.Name,
		}

		err = client.SetSettings(ctx, np)
		if err != nil {
			general.LogErrorWithEvent(client, event, err)
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

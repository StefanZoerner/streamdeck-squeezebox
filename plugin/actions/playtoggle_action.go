package actions

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

func (pmo PlayModeObserver) PlaymodeChanged(s string) {
	err := setImageForPlayMode(pmo.ctx, pmo.client, s)
	if err != nil {
		pmo.client.LogMessage(err.Error())
	}
}

func (pmo PlayModeObserver) AlbumArtChanged(_ string) {
}

func (pmo PlayModeObserver) GetID() string {
	return sdcontext.Context(pmo.ctx)
}

func (pmo PlayModeObserver) String() string {
	return "PlayModeObserver " + pmo.GetID()[:5] + "..."
}

func SetupPlaytoggleAction(client *streamdeck.Client) {

	// Play Toggle
	//
	playtoggleaction := client.Action("de.szoerner.streamdeck.squeezebox.actions.playtoggle")

	playtoggleaction.RegisterHandler(streamdeck.WillAppear, general.WillAppearRequestGlobalSettingsHandler)
	playtoggleaction.RegisterHandler(streamdeck.KeyDown, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		general.LogEvent(client, event)

		settings, err := GetPlayerSettingsFromKeyDownEvent(event)
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

		gs := general.GetPluginGlobalSettings()
		mode, err := squeezebox.TogglePlayerMode(gs.ConnectionProps(), settings.PlayerId)
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

		settings, err := GetPlayerSettingsFromWillAppearEvent(event)
		if err != nil {
			general.LogErrorWithEvent(client, event, err)
			return err
		}

		if settings.PlayerId == "" {
			general.LogErrorWithEvent(client, event, errors.New("no player configured"))
			return err
		}

		gs := general.GetPluginGlobalSettings()
		status, err := squeezebox.GetPlayerMode(gs.ConnectionProps(), settings.PlayerId)
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
		count := general.AddOberserverForPlayer(settings.PlayerId, pmo)
		client.LogMessage(fmt.Sprintf("added %s for player %s, now total %d", pmo, settings.PlayerId, count))

		return nil
	})

	playtoggleaction.RegisterHandler(streamdeck.WillDisappear, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		general.LogEvent(client, event)

		settings, err := GetPlayerSettingsFromWillDisappearEvent(event)
		if err != nil {
			general.LogErrorWithEvent(client, event, err)
			return err
		}

		pmo := PlayModeObserver{
			client: client,
			ctx:    ctx,
		}
		count := general.RemoveOberserverForPlayer(settings.PlayerId, pmo)
		client.LogMessage(fmt.Sprintf("remove %s for player %s, now total %d", pmo, settings.PlayerId, count))

		return nil
	})

	playtoggleaction.RegisterHandler(streamdeck.WillAppear, SelectPlayerHandlerWillAppear)
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

	globalSettings := general.GetPluginGlobalSettings()

	if fromPI.Command == "getPlayerSelectionOptions" {

		payload, err := getPlayerSelection()
		if err == nil {
			err = client.SendToPropertyInspector(ctx, &payload)
		}

		if err != nil {
			general.LogErrorWithEvent(client, event, err)
			return err
		}
	} else if fromPI.Command == "setSelectedPlayer" {

		pmo := PlayModeObserver{
			client: client,
			ctx:    ctx,
		}
		general.RemoveOberserverForAllPlayers(pmo)
		client.LogMessage(fmt.Sprintf("remove PlayerObserver for all players"))

		playerID := fromPI.Value
		count := general.AddOberserverForPlayer(playerID, pmo)
		client.LogMessage(fmt.Sprintf("add PlayerObserver for player %s, now %d", playerID, count))

		conProps := globalSettings.ConnectionProps()

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

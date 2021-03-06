package actions

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/StefanZoerner/streamdeck-squeezebox/plugin/general"
	"github.com/StefanZoerner/streamdeck-squeezebox/plugin/keyimages"
	"github.com/StefanZoerner/streamdeck-squeezebox/squeezebox"
	"github.com/samwho/streamdeck"
	sdcontext "github.com/samwho/streamdeck/context"
)

type playModeObserver struct {
	client *streamdeck.Client
	ctx    context.Context
}

func (pmo playModeObserver) PlaymodeChanged(s string) {
	err := setImageForPlayMode(pmo.ctx, pmo.client, s)
	if err != nil {
		_ = pmo.client.LogMessage(err.Error())
	}
}

func (pmo playModeObserver) AlbumArtChanged(_ string) {
}

func (pmo playModeObserver) GetID() string {
	return sdcontext.Context(pmo.ctx)
}

func (pmo playModeObserver) String() string {
	return "playModeObserver " + pmo.GetID()[:5] + "..."
}

// SetupPlaytoggleAction adds the playtoggle action to the plugin.
//
func SetupPlaytoggleAction(client *streamdeck.Client) {

	playtoggleaction := client.Action("de.szoerner.streamdeck.squeezebox.actions.playtoggle")

	playtoggleaction.RegisterHandler(streamdeck.WillAppear, general.WillAppearRequestGlobalSettingsHandler)
	playtoggleaction.RegisterHandler(streamdeck.KeyDown, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		general.LogEvent(client, event)

		settings, err := GetPlayerSettingsFromKeyDownEvent(event)
		if err != nil {
			general.LogErrorWithEvent(client, event, err)
			return err
		}

		playerID := settings.PlayerId
		globalSettings := general.GetPluginGlobalSettings()
		if playerID == "" {
			playerID = globalSettings.DefaultPlayerID
		}

		if playerID == "" {
			_ = client.ShowAlert(ctx)
			return errors.New("No player configured")
		}

		mode, err := squeezebox.TogglePlayerMode(globalSettings.ConnectionProps(), playerID)
		if err != nil {
			_ = client.ShowAlert(ctx)
		} else {
			err = setImageForPlayMode(ctx, client, mode)
		}

		if err != nil {
			general.LogErrorWithEvent(client, event, err)
		}
		return err
	})

	playtoggleaction.RegisterHandler(streamdeck.WillAppear, SelectPlayerHandlerWillAppear)
	playtoggleaction.RegisterHandler(streamdeck.WillAppear, playtoggleHandlerWillAppear)
	playtoggleaction.RegisterHandler(streamdeck.WillDisappear, playtoggleHandlerWillDisappear)
	playtoggleaction.RegisterHandler(streamdeck.SendToPlugin, playtoggleHandlerSendToPlugin)
}

func playtoggleHandlerWillAppear(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
	general.LogEvent(client, event)

	settings, err := GetPlayerSettingsFromWillAppearEvent(event)
	if err != nil {
		general.LogErrorWithEvent(client, event, err)
		return err
	}

	pmo := playModeObserver{
		client: client,
		ctx:    ctx,
	}
	general.AddOberserverForPlayer(settings.PlayerId, pmo)

	gs := general.GetPluginGlobalSettings()
	conProps := gs.ConnectionProps()
	if conProps.NotEmpty() {

		pid := settings.PlayerId
		if pid == "" {
			pid = gs.DefaultPlayerID
		}

		if pid != "" {

			status, err := squeezebox.GetPlayerMode(gs.ConnectionProps(), pid)
			if err != nil {
				general.LogErrorWithEvent(client, event, err)
			} else {
				err = setImageForPlayMode(ctx, client, status)
				if err != nil {
					general.LogErrorWithEvent(client, event, err)
				}
			}
		}
	}

	return nil
}

func playtoggleHandlerWillDisappear(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
	general.LogEvent(client, event)

	settings, err := GetPlayerSettingsFromWillDisappearEvent(event)
	if err != nil {
		general.LogErrorWithEvent(client, event, err)
		return err
	}

	pmo := playModeObserver{
		client: client,
		ctx:    ctx,
	}
	general.RemoveOberserverForPlayer(settings.PlayerId, pmo)

	return nil
}

func playtoggleHandlerSendToPlugin(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
	general.LogEvent(client, event)

	var err error

	fromPI := DataFromPlayerSelectionPI{}
	err = json.Unmarshal(event.Payload, &fromPI)
	if err != nil {
		general.LogErrorWithEvent(client, event, err)
		return err
	}

	if fromPI.Command == "getPlayerSelectionOptions" {
		var payload PlayerSelection

		payload, err = getPlayerSelection()
		if err == nil {
			err = client.SendToPropertyInspector(ctx, &payload)
		}
	} else if fromPI.Command == "setSelectedPlayer" {

		playerID := fromPI.Value

		pmo := playModeObserver{
			client: client,
			ctx:    ctx,
		}
		general.RemoveOberserverForAllPlayers(pmo)
		general.AddOberserverForPlayer(playerID, pmo)

		globalSettings := general.GetPluginGlobalSettings()
		conProps := globalSettings.ConnectionProps()
		var pinfo *squeezebox.PlayerInfo

		pinfo, err = squeezebox.GetPlayerInfo(conProps, playerID)
		if err == nil {
			np := PlayerSettings{
				PlayerId:   playerID,
				PlayerName: pinfo.Name,
			}
			err = client.SetSettings(ctx, np)
		}
	}

	if err != nil {
		general.LogErrorWithEvent(client, event, err)
	}

	return err
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
		icon = "play_pause"
	}

	if err == nil {
		var image string
		image, err = keyimages.GetStreamDeckImageForIcon(icon)
		if err == nil {
			err = client.SetImage(ctx, image, streamdeck.HardwareAndSoftware)
		}
	}

	return err
}

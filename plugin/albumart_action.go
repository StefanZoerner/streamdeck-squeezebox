package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/StefanZoerner/streamdeck-squeezebox/plugin/keyimages"
	sdcontext "github.com/samwho/streamdeck/context"
	"image"

	"github.com/StefanZoerner/streamdeck-squeezebox/squeezebox"
	"github.com/samwho/streamdeck"
)

type AlbumArtActionSettings struct {
	PlayerSettings
	Dimension  string `json:"albumart_dimension"`
	TileNumber int    `json:"albumart_tile_number"`
}

type AlbumArtFromPI struct {
	Command  string                 `json:"command"`
	Settings AlbumArtActionSettings `json:"settings"`
}

type AlbumArtObserver struct {
	client    *streamdeck.Client
	ctx       context.Context
	dimension string
	tile      int
}

func (aao AlbumArtObserver) playmodeChanged(_ string) {
}

func (aao AlbumArtObserver) albumArtChanged(newURL string) {
	err := showAlbumArtImage(aao.ctx, aao.client, newURL, aao.dimension, aao.tile)
	if err != nil {
		logErrorNoEvent(aao.client, err)
	}
}

func (aao AlbumArtObserver) getID() string {
	return sdcontext.Context(aao.ctx)
}

func (aao AlbumArtObserver) String() string {
	return "AlbumArtObserver " + aao.getID()[:5] + "..."
}

func setupAlbumArtAction(client *streamdeck.Client) {
	albumArtAction := client.Action("de.szoerner.streamdeck.squeezebox.actions.albumart")
	albumArtAction.RegisterHandler(streamdeck.WillAppear, WillAppearRequestGlobalSettingsHandler)

	albumArtAction.RegisterHandler(streamdeck.SendToPlugin, albumArtSendToPlugin)
	albumArtAction.RegisterHandler(streamdeck.WillAppear, albumArtWillAppear)
	albumArtAction.RegisterHandler(streamdeck.WillDisappear, albumArtWillDisappear)
}

func albumArtWillAppear(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
	logEvent(client, event)

	payload := streamdeck.WillAppearPayload{}
	err := json.Unmarshal(event.Payload, &payload)
	if err != nil {
		logError(client, event, err)
		return err
	}

	settings := AlbumArtActionSettings{}
	err = json.Unmarshal(payload.Settings, &settings)
	if err != nil {
		logError(client, event, err)
		return err
	}

	modified := false

	if settings.PlayerId == "" {
		settings.PlayerName = "(None)"
		modified = true
	}
	if settings.Dimension == "" {
		settings.Dimension = keyimages.ALBUM_ART_1x1
		modified = true
	}
	if settings.TileNumber == 0 {
		settings.TileNumber = 1
		modified = true
	}

	if modified {
		err = client.SetSettings(ctx, settings)
		if err != nil {
			logError(client, event, err)
			return err
		}
	}

	aao := AlbumArtObserver{
		client:    client,
		ctx:       ctx,
		dimension: settings.Dimension,
		tile:      settings.TileNumber,
	}
	count := addOberserverForPlayer(settings.PlayerId, aao)
	client.LogMessage(fmt.Sprintf("added %s for player %s, now total %d", aao, settings.PlayerId, count))

	conProps := GetPluginGlobalSettings().connectionProps()
	url, err := squeezebox.GetCurrentArtworkURL(conProps, settings.PlayerId)
	if err != nil {
		logError(client, event, err)
		return err
	}

	err = showAlbumArtImage(ctx, client, url, settings.Dimension, settings.TileNumber)
	if err != nil {
		logError(client, event, err)
	}

	return err
}

func albumArtWillDisappear(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
	logEvent(client, event)

	settings, err := getPlayerSettingsFromWillDisappearEvent(event)
	if err != nil {
		logError(client, event, err)
		return err
	}

	aao := AlbumArtObserver{
		client: client,
		ctx:    ctx,
	}
	count := removeOberserverForPlayer(settings.PlayerId, aao)
	client.LogMessage(fmt.Sprintf("remove %s for player %s, now total %d", aao, settings.PlayerId, count))

	return nil

}

func albumArtSendToPlugin(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
	logEvent(client, event)

	fromPI := AlbumArtFromPI{}
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
	} else if fromPI.Command == "sendFormData" {

		err = client.SetSettings(ctx, fromPI.Settings)
		if err != nil {
			logError(client, event, err)
			return err
		}

		aao := AlbumArtObserver{
			client:    client,
			ctx:       ctx,
			dimension: fromPI.Settings.Dimension,
			tile:      fromPI.Settings.TileNumber,
		}
		removeOberserverForAllPlayers(aao)
		client.LogMessage(fmt.Sprintf("removed %s for all players", aao))

		count := addOberserverForPlayer(fromPI.Settings.PlayerId, aao)
		client.LogMessage(fmt.Sprintf("added %s for player %s, now total %d", aao, fromPI.Settings.PlayerId, count))

		cp := GetPluginGlobalSettings().connectionProps()
		url, err := squeezebox.GetCurrentArtworkURL(cp, fromPI.Settings.PlayerId)
		if err != nil {
			logError(client, event, err)
			return err
		}

		err = showAlbumArtImage(ctx, client, url, fromPI.Settings.Dimension, fromPI.Settings.TileNumber)
		if err != nil {
			logError(client, event, err)
			return err
		}

	}

	return nil
}

func showAlbumArtImage(ctx context.Context, client *streamdeck.Client, url, dim string, tile int) error {

	var img image.Image
	var err error

	if url == "" {
		img, err = keyimages.GetImageByFilename("./assets/images/album_art_default.png")
	} else {
		img, err = keyimages.GetImageByUrl(url)
	}

	if err != nil {
		logErrorNoEvent(client, err)
		return err
	}

	tileImage, err := keyimages.ResizeAndCropImage(img, dim, tile)
	if err != nil {
		logErrorNoEvent(client, err)
		return err
	}

	s, _ := streamdeck.Image(tileImage)
	err = client.SetImage(ctx, s, streamdeck.HardwareAndSoftware)
	if err != nil {
		logErrorNoEvent(client, err)
		return err
	}

	return nil
}

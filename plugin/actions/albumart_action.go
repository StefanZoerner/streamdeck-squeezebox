package actions

import (
	"context"
	"encoding/json"
	"github.com/StefanZoerner/streamdeck-squeezebox/plugin/general"
	"github.com/StefanZoerner/streamdeck-squeezebox/plugin/keyimages"
	sdcontext "github.com/samwho/streamdeck/context"
	"image"

	"github.com/StefanZoerner/streamdeck-squeezebox/squeezebox"
	"github.com/samwho/streamdeck"
)

type albumArtActionSettings struct {
	PlayerSettings
	Dimension  string `json:"albumart_dimension"`
	TileNumber int    `json:"albumart_tile_number"`
}

type albumArtFromPI struct {
	Command  string                 `json:"command"`
	Settings albumArtActionSettings `json:"settings"`
}

type albumArtObserver struct {
	client    *streamdeck.Client
	ctx       context.Context
	dimension string
	tile      int
}

func (aao albumArtObserver) PlaymodeChanged(_ string) {
}

func (aao albumArtObserver) AlbumArtChanged(newURL string) {
	err := showAlbumArtImage(aao.ctx, aao.client, newURL, aao.dimension, aao.tile)
	if err != nil {
		general.LogErrorNoEvent(aao.client, err)
	}
}

func (aao albumArtObserver) GetID() string {
	return sdcontext.Context(aao.ctx)
}

func (aao albumArtObserver) String() string {
	return "albumArtObserver " + aao.GetID()[:5] + "..."
}

// SetupAlbumArtAction adds the albumart action to the plugin.
//
func SetupAlbumArtAction(client *streamdeck.Client) {
	albumArtAction := client.Action("de.szoerner.streamdeck.squeezebox.actions.albumart")
	albumArtAction.RegisterHandler(streamdeck.WillAppear, general.WillAppearRequestGlobalSettingsHandler)

	albumArtAction.RegisterHandler(streamdeck.SendToPlugin, albumArtSendToPlugin)
	albumArtAction.RegisterHandler(streamdeck.WillAppear, albumArtWillAppear)
	albumArtAction.RegisterHandler(streamdeck.WillDisappear, albumArtWillDisappear)
}

func albumArtWillAppear(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
	general.LogEvent(client, event)

	var err error

	payload := streamdeck.WillAppearPayload{}
	err = json.Unmarshal(event.Payload, &payload)
	if err != nil {
		general.LogErrorWithEvent(client, event, err)
		return err
	}

	settings := albumArtActionSettings{}
	err = json.Unmarshal(payload.Settings, &settings)
	if err != nil {
		general.LogErrorWithEvent(client, event, err)
		return err
	}

	modified := false

	if settings.PlayerId == "" {
		settings.PlayerName = "(None)"
		modified = true
	}
	if settings.Dimension == "" {
		settings.Dimension = keyimages.AlbumArt1x1
		modified = true
	}
	if settings.TileNumber == 0 {
		settings.TileNumber = 1
		modified = true
	}

	if modified {
		err = client.SetSettings(ctx, settings)
		if err != nil {
			general.LogErrorWithEvent(client, event, err)
			return err
		}
	}

	aao := albumArtObserver{
		client:    client,
		ctx:       ctx,
		dimension: settings.Dimension,
		tile:      settings.TileNumber,
	}
	general.AddOberserverForPlayer(settings.PlayerId, aao)

	conProps := general.GetPluginGlobalSettings().ConnectionProps()

	var url string
	if conProps.NotEmpty() {
		url, err = squeezebox.GetCurrentArtworkURL(conProps, settings.PlayerId)
		if err != nil {
			general.LogErrorWithEvent(client, event, err)
			return err
		}
	}

	err = showAlbumArtImage(ctx, client, url, settings.Dimension, settings.TileNumber)
	if err != nil {
		general.LogErrorWithEvent(client, event, err)
	}

	return err
}

func albumArtWillDisappear(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
	general.LogEvent(client, event)

	var err error
	var settings PlayerSettings

	settings, err = GetPlayerSettingsFromWillDisappearEvent(event)
	if err == nil {
		aao := albumArtObserver{
			client: client,
			ctx:    ctx,
		}
		general.RemoveOberserverForPlayer(settings.PlayerId, aao)
	}

	if err != nil {
		general.LogErrorWithEvent(client, event, err)
	}
	return err
}

func albumArtSendToPlugin(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
	general.LogEvent(client, event)

	fromPI := albumArtFromPI{}
	err := json.Unmarshal(event.Payload, &fromPI)
	if err != nil {
		general.LogErrorWithEvent(client, event, err)
		return err
	}

	if fromPI.Command == "getPlayerSelectionOptions" {

		payload, err := getPlayerSelection()
		if err == nil {
			err = client.SendToPropertyInspector(ctx, &payload)
		}

		if err != nil {
			general.LogErrorWithEvent(client, event, err)
			return err
		}
	} else if fromPI.Command == "sendFormData" {

		err = client.SetSettings(ctx, fromPI.Settings)
		if err != nil {
			general.LogErrorWithEvent(client, event, err)
			return err
		}

		aao := albumArtObserver{
			client:    client,
			ctx:       ctx,
			dimension: fromPI.Settings.Dimension,
			tile:      fromPI.Settings.TileNumber,
		}
		general.RemoveOberserverForAllPlayers(aao)
		general.AddOberserverForPlayer(fromPI.Settings.PlayerId, aao)

		cp := general.GetPluginGlobalSettings().ConnectionProps()
		url, err := squeezebox.GetCurrentArtworkURL(cp, fromPI.Settings.PlayerId)
		if err != nil {
			general.LogErrorWithEvent(client, event, err)
			return err
		}

		err = showAlbumArtImage(ctx, client, url, fromPI.Settings.Dimension, fromPI.Settings.TileNumber)
		if err != nil {
			general.LogErrorWithEvent(client, event, err)
			return err
		}

	}

	return nil
}

func showAlbumArtImage(ctx context.Context, client *streamdeck.Client, url, dim string, tile int) error {

	var img image.Image
	var err error

	img, err = keyimages.GetImageByUrl(url)
	if err != nil {
		general.LogErrorNoEvent(client, err)
	}

	// Ignore error (if any)  and try to render (default) image, if availaible
	if img != nil {

		var tileImage image.Image

		tileImage, err = keyimages.ResizeAndCropImage(img, dim, tile)
		if err == nil {
			s, _ := streamdeck.Image(tileImage)
			err = client.SetImage(ctx, s, streamdeck.HardwareAndSoftware)
		}
	}

	if err != nil {
		general.LogErrorNoEvent(client, err)
	}
	return err
}

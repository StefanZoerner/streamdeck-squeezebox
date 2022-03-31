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

	albumArtAction.RegisterHandler(streamdeck.SendToPlugin, albumartHandlerSendToPlugin)
	albumArtAction.RegisterHandler(streamdeck.WillAppear, albumartHandlerWillAppear)
	albumArtAction.RegisterHandler(streamdeck.WillDisappear, albumartHandlerWillDisappear)
}

func albumartHandlerWillAppear(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) (err error) {
	general.LogEvent(client, event)

	defer func() {
		if err != nil {
			general.LogErrorWithEvent(client, event, err)
		}
	}()

	payload := streamdeck.WillAppearPayload{}
	err = json.Unmarshal(event.Payload, &payload)

	if err == nil {
		settings := albumArtActionSettings{}
		err = json.Unmarshal(payload.Settings, &settings)

		if err == nil {

			modified := false
			if settings.PlayerId == "" {
				settings.PlayerName = "(Default)"
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

				pid := settings.PlayerId
				if pid == "" {
					pid = general.GetPluginGlobalSettings().DefaultPlayerID
				}

				if pid != "" {
					url, err = squeezebox.GetCurrentArtworkURL(conProps, pid)
					if err != nil {
						return err
					}
				}
			}

			err = showAlbumArtImage(ctx, client, url, settings.Dimension, settings.TileNumber)
		}
	}

	return err
}

func albumartHandlerWillDisappear(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
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

func albumartHandlerSendToPlugin(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
	general.LogEvent(client, event)

	var err error

	fromPI := albumArtFromPI{}
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
	} else if fromPI.Command == "sendFormData" {
		err = client.SetSettings(ctx, fromPI.Settings)
		if err == nil {
			aao := albumArtObserver{
				client:    client,
				ctx:       ctx,
				dimension: fromPI.Settings.Dimension,
				tile:      fromPI.Settings.TileNumber,
			}
			general.RemoveOberserverForAllPlayers(aao)
			general.AddOberserverForPlayer(fromPI.Settings.PlayerId, aao)

			var url string
			cp := general.GetPluginGlobalSettings().ConnectionProps()

			pid := fromPI.Settings.PlayerId
			if pid == "" {
				pid = general.GetPluginGlobalSettings().DefaultPlayerID
			}

			url, err = squeezebox.GetCurrentArtworkURL(cp, pid)
			if err == nil {
				err = showAlbumArtImage(ctx, client, url, fromPI.Settings.Dimension, fromPI.Settings.TileNumber)
			}
		}
	}

	return err
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

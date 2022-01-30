package plugin

import (
	"context"
	"encoding/json"
	"github.com/StefanZoerner/streamdeck-squeezebox/plugin/keyimages"
	"image"

	"github.com/StefanZoerner/streamdeck-squeezebox/squeezebox"
	"github.com/samwho/streamdeck"
)

type AlbumArtActionSettings struct {
	PlayerSettings
	Dimension string `json:"albumart_dimension"`
	TileNumber int   `json:"albumart_tile_number"`
}

type AlbumArtFromPI struct {
	Command  string                 `json:"command"`
	Settings AlbumArtActionSettings `json:"settings"`
}


func setupAlbumArtAction(client *streamdeck.Client) {
	albumArtAction := client.Action("de.szoerner.streamdeck.squeezebox.actions.albumart")
	albumArtAction.RegisterHandler(streamdeck.WillAppear, WillAppearRequestGlobalSettingsHandler)

	albumArtAction.RegisterHandler(streamdeck.SendToPlugin, albumArtSendToPlugin)
	albumArtAction.RegisterHandler(streamdeck.WillAppear, albumArtWillAppear)
}


func albumArtWillAppear (ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
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

	// TODO: get From Global Props
	conProps := squeezebox.NewConnectionProperties("elfman", 9002, 9090)

	url, err := squeezebox.GetCurrentArtworkUrl(conProps, settings.PlayerId)
	if err != nil {
		logError(client, event, err)
		return err
	}


	err = showAlbumArtImage(ctx, client, event, url, settings.Dimension, settings.TileNumber)
	if err != nil {
		logError(client, event, err)
	}


	return err
}

func albumArtSendToPlugin (ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
	logEvent(client, event)

	fromPI := AlbumArtFromPI{}
	err := json.Unmarshal(event.Payload, &fromPI)
	if err != nil {
		logError(client, event, err)
		return err
	}

	globalSettings := GetPluginGlobalSettings()

	if fromPI.Command == "getPlayerSelectionOptions" {

		players, err := squeezebox.GetPlayers(globalSettings.Hostname, globalSettings.CliPort)
		if (err != nil) {
			logError(client, event, err)
			return err
		}

		playerSettings := []PlayerSettings{}
		for _, p := range players {
			np := PlayerSettings{
				PlayerId:   p.Id,
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
		} else {
			err = showAlbumArtImage(ctx, client, event, "", fromPI.Settings.Dimension, fromPI.Settings.TileNumber)
			if err != nil {
				logError(client, event, err)
				return err
			}
		}

	}

	return nil
}

func showAlbumArtImage(ctx context.Context, client *streamdeck.Client, event streamdeck.Event, url, dim string, tile int) error {

	var img image.Image
	var err error

	if url == "" {
		img, err = keyimages.GetImageByFilename("./assets/images/album_art_default.png")
	} else {
		img, err = keyimages.GetImageByUrl(url)
	}

	if err != nil {
		logError(client, event, err)
		return err
	}

	tileImage, err := keyimages.ResizeAndCropImage(img, dim, tile)
	if err != nil {
		logError(client, event, err)
		return err
	}

	s, _ := streamdeck.Image(tileImage)
	err = client.SetImage(ctx, s, streamdeck.HardwareAndSoftware)
	if err != nil {
		logError(client, event, err)
		return err
	}

	return nil
}
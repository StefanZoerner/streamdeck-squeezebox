package plugin

import (
	"context"
	"github.com/StefanZoerner/streamdeck-squeezebox/plugin/general"
	"github.com/samwho/streamdeck"
	"os"
)

func Run(ctx context.Context) error {
	params, err := streamdeck.ParseRegistrationParams(os.Args)
	if err != nil {
		return err
	}

	general.PluginUUID = params.PluginUUID
	client := streamdeck.NewClient(ctx, params)

	client.RegisterHandler(streamdeck.DidReceiveGlobalSettings, general.DidReceiveGlobalSettingsHandler)
	setup(client)

	go general.StartTicker()

	return client.Run()
}

func setup(client *streamdeck.Client) {
	setupConfigurationAction(client)
	setupVolumeAction(client)
	setupPlaytoggleAction(client)
	setupTrackActions(client)
	setupAlbumArtAction(client)
}

package plugin

import (
	"github.com/samwho/streamdeck"
	"strings"
)

func logError(client *streamdeck.Client, event streamdeck.Event, err error) {
	action_name := getActionNameFromEvent(event)
	client.LogMessage("Error in "+action_name + ": "+err.Error())
}

func logEvent(client *streamdeck.Client, event streamdeck.Event) {

	// Determine last part of dot divided action name
	action_name := getActionNameFromEvent(event)

	msg := action_name + " " + event.Event + " "
	client.LogMessage("Event : " +msg)
	pl, _ := event.Payload.MarshalJSON()

	client.LogMessage("Payload: "+string(pl)+"\n")
}

func getActionNameFromEvent(event streamdeck.Event) string {
	action_name := "???"
	actionParts := strings.Split(event.Action, ".");
	if len(actionParts) > 0 {
		action_name = actionParts[len(actionParts)-1]
	}
	return action_name
}
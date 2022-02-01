package plugin

import (
	"github.com/samwho/streamdeck"
	"strings"
)

func logError(client *streamdeck.Client, event streamdeck.Event, err error) {
	actionName := getActionNameFromEvent(event)
	client.LogMessage("Error in " + actionName + ": " + err.Error())
}

func logEvent(client *streamdeck.Client, event streamdeck.Event) {

	// Determine last part of dot divided action name
	actionName := getActionNameFromEvent(event)

	msg := actionName + " " + event.Event + " "
	client.LogMessage("Event : " + msg)
	pl, _ := event.Payload.MarshalJSON()

	client.LogMessage("Payload: " + string(pl) + "\n")
}

func getActionNameFromEvent(event streamdeck.Event) string {
	actionName := "???"
	actionParts := strings.Split(event.Action, ".")
	if len(actionParts) > 0 {
		actionName = actionParts[len(actionParts)-1]
	}
	return actionName
}

package processor

import (
	dialogFlow "cloud.google.com/go/dialogflow/apiv2"
	"context"
)

type DialogFlowProcessor struct {
	projectID string
	authJSONFilePath string
	lang string
	timeZone string
	sessionClient *dialogFlow.SessionsClient
	ctx context.Context
}

type NLPResponse struct {
	Intent string
	Confidence float32
	Entities map[string]string
}

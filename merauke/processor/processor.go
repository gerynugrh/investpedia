package processor

import (
	dialogFlow "cloud.google.com/go/dialogflow/apiv2"
	"context"
	"fmt"
	structPB "github.com/golang/protobuf/ptypes/struct"
	"google.golang.org/api/option"
	dialogFlowPB "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
	"log"
	"strconv"
)

func (dp *DialogFlowProcessor) init(projectID string, authJSONFilePath string, lang string, timeZone string) (err error) {
	dp.projectID = projectID
	dp.authJSONFilePath = authJSONFilePath
	dp.lang = lang
	dp.timeZone = timeZone

	dp.ctx = context.Background()
	sessionClient, err := dialogFlow.NewSessionsClient(dp.ctx, option.WithCredentialsFile(authJSONFilePath))
	if err != nil {
		log.Fatal("Error in auth with dialog flow")
		return
	}
	dp.sessionClient = sessionClient
	return
}

func (dp *DialogFlowProcessor) processNLP(rawMessage string, username string) (r NLPResponse) {
	sessionID := username
	request := dialogFlowPB.DetectIntentRequest{
		Session: fmt.Sprintf("projects/%s/agent/sessions/%s", dp.projectID, sessionID),
		QueryInput: &dialogFlowPB.QueryInput{
			Input: &dialogFlowPB.QueryInput_Text{
				Text: &dialogFlowPB.TextInput{
					Text:         rawMessage,
					LanguageCode: dp.lang,
				},
			},
		},
		QueryParams: &dialogFlowPB.QueryParameters{
			TimeZone: dp.timeZone,
		},
	}
	response, err := dp.sessionClient.DetectIntent(dp.ctx, &request)
	if err != nil {
		log.Fatalf("Error in communication with Dialogflow %s", err.Error())
		return
	}
	queryResult := response.GetQueryResult()
	if queryResult.Intent != nil {
		r.Intent = queryResult.Intent.DisplayName
		r.Confidence = queryResult.IntentDetectionConfidence
	}
	r.Entities = make(map[string]string)
	params := queryResult.Parameters.GetFields()
	if len(params) > 0 {
		for paramName, p := range params {
			fmt.Printf("Param %s: %s (%s)", paramName, p.GetStringValue(), p.String())
			extractedValue := extractDialogFlowEntities(p)
			r.Entities[paramName] = extractedValue
		}
	}
	return
}

func extractDialogFlowEntities(p *structPB.Value) (extractedEntity string) {
	kind := p.GetKind()
	switch kind.(type) {
	case *structPB.Value_StringValue:
		return p.GetStringValue()
	case *structPB.Value_NumberValue:
		return strconv.FormatFloat(p.GetNumberValue(), 'f', 6, 64)
	case *structPB.Value_BoolValue:
		return strconv.FormatBool(p.GetBoolValue())
	case *structPB.Value_StructValue:
		s := p.GetStructValue()
		fields := s.GetFields()
		extractedEntity = ""
		for key, value := range fields {
			if key == "amount" {
				extractedEntity = fmt.Sprintf("%s%s", extractedEntity, strconv.FormatFloat(value.GetNumberValue(), 'f', 6, 64))
			}
			if key == "unit" {
				extractedEntity = fmt.Sprintf("%s%s", extractedEntity, value.GetStringValue())
			}
			if key == "date_time" {
				extractedEntity = fmt.Sprintf("%s%s", extractedEntity, value.GetStringValue())
			}
			// @TODO: Other entity types can be added here
		}
		return extractedEntity
	case *structPB.Value_ListValue:
		list := p.GetListValue()
		if len(list.GetValues()) > 1 {
			// @TODO: Extract more values
		}
		extractedEntity = extractDialogFlowEntities(list.GetValues()[0])
		return extractedEntity
	default:
		return ""
	}
}

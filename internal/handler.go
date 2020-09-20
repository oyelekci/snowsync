package snowsync

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// SQSAPI writes to SQS
type SQSAPI interface {
	SendMessageWithContext(aws.Context, *sqs.SendMessageInput, ...request.Option) (*sqs.SendMessageOutput, error)
}

// ticket can be an incident or change
type ticket interface {
	execute(context.Context) (*Response, error)
}

// Handler respresents the handler type
type Handler struct {
	sqs SQSAPI
}

// Response is returned after capturing the request
type Response struct {
	ResponseType string `json:"response_type"`
	Text         string `json:"text"`
}

// NewHandler returns a new Handler
func NewHandler(s SQSAPI) *Handler {
	return &Handler{sqs: s}
}

// Handle deals with the incoming request
func (h *Handler) Handle(ctx context.Context, request *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	tk, err := h.parseTicket(ctx, request.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}, nil
	}

	rsp, err := tk.execute(ctx)
	if err != nil {
		log.Println(err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}

	body, err := json.Marshal(rsp)
	if err != nil {
		log.Println(err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	headers := map[string]string{"Content-Type": "application/json"}
	return events.APIGatewayProxyResponse{
		Headers:    headers,
		Body:       string(body),
		StatusCode: http.StatusOK,
	}, nil
}

func (h *Handler) parseTicket(ctx context.Context, input string) (ticket, error) {

	ia, err := parseIncident(ctx, input)
	if err == nil {
		ia.sqs = h.sqs
		return ia, err
	}

	// ca, err := parseChange(ctx, input)
	// if err == nil {
	// 	ca.sqs = h.sqs
	// 	return ca, err
	// }

	return nil, fmt.Errorf("failed to parse the ticket")
}

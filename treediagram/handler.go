package treediagram

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/jukeizu/contacts/api/protobuf-spec/contactspb"
	"github.com/jukeizu/contract"
	"github.com/rs/zerolog"
)

type Handler interface {
	SetAddress(contract.Request) (*contract.Response, error)
	SetPhone(contract.Request) (*contract.Response, error)
	Query(contract.Request) (*contract.Response, error)
	RemoveContact(contract.Request) (*contract.Response, error)
	Start() error
	Stop() error
}

type handler struct {
	logger     zerolog.Logger
	client     contactspb.ContactsClient
	httpServer *http.Server
}

func NewHandler(logger zerolog.Logger, client contactspb.ContactsClient, addr string) Handler {
	logger = logger.With().Str("component", "handler").Logger()

	httpServer := http.Server{
		Addr: addr,
	}

	return &handler{logger, client, &httpServer}
}

func (h *handler) SetAddress(request contract.Request) (*contract.Response, error) {
	name, address := parseNameValue("!setaddress", request.Content)

	setAddressRequest := contactspb.SetAddressRequest{
		ServerId: request.ServerId,
		Name:     name,
		Address:  address,
	}

	reply, err := h.client.SetAddress(context.Background(), &setAddressRequest)
	if err != nil {
		return nil, err
	}

	message := contract.Message{
		Content: formatContact(reply.Contact),
	}

	return &contract.Response{Messages: []*contract.Message{&message}}, nil
}

func (h *handler) SetPhone(request contract.Request) (*contract.Response, error) {
	name, phone := parseNameValue("!setphone", request.Content)

	setPhoneRequest := contactspb.SetPhoneRequest{
		ServerId: request.ServerId,
		Name:     name,
		Phone:    phone,
	}

	reply, err := h.client.SetPhone(context.Background(), &setPhoneRequest)
	if err != nil {
		return nil, err
	}

	message := contract.Message{
		Content: formatContact(reply.Contact),
	}

	return &contract.Response{Messages: []*contract.Message{&message}}, nil
}

func (h *handler) Query(request contract.Request) (*contract.Response, error) {
	queryRequest := contactspb.QueryRequest{
		ServerId: request.ServerId,
	}

	reply, err := h.client.Query(context.Background(), &queryRequest)
	if err != nil {
		return nil, err
	}

	contacts := reply.Contacts

	if len(contacts) == 0 {
		message := contract.Message{
			Content: "no results :cry:",
		}

		return &contract.Response{Messages: []*contract.Message{&message}}, nil
	}

	buffer := bytes.Buffer{}
	for _, contact := range contacts {
		buffer.WriteString(formatContact(contact))
	}

	message := contract.Message{
		Content: buffer.String(),
	}

	return &contract.Response{Messages: []*contract.Message{&message}}, nil
}

func (h *handler) RemoveContact(request contract.Request) (*contract.Response, error) {
	input := strings.SplitAfterN(request.Content, "!removecontact ", 2)[1]
	name := strings.Split(input, "'")[1]

	removeContactRequest := contactspb.RemoveContactRequest{
		ServerId: request.ServerId,
		Name:     name,
	}

	reply, err := h.client.RemoveContact(context.Background(), &removeContactRequest)
	if err != nil {
		return nil, err
	}

	removeResponse := ""
	if reply.Removed {
		removeResponse = "removed"
	} else {
		removeResponse = "could not remove"
	}

	message := contract.Message{
		Content: fmt.Sprintf("%s '%s'", removeResponse, name),
	}

	return &contract.Response{Messages: []*contract.Message{&message}}, nil
}

func (h *handler) Start() error {
	h.logger.Info().Msg("starting")

	mux := http.NewServeMux()
	mux.HandleFunc("/setaddress", h.makeLoggingHttpHandlerFunc("setaddress", h.SetAddress))
	mux.HandleFunc("/setphone", h.makeLoggingHttpHandlerFunc("setname", h.SetPhone))
	mux.HandleFunc("/query", h.makeLoggingHttpHandlerFunc("query", h.Query))
	mux.HandleFunc("/removecontact", h.makeLoggingHttpHandlerFunc("removecontact", h.RemoveContact))

	h.httpServer.Handler = mux

	return h.httpServer.ListenAndServe()
}

func (h *handler) Stop() error {
	h.logger.Info().Msg("stopping")

	return h.httpServer.Shutdown(context.Background())
}

func (h *handler) makeLoggingHttpHandlerFunc(name string, f func(contract.Request) (*contract.Response, error)) http.HandlerFunc {
	contractHandlerFunc := contract.MakeHttpHandlerFunc(f)

	return func(w http.ResponseWriter, r *http.Request) {
		defer func(begin time.Time) {
			h.logger.Info().
				Str("intent", name).
				Str("took", time.Since(begin).String()).
				Msg("called")
		}(time.Now())

		contractHandlerFunc.ServeHTTP(w, r)
	}
}

func parseNameValue(command string, content string) (string, string) {
	input := strings.SplitAfterN(content, command, 2)[1]
	split := strings.SplitN(input, "'", 3)
	name, value := split[1], strings.TrimSpace(split[2])

	return name, value
}

func formatContact(contact *contactspb.Contact) string {
	if contact == nil {
		return ""
	}

	buffer := bytes.Buffer{}
	buffer.WriteString(fmt.Sprintf("**%s**\n", contact.Name))
	buffer.WriteString(fmt.Sprintf(":house: %s", contact.Address))
	buffer.WriteString(fmt.Sprintf("\n\n:iphone: %s", contact.Phone))
	buffer.WriteString("\n\n\n")

	return buffer.String()
}

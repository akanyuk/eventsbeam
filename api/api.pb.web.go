// Code generated by protoc-gen-web. DO NOT EDIT.
// source: api.proto

package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/micro/go-micro/v2/metadata"
	"net/http"
	"strconv"
)

type EventsConfiguratorWebHandler interface {
	EventsConfiguratorHandler
}

type eventsConfiguratorWebHandler struct {
	handler EventsConfiguratorWebHandler
}

func RegisterWebHandlers(handler EventsConfiguratorWebHandler, urlsGroup string) {
	h := &eventsConfiguratorWebHandler{
		handler: handler,
	}

	server.AddHandle("/event_levels", h.EventLevels, "Список уровней событий", urlsGroup) // method: GET
}

// swagger:operation GET /event_levels generated EventsConfiguratorEventLevels
// Список уровней событий
//
// Список уровней событий
// ---
// consumes:
// - application/json
// parameters:
// produces:
// - application/json
// responses:
//   '200':
//     description: success
//     schema:
//       "$ref": "#/definitions/EventLevelsResponse"
//   '400':
//     description: bad request
//     schema:
//       "$ref": "#/definitions/ErrorMessage"
func (h *eventsConfiguratorWebHandler) EventLevels(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		response.DetailedError(w, api_error.NewBadRequestError("wrong request method: %s instead %s", r.Method, "GET"))
		return
	}

	userFullName, err := auth.GetUserFullName(r)
	if err != nil {
		response.DetailedError(w, api_error.BadRequestError(err))
		return
	}

	ctx := metadata.Set(r.Context(), "user_full_name", userFullName)
	ctx = metadata.Set(ctx, "lang", server.RequestLanguage(r))
	ctx = metadata.Set(ctx, "is_web_request", "1")

	in := EmptyRequest{}
	out := EventLevelsResponse{}

	if err := h.handler.EventLevels(ctx, &in, &out); err != nil {
		response.DetailedError(w, err)
		return
	}

	result, err := json.Marshal(out)
	if err != nil {
		response.DetailedError(w, api_error.InternalError(err))
		return
	}

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	_, _ = w.Write(result)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ strconv.NumError
var _ mux.Router
var _ json.Decoder

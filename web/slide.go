package web

import (
	"TEST-LOCAL/eventsbeam/web/response"
	"net/http"
)

// swagger:operation POST /setup/slide/create slides handleSlideCreate
// Create slide
//
// Создать новый слайд.
// ---
// parameters:
//   - name: params
//     in: body
//     description: Slide parameters
//     schema:
//       "$ref": "#/definitions/Slide"
// produces:
//   - application/json
// responses:
//   '200':
//     description: success
//     schema:
//       "$ref": "#/definitions/SuccessMessage"
//     examples:
//       application/json: { "success": true, "message": "success", "payload": { } }
//   '400':
//     description: bad request
//     schema:
//       "$ref": "#/definitions/ErrorMessage"
//     examples:
//       application/json: { "success": false, "message": "can not extract slide", "errors": { } }
func (h *handler) handleSlideCreate(w http.ResponseWriter, r *http.Request) {
	slide, err := ExtractSlide(r)
	if err != nil {
		response.WriteErrorResponse(w, http.StatusBadRequest, nil, err.Error())
		return
	}

	validationErrors := h.shower.Slider().Validate(slide)
	if len(validationErrors) > 0 {
		response.WriteErrorResponse(w, http.StatusOK, validationErrors, "validation error")
		return
	}

	if err := h.shower.Slider().Create(slide); err != nil {
		response.WriteErrorResponse(w, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.WriteSuccessResponse(w, nil, "slide created")
}

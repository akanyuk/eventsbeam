package web

import (
	"TEST-LOCAL/eventsbeam/web/response"
	"github.com/gorilla/mux"
	"net/http"
)

// swagger:operation GET /setup/slides/{compo} slide handleSlides
// Read slides
//
// Чтение списка всех слайдов для заданного компо.
// ---
// parameters:
//   - name: compo
//     in: path
//     type: string
//     required: true
//     description: Идентификатор компо
// produces:
//   - application/json
// responses:
//   '200':
//     description: success
//     schema:
//       type: array
//       description: Массив всех компо
//       items:
//         "$ref": "#/definitions/Slide"
//     examples:
//       application/json: [{"id":1,"position":1,"compo":"zxdemo","template":"test1","params":{}},{"id":2,"position":2,"compo":"zxdemo","template":"test1","params":{}},{"id":3,"position":1,"compo":"zxintro","template":"test2","params":{}}]
func (h *handler) handleSlides(w http.ResponseWriter, r *http.Request) {
	compo := mux.Vars(r)["compo"]

	response.WriteDataResponse(w, h.shower.Slider().Slides(compo))
}

// swagger:operation POST /setup/slide/create slide handleSlideCreate
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

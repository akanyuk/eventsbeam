package web

import (
	"TEST-LOCAL/eventsbeam/web/response"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
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
//       application/json: [{"id":1,"compo":"zxdemo","template":"test1","params":{}},{"id":2,"compo":"zxdemo","template":"test1","params":{}},{"id":3,"compo":"zxintro","template":"test2","params":{}}]
func (h *handler) handleSlides(w http.ResponseWriter, r *http.Request) {
	compo := mux.Vars(r)["compo"]

	response.WriteDataResponse(w, h.slider.Slides(compo))
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

	validationErrors := h.slider.Validate(slide)
	if len(validationErrors) > 0 {
		response.WriteErrorResponse(w, http.StatusOK, validationErrors, "validation error")
		return
	}

	if err := h.slider.Create(slide); err != nil {
		response.WriteErrorResponse(w, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.WriteSuccessResponse(w, nil, "slide created")
}

// swagger:operation GET /setup/slide/read/{id} slide handleSlideRead
// Read slide
//
// Чтение слайда.
// ---
// parameters:
//   - name: id
//     in: path
//     type: integer
//     required: true
//     description: Идентификатор слайда
// produces:
//   - application/json
// responses:
//   '200':
//     description: success
//     schema:
//       allOf:
//       - "$ref": "#/definitions/SuccessMessage"
//       - type: object
//         properties:
//           payload:
//             description: Параметры слайда
//             "$ref": "#/definitions/Slide"
//     examples:
//       application/json: {"success":true,"message":"success","errors":null,"payload":{"id":1,"compo":"zxdemo","template":"test2","params":{}}}
//   '400':
//     description: bad request
//     schema:
//       "$ref": "#/definitions/ErrorMessage"
//     examples:
//       application/json: { "success": false, "message": "wrong slide ID", "errors": { } }
//   '404':
//     description: not found
//     schema:
//       "$ref": "#/definitions/ErrorMessage"
//     examples:
//       application/json: { "success": false, "message": "slide not found", "errors": { } }
func (h *handler) handleSlideRead(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		response.WriteErrorResponse(w, http.StatusBadRequest, nil, "wrong slide ID: %v", err.Error())
		return
	}

	compo, err := h.slider.Read(id)
	if err != nil {
		response.WriteErrorResponse(w, http.StatusNotFound, nil, err.Error())
		return
	}

	response.WriteSuccessResponse(w, compo, "success")
}

// swagger:operation POST /setup/slide/update/{id} slide handleSlideUpdate
// Update slide
//
// Редактирование слайда.
// ---
// parameters:
//   - name: id
//     in: path
//     type: integer
//     required: true
//     description: Идентификатор слайда
//   - name: params
//     in: body
//     description: Свойства слайда
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
//       application/json: {"success":false,"message":"validation error","errors":[{"code":"alias","message":"alias already exists"},{"code":"title","message":"title can not be empty"}]}
//   '400':
//     description: bad request
//     schema:
//       "$ref": "#/definitions/ErrorMessage"
//     examples:
//       application/json: { "success": false, "message": "can not extract slide", "errors": { } }
//   '404':
//     description: not found
//     schema:
//       "$ref": "#/definitions/ErrorMessage"
//     examples:
//       application/json: { "success": false, "message": "slide not found", "errors": { } }
//   '500':
//     description: internal error
//     schema:
//       "$ref": "#/definitions/ErrorMessage"
//     examples:
//       application/json: { "success": false, "message": "unable to save slide", "errors": { } }
func (h *handler) handleSlideUpdate(w http.ResponseWriter, r *http.Request) {
	slide, err := ExtractSlide(r)
	if err != nil {
		response.WriteErrorResponse(w, http.StatusBadRequest, nil, "wrong request params: %v", err.Error())
		return
	}

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		response.WriteErrorResponse(w, http.StatusBadRequest, nil, "wrong slide ID: %v", err.Error())
		return
	}

	validationErrors := h.slider.Validate(slide)
	if len(validationErrors) > 0 {
		response.WriteErrorResponse(w, http.StatusOK, validationErrors, "validation error")
		return
	}

	if err := h.slider.Update(id, slide); err != nil {
		response.WriteErrorResponse(w, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.WriteSuccessResponse(w, nil, "slide updated")
}

// swagger:operation POST /setup/slide/delete/{id} slide handleSlideDelete
// Delete slide
//
// Удаление слайда.
// ---
// parameters:
// - name: id
//   in: path
//   type: integer
//   required: true
//   description: Идентификатор слайда
// requestBody:
//   required: true
//   content:
//     application/x-www-form-urlencoded:
//       type: object
//       properties:
//         name: confirm
//         type: boolean
//         description: Подтверждение удаления (защита от прямых GET-запросов)
//         required: true
//         enum: [true]
// produces:
// - application/json
// responses:
//   '200':
//     description: success
//     schema:
//       "$ref": "#/definitions/SuccessMessage"
//     examples:
//       application/json: { "success": true, "message": "slide deleted", "payload": { } }
//   '400':
//     description: bad request
//     schema:
//       "$ref": "#/definitions/ErrorMessage"
//     examples:
//       application/json: { "success": false, "message": "request not confirmed", "errors": { } }
//   '404':
//     description: not found
//     schema:
//       "$ref": "#/definitions/ErrorMessage"
//     examples:
//       application/json: { "success": false, "message": "slide not found", "errors": { } }
//   '500':
//     description: internal error
//     schema:
//       "$ref": "#/definitions/ErrorMessage"
//     examples:
//       application/json: { "success": false, "message": "unable to save slides", "errors": { } }
func (h *handler) handleSlideDelete(w http.ResponseWriter, r *http.Request) {
	isConfirmed, err := strconv.ParseBool(r.FormValue("confirm"))
	if err != nil || !isConfirmed {
		response.WriteErrorResponse(w, http.StatusBadRequest, nil, "request not confirmed")
		return
	}

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		response.WriteErrorResponse(w, http.StatusBadRequest, nil, "wrong slide ID: %v", err.Error())
		return
	}

	_, err = h.slider.Read(id)
	if err != nil {
		response.WriteErrorResponse(w, http.StatusNotFound, nil, err.Error())
		return
	}

	if err := h.slider.Delete(id); err != nil {
		response.WriteErrorResponse(w, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.WriteSuccessResponse(w, nil, "slide deleted")
}

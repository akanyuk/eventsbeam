package web

import (
	"TEST-LOCAL/events_beam/web/response"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// swagger:operation GET /setup/compos compo handleCompos
// Read all compos
//
// Чтение списка всех компо.
// ---
// produces:
//   - application/json
// responses:
//   '200':
//     description: success
//     schema:
//       type: array
//       description: Массив всех компо
//       items:
//         "$ref": "#/definitions/Compo"
//     examples:
//       application/json: [ { "alias": "zxdemo", "title": "ZX Spectrum demo", "slides": [] }, { "alias": "intro", "title": "ZX Spectrum 256 bytes intro", "slides": [] } ]
func (h *handler) handleCompos(w http.ResponseWriter, r *http.Request) {
	response.WriteDataResponse(w, h.shower.Comper().Compos())
}

// swagger:operation POST /setup/compo/create compo handleCompoCreate
// Create compo
//
// Создать новое компо.
// ---
// parameters:
//   - name: params
//     in: body
//     description: Свойства компо
//     schema:
//       "$ref": "#/definitions/Compo"
// produces:
//   - application/json
// responses:
//   '200':
//     description: success
//     schema:
//       "$ref": "#/definitions/SuccessOrError"
//     examples:
//       application/json: {"success":false,"message":"validation error","errors":[{"code":"alias","message":"alias already exists"},{"code":"title","message":"title can not be empty"}]}
//   '400':
//     description: bad request
//     schema:
//       "$ref": "#/definitions/ErrorMessage"
//     examples:
//       application/json: { "success": false, "message": "can not extract compo", "errors": { } }
//   '500':
//     description: internal error
//     schema:
//       "$ref": "#/definitions/ErrorMessage"
//     examples:
//       application/json: { "success": false, "message": "unable to save compo", "errors": { } }
func (h *handler) handleCompoCreate(w http.ResponseWriter, r *http.Request) {
	compo, err := ExtractCompo(r)
	if err != nil {
		response.WriteErrorResponse(w, http.StatusBadRequest, nil, err.Error())
		return
	}

	validationErrors := h.shower.Comper().Validate(compo)
	if len(validationErrors) > 0 {
		response.WriteErrorResponse(w, http.StatusOK, validationErrors, "validation error")
		return
	}

	if err := h.shower.Comper().Create(compo); err != nil {
		response.WriteErrorResponse(w, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.WriteSuccessResponse(w, nil, "compo created")
}

// swagger:operation GET /setup/compo/read/{alias} compo handleCompoRead
// Read compo
//
// Чтение компо.
// ---
// parameters:
//   - name: alias
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
//       AllOf:
//       - "$ref": "#/definitions/SuccessMessage"
//       - type: object
//         properties:
//           payload:
//             description: Параметры компо
//             "$ref": "#/definitions/Compo"
//     examples:
//       application/json: { "success": true, "message": "success", "payload":  { "alias": "zxdemo", "title": "ZX Spectrum demo", "slides": null } }
//   '404':
//     description: not found
//     schema:
//       "$ref": "#/definitions/ErrorMessage"
//     examples:
//       application/json: { "success": false, "message": "compo not found", "errors": { } }
func (h *handler) handleCompoRead(w http.ResponseWriter, r *http.Request) {
	alias := mux.Vars(r)["alias"]

	compo, err := h.shower.Comper().Read(alias)
	if err != nil {
		response.WriteErrorResponse(w, http.StatusNotFound, nil, err.Error())
		return
	}

	response.WriteSuccessResponse(w, compo, "success")
}

// swagger:operation POST /setup/compo/update/{alias} compo handleCompoUpdate
// Update compo
//
// Редактирование компо.
// ---
// parameters:
//   - name: alias
//     in: path
//     type: string
//     required: true
//     description: Идентификатор компо
//   - name: params
//     in: body
//     description: Свойства компо
//     schema:
//       "$ref": "#/definitions/Compo"
// produces:
//   - application/json
// responses:
//   '200':
//     description: success
//     schema:
//       "$ref": "#/definitions/SuccessOrError"
//     examples:
//       application/json: {"success":false,"message":"validation error","errors":[{"code":"alias","message":"alias already exists"},{"code":"title","message":"title can not be empty"}]}
//   '400':
//     description: bad request
//     schema:
//       "$ref": "#/definitions/ErrorMessage"
//     examples:
//       application/json: { "success": false, "message": "can not extract compo", "errors": { } }
//   '404':
//     description: not found
//     schema:
//       "$ref": "#/definitions/ErrorMessage"
//     examples:
//       application/json: { "success": false, "message": "compo not found", "errors": { } }
//   '500':
//     description: internal error
//     schema:
//       "$ref": "#/definitions/ErrorMessage"
//     examples:
//       application/json: { "success": false, "message": "unable to save compo", "errors": { } }
func (h *handler) handleCompoUpdate(w http.ResponseWriter, r *http.Request) {
	compo, err := ExtractCompo(r)
	if err != nil {
		response.WriteErrorResponse(w, http.StatusBadRequest, nil, err.Error())
		return
	}

	alias := mux.Vars(r)["alias"]
	_, err = h.shower.Comper().Read(alias)
	if err != nil {
		response.WriteErrorResponse(w, http.StatusNotFound, nil, err.Error())
		return
	}

	validationErrors := h.shower.Comper().Validate(compo)
	if len(validationErrors) > 0 {
		response.WriteErrorResponse(w, http.StatusOK, validationErrors, "validation error")
		return
	}

	if err := h.shower.Comper().Update(alias, compo); err != nil {
		response.WriteErrorResponse(w, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.WriteSuccessResponse(w, nil, "compo updated")
}

// swagger:operation POST /setup/compo/delete/{alias} compo handleCompoDelete
// Delete compo
//
// Удаление компо.
// ---
// parameters:
// - name: alias
//   in: path
//   type: string
//   required: true
//   description: Идентификатор компо
// - in: formData
//   name: confirm
//   type: boolean
//   description: Подтверждение удаления (защита от прямых GET-запросов)
//   required: true
//   enum: [true]
// produces:
//   - application/json
// responses:
//   '200':
//     description: success
//     schema:
//       "$ref": "#/definitions/SuccessMessage"
//     examples:
//       application/json: { "success": true, "message": "compo deleted", "payload": { } }
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
//       application/json: { "success": false, "message": "compo not found", "errors": { } }
//   '500':
//     description: internal error
//     schema:
//       "$ref": "#/definitions/ErrorMessage"
//     examples:
//       application/json: { "success": false, "message": "unable to save compo", "errors": { } }
func (h *handler) handleCompoDelete(w http.ResponseWriter, r *http.Request) {
	isConfirmed, err := strconv.ParseBool(r.FormValue("confirm"))
	if err != nil || !isConfirmed {
		response.WriteErrorResponse(w, http.StatusBadRequest, nil, "request not confirmed")
		return
	}

	alias := mux.Vars(r)["alias"]
	_, err = h.shower.Comper().Read(alias)
	if err != nil {
		response.WriteErrorResponse(w, http.StatusNotFound, nil, err.Error())
		return
	}

	if err := h.shower.Comper().Delete(alias); err != nil {
		response.WriteErrorResponse(w, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.WriteSuccessResponse(w, nil, "compo deleted")
}

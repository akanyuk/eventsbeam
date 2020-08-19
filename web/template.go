package web

import (
	"bitbucket.org/nyuk/eventsbeam/web/response"
	"net/http"
)

// swagger:operation GET /setup/templates template handleTemplates
// Read templates
//
// Чтение списка всех шаблонов, доступных в каталоге templates.
// ---
// produces:
//   - application/json
// responses:
//   '200':
//     description: success
//     schema:
//       type: array
//       description: Массив всех шаблонов
//       items:
//         "$ref": "#/definitions/Template"
//     examples:
//       application/json: [ { "name": "picture" }, { "name": "ascii" }, { "name": "prizegiving" } ]
func (h *handler) handleTemplates(w http.ResponseWriter, r *http.Request) {
	response.WriteDataResponse(w, h.beamer.Templates())
}

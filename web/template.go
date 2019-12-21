package web

import (
	"TEST-LOCAL/eventsbeam/web/response"
	"net/http"
)

// swagger:operation GET /setup/templates template handleTemplates
// Read all templates
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
//       description: Массив всех компо
//       items:
//         "$ref": "#/definitions/Compo"
//     examples:
//       application/json: [ { "alias": "zxdemo", "title": "ZX Spectrum demo", "slides": [] }, { "alias": "intro", "title": "ZX Spectrum 256 bytes intro", "slides": [] } ]
func (h *handler) handleTemplates(w http.ResponseWriter, r *http.Request) {
	response.WriteDataResponse(w, h.shower.Comper().Compos())
}

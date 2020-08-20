package web

import (
	"github.com/akanyuk/eventsbeam/cmd/protoc-gen-web/generator"
	webProto "github.com/akanyuk/eventsbeam/cmd/protoc-gen-web/proto/web"
	pb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"google.golang.org/protobuf/proto"
	"strings"
)

var imports = []string{
	"github.com/micro/go-micro/v2/metadata",
	"github.com/gorilla/mux",
	"encoding/json",
	"net/http",
	"strconv",
}

const pluginName = "web"

var (
	pkgImports map[string]bool
)

func init() {
	generator.RegisterPlugin(new(web))
}

func unExport(s string) string {
	if len(s) == 0 {
		return ""
	}
	name := strings.ToLower(s[:1]) + s[1:]
	if pkgImports[name] {
		return name + "_"
	}
	return name
}

// web is an implementation of the Go protocol buffer compiler's
// plugin architecture.  It generates bindings for go-web support.
type web struct {
	gen              *generator.Generator
	withFileDownload bool
}

// Name returns the name of this plugin, "web".
func (g *web) Name() string {
	return pluginName
}

// Init initializes the plugin.
func (g *web) Init(gen *generator.Generator) {
	g.gen = gen
}

// Given a type name defined in a .proto, return its object.
// Also record that we're using it, to guarantee the associated import.
func (g *web) objectNamed(name string) generator.Object {
	g.gen.RecordTypeUse(name)
	return g.gen.ObjectNamed(name)
}

// Given a type name defined in a .proto, return its name as we will print it.
func (g *web) typeName(str string) string {
	return g.gen.TypeName(g.objectNamed(str))
}

// P forwards to g.gen.P.
func (g *web) P(args ...interface{}) { g.gen.P(args...) }

// Generate generates code for the services in the given file.
func (g *web) Generate(file *generator.FileDescriptor) {
	if len(file.FileDescriptorProto.Service) == 0 {
		return
	}

	for i, service := range file.FileDescriptorProto.Service {
		g.generateService(file, service, i)
	}

	g.generateInternal()
}

// GenerateImports generates the import declaration for this file.
func (g *web) GenerateImports(file *generator.FileDescriptor) {
	if len(file.FileDescriptorProto.Service) == 0 {
		return
	}
	g.P("import (")
	for _, i := range imports {
		g.P("\t\"" + i + "\"")
	}
	g.P(")")
	g.P()

	// We need to keep track of imported packages to make sure we don't produce
	// a name collision when generating types.
	pkgImports = make(map[string]bool)
	for _, name := range imports {
		pkgImports[name] = true
	}
}

// generateService generates all the code for the named service.
func (g *web) generateService(_ *generator.FileDescriptor, service *pb.ServiceDescriptorProto, _ int) {
	origServiceName := service.GetName()
	serviceName := generator.CamelCase(origServiceName)

	g.generateRegisterFunction(serviceName, service.Method)

	for _, method := range service.Method {
		g.generateServerMethodSwagger(serviceName, method)
		g.generateServerMethod(serviceName, method)
	}
}

func (g *web) generateRegisterFunction(serviceName string, methods []*pb.MethodDescriptorProto) {
	g.P("type " + serviceName + "WebHandler interface {")
	g.P(serviceName + "Handler")
	g.P("}")
	g.P()
	g.P("type ", unExport(serviceName), "WebHandler struct {")
	g.P("handler " + serviceName + "WebHandler")
	g.P("}")
	g.P()
	g.P("func RegisterWebHandlers(handler " + serviceName + "WebHandler, urlsGroup string) {")
	g.P("h := &" + unExport(serviceName) + "WebHandler{")
	g.P("handler: handler,")
	g.P("}")
	g.P()

	for _, method := range methods {
		methodName := generator.CamelCase(method.GetName())

		r := proto.GetExtension(method.Options, webProto.E_Handler)
		params := r.(*webProto.Handler)

		requestMethod, path := g.handlerMethodPath(params)
		g.P("server.AddHandle(\"" + path + "\", h." + methodName + ", \"" + params.GetDescription() + "\", urlsGroup) // method: " + requestMethod)
	}

	g.P("}")
	g.P()
}

func (g *web) handlerMethodPath(params *webProto.Handler) (string, string) {
	switch {
	case len(params.GetGet()) > 0:
		return "GET", params.GetGet()
	case len(params.GetPost()) > 0:
		return "POST", params.GetPost()
	case len(params.GetDelete()) > 0:
		return "DELETE", params.GetDelete()
	case len(params.GetPatch()) > 0:
		return "PATCH", params.GetPatch()
	case len(params.GetPut()) > 0:
		return "PUT", params.GetPut()
	}

	g.gen.Fail("wrong request method")
	return "", ""
}

func (g *web) generateServerMethod(serviceName string, method *pb.MethodDescriptorProto) {
	methodName := generator.CamelCase(method.GetName())

	inType := g.typeName(method.GetInputType())
	outType := g.typeName(method.GetOutputType())

	r := proto.GetExtension(method.Options, webProto.E_Handler)
	params := r.(*webProto.Handler)

	requestMethod, _ := g.handlerMethodPath(params)

	g.P("func (h *", unExport(serviceName), "WebHandler) ", methodName, "(w http.ResponseWriter, r *http.Request) {")
	g.P("if r.Method != \"" + requestMethod + "\" {")
	g.P("response.DetailedError(w, api_error.NewBadRequestError(\"wrong request method: %s instead %s\", r.Method, \"" + requestMethod + "\"))")
	g.P("return")
	g.P("}")
	g.P("")
	g.P("userFullName, err := auth.GetUserFullName(r)")
	g.P("if err != nil {")
	g.P("response.DetailedError(w, api_error.BadRequestError(err))")
	g.P("return")
	g.P("}")
	g.P()
	g.P("ctx := metadata.Set(r.Context(), \"user_full_name\", userFullName)")
	g.P("ctx = metadata.Set(ctx, \"lang\", server.RequestLanguage(r))")
	g.P("ctx = metadata.Set(ctx, \"is_web_request\", \"1\")")
	g.P()
	g.P("in := " + inType + "{}")
	g.P("out := " + outType + "{}")
	g.P()

	// `body` always first
	for _, p := range params.GetParameters() {
		if g.parameterIn(p) == "body" {
			g.requestParameterBody()
			break
		}
	}

	for _, p := range params.GetParameters() {
		if g.parameterIn(p) == "path" {
			g.requestParameterPath(method, p)
		}
	}

	g.P("if err := h.handler." + methodName + "(ctx, &in, &out); err != nil {")
	g.P("response.DetailedError(w, err)")
	g.P("return")
	g.P("}")
	g.P()

	g.P("result, err := json.Marshal(out)")
	g.P("if err != nil {")
	g.P("response.DetailedError(w, api_error.InternalError(err))")
	g.P("return")
	g.P("}")
	g.P()
	g.P("w.Header().Set(\"Content-Type\", \"application/json;charset=utf-8\")")
	g.P("_, _ = w.Write(result)")

	g.P("}")
	g.P()
}

func (g *web) requestParameterBody() {
	g.P("decoder := json.NewDecoder(r.Body)")
	g.P("defer func() {")
	g.P("_ = r.Body.Close()")
	g.P("}()")
	g.P()
	g.P("if err := decoder.Decode(&in); err != nil {")
	g.P("response.DetailedError(w, api_error.BadRequestError(err))")
	g.P("return")
	g.P("}")
	g.P()
}

func (g *web) requestParameterPath(method *pb.MethodDescriptorProto, parameter *webProto.Parameter) {
	fieldName := generator.CamelCase(parameter.GetName())
	goType := g.goType(g.gen.ObjectNamed(method.GetInputType()), parameter.GetName())

	switch goType {
	case "":
		g.gen.Fail("not found go type for " + method.GetName() + "." + fieldName)
	case "int32", "int64":
		g.requestParameterPathInt(parameter.GetName(), fieldName, goType)
	default:
		g.gen.Fail("parameter extractor not implemented. type: " + goType + ", in: path")
	}
}

func (g *web) requestParameterPathInt(name string, fieldName string, goType string) {
	g.P(name + ", err := strconv.Atoi(mux.Vars(r)[\"" + name + "\"])")
	g.P("if err != nil {")
	g.P("response.DetailedError(w, err)")
	g.P("return")
	g.P("}")
	g.P("in." + fieldName + " = " + goType + "(" + name + ")")
	g.P()
}

func (g *web) parameterIn(parameter *webProto.Parameter) string {
	if parameter.GetBody() != "" {
		return "body"
	}

	if parameter.GetPath() != "" {
		return "path"
	}

	g.gen.Fail("parameter miss in:description fields")
	return ""
}

func (g *web) parameterDescription(parameter *webProto.Parameter) string {
	if parameter.GetBody() != "" {
		return parameter.GetBody()
	}

	if parameter.GetPath() != "" {
		return parameter.GetPath()
	}

	g.gen.Fail("parameter miss in:description fields")
	return ""
}

func (g *web) goType(obj generator.Object, fieldName string) string {
	d, ok := obj.(*generator.Descriptor)
	if !ok {
		g.gen.Fail("unable to convert object to descriptor")
	}

	for _, field := range d.GetField() {
		if field.GetName() == fieldName {
			goType, _ := g.gen.GoType(d, field)
			return goType
		}
	}

	g.gen.Fail("field name not found in descriptor")
	return ""
}

func (g *web) generateInternal() {
	g.P("// Reference imports to suppress errors if they are not otherwise used.")
	g.P("var _ strconv.NumError")
	g.P("var _ mux.Router")
	g.P("var _ json.Decoder")
	g.P()
}

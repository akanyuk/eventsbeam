package web

import (
	"github.com/akanyuk/eventsbeam/cmd/protoc-gen-web/generator"
	webProto "github.com/akanyuk/eventsbeam/cmd/protoc-gen-web/proto/web"
	pb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"google.golang.org/protobuf/proto"
)

func (g *web) generateServerMethodSwagger(serviceName string, method *pb.MethodDescriptorProto) {
	methodName := generator.CamelCase(method.GetName())

	r := proto.GetExtension(method.Options, webProto.E_Handler)
	params := r.(*webProto.Handler)

	requestMethod, path := g.handlerMethodPath(params)

	g.P("// swagger:operation " + requestMethod + " " + path + " generated " + serviceName + methodName)
	g.P("// " + params.GetDescription())
	g.P("//")
	g.P("// " + params.GetDescription())
	g.P("// ---")
	g.P("// consumes:")
	g.P("// - application/json")
	g.P("// parameters:")
	for _, p := range params.GetParameters() {
		if g.parameterIn(p) == "body" {
			g.P("// - in: body")
			g.P("//   name: requestBody")
			g.P("//   required: true")
			g.P("//   description: " + g.parameterDescription(p))
			g.P("//   schema:")
			g.P("//     \"$ref\": \"#/definitions/" + g.typeName(method.GetInputType()) + "\"")
			continue
		}

		goType := g.goType(g.gen.ObjectNamed(method.GetInputType()), p.GetName())

		g.P("// - in: " + g.parameterIn(p))
		g.P("//   name: " + p.GetName())
		g.P("//   type: " + swaggerType(goType))
		g.P("//   description: " + g.parameterDescription(p))

		if g.parameterIn(p) == "path" {
			// `path` parameters always required
			g.P("//   required: true")
		}
	}

	g.P("// produces:")
	g.P("// - application/json")
	g.P("// responses:")
	g.P("//   '200':")
	g.P("//     description: success")
	g.P("//     schema:")
	g.P("//       \"$ref\": \"#/definitions/" + g.typeName(method.GetOutputType()) + "\"")
	g.P("//   '400':")
	g.P("//     description: bad request")
	g.P("//     schema:")
	g.P("//       \"$ref\": \"#/definitions/ErrorMessage\"")
}

func (g *web) generateSwaggerDownloadResponse() {
	g.P("// produces:")
	g.P("// - application/json")
	g.P("// - application/octet-stream")
	g.P("// responses:")
	g.P("//   '200':")
	g.P("//     description: success")
	g.P("//     schema:")
	g.P("//       type: file")
	g.P("//   '400':")
	g.P("//     description: bad request")
	g.P("//     schema:")
	g.P("//       \"$ref\": \"#/definitions/ErrorMessage\"")
}

func swaggerType(goType string) string {
	switch goType {
	case "int", "int32", "int64":
		return "integer"
	case "float32", "float64":
		return "number"
	default:
		return goType
	}

}

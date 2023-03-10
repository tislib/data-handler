package rest

import (
	"context"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/http-swagger"
	"github.com/tislib/data-handler/pkg/abs"
	"github.com/tislib/data-handler/pkg/errors"
	"github.com/tislib/data-handler/pkg/model"
	"github.com/tislib/data-handler/pkg/service/annotations"
	"github.com/tislib/data-handler/pkg/types"
	"github.com/tislib/data-handler/pkg/util"
	"io"
	"net/http"
)

type SwaggerApi interface {
	ConfigureRouter(r *mux.Router)
}

type swaggerApi struct {
	resourceService abs.ResourceService
}

func (s *swaggerApi) ConfigureRouter(r *mux.Router) {
	swaggerFiles.Handler.Prefix = "/docs/"

	file, err := statikFS.Open("/openapi.yaml")

	if err != nil {
		log.Fatal(err)
	}

	openApiData, err := io.ReadAll(file)

	r.HandleFunc("/docs/api.json", func(w http.ResponseWriter, req *http.Request) {
		doc, serviceErr := s.prepareDoc(req.Context(), openApiData)

		if serviceErr != nil {
			http.Error(w, serviceErr.GetFullMessage(), 500)
			return
		}

		data, err := doc.MarshalJSON()

		if err != nil {
			http.Error(w, serviceErr.GetFullMessage(), 400)
			return
		}

		_, err = w.Write(data)

		if err != nil {
			log.Error(err)
		}
	})

	if err != nil {
		log.Fatal(err)
	}

	r.HandleFunc("/docs/resources/{namespace}/{resourceName}.json", func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		namespace := vars["namespace"]
		resourceName := vars["resourceName"]

		resource := s.resourceService.GetResourceByName(req.Context(), namespace, resourceName)

		if resource == nil {
			http.Error(w, errors.ResourceNotFoundError.GetFullMessage(), 404)
			return
		}

		doc := s.prepareResourceSchema(resource)

		data, err := doc.MarshalJSON()

		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		_, err = w.Write(data)

		if err != nil {
			log.Error(err)
		}

	})

	r.HandleFunc("/docs", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Location", "/docs/index.html")
		writer.WriteHeader(301)
	})

	r.PathPrefix("/docs").HandlerFunc(httpSwagger.Handler(
		httpSwagger.URL("/docs/api.json"), //The url pointing to API definition
	))
}

func (s *swaggerApi) prepareDoc(ctx context.Context, openApiData []byte) (*openapi3.T, errors.ServiceError) {
	loader := openapi3.NewLoader()

	doc, err := loader.LoadFromData(openApiData)

	if err != nil {
		panic(err)
	}

	list := s.resourceService.List(ctx)

	for _, item := range list {
		if item.Namespace != "system" {
			s.appendResourceApis(ctx, doc, item)
		}
	}

	// post processing
	var security = &openapi3.SecurityRequirements{
		{
			"bearerAuth": []string{},
		},
	}

	for pathKey, path := range doc.Paths {
		for operationKey, operation := range path.Operations() {
			if pathKey == "/authentication/token" && operationKey == "POST" {
				continue
			}

			operation.Security = security
		}
	}

	return doc, nil
}

func (s *swaggerApi) appendResourceApis(ctx context.Context, doc *openapi3.T, resource *model.Resource) {
	jsonSchemaRef := "/docs/resources/" + resource.Namespace + "/" + resource.Name + ".json"

	var tags []string
	if resource.GetNamespace() == "default" {
		tags = []string{resource.GetName()}
	} else {
		tags = []string{resource.GetNamespace() + " / " + resource.GetName()}
	}

	title := resource.GetTitle()
	description := resource.GetDescription()

	if title == "" {
		title = resource.GetName()
	}

	if description == "" {
		description = "Api for " + resource.GetName()
	}

	doc.Paths["/"+s.getResourceFQN(resource)] = &openapi3.PathItem{
		Summary:     title,
		Description: description,
		Get: &openapi3.Operation{
			Tags:        tags,
			Summary:     title,
			Description: description + " - get list",
			Responses: map[string]*openapi3.ResponseRef{
				"200": {
					Value: &openapi3.Response{
						Content: openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{
							Ref: "#/components/schemas/list-" + s.getResourceFQN(resource),
						}),
					},
				},
			},
		},
		Post: &openapi3.Operation{
			Tags:        tags,
			Summary:     title,
			Description: description + " - create item",
			RequestBody: &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					Required: true,
					Content: openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{
						Ref: jsonSchemaRef,
					}),
				},
			},
			Responses: map[string]*openapi3.ResponseRef{
				"200": {
					Value: &openapi3.Response{
						Content: openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{
							Ref: jsonSchemaRef,
						}),
					},
				},
			},
		},
	}

	doc.Paths["/"+s.getResourceFQN(resource)+"/{id}"] = &openapi3.PathItem{
		Summary:     title,
		Description: description,
		Get: &openapi3.Operation{
			Tags:        tags,
			Summary:     title,
			Description: description + " - get item",
			Responses: map[string]*openapi3.ResponseRef{
				"200": {
					Value: &openapi3.Response{
						Content: openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{
							Ref: "#/components/schemas/item-" + s.getResourceFQN(resource),
						}),
					},
				},
			},
		},
		Delete: &openapi3.Operation{
			Tags:        tags,
			Summary:     title,
			Description: description + " - delete item",
			Responses: map[string]*openapi3.ResponseRef{
				"200": {
					Value: &openapi3.Response{
						Content: openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{
							Ref: "#/components/schemas/item-" + s.getResourceFQN(resource),
						}),
					},
				},
			},
		},
		Put: &openapi3.Operation{
			Tags:        tags,
			Summary:     title,
			Description: description + " - update item",
			RequestBody: &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					Required: true,
					Content: openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{
						Ref: jsonSchemaRef,
					}),
				},
			},
			Responses: map[string]*openapi3.ResponseRef{
				"200": {
					Value: &openapi3.Response{
						Content: openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{
							Ref: jsonSchemaRef,
						}),
					},
				},
			},
		},
	}

	doc.Components.Schemas["list-"+s.getResourceFQN(resource)] = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Properties: map[string]*openapi3.SchemaRef{
				"content": {
					Value: &openapi3.Schema{
						Type: openapi3.TypeArray,
						Items: &openapi3.SchemaRef{
							Ref: jsonSchemaRef,
						},
					},
				},
			},
		},
	}

	doc.Components.Schemas["item-"+s.getResourceFQN(resource)] = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Properties: map[string]*openapi3.SchemaRef{
				"content": {
					Ref: jsonSchemaRef,
				},
			},
		},
	}

	for _, tag := range tags {
		doc.Tags = append(doc.Tags, &openapi3.Tag{
			Name:        tag,
			Description: "",
		})
	}
}

func (s *swaggerApi) getResourceFQN(resource *model.Resource) string {
	if resource.Namespace == "default" {
		return util.ToDashCase(resource.Name)
	} else {
		return util.ToDashCase(resource.Namespace) + "-" + util.ToDashCase(resource.Name)
	}
}

func (s *swaggerApi) prepareResourceSchema(resource *model.Resource) *openapi3.Schema {
	var requiredItems []string

	propertiesSchema := &openapi3.Schema{
		Properties: map[string]*openapi3.SchemaRef{},
	}

	for _, property := range resource.Properties {
		propSchema := &openapi3.Schema{
			Type: types.ResourcePropertyTypeToJsonSchemaType(property.Type),
		}

		if property.ExampleValue != nil {
			propSchema.Example = property.ExampleValue.AsInterface()
		}

		if property.DefaultValue != nil {
			propSchema.Default = property.DefaultValue.AsInterface()
		}

		propertiesSchema.Properties[property.Name] = &openapi3.SchemaRef{
			Value: propSchema,
		}

		if property.Required {
			requiredItems = append(requiredItems, property.Name)
		}
	}

	propertiesSchema.Required = requiredItems
	propertiesSchema.Title = resource.GetTitle()
	propertiesSchema.Description = resource.GetDescription()

	recordSchema := &openapi3.Schema{
		Properties: map[string]*openapi3.SchemaRef{},
	}

	recordSchema.Properties["properties"] = &openapi3.SchemaRef{
		Value: propertiesSchema,
	}

	recordSchema.Properties["id"] = &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "string", Format: "uuid"}}

	if !annotations.IsEnabled(resource, annotations.DisableVersion) {
		recordSchema.Properties["version"] = &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "number", ReadOnly: true}}
	}

	if !annotations.IsEnabled(resource, annotations.DisableAudit) {
		recordSchema.Properties["auditData"] = &openapi3.SchemaRef{
			Value: &openapi3.Schema{
				Properties: map[string]*openapi3.SchemaRef{
					"createdOn": {Value: &openapi3.Schema{Type: "string", Format: "date"}},
					"updatedOn": {Value: &openapi3.Schema{Type: "string", Format: "date"}},
					"createdBy": {Value: &openapi3.Schema{Type: "string"}},
					"updatedBy": {Value: &openapi3.Schema{Type: "string"}},
				},
				ReadOnly: true,
			},
		}
	}

	return recordSchema
}

func NewSwaggerApi(resourceService abs.ResourceService) SwaggerApi {
	return &swaggerApi{
		resourceService: resourceService,
	}
}

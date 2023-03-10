package generator

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tislib/data-handler/pkg/model"
	"github.com/tislib/data-handler/pkg/service/annotations"
	"github.com/tislib/data-handler/pkg/types"
	"go/format"
	"reflect"
	"strings"
)

type GenerateResourceCodeParams struct {
	Package   string
	Resources []*model.Resource
}

func GenerateGoResourceCode(resource *model.Resource, params GenerateResourceCodeParams) string {
	var sb strings.Builder

	// scan properties for needed packages
	uuidNeeded := false

	for _, prop := range resource.Properties {
		if prop.Type == model.ResourceProperty_UUID {
			uuidNeeded = true
		}
	}

	sb.WriteString(fmt.Sprintf("package %s\n", params.Package))
	sb.WriteRune('\n')
	sb.WriteString("import \"time\" \n")
	sb.WriteString("import \"github.com/tislib/data-handler/pkg/model\" \n")
	sb.WriteString("import \"github.com/tislib/data-handler/pkg/client\" \n")
	if uuidNeeded {
		sb.WriteString("import \"github.com/google/uuid\" \n")
	}
	sb.WriteString("import \"github.com/tislib/data-handler/pkg/types\" \n")
	sb.WriteString("import \"google.golang.org/protobuf/types/known/structpb\" \n")
	sb.WriteRune('\n')

	writeResourceStruct(&sb, resource, params)
	sb.WriteRune('\n')
	writeResourceStructGetIdFunc(&sb, resource, params)
	sb.WriteRune('\n')
	writeResourceStructToRecordFunc(&sb, resource, params)
	sb.WriteRune('\n')
	writeResourceStructFromRecordFunc(&sb, resource, params)
	sb.WriteRune('\n')
	writeResourceStructFromPropertiesFunc(&sb, resource, params)
	sb.WriteRune('\n')
	writeResourceStructToPropertiesFunc(&sb, resource, params)
	sb.WriteRune('\n')
	writeResourceGetResourceNameFunc(&sb, resource, params)
	sb.WriteRune('\n')
	writeResourceGetNamespaceNameFunc(&sb, resource, params)
	sb.WriteRune('\n')
	writeResourceNewRepository(&sb, resource, params)
	sb.WriteRune('\n')

	formatted, err := format.Source([]byte(sb.String()))
	if err != nil {
		log.Print(sb.String())
		panic(err)
	}

	return string(formatted)
}

func writeResourceNewRepository(sb *strings.Builder, resource *model.Resource, params GenerateResourceCodeParams) {
	resourceStructName := dashCaseToCamelCase(resource.Name)

	sb.WriteString(fmt.Sprintf("func New%sRepository(dhClient client.DhClient) client.Repository[*%s] {\n", resourceStructName, resourceStructName))
	sb.WriteString(fmt.Sprintf("return client.NewRepository[*%s](dhClient, client.RepositoryParams[*%s]{Instance: new(%s)}) \n", resourceStructName, resourceStructName, resourceStructName))
	sb.WriteString("}\n")
}

func writeResourceGetResourceNameFunc(sb *strings.Builder, resource *model.Resource, params GenerateResourceCodeParams) {
	sb.WriteString(fmt.Sprintf("func (s*%s) GetResourceName() string {\n", dashCaseToCamelCase(resource.Name)))
	sb.WriteString(fmt.Sprintf("return \"%s\"\n", resource.Name))
	sb.WriteString("}\n")
}

func writeResourceGetNamespaceNameFunc(sb *strings.Builder, resource *model.Resource, params GenerateResourceCodeParams) {
	sb.WriteString(fmt.Sprintf("func (s*%s) GetNamespace() string {\n", dashCaseToCamelCase(resource.Name)))
	sb.WriteString(fmt.Sprintf("return \"%s\"\n", resource.Namespace))
	sb.WriteString("}\n")
}

func writeResourceStructGetIdFunc(sb *strings.Builder, resource *model.Resource, params GenerateResourceCodeParams) {
	sb.WriteString(fmt.Sprintf("func (s*%s) GetId() string {\n", dashCaseToCamelCase(resource.Name)))
	sb.WriteString("return s.Id\n")
	sb.WriteString("}\n")
}

func writeResourceStructToRecordFunc(sb *strings.Builder, resource *model.Resource, params GenerateResourceCodeParams) {
	sb.WriteString(fmt.Sprintf("func (s*%s) ToRecord() *model.Record {\n", dashCaseToCamelCase(resource.Name)))
	sb.WriteString(" var rec = &model.Record{} \n")
	sb.WriteString(" rec.Id = s.Id \n")
	sb.WriteString(" rec.Properties = s.ToProperties(false) \n")
	sb.WriteRune('\n')
	sb.WriteString("return rec\n")
	sb.WriteString("}\n")
}

func writeResourceStructToPropertiesFunc(sb *strings.Builder, resource *model.Resource, params GenerateResourceCodeParams) {
	sb.WriteString(fmt.Sprintf("func (s*%s) ToProperties(includeTopProperties bool) map[string]*structpb.Value {\n", dashCaseToCamelCase(resource.Name)))
	sb.WriteString(" var properties = make(map[string]*structpb.Value)\n")
	sb.WriteString(" if includeTopProperties {\n")
	sb.WriteString(" properties[\"id\"] = structpb.NewStringValue(s.Id) \n")
	sb.WriteString(" }\n")
	sb.WriteRune('\n')

	for i, prop := range resource.Properties {
		if !prop.Required || prop.Type == model.ResourceProperty_REFERENCE {
			sb.WriteString(fmt.Sprintf("if s.%s != nil { \n", dashCaseToCamelCase(prop.Name)))
		}

		if prop.Type == model.ResourceProperty_REFERENCE {
			sb.WriteString(fmt.Sprintf("properties[\"%s\"] = structpb.NewStructValue(&structpb.Struct{Fields: s.%s.ToProperties(true)})\n", prop.Name, dashCaseToCamelCase(prop.Name)))
		} else {
			sb.WriteString(fmt.Sprintf("val%d, _ := types.ByResourcePropertyType(model.ResourcePropertyType_%s).Pack(s.%s) \n", i, prop.Type.String(), dashCaseToCamelCase(prop.Name)))
			sb.WriteString(fmt.Sprintf("properties[\"%s\"] = val%d\n", prop.Name, i))
		}

		if !prop.Required || prop.Type == model.ResourceProperty_REFERENCE {
			sb.WriteString("}\n")
		}
		sb.WriteString("\n")
	}
	sb.WriteString("return properties\n")
	sb.WriteString("}\n")
}

func writeResourceStructFromRecordFunc(sb *strings.Builder, resource *model.Resource, params GenerateResourceCodeParams) {
	sb.WriteString(fmt.Sprintf("func (s*%s) FromRecord(record *model.Record) {\n", dashCaseToCamelCase(resource.Name)))
	sb.WriteString(" s.Id = record.Id \n")
	sb.WriteRune('\n')
	sb.WriteString("s.FromProperties(record.Properties)")

	sb.WriteString("}\n")
}

func writeResourceStructFromPropertiesFunc(sb *strings.Builder, resource *model.Resource, params GenerateResourceCodeParams) {
	sb.WriteString(fmt.Sprintf("func (s*%s) FromProperties(properties map[string]*structpb.Value) {\n", dashCaseToCamelCase(resource.Name)))
	sb.WriteString(" if properties[\"id\"] != nil { \n s.Id = properties[\"id\"].GetStringValue() \n } \n")
	sb.WriteRune('\n')

	for i, prop := range resource.Properties {
		sb.WriteString(fmt.Sprintf("if properties[\"%s\"] != nil { \n", prop.Name))

		if prop.Type == model.ResourceProperty_REFERENCE {
			sb.WriteString(fmt.Sprintf("s.%s = new(%s)\n", dashCaseToCamelCase(prop.Name), dashCaseToCamelCase(prop.Reference.ReferencedResource)))
			sb.WriteString(fmt.Sprintf("s.%s.FromProperties(properties[\"%s\"].GetStructValue().Fields)\n", dashCaseToCamelCase(prop.Reference.ReferencedResource), prop.Name))
		} else {
			sb.WriteString(fmt.Sprintf("val%d, _ := types.ByResourcePropertyType(model.ResourcePropertyType_%s).UnPack(properties[\"%s\"]) \n", i, prop.Type.String(), prop.Name))
			if prop.Required {
				sb.WriteString(fmt.Sprintf("s.%s = val%d.(%s)\n", dashCaseToCamelCase(prop.Name), i, getGoType(prop.Type)))
			} else {

				sb.WriteString(fmt.Sprintf("s.%s = new(%s)\n", dashCaseToCamelCase(prop.Name), getGoType(prop.Type)))
				sb.WriteString(fmt.Sprintf("*s.%s = val%d.(%s)\n", dashCaseToCamelCase(prop.Name), i, getGoType(prop.Type)))
			}
		}
		sb.WriteString("}\n\n")
	}

	sb.WriteString("}\n")
}

func writeResourceStruct(sb *strings.Builder, resource *model.Resource, params GenerateResourceCodeParams) {
	sb.WriteString(fmt.Sprintf("type %s struct {\n", dashCaseToCamelCase(resource.Name)))

	if !annotations.IsEnabled(resource, annotations.DoPrimaryKeyLookup) {
		sb.WriteString("    Id string\n")
	}

	for _, prop := range resource.Properties {
		typeDef := getGoType(prop.Type)
		if !prop.Required {
			typeDef = "*" + typeDef
		}

		if prop.Type == model.ResourceProperty_REFERENCE {
			typeDef = "*" + dashCaseToCamelCase(prop.Reference.ReferencedResource)
		}

		sb.WriteString(fmt.Sprintf("    %s %s\n", dashCaseToCamelCase(prop.Name), typeDef))
	}

	if !annotations.IsEnabled(resource, annotations.DisableAudit) {
		sb.WriteString("    CreatedBy string\n")
		sb.WriteString("    UpdatedBy *string\n")
		sb.WriteString("    CreatedOn time.Time\n")
		sb.WriteString("    UpdatedOn *time.Time\n")
	}

	if !annotations.IsEnabled(resource, annotations.DisableVersion) {
		sb.WriteString("    Version uint64\n")
	}

	sb.WriteString("}\n")
}

func dashCaseToCamelCase(inputUnderScoreStr string) (camelCase string) {
	//snake_case to camelCase

	isToUpper := false

	for k, v := range inputUnderScoreStr {
		if k == 0 {
			camelCase = strings.ToUpper(string(inputUnderScoreStr[0]))
		} else {
			if isToUpper {
				camelCase += strings.ToUpper(string(v))
				isToUpper = false
			} else {
				if v == '-' {
					isToUpper = true
				} else if v == '_' {
					isToUpper = true
				} else {
					camelCase += string(v)
				}
			}
		}
	}
	return

}

func getGoType(propertyType model.ResourceProperty_Type) string {
	typ := types.ByResourcePropertyType(propertyType)

	return reflect.TypeOf(typ.Default()).String()
}

package postgres

import (
	"data-handler/stub/model"
	"database/sql"
	"github.com/huandu/go-sqlbuilder"
	"strconv"
	"time"
)

type QueryRunner interface {
	QueryRow(query string, args ...any) *sql.Row
	Exec(query string, args ...any) (sql.Result, error)
}

func resourceCountsByName(runner QueryRunner, resourceName string) (int, error) {
	res := runner.QueryRow("select count(*) as count from resource where name = $1", resourceName)

	var count = new(int)
	err := res.Scan(count)

	return *count, err
}

func resourceCreateTable(runner QueryRunner, resource *model.Resource) error {
	builder := sqlbuilder.CreateTable(resource.SourceConfig.Mapping)

	builder.IfNotExists()

	builder.Define("id", "uuid", "NOT NULL", "PRIMARY KEY")

	for _, property := range resource.Properties {
		if sourceConfig, ok := property.SourceConfig.(*model.ResourceProperty_Mapping); ok {
			nullModifier := "NULL"
			if property.Required {
				nullModifier = "NOT NULL"
			}
			sqlType := getPsqlType(property.Type, property.Length)
			builder.Define(sourceConfig.Mapping, sqlType, nullModifier)
		}
	}

	// audit
	builder.Define("created_on", "timestamp", "NOT NULL")
	builder.Define("updated_on", "timestamp", "NULL")
	builder.Define("created_by", DbNameType, "NOT NULL")
	builder.Define("updated_by", DbNameType, "NULL")
	// version
	builder.Define("version", "int2", "NOT NULL")

	sqlQuery, _ := builder.Build()
	_, err := runner.Exec(sqlQuery)

	return err
}

func resourceInsert(runner QueryRunner, resource *model.Resource) error {
	if resource.Flags == nil {
		resource.Flags = &model.ResourceFlags{}
	}

	insertBuilder := sqlbuilder.InsertInto("resource")
	insertBuilder.SetFlavor(sqlbuilder.PostgreSQL)
	insertBuilder.Cols(
		"name",
		"workspace",
		"type",
		"source_data_source",
		"source_mapping",
		"read_only_records",
		"unique_record",
		"keep_history",
		"created_on",
		"updated_on",
		"created_by",
		"updated_by",
		"version")
	insertBuilder.Values(
		resource.Name,
		resource.Workspace,
		resource.Type.Number(),
		resource.SourceConfig.DataSource,
		resource.SourceConfig.Mapping,
		resource.Flags.ReadOnlyRecords,
		resource.Flags.UniqueRecord,
		resource.Flags.KeepHistory,
		time.Now(),
		nil,
		"test-usr",
		nil,
		1,
	)

	sqlQuery, args := insertBuilder.Build()

	_, err := runner.Exec(sqlQuery, args...)

	return err
}

func resourceInsertProperties(runner QueryRunner, resource *model.Resource) error {
	for _, property := range resource.Properties {
		propertyInsertBuilder := sqlbuilder.InsertInto("resource_property")
		propertyInsertBuilder.SetFlavor(sqlbuilder.PostgreSQL)
		propertyInsertBuilder.Cols(
			"resource_name",
			"property_name",
			"type",
			"source_type",
			"source_mapping",
			"required",
			"length",
		)
		sourceType := 0
		propertyInsertBuilder.Values(
			resource.Name,
			property.Name,
			property.Type,
			sourceType,
			property.SourceConfig.(*model.ResourceProperty_Mapping).Mapping,
			property.Required,
			property.Length,
		)

		sql, args := propertyInsertBuilder.Build()

		_, err := runner.Exec(sql, args...)

		if err != nil {
			return err
		}
	}

	return nil
}

func getPsqlType(propertyType model.ResourcePropertyType, length uint32) string {
	switch propertyType {
	case model.ResourcePropertyType_INT32:
		return "INT"
	case model.ResourcePropertyType_STRING:
		return "VARCHAR(" + strconv.Itoa(int(length)) + ")"
	default:

		panic("unknown property type")
	}
}

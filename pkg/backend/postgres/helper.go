package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/tislib/data-handler/pkg/errors"
	"github.com/tislib/data-handler/pkg/logging"
	"github.com/tislib/data-handler/pkg/model"
	"github.com/tislib/data-handler/pkg/service/annotations"
	"github.com/tislib/data-handler/pkg/types"
	"google.golang.org/protobuf/types/known/structpb"
	"net"
	"runtime/debug"
)

func locatePropertyByName(resource *model.Resource, propertyName string) *model.ResourceProperty {
	for _, property := range resource.Properties {
		if property.Name == propertyName {
			return property
		}
	}

	return nil
}

func handleDbError(ctx context.Context, err error) errors.ServiceError {
	logger := log.WithFields(logging.CtxFields(ctx))

	if err == nil {
		return nil
	}

	if err == sql.ErrNoRows {
		return errors.RecordNotFoundError
	}

	logger.Errorf("Db Error: %s", err)
	logger.Tracef("Stack: " + string(debug.Stack()))

	if err == sql.ErrTxDone {
		logger.Panic("Illegal situation")
	}

	if _, ok := err.(errors.ServiceError); ok {
		logger.Panic("database error is expected: ", err)
	}

	if pqErr, ok := err.(*pq.Error); ok {
		return handlePqErr(pqErr)
	}

	if netErr, ok := err.(*net.OpError); ok {
		return errors.InternalError.WithDetails(netErr.Error())
	}

	if err.Error() == "context cancelled" {
		return errors.InternalError.WithDetails(err.Error())
	}

	logger.Print("Unhandled Error: ", err)
	return errors.InternalError.WithDetails(err.Error())
}

func handlePqErr(err *pq.Error) errors.ServiceError {
	switch err.Code {
	case "28000":
		return errors.BackendConnectionAuthenticationError.WithMessage(err.Message)
	case "28P01":
		return errors.BackendConnectionAuthenticationError.WithMessage(err.Message)
	case "23505":
		return errors.UniqueViolation.WithDetails(err.Message)
	case "23503":
		return errors.ForeignKeyViolation.WithDetails(err.Message)
	default:
		return errors.InternalError.WithMessage(err.Message)
	}
}

func DbEncode(property *model.ResourceProperty, packedVal *structpb.Value) (interface{}, errors.ServiceError) {
	propertyType := types.ByResourcePropertyType(property.Type)
	var val interface{}

	if property.Type == model.ResourcePropertyType_TYPE_OBJECT || property.Type == model.ResourcePropertyType_TYPE_ENUM || property.Type == model.ResourcePropertyType_TYPE_MAP || property.Type == model.ResourcePropertyType_TYPE_LIST {
		var err error
		val, err = json.Marshal(packedVal.AsInterface())

		if err != nil {
			return nil, errors.InternalError.WithDetails(err.Error())
		}
		val = string(val.([]byte))
	} else {
		var err error
		val, err = propertyType.UnPack(packedVal)

		if err != nil {
			return nil, errors.InternalError.WithDetails(err.Error())
		}
	}
	return val, nil
}

func getTableName(sourceConfig *model.ResourceSourceConfig, history bool) string {
	def := ""

	if sourceConfig.Catalog != "" {
		def = fmt.Sprintf("\"%s\".", sourceConfig.Catalog)
	}

	def = fmt.Sprintf("\"%s\"", sourceConfig.Entity)

	if history {
		def = fmt.Sprintf("\"%s_h\"", sourceConfig.Entity)
	} else {
		def = fmt.Sprintf("\"%s\"", sourceConfig.Entity)
	}

	return def
}

func checkHasOwnId(resource *model.Resource) bool {
	return !annotations.IsEnabled(resource, annotations.DoPrimaryKeyLookup)
}

func locatePrimaryKey(resource *model.Resource) (string, errors.ServiceError) {
	for _, property := range resource.Properties {
		if property.Primary {
			return property.Mapping, nil
		}
	}

	return "", errors.UnableToLocatePrimaryKey
}

func prepareResourceRecordCols(resource *model.Resource) []string {
	var cols []string

	if checkHasOwnId(resource) {
		cols = append(cols, "id")
	}

	for _, property := range resource.Properties {
		col := fmt.Sprintf("%s", property.Mapping)
		cols = append(cols, col)
	}

	// referenced columns

	if !annotations.IsEnabled(resource, annotations.DisableAudit) {
		cols = append(cols, "created_on")
		cols = append(cols, "updated_on")
		cols = append(cols, "created_by")
		cols = append(cols, "updated_by")
		cols = append(cols, "version")
	}
	return cols
}
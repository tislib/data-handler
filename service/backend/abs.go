package backend

import (
	"context"
	"data-handler/model"
	"data-handler/service/errors"
)

type DataSourceConnectionDetails interface {
}

type AddResourceParams struct {
	Resource       *model.Resource
	IgnoreIfExists bool
	Migrate        bool
	ForceMigrate   bool
}

type BulkRecordsParams struct {
	Resource       *model.Resource
	Records        []*model.Record
	CheckVersion   bool
	IgnoreIfExists bool
}

type ListRecordParams struct {
	Resource          *model.Resource
	Query             *model.BooleanExpression
	Limit             uint32
	Offset            uint64
	UseHistory        bool
	ResolveReferences bool
}

type Constructor func(dataSource DataSourceConnectionDetails) Backend

type GenericInterface interface {
	GetStatus(ctx context.Context) (connectionAlreadyInitiated bool, testConnection bool, err errors.ServiceError)
	DestroyDataSource(ctx context.Context)
}

type RecordsInterface interface {
	AddRecords(ctx context.Context, params BulkRecordsParams) ([]*model.Record, bool, errors.ServiceError)
	UpdateRecords(ctx context.Context, params BulkRecordsParams) ([]*model.Record, errors.ServiceError)
	GetRecord(ctx context.Context, resource *model.Resource, id string) (*model.Record, errors.ServiceError)
	DeleteRecords(ctx context.Context, resource *model.Resource, list []string) errors.ServiceError
	ListRecords(ctx context.Context, params ListRecordParams) ([]*model.Record, uint32, errors.ServiceError)
}

type SchemaInterface interface {
	ListEntities(ctx context.Context) ([]string, errors.ServiceError)
	PrepareResourceFromEntity(ctx context.Context, catalog, entity string) (*model.Resource, errors.ServiceError)
	UpgradeResource(ctx context.Context, resource *model.Resource, forceMigration bool) errors.ServiceError
	DowngradeResource(ctx context.Context, resource *model.Resource, migration bool) errors.ServiceError
}

type Backend interface {
	GenericInterface
	RecordsInterface
	SchemaInterface
}

package abs

import (
	"context"
	"github.com/tislib/data-handler/pkg/errors"
	"github.com/tislib/data-handler/pkg/model"
)

type AuthenticationService interface {
	Init(data *model.InitData)
	Authenticate(ctx context.Context, username string, password string, term model.TokenTerm) (*model.Token, errors.ServiceError)
	RenewToken(ctx context.Context, token string, term model.TokenTerm) (*model.Token, errors.ServiceError)
	ParseAndVerifyToken(token string) (*UserDetails, errors.ServiceError)
}

type BackendProviderService interface {
	Init(data *model.InitData)
	GetSystemBackend(ctx context.Context) Backend
	GetBackendByDataSourceId(ctx context.Context, dataSourceId string) (Backend, errors.ServiceError)
	GetBackendByDataSourceName(ctx context.Context, dataSourceId string) (Backend, errors.ServiceError)
	DestroyBackend(ctx context.Context, id string) error
}

type DataSourceService interface {
	Init(*model.InitData)
	ListEntities(ctx context.Context, id string) ([]*model.DataSourceCatalog, errors.ServiceError)
	List(ctx context.Context) ([]*model.DataSource, errors.ServiceError)
	GetStatus(ctx context.Context, id string) (connectionAlreadyInitiated bool, testConnection bool, err errors.ServiceError)
	Create(ctx context.Context, sources []*model.DataSource) ([]*model.DataSource, errors.ServiceError)
	Update(ctx context.Context, sources []*model.DataSource) ([]*model.DataSource, errors.ServiceError)
	PrepareResourceFromEntity(ctx context.Context, dataSourceId string, catalog, entity string) (*model.Resource, errors.ServiceError)
	Get(ctx context.Context, id string) (*model.DataSource, errors.ServiceError)
	Delete(ctx context.Context, ids []string) errors.ServiceError
}

type PluginService interface {
	Init(data *model.InitData)
}

type RecordService interface {
	PrepareQuery(resource *model.Resource, queryMap map[string]interface{}) (*model.BooleanExpression, errors.ServiceError)
	GetRecord(ctx context.Context, namespace, resourceName, id string) (*model.Record, errors.ServiceError)
	FindBy(ctx context.Context, namespace, resourceName, propertyName string, value interface{}) (*model.Record, errors.ServiceError)

	Init(data *model.InitData)

	List(ctx context.Context, params RecordListParams) ([]*model.Record, uint32, errors.ServiceError)
	Create(ctx context.Context, params RecordCreateParams) ([]*model.Record, []bool, errors.ServiceError)
	Update(ctx context.Context, params RecordUpdateParams) ([]*model.Record, errors.ServiceError)
	Get(ctx context.Context, params RecordGetParams) (*model.Record, errors.ServiceError)
	Delete(ctx context.Context, params RecordDeleteParams) errors.ServiceError
}

type ResourceService interface {
	Init(data *model.InitData)
	CheckResourceExists(ctx context.Context, namespace, name string) bool
	GetResourceByName(ctx context.Context, namespace, resource string) *model.Resource
	GetSystemResourceByName(ctx context.Context, resourceName string) *model.Resource
	Create(ctx context.Context, resource *model.Resource, doMigration bool, forceMigration bool) (*model.Resource, errors.ServiceError)
	Update(ctx context.Context, resource *model.Resource, doMigration bool, forceMigration bool) errors.ServiceError
	Delete(ctx context.Context, ids []string, doMigration bool, forceMigration bool) errors.ServiceError
	List(ctx context.Context) []*model.Resource
	ReloadSchema(ctx context.Context) errors.ServiceError
	Get(ctx context.Context, id string) *model.Resource
	GetSchema() *Schema
	PrepareResourceMigrationPlan(ctx context.Context, resources []*model.Resource, prepareFromDataSource bool) ([]*model.ResourceMigrationPlan, errors.ServiceError)
}

type ResourceMigrationService interface {
	PreparePlan(ctx context.Context, existingResource *model.Resource, resource *model.Resource) (*model.ResourceMigrationPlan, errors.ServiceError)
}

type UserService interface {
	Init(data *model.InitData)
	Create(ctx context.Context, users []*model.User) ([]*model.User, errors.ServiceError)
	Update(ctx context.Context, users []*model.User) ([]*model.User, errors.ServiceError)
	Delete(ctx context.Context, ids []string) errors.ServiceError
	Get(ctx context.Context, id string) (*model.User, errors.ServiceError)
	List(ctx context.Context, query *model.BooleanExpression, limit uint32, offset uint64) ([]*model.User, errors.ServiceError)
	InjectBackendProviderService(service BackendProviderService)
}

type WatchService interface {
	Watch(ctx context.Context, params WatchParams) <-chan *model.WatchMessage
}

type NamespaceService interface {
	Init(data *model.InitData)
	Create(ctx context.Context, namespaces []*model.Namespace) ([]*model.Namespace, errors.ServiceError)
	Update(ctx context.Context, namespaces []*model.Namespace) ([]*model.Namespace, errors.ServiceError)
	Delete(ctx context.Context, ids []string) errors.ServiceError
	Get(ctx context.Context, id string) (*model.Namespace, errors.ServiceError)
	List(ctx context.Context) ([]*model.Namespace, errors.ServiceError)
}

type ExtensionService interface {
	Init(*model.InitData)
	List(ctx context.Context) ([]*model.RemoteExtension, errors.ServiceError)
	Create(ctx context.Context, sources []*model.RemoteExtension) ([]*model.RemoteExtension, errors.ServiceError)
	Update(ctx context.Context, sources []*model.RemoteExtension) ([]*model.RemoteExtension, errors.ServiceError)
	Get(ctx context.Context, id string) (*model.RemoteExtension, errors.ServiceError)
	Delete(ctx context.Context, ids []string) errors.ServiceError
	RegisterExtension(extension Extension)
	UnRegisterExtension(extension Extension)
}

type WatchParams struct {
	Namespace  string
	Resource   string
	Query      *model.BooleanExpression
	BufferSize int
}

type RecordListParams struct {
	Query             *model.BooleanExpression
	Namespace         string
	Resource          string
	Limit             uint32
	Offset            uint64
	UseHistory        bool
	ResolveReferences []string
	ResultChan        chan<- *model.Record
	PackRecords       bool
}

type RecordCreateParams struct {
	Namespace      string
	Resource       string
	Records        []*model.Record
	IgnoreIfExists bool
}

type RecordUpdateParams struct {
	Namespace    string
	Resource     string
	Records      []*model.Record
	CheckVersion bool
}

type RecordGetParams struct {
	Namespace string
	Resource  string
	Id        string
}

type RecordDeleteParams struct {
	Namespace string
	Resource  string
	Ids       []string
}

type UserDetails struct {
	Username        string
	SecurityContext *model.SecurityContext
}

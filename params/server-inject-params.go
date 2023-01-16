package params

import (
	"data-handler/service"
)

type ServerInjectionConstructorParams struct {
	ResourceService       service.ResourceService
	RecordService         service.RecordService
	AuthenticationService service.AuthenticationService
	DataSourceService     service.DataSourceService
	NamespaceService      service.NamespaceService
	UserService           service.UserService
	WatchService          service.WatchService
}

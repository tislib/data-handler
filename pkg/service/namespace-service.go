package service

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/tislib/data-handler/pkg/abs"
	"github.com/tislib/data-handler/pkg/errors"
	"github.com/tislib/data-handler/pkg/model"
	"github.com/tislib/data-handler/pkg/resources"
	mapping2 "github.com/tislib/data-handler/pkg/resources/mapping"
	"github.com/tislib/data-handler/pkg/service/security"
)

type namespaceService struct {
	recordService          abs.RecordService
	serviceName            string
	resourceService        abs.ResourceService
	backendProviderService abs.BackendProviderService
}

func (d *namespaceService) Create(ctx context.Context, namespaces []*model.Namespace) ([]*model.Namespace, errors.ServiceError) {
	// insert records via resource service
	records := mapping2.MapToRecord(namespaces, mapping2.NamespaceToRecord)

	result, _, err := d.recordService.Create(ctx, abs.RecordCreateParams{
		Namespace: resources.NamespaceResource.Namespace,
		Resource:  resources.NamespaceResource.Name,
		Records:   records,
	})

	if err != nil {
		return nil, err
	}

	return mapping2.MapFromRecord(result, mapping2.NamespaceFromRecord), nil
}

func (d *namespaceService) Update(ctx context.Context, namespaces []*model.Namespace) ([]*model.Namespace, errors.ServiceError) {
	// insert records via resource service
	records := mapping2.MapToRecord(namespaces, mapping2.NamespaceToRecord)

	result, err := d.recordService.Update(ctx, abs.RecordUpdateParams{
		Namespace: resources.NamespaceResource.Namespace,
		Resource:  resources.NamespaceResource.Name,
		Records:   records,
	})

	if err != nil {
		return nil, err
	}

	return mapping2.MapFromRecord(result, mapping2.NamespaceFromRecord), nil
}

func (d *namespaceService) Delete(ctx context.Context, ids []string) errors.ServiceError {
	return d.recordService.Delete(ctx, abs.RecordDeleteParams{
		Namespace: resources.NamespaceResource.Namespace,
		Resource:  resources.NamespaceResource.Name,
		Ids:       ids,
	})
}

func (d *namespaceService) Get(ctx context.Context, id string) (*model.Namespace, errors.ServiceError) {

	record, err := d.recordService.Get(ctx, abs.RecordGetParams{
		Namespace: resources.NamespaceResource.Namespace,
		Resource:  resources.NamespaceResource.Name,
		Id:        id,
	})

	if err != nil {
		return nil, err
	}

	return mapping2.NamespaceFromRecord(record), nil
}

func (d *namespaceService) List(ctx context.Context) ([]*model.Namespace, errors.ServiceError) {

	result, _, err := d.recordService.List(ctx, abs.RecordListParams{
		Namespace: resources.NamespaceResource.Namespace,
		Resource:  resources.NamespaceResource.Name,
	})

	if err != nil {
		return nil, err
	}

	return mapping2.MapFromRecord(result, mapping2.NamespaceFromRecord), err
}

func (d *namespaceService) Init(data *model.InitData) {
	if len(data.InitNamespaces) > 0 {
		_, _, err := d.recordService.Create(security.SystemContext, abs.RecordCreateParams{
			Namespace:      resources.NamespaceResource.Namespace,
			Resource:       resources.NamespaceResource.Name,
			Records:        mapping2.MapToRecord(data.InitNamespaces, mapping2.NamespaceToRecord),
			IgnoreIfExists: true,
		})

		if err != nil {
			log.Fatal(err)
		}
	}

	_, _, err := d.recordService.Create(security.SystemContext, abs.RecordCreateParams{
		Namespace: resources.NamespaceResource.Namespace,
		Resource:  resources.NamespaceResource.Name,
		Records: []*model.Record{mapping2.NamespaceToRecord(&model.Namespace{
			Name:        "default",
			Description: "default namespace",
		})},
		IgnoreIfExists: true,
	})

	if err != nil {
		log.Fatal(err)
	}
}

func NewNamespaceService(resourceService abs.ResourceService, recordService abs.RecordService, backendProviderService abs.BackendProviderService) abs.NamespaceService {
	return &namespaceService{
		serviceName:            "NamespaceService",
		resourceService:        resourceService,
		recordService:          recordService,
		backendProviderService: backendProviderService,
	}
}

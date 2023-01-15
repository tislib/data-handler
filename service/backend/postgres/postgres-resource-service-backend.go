package postgres

import (
	"context"
	"data-handler/model"
	"data-handler/service/backend"
	"data-handler/service/errors"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

const DbNameType = "VARCHAR(64)"

type postgresResourceServiceBackend struct {
	connectionDetails *model.DataSource_PostgresqlParams
	connection        *sql.DB
}

func (p *postgresResourceServiceBackend) ListResources(ctx context.Context) (result []*model.Resource, err errors.ServiceError) {
	err = p.withBackend(ctx, true, func(tx *sql.Tx) errors.ServiceError {
		result, err = resourceList(ctx, tx)

		for _, resource := range result {
			if err = resourceLoadProperties(tx, resource, resource.Workspace, resource.Name); err != nil {
				log.Errorf("Unable to load resource properties for %s/%s Err: %s", resource.Workspace, resource.Name, err)
				return err
			}

			if err = resourceLoadReferences(tx, resource, resource.Workspace, resource.Name); err != nil {
				log.Errorf("Unable to load resource references for %s/%s Err: %s", resource.Workspace, resource.Name, err)
				return err
			}
		}

		return err
	})

	return
}

func (p *postgresResourceServiceBackend) ListEntities(ctx context.Context) (result []string, err errors.ServiceError) {
	err = p.withBackend(ctx, true, func(tx *sql.Tx) errors.ServiceError {
		result, err = resourceListEntities(ctx, tx)

		return err
	})

	return
}

func (p *postgresResourceServiceBackend) DestroyDataSource(ctx context.Context) {
	if p.connection != nil {
		p.connection.Close()

		p.connection = nil
	}
}

func (p *postgresResourceServiceBackend) GetStatus(ctx context.Context) (connectionAlreadyInitiated bool, testConnection bool, err errors.ServiceError) {
	conn, err := p.acquireConnection(ctx)

	if err != nil {
		return
	}

	err = handleDbError(conn.Ping())

	testConnection = err == nil

	return
}

func (p *postgresResourceServiceBackend) PrepareResourceFromEntity(ctx context.Context, entity string) (resource *model.Resource, err errors.ServiceError) {
	err = p.withBackend(ctx, false, func(tx *sql.Tx) errors.ServiceError {
		if resource, err = resourcePrepareResourceFromEntity(ctx, tx, entity); err != nil {
			log.Errorf("[PrepareResourceFromEntity] Unable to load resource details for %s Err: %s", entity, err)
			return err
		}

		resource.SourceConfig = &model.ResourceSourceConfig{
			Mapping: entity,
		}

		return nil
	})

	if err != nil {
		log.Errorf("Unable to load resource for %s Err: %s", entity, err)
		return nil, err
	}

	return resource, nil
}

func (p *postgresResourceServiceBackend) GetResource(ctx context.Context, workspace string, id string) (resource *model.Resource, err errors.ServiceError) {
	if workspace == "" {
		workspace = "default"
	}

	resource = new(model.Resource)

	err = p.withBackend(ctx, true, func(tx *sql.Tx) errors.ServiceError {
		if err = resourceLoadDetails(tx, resource, workspace, id); err != nil {
			log.Errorf("Unable to load resource details for %s/%s Err: %s", workspace, id, err)
			return err
		}

		if err = resourceLoadProperties(tx, resource, workspace, resource.Name); err != nil {
			log.Errorf("Unable to load resource properties for %s/%s Err: %s", workspace, id, err)
			return err
		}

		if err = resourceLoadReferences(tx, resource, workspace, resource.Name); err != nil {
			log.Errorf("Unable to load resource references for %s/%s Err: %s", workspace, id, err)
			return err
		}

		return nil
	})

	if err != nil {
		log.Errorf("Unable to load resource for %s/%s Err: %s", workspace, id, err)
		return nil, err
	}

	return resource, nil
}

func (p *postgresResourceServiceBackend) GetResourceByName(ctx context.Context, workspace string, resourceName string) (resource *model.Resource, err errors.ServiceError) {
	if workspace == "" {
		workspace = "default"
	}

	resource = new(model.Resource)

	err = p.withBackend(ctx, true, func(tx *sql.Tx) errors.ServiceError {
		if err = resourceLoadDetailsByName(tx, resource, workspace, resourceName); err != nil {
			log.Errorf("Unable to load resource details for %s/%s Err: %s", workspace, resourceName, err)
			return err
		}

		if err = resourceLoadProperties(tx, resource, workspace, resourceName); err != nil {
			log.Errorf("Unable to load resource properties for %s/%s Err: %s", workspace, resourceName, err)
			return err
		}

		if err = resourceLoadReferences(tx, resource, workspace, resourceName); err != nil {
			log.Errorf("Unable to load resource references for %s/%s Err: %s", workspace, resourceName, err)
			return err
		}

		return nil
	})

	if err != nil {
		log.Errorf("Unable to load resource for %s/%s Err: %s", workspace, resourceName, err)
		return nil, err
	}

	return resource, nil
}

func (p *postgresResourceServiceBackend) AddResource(ctx context.Context, params backend.AddResourceParams) (*model.Resource, errors.ServiceError) {
	if params.Resource.Workspace == "" {
		params.Resource.Workspace = "default"
	}

	err := p.withBackend(ctx, false, func(tx *sql.Tx) errors.ServiceError {
		// check if resource exists

		if existingCount, err := resourceCountsByName(tx, params.Resource.Workspace, params.Resource.Name); err != nil {
			return err
		} else if existingCount > 0 {
			if params.IgnoreIfExists {
				return nil
			}
			return errors.AlreadyExistsError
		}

		newId, err := uuid.NewRandom()

		if err != nil {
			return errors.InternalError.WithDetails(err.Error())
		}

		params.Resource.Id = newId.String()

		if err := resourceInsert(tx, params.Resource); err != nil {
			log.Error("Unable to insert resource: ", err)
			return err
		}

		if err := resourceUpsertProperties(tx, params.Resource); err != nil {
			log.Error("Unable to insert resource properties: ", err)
			return err
		}

		if err := resourceUpsertReferences(tx, params.Resource); err != nil {
			log.Error("Unable to insert resource properties: ", err)
			return err
		}

		//if params.Migrate {
		//	err := p.withBackend(ctx, params.Resource.SourceConfig.DataSource, false, func(tx *sql.Tx) errors.ServiceError {
		//		err := resourceCreateTable(tx, params.Resource)
		//		if err != nil {
		//			return err
		//		}
		//
		//		for _, reference := range params.Resource.References {
		//			if err = resourceCreateForeignKey(tx, params.Resource, reference); err != nil {
		//				if err != nil {
		//					return err
		//				}
		//			}
		//		}
		//
		//		if params.Resource.Flags.KeepHistory {
		//			if params.Resource.Flags.DisableAudit {
		//				return errors.LogicalError.WithMessage("history cannot created while audit is disabled")
		//			}
		//			err = resourceCreateHistoryTable(tx, params.Resource)
		//			if err != nil {
		//				return err
		//			}
		//		}
		//		return nil
		//	})
		//
		//	if err != nil {
		//		return err
		//	}
		//}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return params.Resource, nil
}

func (p *postgresResourceServiceBackend) UpdateResource(ctx context.Context, resource *model.Resource, doMigration bool, forceMigration bool) errors.ServiceError {
	var err errors.ServiceError
	return p.withBackend(ctx, false, func(tx *sql.Tx) errors.ServiceError {
		if err = resourceUpdate(ctx, tx, resource); err != nil {
			return err
		}

		doResourceCleanup(resource)

		if err = resourceUpsertProperties(tx, resource); err != nil {
			return err
		}

		if err = resourceUpsertReferences(tx, resource); err != nil {
			return err
		}

		//if doMigration {
		//	err = p.withBackend(ctx, resource.SourceConfig.DataSource, false, func(tx *sql.Tx) errors.ServiceError {
		//		if err = resourceCreateTable(tx, resource); err != nil {
		//			return err
		//		}
		//
		//		if err = resourceMigrateTable(ctx, tx, resource, forceMigration, false); err != nil {
		//			return err
		//		}
		//
		//		if resource.Flags.KeepHistory {
		//			if err = resourceCreateHistoryTable(tx, resource); err != nil {
		//				return err
		//			}
		//
		//			if err = resourceMigrateTable(ctx, tx, resource, forceMigration, true); err != nil {
		//				return err
		//			}
		//		}
		//
		//		return nil
		//	})
		//
		//	if err != nil {
		//		return err
		//	}
		//}

		return nil
	})
}

func (p *postgresResourceServiceBackend) DeleteResources(ctx context.Context, workspace string, ids []string, doMigration bool, forceMigration bool) errors.ServiceError {
	if workspace == "" {
		workspace = "default"
	}

	var sources []*model.ResourceSourceConfig

	for _, id := range ids {
		resource, err := p.GetResource(ctx, workspace, id)

		if err != nil {
			return err
		}

		if resource == nil {
			return errors.NotFoundError
		}

		sources = append(sources, resource.SourceConfig)
	}
	err := p.withBackend(ctx, false, func(tx *sql.Tx) errors.ServiceError {
		return resourceDelete(ctx, tx, ids)
	})

	if err != nil {
		return err
	}

	//if doMigration {
	//	for _, source := range sources {
	//		err = p.withBackend(ctx, source.DataSource, false, func(tx *sql.Tx) errors.ServiceError {
	//			return resourceDropTable(tx, source.Mapping)
	//		})
	//
	//		if err != nil {
	//			return err
	//		}
	//	}
	//}

	return nil
}

func (p *postgresResourceServiceBackend) acquireConnection(ctx context.Context) (*sql.DB, errors.ServiceError) {
	if p.connection == nil {

		params := p.connectionDetails.PostgresqlParams

		connStr := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable", params.Username, params.Password, params.Host, params.Port, params.DbName)
		// Connect to database
		conn, sqlErr := sql.Open("postgres", connStr)
		err := handleDbError(sqlErr)

		if err != nil {
			return nil, err
		}

		p.connection = conn

		log.Infof("Connected to Datasource: %s@%s:%d/%s", params.Username, params.Host, params.Port, params.DefaultSchema)
	}

	return p.connection, nil
}

func NewPostgresResourceServiceBackend(connectionDetails backend.DataSourceConnectionDetails) backend.Backend {
	return &postgresResourceServiceBackend{
		connectionDetails: connectionDetails.(*model.DataSource_PostgresqlParams),
	}
}

package postgres

import (
	"context"
	"data-handler/model"
	"data-handler/service/backend"
	"data-handler/service/errors"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

const DbNameType = "VARCHAR(64)"

type postgresResourceServiceBackend struct {
	connectionMap     map[string]*sql.DB
	systemBackend     backend.DataSourceBackend
	dataSourceService backend.DataSourceLocator
}

func (p *postgresResourceServiceBackend) ListResources(ctx context.Context) (result []*model.Resource, err errors.ServiceError) {
	err = p.withSystemBackend(true, func(tx *sql.Tx) errors.ServiceError {
		result, err = resourceList(ctx, tx)

		return err
	})

	return
}

func (p *postgresResourceServiceBackend) ListEntities(ctx context.Context, dataSourceId string) (result []string, err errors.ServiceError) {
	err = p.withBackend(dataSourceId, true, func(tx *sql.Tx) errors.ServiceError {
		result, err = resourceListEntities(ctx, tx)

		return err
	})

	return
}

func (p *postgresResourceServiceBackend) InjectDataSourceService(dataSourceService backend.DataSourceLocator) {
	p.dataSourceService = dataSourceService
}

func (p *postgresResourceServiceBackend) DestroyDataSource(dataSourceId string) {
	if p.connectionMap[dataSourceId] != nil {
		err := p.connectionMap[dataSourceId].Close()

		if err != nil {
			log.Println("Cannot Close destroyed datasource connection: ", err)
		}

		delete(p.connectionMap, dataSourceId)
	}
}

func (p *postgresResourceServiceBackend) Init() {
	p.systemBackend = p.dataSourceService.GetSystemDataSourceBackend()

	err := p.withBackend(p.systemBackend.GetDataSourceId(), false, func(tx *sql.Tx) errors.ServiceError {
		return resourceSetupTables(tx)
	})

	if err != nil {
		panic(err)
	}
}

func (p *postgresResourceServiceBackend) GetStatus(dataSourceId string) (connectionAlreadyInitiated bool, testConnection bool, err errors.ServiceError) {
	connectionAlreadyInitiated = p.connectionMap[dataSourceId] != nil

	conn, err := p.acquireConnection(dataSourceId)

	if err != nil {
		return
	}

	err = handleDbError(conn.Ping())

	testConnection = err == nil

	return
}

func (p *postgresResourceServiceBackend) PrepareResourceFromEntity(ctx context.Context, dataSourceId string, entity string) (resource *model.Resource, err errors.ServiceError) {
	err = p.withBackend(dataSourceId, false, func(tx *sql.Tx) errors.ServiceError {
		if resource, err = resourcePrepareResourceFromEntity(ctx, tx, entity); err != nil {
			log.Error("Unable to load resource details", err)
			return err
		}

		resource.SourceConfig = &model.ResourceSourceConfig{
			DataSource: dataSourceId,
			Mapping:    entity,
		}

		return nil
	})

	if err != nil {
		log.Error("Unable load resource", err)
		return nil, err
	}

	return resource, nil
}

func (p *postgresResourceServiceBackend) GetResource(ctx context.Context, workspace string, id string) (resource *model.Resource, err errors.ServiceError) {
	resource = new(model.Resource)

	err = p.withSystemBackend(true, func(tx *sql.Tx) errors.ServiceError {
		if err = resourceLoadDetails(tx, resource, workspace, id); err != nil {
			log.Error("Unable to load resource details", err)
			return err
		}

		if err = resourceLoadProperties(tx, resource, workspace, resource.Name); err != nil {
			log.Error("Unable to load resource properties", err)
			return err
		}

		if err = resourceLoadReferences(tx, resource, workspace, resource.Name); err != nil {
			log.Error("Unable to load resource properties", err)
			return err
		}

		return nil
	})

	if err != nil {
		log.Error("Unable load resource", err)
		return nil, err
	}

	return resource, nil
}

func (p *postgresResourceServiceBackend) GetResourceByName(ctx context.Context, workspace string, resourceName string) (resource *model.Resource, err errors.ServiceError) {
	resource = new(model.Resource)

	err = p.withSystemBackend(true, func(tx *sql.Tx) errors.ServiceError {
		if err = resourceLoadDetailsByName(tx, resource, workspace, resourceName); err != nil {
			log.Error("Unable to load resource details", err)
			return err
		}

		if err = resourceLoadProperties(tx, resource, workspace, resourceName); err != nil {
			log.Error("Unable to load resource properties", err)
			return err
		}

		if err = resourceLoadReferences(tx, resource, workspace, resourceName); err != nil {
			log.Error("Unable to load resource properties", err)
			return err
		}

		return nil
	})

	if err != nil {
		log.Error("Unable load resource", err)
		return nil, err
	}

	return resource, nil
}

func (p *postgresResourceServiceBackend) AddResource(params backend.AddResourceParams) (*model.Resource, errors.ServiceError) {
	if params.Resource.Workspace == "" {
		params.Resource.Workspace = "default"
	}

	err := p.withSystemBackend(false, func(tx *sql.Tx) errors.ServiceError {
		// check if resource exists

		if existingCount, err := resourceCountsByName(tx, params.Resource.Workspace, params.Resource.Name); err != nil {
			return err
		} else if existingCount > 0 {
			if params.IgnoreIfExists {
				return nil
			}
			return errors.AlreadyExistsError
		}

		if err := resourceInsert(tx, params.Resource); err != nil {
			log.Error("Unable to insert resource", err)
			return err
		}

		if err := resourceUpsertProperties(tx, params.Resource); err != nil {
			log.Error("Unable to insert resource properties", err)
			return err
		}

		if err := resourceUpsertReferences(tx, params.Resource); err != nil {
			log.Error("Unable to insert resource properties", err)
			return err
		}

		if params.Migrate {
			err := p.withBackend(params.Resource.SourceConfig.DataSource, false, func(tx *sql.Tx) errors.ServiceError {
				err := resourceCreateTable(tx, params.Resource)
				if err != nil {
					return err
				}

				for _, reference := range params.Resource.References {
					if err = resourceCreateForeignKey(tx, params.Resource, reference); err != nil {
						if err != nil {
							return err
						}
					}
				}

				if params.Resource.Flags.KeepHistory {
					if params.Resource.Flags.DisableAudit {
						return errors.LogicalError.WithMessage("history cannot created while audit is disabled")
					}
					err = resourceCreateHistoryTable(tx, params.Resource)
					if err != nil {
						return err
					}
				}
				return nil
			})

			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return params.Resource, nil
}

func (p *postgresResourceServiceBackend) UpdateResource(ctx context.Context, resource *model.Resource, doMigration bool, forceMigration bool) errors.ServiceError {
	var err errors.ServiceError
	return p.withSystemBackend(false, func(tx *sql.Tx) errors.ServiceError {
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

		if doMigration {
			err = p.withBackend(resource.SourceConfig.DataSource, false, func(tx *sql.Tx) errors.ServiceError {
				if err = resourceCreateTable(tx, resource); err != nil {
					return err
				}

				if err = resourceMigrateTable(ctx, tx, resource, forceMigration, false); err != nil {
					return err
				}

				if resource.Flags.KeepHistory {
					if err = resourceCreateHistoryTable(tx, resource); err != nil {
						return err
					}

					if err = resourceMigrateTable(ctx, tx, resource, forceMigration, true); err != nil {
						return err
					}
				}

				return nil
			})

			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (p *postgresResourceServiceBackend) DeleteResources(ctx context.Context, workspace string, ids []string, doMigration bool, forceMigration bool) errors.ServiceError {
	var sources []*model.ResourceSourceConfig

	if workspace == "" {
		workspace = "default"
	}

	for _, id := range ids {
		resource, err := p.GetResourceByName(nil, workspace, id)

		if err != nil {
			return err
		}

		if resource == nil {
			return errors.NotFoundError
		}

		sources = append(sources, resource.SourceConfig)
	}
	err := p.withSystemBackend(false, func(tx *sql.Tx) errors.ServiceError {
		return resourceDelete(ctx, tx, ids)
	})

	if err != nil {
		return err
	}

	if doMigration {
		for _, source := range sources {
			err = p.withBackend(source.DataSource, false, func(tx *sql.Tx) errors.ServiceError {
				return resourceDropTable(tx, source.Mapping)
			})

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *postgresResourceServiceBackend) acquireConnection(dataSourceId string) (*sql.DB, errors.ServiceError) {
	if p.connectionMap[dataSourceId] == nil {
		dsBck, err := p.dataSourceService.GetDataSourceBackendById(dataSourceId)

		if err != nil {
			return nil, err
		}

		bck := dsBck.(*postgresDataSourceBackend)

		connStr := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable", bck.Options.Username, bck.Options.Password, bck.Options.Host, bck.Options.Port, bck.Options.DbName)
		// Connect to database
		conn, sqlErr := sql.Open("postgres", connStr)
		err = handleDbError(sqlErr)

		if err != nil {
			return nil, err
		}

		p.connectionMap[dataSourceId] = conn

		log.Info("Connected to Datasource: ", dataSourceId)
	}

	return p.connectionMap[dataSourceId], nil
}

func NewPostgresResourceServiceBackend() backend.ResourceServiceBackend {
	return &postgresResourceServiceBackend{
		connectionMap: make(map[string]*sql.DB),
	}
}

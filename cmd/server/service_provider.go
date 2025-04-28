package main

import (
	"context"
	"log"

	configDef "github.com/Ghaarp/auth/internal/config"
	repositoryDef "github.com/Ghaarp/auth/internal/repository"

	pgRepository "github.com/Ghaarp/auth/internal/repository/auth"
	pgRepositoryConverter "github.com/Ghaarp/auth/internal/repository/auth/converter"

	serviceDef "github.com/Ghaarp/auth/internal/service"
	authService "github.com/Ghaarp/auth/internal/service/auth"
)

type serviceProvider struct {
	dbConfig     configDef.DBConfig
	serverConfig configDef.AuthConfig

	repository          repositoryDef.AuthRepository
	repositoryConverter repositoryDef.RepoConverter

	service serviceDef.AuthService
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (sp *serviceProvider) DBConfig() configDef.DBConfig {
	if sp.dbConfig == nil {
		cfg, err := configDef.NewDBConfig()
		if err != nil {
			panic(err)
		}
		sp.dbConfig = cfg
	}
	return sp.dbConfig
}

func (sp *serviceProvider) ServerConfig(configPath string) configDef.AuthConfig {
	if sp.serverConfig == nil {
		err := configDef.Load(configPath)
		if err != nil {
			log.Print("Unable to load .env")
		}

		cfg, err := configDef.NewAuthConfig()
		if err != nil {
			panic(err)
		}
		sp.serverConfig = cfg
	}
	return sp.serverConfig
}

func (sp *serviceProvider) RepositoryConverter() repositoryDef.RepoConverter {
	if sp.repositoryConverter == nil {
		sp.repositoryConverter = pgRepositoryConverter.CreateConverter()
	}

	return sp.repositoryConverter
}

func (sp *serviceProvider) Repository(ctx context.Context) repositoryDef.AuthRepository {
	if sp.repository == nil {
		var err error
		sp.repository, err = pgRepository.CreateRepository(ctx, sp.DBConfig().DSN())
		if err != nil {
			log.Fatal(err)
		}
	}
	return sp.repository
}

func (sp *serviceProvider) Service(ctx context.Context) serviceDef.AuthService {
	if sp.service == nil {
		sp.service = authService.CreateService(sp.Repository(ctx), sp.RepositoryConverter())
	}
	return sp.service
}

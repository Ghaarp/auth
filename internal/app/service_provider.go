package app

import (
	"context"
	"flag"
	"log"

	configDef "github.com/Ghaarp/auth/internal/config"
	repositoryDef "github.com/Ghaarp/auth/internal/repository"

	pgRepository "github.com/Ghaarp/auth/internal/repository/auth"
	pgRepositoryConverter "github.com/Ghaarp/auth/internal/repository/auth/converter"

	serviceDef "github.com/Ghaarp/auth/internal/service"
	authService "github.com/Ghaarp/auth/internal/service/auth"

	authImplementation "github.com/Ghaarp/auth/internal/api/auth"
)

type serviceProvider struct {
	dbConfig   configDef.DBConfig
	grpcConfig configDef.Config
	httpConfig configDef.Config

	repository          repositoryDef.AuthRepository
	repositoryConverter repositoryDef.RepoConverter

	service  serviceDef.AuthService
	authImpl *authImplementation.AuthImplementation
}

func newServiceProvider() *serviceProvider {
	var configPath string
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")

	err := configDef.Load(configPath)
	if err != nil {
		log.Print("Unable to load .env")
	}
	return &serviceProvider{}
}

func (sp *serviceProvider) AuthImplementation(ctx context.Context) *authImplementation.AuthImplementation {
	if sp.authImpl == nil {
		sp.authImpl = authImplementation.NewAuthImplementation(sp.Service(ctx))
	}
	return sp.authImpl
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

func (sp *serviceProvider) GRPCConfig() configDef.Config {
	if sp.grpcConfig == nil {
		cfg, err := configDef.NewAuthConfig()
		if err != nil {
			panic(err)
		}
		sp.grpcConfig = cfg
	}
	return sp.grpcConfig
}

func (sp *serviceProvider) HttpConfig() configDef.Config {
	if sp.httpConfig == nil {
		cfg, err := configDef.NewHttpConfig()
		if err != nil {
			panic(err)
		}
		sp.httpConfig = cfg
	}
	return sp.httpConfig
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

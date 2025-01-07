package extern

import (
	"io"

	"github.com/AAguilar0x0/txapp/core/services"
	"github.com/AAguilar0x0/txapp/extern/db/psql"
	"github.com/AAguilar0x0/txapp/extern/env"
	"github.com/AAguilar0x0/txapp/extern/hash/chash"
	"github.com/AAguilar0x0/txapp/extern/idgen/ksuid"
	"github.com/AAguilar0x0/txapp/extern/jwt/golangjwt"
	"github.com/AAguilar0x0/txapp/extern/validator/validatorv10"
)

type DefaultServiceProvider struct {
	environment services.Environment
	database    services.DatabaseManager
	validator   services.Validator
	jwt         services.JWTokenizer
	hash        services.Hash
	idGenerator services.IDGenerator
	closable    []io.Closer
}

func New() *DefaultServiceProvider {
	return &DefaultServiceProvider{
		environment: env.New(),
	}
}

func (d *DefaultServiceProvider) cleanup(c io.Closer) {
	d.closable = append(d.closable, c)
}

func (d *DefaultServiceProvider) Environment() (services.Environment, error) {
	if d.environment == nil {
		d.environment = env.New()
	}
	return d.environment, nil
}

func (d *DefaultServiceProvider) databaseManager() (services.DatabaseManager, error) {
	if d.database == nil {
		idGen, err := d.IDGenerator()
		if err != nil {
			return nil, err
		}
		data, err := psql.New(d.environment, idGen)
		if err != nil {
			return nil, err
		}
		d.database = data
		d.cleanup(data)
	}
	return d.database, nil
}

func (d *DefaultServiceProvider) Database() (services.Database, error) {
	return d.databaseManager()
}

func (d *DefaultServiceProvider) Migrator() (services.Migrator, error) {
	return d.databaseManager()
}

func (d *DefaultServiceProvider) Validator() (services.Validator, error) {
	if d.validator == nil {
		data, err := validatorv10.New(d.environment)
		if err != nil {
			return nil, err
		}
		d.validator = data
		d.cleanup(data)
	}
	return d.validator, nil
}

func (d *DefaultServiceProvider) JWTokenizer() (services.JWTokenizer, error) {
	if d.jwt == nil {
		idGen, err := d.IDGenerator()
		if err != nil {
			return nil, err
		}
		data, err := golangjwt.New(d.environment, idGen)
		if err != nil {
			return nil, err
		}
		// d.jwt = data
		d.cleanup(data)
	}
	return d.jwt, nil
}

func (d *DefaultServiceProvider) Hash() (services.Hash, error) {
	if d.hash == nil {
		data, err := chash.New()
		if err != nil {
			return nil, err
		}
		d.hash = data
	}
	return d.hash, nil
}

func (d *DefaultServiceProvider) IDGenerator() (services.IDGenerator, error) {
	if d.idGenerator == nil {
		data, err := ksuid.New()
		if err != nil {
			return nil, err
		}
		d.idGenerator = data
		d.cleanup(data)
	}
	return d.idGenerator, nil
}

func (d *DefaultServiceProvider) Close() error {
	for i := len(d.closable) - 1; i >= 0; i-- {
		if err := d.closable[i].Close(); err != nil {
			return err
		}
	}
	return nil
}

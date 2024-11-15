package extern

import (
	"io"

	"github.com/AAguilar0x0/txapp/core/models"
	"github.com/AAguilar0x0/txapp/core/services"
	authcustom "github.com/AAguilar0x0/txapp/extern/auth/custom"
	"github.com/AAguilar0x0/txapp/extern/db/psql"
	"github.com/AAguilar0x0/txapp/extern/env"
	"github.com/AAguilar0x0/txapp/extern/hash/chash"
	"github.com/AAguilar0x0/txapp/extern/idgen/ksuid"
	"github.com/AAguilar0x0/txapp/extern/validator/validatorv10"
)

type DefaultServiceProvider struct {
	environment   services.Environment
	database      models.DatabaseManager
	validator     services.Validator
	authenticator services.Authenticator
	hash          services.Hash
	idGenerator   services.IDGenerator
	closable      []io.Closer
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

func (d *DefaultServiceProvider) databaseManager() (models.DatabaseManager, error) {
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

func (d *DefaultServiceProvider) Database() (models.Database, error) {
	return d.databaseManager()
}

func (d *DefaultServiceProvider) Migrator() (models.Migrator, error) {
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

func (d *DefaultServiceProvider) Authenticator() (services.Authenticator, error) {
	if d.authenticator == nil {
		data, err := authcustom.New(d.environment)
		if err != nil {
			return nil, err
		}
		d.authenticator = data
		d.cleanup(data)
	}
	return d.authenticator, nil
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

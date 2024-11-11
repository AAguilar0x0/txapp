package extern

import (
	"io"

	"github.com/AAguilar0x0/txapp/core/models"
	"github.com/AAguilar0x0/txapp/core/services"
	authcustom "github.com/AAguilar0x0/txapp/extern/auth/custom"
	"github.com/AAguilar0x0/txapp/extern/db/psql"
	"github.com/AAguilar0x0/txapp/extern/env"
	"github.com/AAguilar0x0/txapp/extern/idgen/ksuid"
	"github.com/AAguilar0x0/txapp/extern/validator/validatorv10"
)

type Services struct {
	environment   services.Environment
	database      models.Database
	validator     services.Validator
	authenticator services.Authenticator
	idGenerator   services.IDGenerator
	closable      []io.Closer
}

func New() *Services {
	return &Services{
		environment: env.New(),
	}
}

func (d *Services) cleanup(c io.Closer) {
	d.closable = append(d.closable, c)
}

func (d *Services) Environment() (services.Environment, error) {
	if d.environment == nil {
		d.environment = env.New()
	}
	return d.environment, nil
}

func (d *Services) Database() (models.Database, error) {
	if d.database == nil {
		data, err := psql.New(d.environment)
		if err != nil {
			return nil, err
		}
		d.database = data
		d.cleanup(data)
	}
	return d.database, nil
}

func (d *Services) Validator() (services.Validator, error) {
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

func (d *Services) Authenticator() (services.Authenticator, error) {
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

func (d *Services) IDGenerator() (services.IDGenerator, error) {
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

func (d *Services) Close() error {
	for i := len(d.closable) - 1; i >= 0; i-- {
		if err := d.closable[i].Close(); err != nil {
			return err
		}
	}
	return nil
}

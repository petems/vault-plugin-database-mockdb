package mockdb

import (
	"context"
	"fmt"
	"sync"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/database/dbplugin"
	"github.com/hashicorp/vault/sdk/database/helper/connutil"
	"github.com/hashicorp/vault/sdk/database/helper/dbutil"

	"github.com/mitchellh/mapstructure"
)

// mockdbConnectionProducer implements ConnectionProducer and provides an
// interface for mockdb databases to make connections.
type mockdbConnectionProducer struct {
	Host     string `json:"host" structs:"host" mapstructure:"host"`
	Username string `json:"username" structs:"username" mapstructure:"username"`
	Password string `json:"password" structs:"password" mapstructure:"password"`
	Port     string `json:"port" structs:"port" mapstructure:"port"` //default to 8086

	rawConfig map[string]interface{}

	Initialized bool
	Type        string
	client      string

	sync.Mutex
}

func (i *mockdbConnectionProducer) Initialize(ctx context.Context, conf map[string]interface{}, verifyConnection bool) error {
	_, err := i.Init(ctx, conf, verifyConnection)
	return err
}

func (i *mockdbConnectionProducer) Init(ctx context.Context, conf map[string]interface{}, verifyConnection bool) (map[string]interface{}, error) {
	i.Lock()
	defer i.Unlock()

	i.rawConfig = conf

	err := mapstructure.WeakDecode(conf, i)
	if err != nil {
		return nil, err
	}

	if i.Port == "" {
		i.Port = "8086"
	}

	switch {
	case len(i.Host) == 0:
		return nil, fmt.Errorf("host cannot be empty")
	case len(i.Username) == 0:
		return nil, fmt.Errorf("username cannot be empty")
	case len(i.Password) == 0:
		return nil, fmt.Errorf("password cannot be empty")
	}

	// Set initialized to true at this point since all fields are set,
	// and the connection can be established at a later time.
	i.Initialized = true

	if verifyConnection {
		if _, err := i.Connection(ctx); err != nil {
			return nil, errwrap.Wrapf("error verifying connection: {{err}}", err)
		}
	}

	return conf, nil
}

func (i *mockdbConnectionProducer) Connection(_ context.Context) (interface{}, error) {
	if !i.Initialized {
		return nil, connutil.ErrNotInitialized
	}

	return nil, nil
}

func (i *mockdbConnectionProducer) Close() error {
	// Grab the write lock
	i.Lock()
	defer i.Unlock()

	return nil
}

func (i *mockdbConnectionProducer) secretValues() map[string]interface{} {
	return map[string]interface{}{
		i.Password: "[password]",
	}
}

// SetCredentials uses provided information to set/create a user in the
// database. Unlike CreateUser, this method requires a username be provided and
// uses the name given, instead of generating a name. This is used for creating
// and setting the password of static accounts, as well as rolling back
// passwords in the database in the event an updated database fails to save in
// Vault's storage.
func (i *mockdbConnectionProducer) SetCredentials(ctx context.Context, statements dbplugin.Statements, staticUser dbplugin.StaticUserConfig) (username, password string, err error) {
	return "", "", dbutil.Unimplemented()
}

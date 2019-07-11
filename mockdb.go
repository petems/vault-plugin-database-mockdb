package mockdb

import (
	"context"
	"database/sql"
	"errors"

	// "strings"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/database/dbplugin"
	"github.com/hashicorp/vault/sdk/database/helper/credsutil"
	"github.com/hashicorp/vault/sdk/database/helper/dbutil"
	// "github.com/hashicorp/vault/sdk/helper/dbtxn"
	// "github.com/hashicorp/vault/sdk/helper/strutil"
)

const (
	defaultUserCreationMQL           = `CREATE USER "{{username}}" WITH PASSWORD '{{password}}';`
	defaultUserDeletionMQL           = `DROP USER "{{username}}";`
	defaultRootCredentialRotationMQL = `SET PASSWORD FOR "{{username}}" = '{{password}}';`
	mockdbTypeName                   = "mockdb"
)

type Mockdb struct {
	*mockdbConnectionProducer
	credsutil.CredentialsProducer
}

// New implements builtinplugins.BuiltinFactory
func New() (interface{}, error) {
	db := new()
	// Wrap the plugin with middleware to sanitize errors
	dbType := dbplugin.NewDatabaseErrorSanitizerMiddleware(db, db.secretValues)
	return dbType, nil
}

func new() *Mockdb {
	connProducer := &mockdbConnectionProducer{}
	connProducer.Type = mockdbTypeName

	credsProducer := &credsutil.SQLCredentialsProducer{
		DisplayNameLen: 15,
		RoleNameLen:    15,
		UsernameLen:    100,
		Separator:      "_",
	}

	return &Mockdb{
		mockdbConnectionProducer: connProducer,
		CredentialsProducer:      credsProducer,
	}
}

// Run instantiates an Mockdb object, and runs the RPC server for the plugin
func Run(apiTLSConfig *api.TLSConfig) error {
	dbType, err := New()
	if err != nil {
		return err
	}

	dbplugin.Serve(dbType.(dbplugin.Database), api.VaultPluginTLSProvider(apiTLSConfig))

	return nil
}

func (m *Mockdb) Type() (string, error) {
	return mockdbTypeName, nil
}

func (m *Mockdb) CreateUser(ctx context.Context, statements dbplugin.Statements, usernameConfig dbplugin.UsernameConfig, expiration time.Time) (username string, password string, err error) {
	statements = dbutil.StatementCompatibilityHelper(statements)

	if len(statements.Creation) == 0 {
		return "", "", dbutil.ErrEmptyCreationStatement
	}

	// Grab the lock
	m.Lock()
	defer m.Unlock()

	username, err = m.GenerateUsername(usernameConfig)
	if err != nil {
		return "", "", err
	}

	password, err = m.GeneratePassword()
	if err != nil {
		return "", "", err
	}

	// expirationStr, err := m.GenerateExpiration(expiration)
	// if err != nil {
	// 	return "", "", err
	// }

	// Get the connection
	// db, err := m.getConnection(ctx)
	// if err != nil {
	// 	return "", "", err

	// }

	// // Start a transaction
	// tx, err := db.Begin()
	// if err != nil {
	// 	return "", "", err

	// }
	// defer func() {
	// 	tx.Rollback()
	// }()

	// // Execute each query
	// for _, stmt := range statements.Creation {
	// 	for _, query := range strutil.ParseArbitraryStringSlice(stmt, ";") {
	// 		query = strings.TrimSpace(query)
	// 		if len(query) == 0 {
	// 			continue
	// 		}

	// 		m := map[string]string{
	// 			"name":       username,
	// 			"password":   password,
	// 			"expiration": expirationStr,
	// 		}

	// 		if err := dbtxn.ExecuteTxQuery(ctx, tx, m, query); err != nil {
	// 			return "", "", err
	// 		}
	// 	}
	// }

	// // Commit the transaction
	// if err := tx.Commit(); err != nil {
	// 	return "", "", err

	// }

	// Return the secret
	return username, password, nil
}

func (m *Mockdb) RenewUser(ctx context.Context, statements dbplugin.Statements, username string, expiration time.Time) error {
	return nil // NOOP
}

func (m *Mockdb) RevokeUser(ctx context.Context, statements dbplugin.Statements, username string) error {
	// Grab the lock
	m.Lock()
	defer m.Unlock()

	// // Get the connection
	// db, err := m.getConnection(ctx)
	// if err != nil {
	// 	return err
	// }

	// tx, err := db.Begin()
	// if err != nil {
	// 	return err
	// }
	// defer func() {
	// 	tx.Rollback()
	// }()

	// if err := m.disconnectSession(db, username); err != nil {
	// 	return err
	// }

	// statements = dbutil.StatementCompatibilityHelper(statements)
	// revocationStatements := statements.Revocation
	// if len(revocationStatements) == 0 {
	// 	revocationStatements = []string{defaultRootCredentialRotationMQL}
	// }

	// // We can't use a transaction here, because Mockdb treats DROP USER as a DDL statement, which commits immediately.
	// // Execute each query
	// for _, stmt := range revocationStatements {
	// 	for _, query := range strutil.ParseArbitraryStringSlice(stmt, ";") {
	// 		query = strings.TrimSpace(query)
	// 		if len(query) == 0 {
	// 			continue
	// 		}

	// 		m := map[string]string{
	// 			"name": username,
	// 		}

	// 		if err := dbtxn.ExecuteTxQuery(ctx, tx, m, query); err != nil {
	// 			return err
	// 		}
	// 	}
	// }

	return nil
}

func (m *Mockdb) RotateRootCredentials(ctx context.Context, statements []string) (map[string]interface{}, error) {
	m.Lock()
	defer m.Unlock()

	if len(m.Username) == 0 || len(m.Password) == 0 {
		return nil, errors.New("username and password are required to rotate")
	}

	rotateStatements := statements
	if len(rotateStatements) == 0 {
		rotateStatements = []string{defaultRootCredentialRotationMQL}
	}

	// db, err := m.getConnection(ctx)
	// if err != nil {
	// 	return nil, err
	// }

	// tx, err := db.Begin()
	// if err != nil {
	// 	return nil, err
	// }
	// defer func() {
	// 	tx.Rollback()
	// }()

	password, err := m.GeneratePassword()
	if err != nil {
		return nil, err
	}

	// for _, stmt := range rotateStatements {
	// 	for _, query := range strutil.ParseArbitraryStringSlice(stmt, ";") {
	// 		query = strings.TrimSpace(query)
	// 		if len(query) == 0 {
	// 			continue
	// 		}

	// 		m := map[string]string{
	// 			"username": m.Username,
	// 			"password": password,
	// 		}

	// 		if err := dbtxn.ExecuteTxQuery(ctx, tx, m, query); err != nil {
	// 			return nil, err
	// 		}
	// 	}
	// }

	// if err := tx.Commit(); err != nil {
	// 	return nil, err
	// }

	// if err := db.Close(); err != nil {
	// 	return nil, err
	// }

	m.rawConfig["password"] = password
	return m.rawConfig, nil
}

func (m *Mockdb) disconnectSession(db *sql.DB, username string) error {
	return nil
}

func (m *Mockdb) getConnection(ctx context.Context) (*sql.DB, error) {
	db, err := m.Connection(ctx)
	if err != nil {
		return nil, err
	}

	return db.(*sql.DB), nil
}

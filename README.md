# vault-database-plugin-mockdb

A [Vault](https://www.vaultproject.io) plugin for "MockDB"... which is just a mocked database that doesn't actually do anything!

## Background

This is basically an example of how to write your own database secrets engine. It's largely based on the [vault-plugin-database-oracle repo](https://github.com/hashicorp/vault-plugin-database-oracle) repo, with some other code taken from the [InfluxDB secrets engine](https://github.com/hashicorp/vault/blob/master/plugins/database/influxdb/influxdb.go)

## Installation

The Vault plugin system is documented on the [Vault documentation site](https://www.vaultproject.io/docs/internals/plugins.html).

You will need to define a plugin directory using the `plugin_directory` configuration directive, then place the `mockdb` executable generated above in the directory.

Inside this repo, you can do this with the following dev steps:

```
$ vault server -dev -dev-root-token-id=root -dev-plugin-dir=./vault/plugins &
$ export VAULT_ADDR='http://127.0.0.1:8200'
$ vault secrets enable database
$ vault write sys/plugins/catalog/database/mockdb \
  sha_256=$(MOCKDBSHASUM) \
  command="mockdb"
$ vault write database/config/mockdb \
  plugin_name="mockdb" \
  host=127.0.0.1 \
  username=mockdb-root \
  password=password123 \
  allowed_roles=my-role
$ vault write database/roles/my-role \
  db_name=mockdb \
  creation_statements="CREATE USER \"{{username}}\" WITH PASSWORD '{{password}}'; \
       GRANT ALL ON \"vault\" TO \"{{username}}\";" \
  default_ttl="1h" \
  max_ttl="24h"
$ vault read database/creds/my-role
Key                Value
---                -----
lease_id           database/creds/my-role/tXLr2KV5zhpuoSCszMChgsoQ
lease_duration     1h
lease_renewable    true
password           A1a-vJ10mm1hx0uSAjDe
username           v_token_my-role_HMxsjZcsAH6KQRW4OStz_1562877167
```

Or use the makefile steps:

```
$ make start-vault
```

New Terminal

```
$ make dev-flow
go build -o vault/plugins/mockdb plugin/main.go
vault secrets disable database
Success! Disabled the secrets engine (if it existed) at: database/
vault secrets enable database
Success! Enabled the database secrets engine at: database/
vault write sys/plugins/catalog/database/mockdb \
    sha_256=92b14b650aee1e0719e12e3a7ba423ef0b6316e4a956c87c52bdcc38ff9118ac \
    command="mockdb"
Success! Data written to: sys/plugins/catalog/database/mockdb
vault write database/config/mockdb \
     plugin_name="mockdb" \
     host=127.0.0.1 \
     username=mockdb-root \
     password=password123 \
     allowed_roles=my-role
vault write database/roles/my-role \
  	db_name=mockdb \
  	creation_statements="CREATE USER \"{{username}}\" WITH PASSWORD '{{password}}'; \
  	     GRANT ALL ON \"vault\" TO \"{{username}}\";" \
  	default_ttl="1h" \
  	max_ttl="24h"
Success! Data written to: database/roles/my-role
vault read database/creds/my-role
Key                Value
---                -----
lease_id           database/creds/my-role/Re0cSMdAxC5TojqBTR28kAtm
lease_duration     1h
lease_renewable    true
password           A1a-LOPIN4vNs84ZW4G0
username           v_token_my-role_nai3vEL6UOK2pJXmNu8O_1563449205
```

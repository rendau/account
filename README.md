# Auth

### DB dump:

```
pg_dump --no-owner -Fc -U postgres account -f ./account.custom
```

### DB restore:

```
dropdb -U postgres account
createdb -U postgres account
pg_restore --no-owner -d account -U postgres ./account.custom
```

### Install `migrate` command-tool:

https://github.com/golang-migrate/migrate/tree/master/cmd/migrate

### Create new migration:

```
migrate create -ext sql -dir migrations mg_name
```

### Apply migration:

```
migrate -path migrations -database "postgres://localhost:5432/db_name?sslmode=disable" up
```


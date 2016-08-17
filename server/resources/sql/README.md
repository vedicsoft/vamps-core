## Database migration
Note : applied only for sql schemas. Internally we make use of [mattes/migrate](https://github.com/mattes/migrate)
library, and the migration code is in commons/database.go

## Introduce new schema changes to sql database

1. create a new file 00n_platform.up.sql where n is the latest file number
2. Add the schema changes to that
3. create a new file 00n_platform.down.sql where n is the latest file number.
   This will be use to rollback the migration
4. Update the platform.master.sql with all the changes

## How it works
This is a very simple migration tool which creates a table called `schema_migrations`
in the db and keep track of the migration file numbers.

Migration tool runs at every server startup and will apply any schema changes. Currently we support upgrades only

## platform.master.sql

This sql file contains the latest core database schema

## Migration files

The format of migration files looks like this:

```
001_platform.up.sql     # up migration instructions
001_platform.down.sql   # down migration instructions
002_xxx.up.sql
002_xxx.down.sql
...
```

Why two files? This way you could still do sth like
``psql -f ./db/migrations/001_platform.up.sql`` and there is no
need for any custom markup language to divide up and down migrations. Please note
that the filename extension depends on the driver.

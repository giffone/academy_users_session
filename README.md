# academy_users_session
manage user's pc session at domain.


### Add new migrations
```bash
migrate create -ext sql -dir db/migrations -seq table_name
```

### Run service
```bash
make run 'DATABASE_URL=postgres://user:password@host:port/db-name?sslmode=disable'
```

### Run service locally 
for testing and be able to connect to the local database
```bash
make run_local 'DATABASE_URL=postgres://user:password@host:port/db-name?sslmode=disable'
```
# academy_users_session
manage user's pc session at domain.


add new migrations
```bash
migrate create -ext sql -dir db/migrations -seq table_name
```

run service
```bash
make docker DATABASE_URL=postgres://user:password@host:port/db-name?sslmode=disable
```
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
### APIs

#### Add new users
character varying(50)
```http
POST http://localhost:8080/api/session-manager/users
Content-Type: application/json
[
    {
        "Name": "user_1"
    },
    {
        "Name": "user_2"
    },
    {
        "Name": "user_3"
    }
    // ...
]
```
#### Add new computers
character varying(30)
```http
POST http://localhost:8080/api/session-manager/computers
Content-Type: application/json
[
    {
        "Name": "academie-mac-pink0001"
    },
    {
        "Name": "academie-mac-blue0002"
    },
    {
        "Name": "academie-mac-red0003"
    }
    // ...
]
```
#### Add new session
The computer notifies the running script about the start of a session during user authorization
```http
POST http://localhost:8080/api/session-manager/session
Content-Type: application/json
{
  "id": "5f2c9d6c-2a84-4d63-b64c-6a0f12eb3471",
  "comp_name": "academie-mac-pink0001",
  "ip_addr": "192.168.1.100",
  "login": "user_1",
  "next_ping_sec": 60,
  "date_time": "2023-09-06T12:30:00Z" // current time from pc
}

```
#### Add new activity (after creating a session)
just calculate the session time
```http
POST http://localhost:8080/api/session-manager/activity
Content-Type: application/json
{
  "session_id": "5f2c9d6c-2a84-4d63-b64c-6a0f12eb3471",
  "session_type": "", // event name empty
  "login": "user_1",
  "next_ping_sec": 60,
  "date_time": "2023-09-06T15:30:00Z" // current time from pc
}
```
or if you want to calculate the session time for individual events
```http
POST http://localhost:8080/api/session-manager/activity
Content-Type: application/json
{
  "session_id": "5f2c9d6c-2a84-4d63-b64c-6a0f12eb3471",
  "session_type": "platform zero", // event name
  "login": "user_1",
  "next_ping_sec": 60,
  "date_time": "2023-09-06T15:30:00Z" // current time from pc
}
```
the last notification (session or activity) sent from the computer will mean the end of the session ("date_time" + "next_ping_sec").
#### Get online sessions
```http
GET http://localhost:8080/api/session-manager/dashboard
```
response:
```json
// Content-Type: application/json
{
  "message": "Success",
  "data": [
    {
      "id": "5f2c9d6c-2a84-4d63-b64c-6a0f12eb3471",
      "comp_name": "academie-mac-pink0001",
      "ip_addr": "192.168.1.100",
      "login": "user_1",
      "start_date_time": "2023-09-06T08:00:00Z",
      "end_date_time": "2023-09-06T09:30:00Z" // ends at: last ping + "next_ping_sec"
    },
    {
      "id": "6b8a4f1a-1d09-4d08-8a03-86be3e3b9104",
      "comp_name": "academie-mac-blue0002",
      "ip_addr": "192.168.1.101",
      "login": "user_2",
      "start_date_time": "2023-09-06T10:15:00Z",
      "end_date_time": "2023-09-06T11:45:00Z"
    }
    // ...
  ]
}
```
#### Get user activity
query param
- `session_type` - ***"platform zero"*** or your event
- `login` - ***"user_1"***
- `from_date` - ***"2022-09-01T00:00:00Z"***
- `to_date` - ***"2022-12-31T00:00:00Z"***
- `group_by` - ***"month"*** or ***"date"***
```http
GET http://localhost:8080/api/session-manager/activity?session_type=xxx&login=xxx&from_date=xxx&to_date=xxx&group_by=xxx
```
response - group by month:
```json
// Content-Type: application/json
{
	"message": "Success",
	"data": {
		"id": "user_1",
		"total_hours": 436.63,
		"user_activity": [
			{
				"year": "2023",
				"month_num": "9",
				"hours": 90.616667
			},
			{
				"year": "2023",
				"month_num": "5",
				"hours": 79.016666
			}
            // ...
		]
	}
}
```
response - group by date:
```json
// Content-Type: application/json
{
    "message": "Success",
    "data": {
        "login": "user_1",
        "total_hours": 37.62,
        "user_activity": [
            {
                "date": "2023-01-01T00:00:00Z",
                "hours": 13.266666
            },
            {
                "date": "2023-01-02T00:00:00Z",
                "hours": 7.35
            },
            // ...
        ]
    }
}
```
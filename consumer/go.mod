module github.com/gbouv/queue-performance/consumer

go 1.19

// replace github.com/gbouv/queue-performance/queue => ../queue

require (
	github.com/gbouv/queue-performance/queue v0.0.0-20230808212659-5fdb43b3572b
	github.com/redis/go-redis/v9 v9.0.5
	github.com/sirupsen/logrus v1.9.3
)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.3.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/palantir/stacktrace v0.0.0-20161112013806-78658fd2d177 // indirect
	golang.org/x/crypto v0.8.0 // indirect
	golang.org/x/sys v0.10.0 // indirect
	golang.org/x/text v0.11.0 // indirect
	gorm.io/driver/postgres v1.5.2 // indirect
	gorm.io/gorm v1.25.2 // indirect
)

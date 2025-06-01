package pg

import (
	"database/sql"
	"log"
	"os"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"
)

type PostgresConn struct {
	Conn *sql.DB
	Qb   sq.StatementBuilderType
}

func SetupDB() PostgresConn {
	connStr := os.Getenv("POSTGRES_DSN")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	pg := PostgresConn{
		Conn: db,
		Qb:   sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}

	// migrator()
	return pg
}

// func migrator() {
// 	goose.SetLogger(goose.NopLogger())
//
// 	if len(os.Args) < 2 {
// 		log.Fatal("missing goose command")
// 	}
//
// 	dbstring := os.Getenv("GOOSE_DBSTRING")
// 	db, err := goose.OpenDBWithDriver("postgres", dbstring)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	if err := goose.RunContext(context.Background(), os.Args[1], db, os.Getenv("GOOSE_MIGRATION_DIR")); err != nil {
// 		log.Fatal(err)
// 	}
// }

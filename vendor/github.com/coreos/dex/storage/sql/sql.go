// Package sql provides SQL implementations of the storage interface.
package sql

import (
	"database/sql"
	"fmt"
	"regexp"
	"time"

	"github.com/sirupsen/logrus"

	// import third party drivers
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	// _ "github.com/mattn/go-sqlite3"
)

// flavor represents a specific SQL implementation, and is used to translate query strings
// between different drivers. Flavors shouldn't aim to translate all possible SQL statements,
// only the specific queries used by the SQL storages.
type flavor struct {
	queryReplacers []replacer

	// Optional function to create and finish a transaction. This is mainly for
	// cockroachdb support which requires special retry logic provided by their
	// client package.
	//
	// This will be nil for most flavors.
	//
	// See: https://github.com/cockroachdb/docs/blob/63761c2e/_includes/app/txn-sample.go#L41-L44
	executeTx func(db *sql.DB, fn func(*sql.Tx) error) error

	// Does the flavor support timezones?
	supportsTimezones bool
}

// A regexp with a replacement string.
type replacer struct {
	re   *regexp.Regexp
	with string
}

// Match a postgres query binds. E.g. "$1", "$12", etc.
var bindRegexp = regexp.MustCompile(`\$\d+`)

func matchLiteral(s string) *regexp.Regexp {
	return regexp.MustCompile(`\b` + regexp.QuoteMeta(s) + `\b`)
}

var (
	// The "github.com/lib/pq" driver is the default flavor. All others are
	// translations of this.
	flavorPostgres = flavor{
		// The default behavior for Postgres transactions is consistent reads, not consistent writes.
		// For each transaction opened, ensure it has the correct isolation level.
		//
		// See: https://www.postgresql.org/docs/9.3/static/sql-set-transaction.html
		//
		// NOTE(ericchiang): For some reason using `SET SESSION CHARACTERISTICS AS TRANSACTION` at a
		// session level didn't work for some edge cases. Might be something worth exploring.
		executeTx: func(db *sql.DB, fn func(sqlTx *sql.Tx) error) error {
			tx, err := db.Begin()
			if err != nil {
				return err
			}
			defer tx.Rollback()

			if _, err := tx.Exec(`SET TRANSACTION ISOLATION LEVEL SERIALIZABLE;`); err != nil {
				return err
			}
			if err := fn(tx); err != nil {
				return err
			}
			return tx.Commit()
		},

		supportsTimezones: true,
	}

	flavorSQLite3 = flavor{
		queryReplacers: []replacer{
			{bindRegexp, "?"},
			// Translate for booleans to integers.
			{matchLiteral("true"), "1"},
			{matchLiteral("false"), "0"},
			{matchLiteral("boolean"), "integer"},
			// Translate other types.
			{matchLiteral("bytea"), "blob"},
			{matchLiteral("timestamptz"), "timestamp"},
			// SQLite doesn't have a "now()" method, replace with "date('now')"
			{regexp.MustCompile(`\bnow\(\)`), "date('now')"},
		},
	}

	// Incomplete.
	flavorMySQL = flavor{
		queryReplacers: []replacer{
			{bindRegexp, "?"},
			{matchLiteral("timestamptz"), "datetime(6)"},
			{matchLiteral("bytea"), "blob"},
			{matchLiteral("text"), "varchar(555)"},
			{matchLiteral("keys"), "dex_keys"},
			//{regexp.MustCompile(`\b(keys)\b`), "`$1`"},
			// Change default timestamp to fit datetime.
			{regexp.MustCompile(`0001-01-01 00:00:00 UTC`), "1000-01-01 00:00:00"},
		},

		executeTx: func(db *sql.DB, fn func(sqlTx *sql.Tx) error) error {
			/*
				if _, err := db.Exec(`SET GLOBAL connect_timeout=100;;`); err != nil {
					return err
				}

				tx, err := db.BeginTx(context.Background(), &sql.TxOptions{
					Isolation: sql.LevelSerializable,
				})*/
			/**/

			tx, err := db.Begin()
			if err != nil {
				return err
			}

			defer tx.Rollback()

			//fmt.Println("Wait forrrrrr")
			/*
				if _, err := tx.Exec(`SET TRANSACTION ISOLATION LEVEL SERIALIZABLE;`); err != nil {
					return err
				}
			*/

			if err := fn(tx); err != nil {
				//Error: 1213 SQLSTATE: 40001 (ER_LOCK_DEADLOCK)
				//Message: Deadlock found when trying to get lock; try restarting transaction
				fmt.Println(err.Error())
				return err
			}

			return tx.Commit()
		},
		supportsTimezones: true,
	}

	// Not tested.
	// flavorCockroach = flavor{
	// 	executeTx: crdb.ExecuteTx,
	// }
)

func (f flavor) translate(query string) string {
	// TODO(ericchiang): Heavy cashing.
	for _, r := range f.queryReplacers {
		query = r.re.ReplaceAllString(query, r.with)
	}
	return query
}

// translateArgs translates query parameters that may be unique to
// a specific SQL flavor. For example, standardizing "time.Time"
// types to UTC for clients that don't provide timezone support.
func (c *conn) translateArgs(args []interface{}) []interface{} {
	if c.flavor.supportsTimezones {
		return args
	}

	for i, arg := range args {
		if t, ok := arg.(time.Time); ok {
			args[i] = t.UTC()
		}
	}
	return args
}

// conn is the main database connection.
type conn struct {
	db                 *sql.DB
	flavor             flavor
	logger             logrus.FieldLogger
	alreadyExistsCheck func(err error) bool
}

func (c *conn) Close() error {
	return c.db.Close()
}

// conn implements the same method signatures as encoding/sql.DB.

func (c *conn) Exec(query string, args ...interface{}) (sql.Result, error) {
	query = c.flavor.translate(query)
	return c.db.Exec(query, c.translateArgs(args)...)
}

func (c *conn) Query(query string, args ...interface{}) (*sql.Rows, error) {
	query = c.flavor.translate(query)
	return c.db.Query(query, c.translateArgs(args)...)
}

func (c *conn) QueryRow(query string, args ...interface{}) *sql.Row {
	query = c.flavor.translate(query)
	return c.db.QueryRow(query, c.translateArgs(args)...)
}

// ExecTx runs a method which operates on a transaction.
func (c *conn) ExecTx(fn func(tx *trans) error) error {
	if c.flavor.executeTx != nil {
		return c.flavor.executeTx(c.db, func(sqlTx *sql.Tx) error {
			return fn(&trans{sqlTx, c})
		})
	}

	sqlTx, err := c.db.Begin()
	if err != nil {
		return err
	}
	if err := fn(&trans{sqlTx, c}); err != nil {
		sqlTx.Rollback()
		return err
	}
	return sqlTx.Commit()
}

type trans struct {
	tx *sql.Tx
	c  *conn
}

// trans implements the same method signatures as encoding/sql.Tx.

func (t *trans) Exec(query string, args ...interface{}) (sql.Result, error) {
	query = t.c.flavor.translate(query)
	return t.tx.Exec(query, t.c.translateArgs(args)...)
}

func (t *trans) Query(query string, args ...interface{}) (*sql.Rows, error) {
	query = t.c.flavor.translate(query)
	return t.tx.Query(query, t.c.translateArgs(args)...)
}

func (t *trans) QueryRow(query string, args ...interface{}) *sql.Row {
	query = t.c.flavor.translate(query)
	return t.tx.QueryRow(query, t.c.translateArgs(args)...)
}

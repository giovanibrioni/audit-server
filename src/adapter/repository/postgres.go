package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"strconv"

	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/giovanibrioni/audit-server/audit"
	"github.com/giovanibrioni/audit-server/helper"
)

var (
	PG_URL      = helper.GetEnvOrDefault("POSTGRES_URL", "postgresql://localhost:5432")
	DB_MAX_CONN = helper.GetEnvOrDefault("DB_MAX_CONN", "10")
)

type postgresAuditRepository struct {
	conn   *sql.DB
	ctx    context.Context
	logger *zap.SugaredLogger
}

func NewPostgresAuditRepository(ctx context.Context, logger *zap.SugaredLogger) audit.AuditRepo {
	db := postgresConnect()
	createDatabase(db)
	dbMaxConn, _ := strconv.Atoi(DB_MAX_CONN)
	db.SetMaxOpenConns(dbMaxConn)
	db.SetMaxIdleConns(dbMaxConn)

	return &postgresAuditRepository{
		conn:   db,
		ctx:    ctx,
		logger: logger,
	}
}

func (p *postgresAuditRepository) SaveBatch(auditLogs []*audit.AuditEntity) error {
	for _, auditLog := range auditLogs {
		encoded, err := json.Marshal(auditLog.RawMessage)
		if err != nil {
			p.logger.Fatal("Unable to marshal auditLogs")
			return err
		}
		stmt := `INSERT INTO AuditLogs VALUES ($1, $2, $3)`
		_, err = p.conn.Exec(stmt, auditLog.AuditId, auditLog.JobId, encoded)
		checkIfError(err)
	}
	return nil
}

func postgresConnect() *sql.DB {
	db, err := sql.Open("postgres", PG_URL)
	checkIfError(err)

	return db
}

func checkIfError(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}

func createDatabase(db *sql.DB) {
	stmt := `CREATE TABLE IF NOT EXISTS AuditLogs (
						audit_id uuid NOT NULL PRIMARY KEY,
						job_id uuid NOT NULL,
						raw_info json NOT NULL)`

	_, err := db.Exec(stmt)
	checkIfError(err)

	log.Println(">>>> Successfully created table AuditLogs.")
}

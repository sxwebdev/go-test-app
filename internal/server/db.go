package server

import (
	"database/sql"

	"github.com/tkcrm/modules/cfg"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func (s *Server) newDB() error {

	s.logger.Infof("try to connect to postgres")

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(
		cfg.GetPostgreSqlURL(
			s.config.DBUser,
			s.config.DBPass,
			s.config.DBHost,
			s.config.DBPort,
			s.config.DBName,
		),
	)))
	if err := sqldb.Ping(); err != nil {
		return err
	}

	s.logger.Infof("successfully connected to postgres")

	db := bun.NewDB(sqldb, pgdialect.New())

	s.store = db

	return nil
}

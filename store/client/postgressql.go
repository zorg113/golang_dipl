package client

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // for postgres driver import
	"github.com/rs/zerolog"
	"github.com/zorg113/golang_dipl/atibruteforce/internal/config"
)

type PostgresSql struct {
	Db     *sqlx.DB
	cfg    *config.Config
	logger *zerolog.Logger
}

func NewPostgresSql(log *zerolog.Logger, config *config.Config) *PostgresSql {
	return &PostgresSql{logger: log, cfg: config}
}

func (s *PostgresSql) Open() error {
	dbSrcName := "host=" + s.cfg.DbData.Host + " " + "dbname=" + s.cfg.DbData.DbName +
		" " + "port=" + s.cfg.DbData.Port + " " + "user=" + s.cfg.DbData.DbUser +
		" " + "password=" + s.cfg.DbData.DbPassword + " " + "sslmode=" + s.cfg.DbData.SslMode
	s.logger.Info().Msg(dbSrcName)
	db, err := sqlx.Open("postgres", dbSrcName)
	if err != nil {
		s.logger.Fatal().Err(err).Msg("failed to connect to postgres")
		return err
	}
	err = db.Ping()
	if err != nil {
		s.logger.Fatal().Err(err).Msg("failed to ping postgres")
		return err
	}
	s.Db = db
	s.logger.Info().Msg("connection to db successfully")
	return nil
}

func (s *PostgresSql) Close() error {
	err := s.Db.Close()
	if err != nil {
		s.logger.Fatal().Err(err).Msg("failed to close db connection")
		return err
	}
	s.logger.Info().Msg("db connection closed successfully")
	return nil
}

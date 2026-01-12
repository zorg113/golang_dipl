package client

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // for postgres driver import
	"github.com/rs/zerolog"
	"github.com/zorg113/golang_dipl/atibruteforce/internal/config"
)

type PostgresSQL struct {
	DB     *sqlx.DB
	cfg    *config.Config
	logger *zerolog.Logger
}

func NewPostgresSQL(log *zerolog.Logger, config *config.Config) *PostgresSQL {
	return &PostgresSQL{logger: log, cfg: config}
}

func (s *PostgresSQL) Open() error {
	dbSrcName := "host=" + s.cfg.DBData.Host + " " + "dbname=" + s.cfg.DBData.DBName +
		" " + "port=" + s.cfg.DBData.Port + " " + "user=" + s.cfg.DBData.DBUser +
		" " + "password=" + s.cfg.DBData.DBPassword + " " + "sslmode=" + s.cfg.DBData.SslMode
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
	s.DB = db
	s.logger.Info().Msg("connection to db successfully")
	return nil
}

func (s *PostgresSQL) Close() error {
	err := s.DB.Close()
	if err != nil {
		s.logger.Fatal().Err(err).Msg("failed to close db connection")
		return err
	}
	s.logger.Info().Msg("db connection closed successfully")
	return nil
}

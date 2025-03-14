package pg

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/alextotalk/feline-intelligence/internal/config"
)

func NewPostgres(cfg *config.Config) (*sql.DB, error) {
	const op = "storage.pg.New"

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("не вдалося відкрити підключення до Postgres: %s %w", op, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("не вдалося виконати ping до Postgres: %s %w", op, err)
	}

	return db, nil
}

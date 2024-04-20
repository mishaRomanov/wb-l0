package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/mishaRomanov/wb-l0/internal/config"
	"github.com/mishaRomanov/wb-l0/internal/entities"
)

// Pgdb struct stands for postgres database
type Pgdb struct {
	db *pgx.Conn
}

func CreateDB() (*pgx.Conn, error) {
	cfg, err := config.InitConfig()
	if err != nil {
		return nil, err
	}
	// urlExample = "postgres://username:password@127.0.0.1:5432/database_name"
	connectString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Db)
	conn, err := pgx.Connect(context.Background(), connectString)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (p *Pgdb) WriteData(order entities.Order) error {

}

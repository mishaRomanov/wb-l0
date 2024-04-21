package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mishaRomanov/wb-l0/internal/config"
	"github.com/mishaRomanov/wb-l0/internal/entities"
	"log"
)

// Pgdb struct stands for postgres database
type Pgdb struct {
	//The use of pgxpool over pgx.Conn is justified because pgxpool
	// is concurrent-safe, according to official documentation
	Db *pgxpool.Pool
}

func CreateDB() (*Pgdb, error) {
	//creating context
	ctx := context.Background()
	//parsing config
	cfg, err := config.InitConfig()
	if err != nil {
		return nil, err
	}

	// urlExample = "postgres://username:password@localhost5432/database_name"
	connectString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Db)
	//creating pgxpool
	pool, err := pgxpool.New(ctx, connectString)
	if err != nil {
		return nil, err
	}
	return &Pgdb{
		Db: pool,
	}, nil

}

func (p *Pgdb) WriteData(order entities.Order) error {
	tag, err := p.Db.Exec(context.Background(),
		`INSERT INTO orders 
    (order_uid,
    track_number,
    entry,
    delivery,
    payment,
    items,
    locale_chr,
    internal_signature,
    customer_id,
    delivery_service,
    shardkey,
    sm_id,
    data_created,
    oof_shard) VALUES (
    @order_uid,
    @track_number,
    @entry,
    @delivery,
    @payment,
    @items,
    @locale_chr,
    @internal_signature,
    @customer_id,
    @delivery_service,
    @shardkey,
    @sm_id,
    @data_created,
    @oof_shard)`, pgx.NamedArgs{
			"order_uid":          order.OrderUID,
			"track_number":       order.TrackNumber,
			"entry":              order.Entry,
			"delivery":           order.DeliveryString(),
			"payment":            order.PaymentString(),
			"items":              order.ItemsString(),
			"locale_chr":         order.Locale,
			"internal_signature": order.InternalSignature,
			"customer_id":        order.CustomerID,
			"delivery_service":   order.DeliveryService,
			"shardkey":           order.Shardkey,
			"sm_id":              order.SmID,
			"data_created":       order.DateCreated,
			"oof_shard":          order.OofShard,
		})

	if err != nil {
		return err
	}
	r := tag.RowsAffected()
	log.Printf("Postgres data written. Rows affected %d \n", r)
	return nil
}

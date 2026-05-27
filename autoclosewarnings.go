package database

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type AutoCloseWarnings struct {
	*pgxpool.Pool
}

func newAutoCloseWarnings(db *pgxpool.Pool) *AutoCloseWarnings {
	return &AutoCloseWarnings{
		db,
	}
}

func (a AutoCloseWarnings) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS auto_close_warnings(
	"guild_id" int8 NOT NULL,
	"ticket_id" int4 NOT NULL,
	"last_message_id" int8,
	"sent_at" timestamptz NOT NULL DEFAULT NOW(),
	FOREIGN KEY("guild_id", "ticket_id") REFERENCES tickets("guild_id", "id") ON DELETE CASCADE,
	PRIMARY KEY("guild_id", "ticket_id")
);
`
}

func (a *AutoCloseWarnings) MarkSent(ctx context.Context, guildId uint64, ticketId int, lastMessageId *uint64) (err error) {
	query := `
INSERT INTO auto_close_warnings("guild_id", "ticket_id", "last_message_id", "sent_at")
VALUES($1, $2, $3, NOW())
ON CONFLICT("guild_id", "ticket_id")
DO UPDATE SET "last_message_id" = $3, "sent_at" = NOW();`

	_, err = a.Exec(ctx, query, guildId, ticketId, lastMessageId)
	return
}

func (a *AutoCloseWarnings) Delete(ctx context.Context, guildId uint64, ticketId int) (err error) {
	query := `DELETE FROM auto_close_warnings WHERE "guild_id" = $1 AND "ticket_id" = $2;`
	_, err = a.Exec(ctx, query, guildId, ticketId)
	return
}

package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/viniciusgferreira/viniciusgferreira-rinhaq1-2024-go/data"
	"log"
	"time"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateStatement(ctx context.Context, id int) *data.Client {
	//log.Printf("idle: %v - acquired: %v - total: %v\n", r.db.Stat().IdleConns(), r.db.Stat().AcquiredConns(), r.db.Stat().TotalConns())
	balance := r.FindBalance(ctx, id)
	balance.Date = time.Now()
	transactions := r.GetLastTransactions(ctx, id)

	client := &data.Client{
		Balance:      balance,
		Transactions: transactions,
	}
	return client
}

func (r *Repository) FindBalance(ctx context.Context, id int) *data.Balance {
	statement := fmt.Sprintf("SELECT limite, saldo from clientes where id=$1")
	var balance data.Balance
	err := r.db.QueryRow(ctx, statement, id).Scan(&balance.Limit, &balance.Total)
	if err != nil {
		log.Println("queryrow findbalance", err.Error())
		return nil
	}
	return &balance
}

func (r *Repository) Debit(ctx context.Context, t *data.Transaction, id int) (*data.ResponseTransaction, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		log.Println("begin transaction:", err.Error())
		return nil, err
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		_ = tx.Rollback(ctx)

	}(tx, ctx)

	response := &data.ResponseTransaction{}
	err = tx.QueryRow(ctx, "SELECT limite, saldo FROM clientes WHERE id=$1 FOR UPDATE", id).Scan(&response.Limit, &response.Total)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("no such client")
		}
		log.Println("error select update: ", err.Error())
		return nil, err
	}
	newBalance := response.Total - t.Amount
	if newBalance < (response.Limit * -1) {
		return nil, errors.New("insufficient funds")
	}
	_, err = tx.Exec(ctx, "UPDATE clientes SET saldo=$1 where id=$2", newBalance, id)
	response.Total = newBalance

	_, err = tx.Exec(ctx, "INSERT INTO transacoes (cliente_id, valor, tipo, descricao, realizada_em) VALUES ($1, $2, $3, $4, $5)", id, t.Amount, t.Type, t.Description, t.Date)
	if err != nil {
		log.Println("insert", err.Error())
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Println("erro commit debit", err.Error())
		return nil, err
	}
	return response, nil
}

func (r *Repository) Credit(ctx context.Context, t *data.Transaction, id int) (*data.ResponseTransaction, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		log.Println("begin transaction:", err.Error())
		return nil, err
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		_ = tx.Rollback(ctx)
	}(tx, ctx)

	response := &data.ResponseTransaction{}
	err = tx.QueryRow(ctx, "SELECT limite, saldo FROM clientes WHERE id=$1", id).Scan(&response.Limit, &response.Total)
	if err != nil {
		log.Println("select get saldo/limite:", err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("no such client")
		}
		return nil, err
	}

	_, err = tx.Exec(ctx, "UPDATE clientes SET saldo=saldo+$1 WHERE id=$2", t.Amount, id)
	if err != nil {
		log.Println("update saldo:", err.Error())
		return nil, err
	}

	_, err = tx.Exec(ctx, "INSERT INTO transacoes (cliente_id, valor, tipo, descricao, realizada_em) VALUES ($1, $2, $3, $4, $5)", id, t.Amount, t.Type, t.Description, t.Date)
	if err != nil {
		log.Println("insert transaction:", err.Error())
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Println("erro commit credit", err.Error())
		return nil, err
	}

	response.Total += t.Amount
	return response, nil
}

func (r *Repository) GetLastTransactions(ctx context.Context, id int) []*data.Transaction {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		log.Println("begin transaction extrato:", err.Error())
		return nil
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err == nil {
			log.Println("rollback rolou extrato")
		}
	}(tx, ctx)

	rows, err := r.db.Query(ctx, "SELECT valor, tipo, descricao, realizada_em FROM transacoes WHERE cliente_id=$1 ORDER BY realizada_em DESC LIMIT 10", id)
	if err != nil {
		log.Println("query", err.Error())
	}
	defer rows.Close()

	var transactions []*data.Transaction

	for rows.Next() {
		var valor int
		var tipo string
		var descricao string
		var realizadaEm time.Time

		if err := rows.Scan(&valor, &tipo, &descricao, &realizadaEm); err != nil {
			log.Println("scan", err.Error())
		}
		transaction := &data.Transaction{Amount: valor, Type: tipo, Description: descricao, Date: realizadaEm}
		transactions = append(transactions, transaction)
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Println("erro commit getlast", err.Error())
		return nil
	}
	return transactions
}

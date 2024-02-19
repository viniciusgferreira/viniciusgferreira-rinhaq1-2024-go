package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/viniciusgferreira/viniciusgferreira-rinhaq1-2024-go/data"
	"log"
	"time"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateStatement(id uint8) *data.Client {
	balance := r.FindBalance(id)
	log.Println(balance)
	transactions := r.GetLastTransactions(id)

	client := &data.Client{
		Balance:      balance,
		Transactions: transactions,
	}
	return client
}

func (r *Repository) FindBalance(id uint8) *data.Balance {
	statement := fmt.Sprintf("SELECT limite, saldo from clientes where id=$1")
	var balance data.Balance
	err := r.db.QueryRow(statement, id).Scan(&balance.Limit, &balance.Total)
	if err != nil {
		log.Println("queryrow findbalance", err.Error())
	}
	return &balance
}

func (r *Repository) Debit(t *data.Transaction, id uint8) (*data.Balance, error) {
	tx, err := r.db.Begin()
	if err != nil {
		log.Println("begin transaction:", err.Error())
		return nil, err
	}
	defer tx.Rollback()

	newBalance := &data.Balance{}
	err = tx.QueryRow("SELECT limite, saldo FROM clientes WHERE id=$1", id).Scan(&newBalance.Limit, &newBalance.Total)
	if err != nil {
		log.Println("select", err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("no such client")
		}
		return nil, err
	}

	if newBalance.Total+newBalance.Limit < t.Amount {
		return nil, errors.New("not enough funds or limit")
	}
	_, err = tx.Exec("UPDATE clientes SET saldo=saldo-$1 where id=$2", t.Amount, id)
	newBalance.Total -= t.Amount

	_, err = tx.Exec("INSERT INTO transacoes (cliente_id, valor, tipo, descricao, realizada_em) VALUES ($1, $2, $3, $4, $5)", id, t.Amount, t.Type, t.Description, t.Date)
	if err != nil {
		log.Println("insert", err.Error())
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return newBalance, nil
}

func (r *Repository) Credit(t *data.Transaction, id uint8) (*data.Balance, error) {
	tx, err := r.db.Begin()
	if err != nil {
		log.Println("begin transaction:", err.Error())
		return nil, err
	}
	defer tx.Rollback()

	newBalance := &data.Balance{}
	err = tx.QueryRow("SELECT limite, saldo FROM clientes WHERE id=$1", id).Scan(&newBalance.Limit, &newBalance.Total)
	if err != nil {
		log.Println("select:", err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("no such client")
		}
		return nil, err
	}

	_, err = tx.Exec("UPDATE clientes SET saldo=saldo+$1 WHERE id=$2", t.Amount, id)
	if err != nil {
		log.Println("update:", err.Error())
		return nil, err
	}

	_, err = tx.Exec("INSERT INTO transacoes (cliente_id, valor, tipo, descricao, realizada_em) VALUES ($1, $2, $3, $4, $5)", id, t.Amount, t.Type, t.Description, t.Date)
	if err != nil {
		log.Println("insert:", err.Error())
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		log.Println("commit:", err.Error())
		return nil, err
	}

	newBalance.Total += t.Amount
	return newBalance, nil
}

func (r *Repository) GetLastTransactions(id uint8) []*data.Transaction {
	rows, err := r.db.Query("SELECT valor, tipo, descricao, realizada_em FROM transacoes WHERE cliente_id=$1 ORDER BY realizada_em DESC LIMIT 10", id)
	if err != nil {
		log.Println("query", err.Error())
	}
	defer rows.Close()

	var transactions []*data.Transaction

	for rows.Next() {
		var valor int
		var tipo string
		var descricao string
		var realizadaEm *time.Time

		if err := rows.Scan(&valor, &tipo, &descricao, &realizadaEm); err != nil {
			log.Println("scan", err.Error())
		}
		transaction := &data.Transaction{Amount: valor, Type: tipo, Description: descricao, Date: realizadaEm}
		log.Println(transaction)
		transactions = append(transactions, transaction)
	}
	return transactions
}

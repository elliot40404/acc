package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type transactionRepository struct {
	db *sqlx.DB
}

type Transaction struct {
	ID          int     `db:"id" json:"id"`
	Type        string  `db:"type" json:"type"`
	Description string  `db:"description" json:"description"`
	Amount      float64 `db:"amount" json:"amount"`
	CreatedAt   string  `db:"created_at" json:"created_at"`
	UpdatedAt   string  `db:"updated_at" json:"updated_at"`
}

type TransactionConfig struct {
	Dry     bool
	Verbose bool
	TxType  string
	Page    int
	Limit   int
	All     bool
	Date    string
	Amount  string
	Sort    string
	SortAsc bool
	Desc    string
	Columns []string
	IsHRTime bool
	Format string
	IsPretty bool
}

type DeleteConfig struct {
	Dry     bool
	Verbose bool
	All     bool
	Ids     []string
}

type TQuery struct {
	Query   string
	Config  TransactionConfig
	isCount bool
}

type TransactionRepository interface {
	CreateTransaction(transaction Transaction) error
	GetTransactionsWithConfig(c TransactionConfig) ([]Transaction, error)
	GetTransactionCountWithConfig(c TransactionConfig) (int, error)
	DeleteTransactions(c DeleteConfig) error
}

func NewTransactionRepository() TransactionRepository {
	db, err := GetDB()
	if err != nil {
		panic(err)
	}
	return &transactionRepository{db: db}
}

func (r *transactionRepository) CreateTransaction(transaction Transaction) error {
	_, err := r.db.Exec("INSERT INTO transactions (type, description, amount) VALUES (?, ?, ?)", transaction.Type, transaction.Description, transaction.Amount)
	if err != nil {
		return err
	}
	return nil
}

func (r *transactionRepository) GetTransactionsWithConfig(c TransactionConfig) ([]Transaction, error) {
	var transactions []Transaction
	query, err := buildQuery(c, false)
	if err != nil {
		return nil, err
	}
	if c.Verbose {
		fmt.Println("SELECT =>", query)
	}
	if c.Dry {
		return nil, nil
	}
	err = r.db.Select(&transactions, query)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *transactionRepository) GetTransactionCountWithConfig(c TransactionConfig) (int, error) {
	query, err := buildQuery(c, true)
	if err != nil {
		return 0, err
	}
	if c.Verbose {
		fmt.Println("COUNT =>", query)
	}
	if c.Dry {
		return 0, nil
	}
	var count int
	err = r.db.Get(&count, query)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *transactionRepository) DeleteTransactions(c DeleteConfig) error {
	var query string
	if c.All {
		query = "DELETE FROM transactions"
		if c.Verbose {
			fmt.Println("DELETE =>", query)
		}
		if c.Dry {
			return nil
		}
		_, err := r.db.Exec(query)
		if err != nil {
			return err
		}
		_, err = r.db.Exec("DELETE FROM sqlite_sequence WHERE name='transactions'")
		if err != nil {
			return err
		}
		// VACUUM is used to rebuild the database file, repacking it into a minimal amount of disk space.
		_, err = r.db.Exec("VACUUM")
		if err != nil {
			return err
		}
		return nil
	}
	q, args, err := sqlx.In("DELETE FROM transactions WHERE id IN (?)", c.Ids)
	if err != nil {
		return err
	}
	q = r.db.Rebind(q)
	if c.Verbose {
		fmt.Println("DELETE =>", q, args)
	}
	if c.Dry {
		return nil
	}
	_, err = r.db.Exec(q, args...)
	if err != nil {
		return err
	}
	return nil
}

func NewQuery(t TransactionConfig, isCount bool) TQuery {
	if isCount {
		return TQuery{
			Query:   "SELECT COUNT(*) FROM transactions",
			Config:  t,
			isCount: isCount,
		}
	}
	return TQuery{
		Query:   "SELECT * FROM transactions",
		Config:  t,
		isCount: isCount,
	}
}

func buildQuery(c TransactionConfig, isCount bool) (string, error) {
	q := NewQuery(c, isCount)
	// q.AddColumns()
	q.AddType()
	q.AddDate()
	q.AddAmount()
	q.AddDesc()
	if isCount {
		return q.Build(), nil
	}
	q.AddSort()
	if c.All {
		return q.Build(), nil
	}
	q.AddLimit()
	q.AddOffset()
	return q.Build(), nil
}

package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		Queries: New(db),
		db:      db,
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err %v,rb err%v", err, rbErr)
		}
		return err
	}
	return tx.Commit()

}

type TransferTxParams struct {
	SenderId   int64 `json:"senderId"`
	ReceiverId int64 `json:"receiverId"`
	Amount     int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer        Transfer `json:"transfer"`
	SenderAccount   Account  `json:"senderAccount"`
	ReceiverAccount Account  `json:"receiverAccount"`
	SenderEntry     Entry    `json:"senderEntry"`
	ReceiverEntry   Entry    `json:"receiverEntry"`
}

var txKey = struct{}{}

func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		txName := ctx.Value(txKey)
		fmt.Println(txName, "create transfer")
		result.Transfer, err = q.NewTransfer(ctx, NewTransferParams{
			Amount:     arg.Amount,
			SenderID:   arg.SenderId,
			ReceiverID: arg.ReceiverId,
		})
		if err != nil {
			return err
		}
		fmt.Println(txName, "create entry 1")
		result.SenderEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			Amount: -arg.Amount,
			AccountID: sql.NullInt64{
				Int64: arg.SenderId,
				Valid: true,
			},
		})
		if err != nil {
			return err
		}
		fmt.Println(txName, "create entry 2")
		result.ReceiverEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			Amount: arg.Amount,
			AccountID: sql.NullInt64{
				Int64: arg.ReceiverId,
				Valid: true,
			},
		})
		if err != nil {
			return err
		}
		fmt.Println(txName, "Getting account 1 for update")
		senderAct, err := q.GetAccountForUpdate(ctx, arg.SenderId)
		if senderAct.Balance < 0 {
			return err
		}
		if err != nil {
			return err
		}
		fmt.Println(txName, "Updating account 1")
		result.SenderAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      arg.SenderId,
			Balance: senderAct.Balance - arg.Amount,
		})
		if err != nil {
			return err
		}
		fmt.Println(txName, "Getting account 2 for update")
		receipientAct, err := q.GetAccountForUpdate(ctx, arg.ReceiverId)
		if err != nil {
			return err
		}
		if receipientAct.Balance < 0 {
			return err
		}
		fmt.Println(txName, "Updating account 2")
		result.ReceiverAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      arg.ReceiverId,
			Balance: receipientAct.Balance + arg.Amount,
		})
		return nil
	})
	return result, err
}

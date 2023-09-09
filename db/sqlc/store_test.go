package db

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTransfer(t *testing.T) {
	store := NewStore(testDB)

	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)
	fmt.Println(">>before transaction: ", account1.Balance, ":", account2.Balance)
	//run transfer transaction concurrently
	n := 2
	amount := int64(10)

	transResChan := make(chan TransferTxResult)
	transErr := make(chan error)
	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx %d", i+1)
		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			result, err := store.TransferTx(ctx, TransferTxParams{
				SenderId:   account1.ID,
				ReceiverId: account2.ID,
				Amount:     amount,
			})
			transErr <- err
			transResChan <- result
		}()
	}
	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-transErr
		require.NoError(t, err)

		//get transfer results
		results := <-transResChan
		require.NotEmpty(t, results)

		//check transfer
		transfer := results.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.SenderID)
		require.Equal(t, account2.ID, transfer.ReceiverID)
		require.Equal(t, transfer.Amount, amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		//check entry
		senderEntry := results.SenderEntry
		require.NotZero(t, senderEntry)
		//require.Equal(t, senderEntry.AccountID, account1.ID)
		require.Equal(t, -amount, senderEntry.Amount)
		require.NotZero(t, senderEntry.CreatedAt)
		require.NotZero(t, senderEntry.ID)

		_, err = store.GetEntry(context.Background(), senderEntry.ID)
		require.NoError(t, err)

		receiverEntry := results.ReceiverEntry
		require.NotZero(t, receiverEntry)
		require.NotZero(t, receiverEntry.ID)
		require.NotZero(t, receiverEntry.CreatedAt)
		require.Equal(t, receiverEntry.Amount, amount)
		//require.Equal(t, receiverEntry.AccountID, account2.ID)

		_, err = store.GetEntry(context.Background(), receiverEntry.ID)
		require.NoError(t, err)

		//check account
		senderAct := results.SenderAccount
		require.NotEmpty(t, senderAct)
		require.Equal(t, senderAct.ID, account1.ID)

		receiverAct := results.ReceiverAccount
		require.NotEmpty(t, receiverAct)
		require.Equal(t, receiverAct.ID, account2.ID)

		//check account balance
		fmt.Println(">>tx: ", senderAct.Balance, ":", receiverAct.Balance)
		diff1 := account1.Balance - senderAct.Balance
		diff2 := receiverAct.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	fmt.Println(">>after transaction: ", updatedAccount1.Balance, ":", updatedAccount2.Balance)
	require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance)
}
func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)
	fmt.Println(">>before transaction: ", account1.Balance, ":", account2.Balance)
	//run transfer transaction concurrently
	n := 10
	amount := int64(10)

	transErr := make(chan error)
	for i := 0; i < n; i++ {
		senderId := account1.ID
		receiverId := account2.ID
		if i%2 == 1 {
			senderId = account2.ID
			receiverId = account1.ID
		}
		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				SenderId:   senderId,
				ReceiverId: receiverId,
				Amount:     amount,
			})
			transErr <- err

		}()
	}

	for i := 0; i < n; i++ {
		err := <-transErr
		require.NoError(t, err)

	}

	if account1.ID < account2.ID {
		updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)
		require.NoError(t, err)

		updatedAccount2, err := store.GetAccount(context.Background(), account2.ID)
		require.NoError(t, err)
		fmt.Println(">>after transaction: ", updatedAccount1.Balance, ":", updatedAccount2.Balance)
		require.Equal(t, account1.Balance, updatedAccount1.Balance)
		require.Equal(t, account2.Balance, updatedAccount2.Balance)
	} else {

		updatedAccount2, err := store.GetAccount(context.Background(), account2.ID)
		require.NoError(t, err)

		updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)
		require.NoError(t, err)

		fmt.Println(">>after transaction: ", updatedAccount1.Balance, ":", updatedAccount2.Balance)
		require.Equal(t, account1.Balance, updatedAccount1.Balance)
		require.Equal(t, account2.Balance, updatedAccount2.Balance)
	}

}

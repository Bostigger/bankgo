package db

import (
	"context"
	"database/sql"
	"github.com/bostigger/bankgo/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func CreateRandomTransfer(t *testing.T) Transfer {
	arg := NewTransferParams{
		Amount:     util.RandomMoney(),
		SenderID:   util.RandomActId(),
		ReceiverID: util.RandomActId(),
	}
	transfer, err := testQueries.NewTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)
	require.NotEmpty(t, transfer.ID)
	require.Equal(t, arg.Amount, transfer.Amount)
	require.Equal(t, arg.SenderID, transfer.SenderID)
	require.Equal(t, arg.ReceiverID, transfer.ReceiverID)
	return transfer
}

func TestCreateTransfer(t *testing.T) {
	CreateRandomTransfer(t)
}

func TestUpdateTransfer(t *testing.T) {
	transfer := CreateRandomTransfer(t)

	arg := UpdateTransferParams{
		ID:     transfer.ID,
		Amount: util.RandomMoney(),
	}

	trans, err := testQueries.UpdateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotZero(t, trans.ID)
	require.Equal(t, arg.ID, trans.ID)
	require.Equal(t, transfer.SenderID, trans.SenderID)
	require.Equal(t, transfer.ReceiverID, trans.ReceiverID)
	require.Equal(t, arg.Amount, trans.Amount)
	require.WithinDuration(t, transfer.CreatedAt, trans.CreatedAt, time.Second)

}

func TestGetTransfer(t *testing.T) {
	transfer := CreateRandomTransfer(t)

	transRes, err := testQueries.GetTransfer(context.Background(), transfer.ID)
	require.NoError(t, err)
	require.NotZero(t, transRes.ID)
	require.Equal(t, transfer.ID, transRes.ID)
	require.Equal(t, transRes.Amount, transRes.Amount)
	require.Equal(t, transfer.SenderID, transRes.SenderID)
	require.Equal(t, transfer.ReceiverID, transRes.ReceiverID)
	require.WithinDuration(t, transfer.CreatedAt, transRes.CreatedAt, time.Second)
}

func TestGetTransfers(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateRandomTransfer(t)
	}

	arg := GetTransfersParams{
		Limit:  5,
		Offset: 5,
	}

	transfers, err := testQueries.GetTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}
}

func TestDeleteTransfer(t *testing.T) {
	newTrans := CreateRandomTransfer(t)

	err := testQueries.DeleteTransfer(context.Background(), newTrans.ID)
	require.NoError(t, err)

	transfer, err := testQueries.GetTransfer(context.Background(), newTrans.ID)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, transfer)
}

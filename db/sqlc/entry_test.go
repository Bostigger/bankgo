package db

import (
	"context"
	"database/sql"
	"github.com/bostigger/bankgo/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func CreateRandomEntry(t *testing.T) Entry {
	arg := CreateEntryParams{
		Amount:    util.RandomMoney(),
		AccountID: util.RandomId(),
	}
	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, arg.Amount, entry.Amount)
	require.NotZero(t, entry.CreatedAt)
	require.NotZero(t, entry.ID)

	return entry

}

func TestCreateEntry(t *testing.T) {
	CreateRandomEntry(t)
}

func TestUpdateEntry(t *testing.T) {
	newentry := CreateRandomEntry(t)
	arg := UpdateEntryParams{
		ID:     newentry.ID,
		Amount: util.RandomMoney(),
	}
	entry, err := testQueries.UpdateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotZero(t, entry.ID)
	require.Equal(t, arg.Amount, entry.Amount)
	require.Equal(t, arg.ID, entry.ID)
	require.WithinDuration(t, entry.CreatedAt, newentry.CreatedAt, time.Second)
}

func TestGetSingleEntry(t *testing.T) {
	newentry := CreateRandomEntry(t)
	entry, err := testQueries.GetEntry(context.Background(), newentry.ID)
	require.NoError(t, err)
	require.NotZero(t, entry.ID)
	require.Equal(t, newentry.AccountID, entry.AccountID)
	require.Equal(t, newentry.Amount, entry.Amount)
	require.WithinDuration(t, newentry.CreatedAt, entry.CreatedAt, time.Second)

}

func TestDeleteEntry(t *testing.T) {
	newentry := CreateRandomEntry(t)
	err := testQueries.DeleteEntry(context.Background(), newentry.ID)
	require.NoError(t, err)

	entry, err := testQueries.GetEntry(context.Background(), newentry.ID)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, entry)

}

func TestGetEntries(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateRandomEntry(t)
	}
	arg := GetEntriesParams{
		Limit:  5,
		Offset: 5,
	}
	entries, err := testQueries.GetEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}

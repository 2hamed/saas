package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDataStore(t *testing.T) {

	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	storage, err := NewDataStore()
	assert.NoError(t, err)

	_, err = storage.Fetch("someurl")
	assert.Error(t, err)
	//-------------
	exists, isPending, isFinished, err := storage.FetchStatus("someurl")

	assert.NoError(t, err)
	assert.False(t, exists)
	assert.False(t, isPending)
	assert.False(t, isFinished)
	//-------------
	err = storage.Store("someurl", "somepath")
	assert.NoError(t, err)

	path, err := storage.Fetch("someurl")
	assert.NoError(t, err)

	assert.Equal(t, "somepath", path)

	//----------------
	err = storage.UpdatePath("someurl", "newpath")
	assert.NoError(t, err)

	path, err = storage.Fetch("someurl")
	assert.NoError(t, err)

	assert.Equal(t, "newpath", path)

	//----------------
	exists, isPending, isFinished, err = storage.FetchStatus("someurl")

	assert.NoError(t, err)
	assert.True(t, exists)
	assert.True(t, isPending)
	assert.False(t, isFinished)

	//---------------
	err = storage.SetFinished("someurl")
	assert.NoError(t, err)

	//---------------
	exists, isPending, isFinished, err = storage.FetchStatus("someurl")

	assert.NoError(t, err)
	assert.True(t, exists)
	assert.False(t, isPending)
	assert.True(t, isFinished)

	//---------------
	err = storage.SetFailed("someurl")
	assert.NoError(t, err)

	exists, isPending, isFinished, err = storage.FetchStatus("someurl")

	assert.NoError(t, err)
	assert.True(t, exists)
	assert.False(t, isPending)
	assert.False(t, isFinished)

	err = storage.Delete("someurl")
	assert.NoError(t, err)
}

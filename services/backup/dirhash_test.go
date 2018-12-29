package backup_test

import (
	"testing"

	"github.com/sesheffield/go-blueprints/services/backup"
	"github.com/stretchr/testify/require"
)

func TestDirHash(t *testing.T) {
	backup.Setup(t)
	defer backup.Teardown(t)
	hash1a, err := backup.DirHash("test/output")
	require.NoError(t, err)
	hash1b, err := backup.DirHash("test/output")
	require.NoError(t, err)

	require.Equal(t, hash1a, hash1b, "hash1 and hash1b should be identical")

	hash2, err := backup.DirHash("test/output1")
	require.NoError(t, err)

	require.NotEqual(t, hash1a, hash2, "hash1 and hash2 should not be the same")

}

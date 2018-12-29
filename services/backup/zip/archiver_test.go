package zip_test

import (
	"testing"

	"github.com/sesheffield/go-blueprints/services/backup"
	"github.com/sesheffield/go-blueprints/services/backup/zip"
	"github.com/stretchr/testify/require"
)

func TestZipArchive(t *testing.T) {
	backup.Setup(t)
	defer backup.Teardown(t)

	err := zip.Archive("test/output", "test/output/1.zip")
	//err := backup.ArchiveFunc("test/output", "test/output/1.zip")
	require.NoError(t, err)

	// unzip
	err = zip.Restore("test/output/1.zip", "test/output/restored")
	require.NoError(t, err)

}

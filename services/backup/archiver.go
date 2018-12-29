package backup

// ArchiveFunc s responsbile for arching the source folder and storing it in the
// destination path
type ArchiveFunc func(src, dest string) error

// DestFmtFunc is responsible for handling extension names
type DestFmtFunc func() string

// RestoreFunc is response for restoring archived folder
type RestoreFunc func(src, dest string) error

package zip

import (
	"archive/zip"
	"io"
	"os"
	"path"
	"path/filepath"
	"sync"
)

// DestFmt handles extension for zip
func DestFmt() string {
	return "%d.zip"
}

// Archive will archive a folder as a zip
func Archive(src, dest string) error {
	if err := os.MkdirAll(filepath.Dir(dest), 0777); err != nil {
		return err
	}
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()
	w := zip.NewWriter(out)
	defer w.Close()
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil // skip
		}
		if err != nil {
			return err
		}
		in, err := os.Open(path)
		if err != nil {
			return err
		}
		defer in.Close()
		f, err := w.Create(path)
		if err != nil {
			return err
		}
		_, err = io.Copy(f, in)
		if err != nil {
			return err
		}
		return nil
	})
}

// Restore will restore a zip folder
func Restore(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()
	var w sync.WaitGroup
	var errs []error
	errChan := make(chan error)
	go func() {
		for err := range errChan {
			errs = append(errs, err)
		}
	}()
	for _, f := range r.File {
		w.Add(1)
		go func(f *zip.File) {
			zippedfile, err := f.Open()
			if err != nil {
				errChan <- err
				w.Done()
				return
			}
			toFilename := path.Join(dest, f.Name)
			err = os.MkdirAll(path.Dir(toFilename), 0777)
			if err != nil {
				errChan <- err
				w.Done()
				return
			}
			newFile, err := os.Create(toFilename)
			if err != nil {
				zippedfile.Close()
				errChan <- err
				w.Done()
				return
			}
			_, err = io.Copy(newFile, zippedfile)
			newFile.Close()
			zippedfile.Close()
			if err != nil {
				errChan <- err
				w.Done()
				return
			}
			w.Done()
		}(f)
	}
	w.Wait()
	close(errChan)
	if len(errs) > 0 {
		return errs[0] // return first error
	}
	return nil
}

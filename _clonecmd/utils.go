package main

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type argString []string

func (a argString) Get(i int, args ...string) (r string) {
	if i >= 0 && i < len(a) {
		r = a[i]
	} else if len(args) > 0 {
		r = args[0]
	}
	return
}

func SanitizePath(fileOrDirPath string) string {
	return path.Clean(strings.Replace(fileOrDirPath, "\\", "/", -1))
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func IsPathExcluded(exclusionList_FullPaths *[]string, sourcePath string) bool {
	if exclusionList_FullPaths == nil {
		return false
	}
	for _, p := range *exclusionList_FullPaths {
		if strings.ToLower(SanitizePath(p)) == strings.ToLower(SanitizePath(sourcePath)) {
			return true
		}
	}
	return false
}

//Thanks to https://gist.github.com/jaybill/2876519
// Copies file source to destination dest.
func CopyFile(source string, dest string, exclusionList_FullPaths *[]string) (err error) {
	sf, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sf.Close()

	if IsPathExcluded(exclusionList_FullPaths, source) {
		return nil
	}

	df, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer df.Close()
	_, err = io.Copy(df, sf)
	if err == nil {
		si, err := os.Stat(source)
		if err != nil {
			err = os.Chmod(dest, si.Mode())
		}

	}

	return
}

//Thanks to https://gist.github.com/jaybill/2876519
// Recursively copies a directory tree, attempting to preserve permissions.
// Source directory must exist, destination directory must *not* exist.
func CopyDir(source string, dest string, exclusionList_FullPaths *[]string) []error {
	// get properties of source dir
	fi, err := os.Stat(source)
	if err != nil {
		return []error{err}
	}

	if !fi.IsDir() {
		return []error{errors.New("Source is not a directory")}
	}

	if IsPathExcluded(exclusionList_FullPaths, source) {
		return nil
	}

	_, err = os.Open(dest)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dest, fi.Mode())
		if err != nil {
			return []error{err}
		}
	}

	entries, err := ioutil.ReadDir(source)

	errs := []error{}
	for _, entry := range entries {

		sfp := source + "/" + entry.Name()
		dfp := dest + "/" + entry.Name()
		if entry.IsDir() {
			tmpErrors := CopyDir(sfp, dfp, exclusionList_FullPaths)
			if tmpErrors != nil {
				for _, e := range tmpErrors {
					errs = append(errs, e)
				}
			}
		} else {
			// perform copy
			err = CopyFile(sfp, dfp, exclusionList_FullPaths)
			if err != nil {
				errs = append(errs, err)
			}
		}

	}

	if len(errs) == 0 {
		return nil
	} else {
		return errs
	}
}

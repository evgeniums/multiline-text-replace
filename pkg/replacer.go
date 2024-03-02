package multiline_text_replace

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Options struct {
	PatternFile string `short:"p" long:"pattern-file" description:"File that contents pattern text to search." required:"true"`
	SubstFile   string `short:"s" long:"substitution-file" description:"File that contents substitution text to replace with." required:"true"`
	TargetFile  string `short:"f" long:"file" description:"Target file to replace text in. Either file or dir must be specified."`
	TargetDir   string `short:"d" long:"dir" description:"Target directory containing files to replace text in. Either file or dir must be specified."`
	TargeExt    string `short:"e" long:"ext" description:"Extensions of files in target directory to replace text in, can be multiple separated with comms, e.g. .txt,.csv"`
	Recursive   bool   `short:"r" long:"recursive" description:"Process files recursivley in the target directory."`
}

func ReplaceTextInFile(pattern string, subst string, fileName string) error {

	fmt.Printf("replacing text in file %s...\n", fileName)

	// read file
	dat, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	// make temporary file copy
	tmpFileName := fileName + ".~tmp"
	err = os.WriteFile(tmpFileName, dat, 0755)
	if err != nil {
		return fmt.Errorf("failed to create temporary file copy for %s: %v", fileName, err)
	}

	// replace content
	src := string(dat)
	dst := strings.ReplaceAll(src, pattern, subst)

	// overwrite file
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	err = f.Truncate(0)
	if err != nil {
		return err
	}
	_, err = f.Seek(0, 0)
	if err != nil {
		return err
	}
	_, err = f.Write([]byte(dst))
	if err != nil {
		return err
	}
	err = f.Close()
	if err != nil {
		return err
	}

	// remove tmp file
	os.Remove(tmpFileName)

	// done
	fmt.Printf("successfully replaced text in file %s...\n", fileName)
	return nil
}

func ReplaceText(opts Options) error {

	// check mode
	if opts.TargetFile == "" && opts.TargetDir == "" {
		return errors.New("either target file or dir must be specified")
	}

	// read pattern
	pattern, err := os.ReadFile(opts.PatternFile)
	if err != nil {
		return err
	}

	// read substitution
	subst, err := os.ReadFile(opts.SubstFile)
	if err != nil {
		return err
	}

	// replace text in file
	if opts.TargetFile != "" {
		return ReplaceTextInFile(string(pattern), string(subst), opts.TargetFile)
	}

	// replace text in dir
	dirPath, err := filepath.Abs(opts.TargetDir)
	if err != nil {
		return err
	}
	exts := make([]string, 0)
	if opts.TargeExt != "" {
		exts = strings.Split(opts.TargeExt, ",")
	}
	checkExt := func(filePath string) bool {
		fileExt := filepath.Ext(filePath)
		matchedExt := len(exts) == 0
		if !matchedExt {
			for _, ext := range exts {
				if ext == fileExt {
					matchedExt = true
					break
				}
			}
		}
		return matchedExt
	}
	if opts.Recursive {
		err = filepath.WalkDir(opts.TargetDir, func(path string, d fs.DirEntry, e error) error {
			if e != nil {
				return e
			}
			if !d.IsDir() && checkExt(path) {

				filePath := path
				parts := strings.Split(path, string(os.PathSeparator))
				if len(parts) > 1 {
					filePath = filepath.Join(parts[1:]...)
					filePath = filepath.Join(dirPath, filePath)
					err = ReplaceTextInFile(string(pattern), string(subst), filePath)
					if err != nil {
						return err
					}
				}
			}

			return nil
		})
		if err != nil {
			return err
		}
	} else {
		items, err := os.ReadDir(opts.TargetDir)
		if err != nil {
			return err
		}
		for _, item := range items {
			if !item.IsDir() && checkExt(item.Name()) {
				filePath := filepath.Join(dirPath, item.Name())
				err = ReplaceTextInFile(string(pattern), string(subst), filePath)
				if err != nil {
					return err
				}
			}
		}
	}

	// done
	return nil
}

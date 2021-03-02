package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func dirTreeRec(out io.Writer, path string, printFiles bool, prefix string) error {
	dir, err := os.Open(path)
	if err != nil {
		err = fmt.Errorf("can't open %s: %s", path, err)
		return (err)
	}
	defer dir.Close()

	stat, err := dir.Stat()
	if err != nil {
		err := fmt.Errorf("can't get stat %s: %s", path, err)
		return err
	}

	if !stat.IsDir() && printFiles {
		size := ""
		if stat.Size() == 0 {
			size = " (empty)"
		} else {
			size = " (" + strconv.Itoa(int(stat.Size())) + "b)"
		}
		fmt.Fprintln(out, prefix+filepath.Base(path)+size)
		return nil
	}

	files, err := dir.Readdir(-1)
	if err != nil {
		err = fmt.Errorf("can't read dir%s: %s", path, err)
		return err
	}
	if prefix != "" {
		fmt.Fprintln(out, prefix+filepath.Base(dir.Name()))
	}

	if !printFiles {
		new_files := make([]os.FileInfo, 0)
		for _, file := range files {
			if file.IsDir() {
				new_files = append(new_files, file)
			}
		}
		files = new_files
	}

	prefix = strings.Replace(prefix, "├───", "│\t", 1)
	prefix = strings.Replace(prefix, "└───", "\t", 1)
	prefix = prefix + "├───"
	sort.Slice(files, func(i, j int) bool { return files[i].Name() < files[j].Name() })
	for i, file := range files {
		if i+1 == len(files) {
			prefix = strings.Replace(prefix, "├───", "└───", 1)
		}
		new_path := path + string(os.PathSeparator) + file.Name()
		err := dirTreeRec(out, new_path, printFiles, prefix)
		if err != nil {
			return err
		}
	}
	return nil

}

func dirTree(out io.Writer, path string, printFiles bool) error {
	err := dirTreeRec(out, path, printFiles, "")
	if err != nil {
		err = fmt.Errorf("problem with dirTreeRec: %s", err)
		return err
	}
	return nil
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

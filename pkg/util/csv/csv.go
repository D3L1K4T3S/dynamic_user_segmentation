package csv

import (
	e "dynamic-user-segmentation/pkg/util/errors"
	"encoding/csv"
	"os"
	"path/filepath"
)

func CreateCSV(records [][]string, createPath string) error {
	f, err := os.Create(filepath.Join(createPath, "users.csv"))
	defer func() {
		err = f.Close()
		err = e.WrapIfErr("can't create CSV file: ", err)
	}()

	if err != nil {
		return e.Wrap("can't create a file: ", err)
	}

	w := csv.NewWriter(f)
	err = w.WriteAll(records)

	if err != nil {
		return e.Wrap("can't write data in file: ", err)
	}

	return nil
}

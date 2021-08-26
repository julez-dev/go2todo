package lists

import (
	"context"
	"encoding/json"
	"io"
	"os"

	"github.com/julez-dev/go2todo/fileutil"
)

type InFile struct {
	fileName string
	inMem    *InMemory
}

func NewInFile(path string) (*InFile, error) {
	inMem := NewInMemory()

	var file *os.File

	if !fileutil.FileExists(path) {
		newFile, err := os.Create(path)

		if err != nil {
			return nil, err
		}

		file = newFile
	} else {
		existingFile, err := os.Open(path)

		if err != nil {
			return nil, err
		}

		file = existingFile
	}

	defer file.Close()

	lists := []*List{}

	err := json.NewDecoder(file).Decode(&lists)

	if err != nil && err != io.EOF {
		return nil, err
	}

	for _, list := range lists {
		_, _ = inMem.CreateList(context.Background(), list)
	}

	return &InFile{
		fileName: path,
		inMem:    inMem,
	}, nil
}

func (inFile *InFile) CreateList(ctx context.Context, list *List) (*List, error) {
	f, err := fileutil.OpenAndTruncate(inFile.fileName)

	if err != nil {
		return nil, err
	}

	defer f.Close()

	_, _ = inFile.inMem.CreateList(ctx, list)
	lists, _ := inFile.inMem.GetLists(ctx)

	err = json.NewEncoder(f).Encode(lists)

	if err != nil {
		return nil, err
	}

	return list, nil
}

func (inFile *InFile) GetList(ctx context.Context, search string) (*List, error) {
	return inFile.inMem.GetList(ctx, search)
}

func (inFile *InFile) GetLists(ctx context.Context) ([]*List, error) {
	return inFile.inMem.GetLists(ctx)
}

func (inFile *InFile) DeleteList(ctx context.Context, id string) error {
	f, err := fileutil.OpenAndTruncate(inFile.fileName)

	if err != nil {
		return err
	}

	defer f.Close()

	_ = inFile.inMem.DeleteList(ctx, id)
	lists, _ := inFile.inMem.GetLists(ctx)

	err = json.NewEncoder(f).Encode(lists)

	if err != nil {
		return err
	}

	return nil
}

func (inFile *InFile) DeleteLists(ctx context.Context) error {
	f, err := fileutil.OpenAndTruncate(inFile.fileName)

	if err != nil {
		return err
	}

	defer f.Close()

	_ = inFile.inMem.DeleteLists(ctx)
	lists, _ := inFile.inMem.GetLists(ctx)

	err = json.NewEncoder(f).Encode(lists)

	if err != nil {
		return err
	}

	return nil
}

func (inFile *InFile) UpdateList(ctx context.Context, list *List) (*List, error) {
	f, err := fileutil.OpenAndTruncate(inFile.fileName)

	if err != nil {
		return nil, err
	}

	defer f.Close()

	_, _ = inFile.inMem.UpdateList(ctx, list)
	lists, _ := inFile.inMem.GetLists(ctx)

	err = json.NewEncoder(f).Encode(lists)

	if err != nil {
		return nil, err
	}

	return list, nil
}

package tasks

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

	tasks := []*Task{}

	err := json.NewDecoder(file).Decode(&tasks)

	if err != nil && err != io.EOF {
		return nil, err
	}

	for _, task := range tasks {
		_, _ = inMem.CreateTask(context.Background(), task)
	}

	return &InFile{
		fileName: path,
		inMem:    inMem,
	}, nil
}

func (inFile *InFile) CreateTask(ctx context.Context, task *Task) (*Task, error) {
	f, err := fileutil.OpenAndTruncate(inFile.fileName)

	if err != nil {
		return nil, err
	}

	defer f.Close()

	_, _ = inFile.inMem.CreateTask(ctx, task)
	tasks, _ := inFile.inMem.GetAllTasks(ctx)

	err = json.NewEncoder(f).Encode(tasks)

	if err != nil {
		return nil, err
	}

	return task, nil
}

func (inFile *InFile) GetTask(ctx context.Context, search string) (*Task, error) {
	return inFile.inMem.GetTask(ctx, search)
}

func (inFile *InFile) GetAllTasks(ctx context.Context) ([]*Task, error) {
	return inFile.inMem.GetAllTasks(ctx)
}

func (inFile *InFile) GetTasks(ctx context.Context, listID string) ([]*Task, error) {
	return inFile.inMem.GetTasks(ctx, listID)
}

func (inFile *InFile) DeleteTask(ctx context.Context, id string) error {
	f, err := fileutil.OpenAndTruncate(inFile.fileName)

	if err != nil {
		return err
	}

	defer f.Close()

	_ = inFile.inMem.DeleteTask(ctx, id)
	tasks, _ := inFile.inMem.GetAllTasks(ctx)

	err = json.NewEncoder(f).Encode(tasks)

	if err != nil {
		return err
	}

	return nil
}

func (inFile *InFile) DeleteTasks(ctx context.Context, listID string) error {
	f, err := fileutil.OpenAndTruncate(inFile.fileName)

	if err != nil {
		return err
	}

	defer f.Close()

	_ = inFile.inMem.DeleteTasks(ctx, listID)
	tasks, _ := inFile.inMem.GetAllTasks(ctx)

	err = json.NewEncoder(f).Encode(tasks)

	if err != nil {
		return err
	}

	return nil
}

func (inFile *InFile) UpdateTask(ctx context.Context, task *Task) (*Task, error) {
	f, err := fileutil.OpenAndTruncate(inFile.fileName)

	if err != nil {
		return nil, err
	}

	defer f.Close()

	_, _ = inFile.inMem.UpdateTask(ctx, task)
	tasks, _ := inFile.inMem.GetAllTasks(ctx)

	err = json.NewEncoder(f).Encode(tasks)

	if err != nil {
		return nil, err
	}

	return task, nil
}

package openai //nolint:testpackage // testing private field

import (
	"context"
	"fmt"
	"io"
	"os"
	"testing"

	utils "github.com/sashabaranov/go-openai/internal"
	"github.com/sashabaranov/go-openai/internal/test/checks"
)

func TestFileBytesUploadWithFailingFormBuilder(t *testing.T) {
	config := DefaultConfig("")
	config.BaseURL = ""
	client := NewClientWithConfig(config)
	mockBuilder := &mockFormBuilder{}
	client.createFormBuilder = func(io.Writer) utils.FormBuilder {
		return mockBuilder
	}

	ctx := context.Background()
	req := FileBytesRequest{
		Name:    "foo",
		Bytes:   []byte("foo"),
		Purpose: PurposeAssistants,
	}

	mockError := fmt.Errorf("mockWriteField error")
	mockBuilder.mockWriteField = func(string, string) error {
		return mockError
	}
	_, err := client.CreateFileBytes(ctx, req)
	checks.ErrorIs(t, err, mockError, "CreateFile should return error if form builder fails")

	mockError = fmt.Errorf("mockCreateFormFile error")
	mockBuilder.mockWriteField = func(string, string) error {
		return nil
	}
	mockBuilder.mockCreateFormFileReader = func(string, io.Reader, string) error {
		return mockError
	}
	_, err = client.CreateFileBytes(ctx, req)
	checks.ErrorIs(t, err, mockError, "CreateFile should return error if form builder fails")

	mockError = fmt.Errorf("mockClose error")
	mockBuilder.mockWriteField = func(string, string) error {
		return nil
	}
	mockBuilder.mockCreateFormFileReader = func(string, io.Reader, string) error {
		return nil
	}
	mockBuilder.mockClose = func() error {
		return mockError
	}
	_, err = client.CreateFileBytes(ctx, req)
	checks.ErrorIs(t, err, mockError, "CreateFile should return error if form builder fails")
}

func TestFileUploadWithFailingFormBuilder(t *testing.T) {
	config := DefaultConfig("")
	config.BaseURL = ""
	client := NewClientWithConfig(config)
	mockBuilder := &mockFormBuilder{}
	client.createFormBuilder = func(io.Writer) utils.FormBuilder {
		return mockBuilder
	}

	ctx := context.Background()
	req := FileRequest{
		FileName: "test.go",
		FilePath: "client.go",
		Purpose:  PurposeAssistants,
	}

	// mock WriteFiled with given error
	fieldWith := func(err error) func(string, string) error {
		return func(string, string) error {
			return err
		}
	}
	// mock CreateFormFileReader with given error
	ffrWith := func(err error) func(string, io.Reader, string) error {
		return func(string, io.Reader, string) error {
			return err
		}
	}
	mockError := fmt.Errorf("mockWriteField error")
	mockBuilder.mockWriteField = fieldWith(mockError)
	_, err := client.CreateFile(ctx, req)
	checks.ErrorIs(t, err, mockError, "CreateFile should return error if form builder fails")

	mockError = fmt.Errorf("mockCreateFormFile error")
	mockBuilder.mockWriteField = fieldWith(nil)
	mockBuilder.mockCreateFormFileReader = ffrWith(mockError)
	_, err = client.CreateFile(ctx, req)
	checks.ErrorIs(t, err, mockError, "CreateFile should return error if form builder fails")

	mockError = fmt.Errorf("mockClose error")
	mockBuilder.mockWriteField = fieldWith(nil)
	mockBuilder.mockCreateFormFileReader = ffrWith(nil)
	mockBuilder.mockClose = func() error {
		return mockError
	}
	_, err = client.CreateFile(ctx, req)
	if err == nil {
		t.Fatal("CreateFile should return error if form builder fails")
	}
	checks.ErrorIs(t, err, mockError, "CreateFile should return error if form builder fails")
}

func TestFileUploadWithNonExistentPath(t *testing.T) {
	config := DefaultConfig("")
	config.BaseURL = ""
	client := NewClientWithConfig(config)

	ctx := context.Background()
	req := FileRequest{
		FilePath: "some non existent file path/F616FD18-589E-44A8-BF0C-891EAE69C455",
		Purpose:  PurposeAssistants,
	}

	_, err := client.CreateFile(ctx, req)
	checks.ErrorIs(t, err, os.ErrNotExist, "CreateFile should return error if file does not exist")
}

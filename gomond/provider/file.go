package provider

import (
	"github.com/hpcloud/tail"
	"github.com/juju/errors"
	"io"
	"os"
)

type FileOption struct {
	Height   int    `json:"height"`
	FileName string `json:"file_name"`
}

type FileProvider struct {
	config   FileOption
	lastLine int
	tailFile *tail.Tail
}

func NewFileProvider(config FileOption) *FileProvider {
	return &FileProvider{config: config}
}

func (f *FileProvider) Follow(out chan []byte) error {
	batch := f.config.Height

	defer f.Close()

	record := make([]byte, 0)

	for line := range f.tailFile.Lines {

		if line.Text == "" {
			continue
		}

		batch--
		if batch == 0 {
			record = append(record, []byte(line.Text)...)
			record = append(record, []byte("\n")...)
			out <- record
			batch = f.config.Height
			record = make([]byte, 0)
		} else {
			record = append(record, []byte(line.Text)...)
			record = append(record, []byte("\n")...)
		}

	}

	return nil
}

func (f *FileProvider) Start() error {
	open, err := os.Open(f.config.FileName)

	if err != nil {
		return errors.Annotate(err, "FileProvider file.Open")
	}

	seek, _ := open.Seek(0, io.SeekEnd)

	err = open.Close()

	if err != nil {
		return errors.Annotate(err, "FileProvider file.Close")
	}

	t, err := tail.TailFile(f.config.FileName, tail.Config{
		ReOpen: false,
		Location: &tail.SeekInfo{
			Offset: seek,
			Whence: io.SeekStart,
		},
		MustExist: false,
		Poll:      false,
		Pipe:      false,
		Follow:    true,
		Logger:    nil,
	})
	f.tailFile = t
	return err
}

func (f *FileProvider) Close() error {
	return f.tailFile.Stop()
}

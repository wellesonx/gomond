package provider

import (
	"github.com/hpcloud/tail"
	"github.com/juju/errors"
	"io"
	"os"
)

type FileConfiguration struct {
	Height   int    `json:"height"`
	FileName string `json:"file_name"`
}

type FileProvider struct {
	config   FileConfiguration
	lastLine int
	tailFile *tail.Tail
}

func NewFileProvider(config FileConfiguration) *FileProvider {
	return &FileProvider{config: config, lastLine: 99969910}
}

func (f *FileProvider) Follow(out chan []byte) error {
	batch := f.config.Height

	defer f.Close()

	record := make([]byte, 0)

	for line := range f.tailFile.Lines {
		record = append(record, []byte(line.Text)...)

		if batch == 0 {
			out <- record
			batch = f.config.Height
			record = make([]byte, 0)
		} else {
			record = append(record, []byte("\n")...)
		}
		batch--

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

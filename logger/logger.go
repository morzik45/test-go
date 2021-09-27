package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type RotateWriter struct {
	lock       sync.Mutex
	dateString string
	fp         *os.File
}

var INFO *log.Logger
var ERROR *log.Logger

func init() {
	writer := New()
	INFO = log.New(writer, "INFO:\t", log.Ldate|log.Ltime|log.Lshortfile)
	ERROR = log.New(writer, "ERROR:\t", log.Ldate|log.Ltime|log.Lshortfile)
}

func New() *RotateWriter {
	w := &RotateWriter{}
	w.dateString = time.Now().UTC().Format("2006-01-02")
	err := w.Rotate()
	if err != nil {
		return nil
	}
	return w
}

func (w *RotateWriter) Write(output []byte) (int, error) {
	w.lock.Lock()
	defer w.lock.Unlock()
	currentDateString := time.Now().UTC().Format("2006-01-02")
	if w.dateString != currentDateString {
		INFO.Println("change log date to", currentDateString)
		w.dateString = currentDateString
		w.Rotate()
	}
	return w.fp.Write(output)
}

func (w *RotateWriter) Rotate() (err error) {
	w.lock.Lock()
	defer w.lock.Unlock()
	if w.fp != nil {
		err = w.fp.Close()
		w.fp = nil
		if err != nil {
			return
		}
	}
	absPath, err := filepath.Abs("./log")
	if err != nil {
		fmt.Println("Error reading given path:", err)
	}
	filename := fmt.Sprintf("%s/log_%s.log", absPath, w.dateString)
	w.fp, err = os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Println("Error opening log file", err)
		os.Exit(1)
	}
	return
}

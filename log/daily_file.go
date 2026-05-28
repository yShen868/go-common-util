package log

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// dailyFileWriter 按自然日写入 {dir}/YYYYMMDD.log，跨日自动切换文件。
type dailyFileWriter struct {
	dir     string
	mu      sync.Mutex
	curDate string
	file    *os.File
}

func newDailyFileWriter(dir string) *dailyFileWriter {
	return &dailyFileWriter{dir: dir}
}

func (w *dailyFileWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if err := w.ensureFile(); err != nil {
		return 0, err
	}
	return w.file.Write(p)
}

func (w *dailyFileWriter) Sync() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.file == nil {
		return nil
	}
	return w.file.Sync()
}

func (w *dailyFileWriter) ensureFile() error {
	date := time.Now().Format("20060102")
	if date == w.curDate && w.file != nil {
		return nil
	}
	if w.file != nil {
		_ = w.file.Close()
		w.file = nil
	}
	if err := os.MkdirAll(w.dir, 0o755); err != nil {
		return fmt.Errorf("mkdir log dir %s: %w", w.dir, err)
	}
	path := filepath.Join(w.dir, date+".log")
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("open log file %s: %w", path, err)
	}
	w.file = f
	w.curDate = date
	return nil
}

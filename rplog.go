package rplog

import (
	"log"
	"os"
)

const (
	Unset = iota
	Error
	Warn
	Info
)

const (
	File = iota
	StdOut
	StdErr
)

const (
	v = iota
	vv
	vvv
)

const (
	chanSize = 100
)

type Config struct {
	Mode      int
	File      string
	Verbosity int
}

type Record struct {
	Message string
	Level   int
}

type loger struct {
	chanel chan *Record
	config *Config
	file   *os.File
}

var Loger loger

func (l *loger) Start(c *Config) error {
	if err := l.configure(c); err != nil {
		return err
	}
	l.chanel = make(chan *Record, chanSize)

	go func() {
		for r := range l.chanel {
			l.write(r)
		}
	}()

	return nil
}

func (l *loger) configure(c *Config) error {
	l.config = c
	if l.config.Mode != File {
		return nil
	}
	f, err := os.OpenFile(l.config.File, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error open log file: %v", err)
	}
	l.file = f
	return err
}

func (l *loger) Stop() {
	l.file.Close()
}

func (l *loger) Write(r *Record) {
	l.chanel <- r
}

func (l *loger) write(r *Record) {
	switch l.config.Mode {
	case File:
		l.writeToFile(r)
	case StdOut:
	case StdErr:
	default:
	}
}

func (l *loger) writeToFile(r *Record) {
	log.SetOutput(l.file)
	log.Println(r.Message)
}

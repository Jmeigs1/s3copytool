package main

import (
	"fmt"
	"io"
	"sync/atomic"
)

type progressWriter struct {
	written int64
	writer  io.WriterAt
	size    int64
}

func (pw *progressWriter) WriteAt(p []byte, off int64) (int, error) {
	atomic.AddInt64(&pw.written, int64(len(p)))

	percentageDownloaded := float32(pw.written*100) / float32(pw.size)

	fmt.Printf("File size:%s downloaded:%s percentage:%.2f%%\r", byteCountDecimal(pw.size), byteCountDecimal(pw.written), percentageDownloaded)

	return pw.writer.WriteAt(p, off)
}
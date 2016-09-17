package main

import (
	"log"
	"os"
	"errors"
	"github.com/eugene-eeo/psync/blockfs"
	"net/http"
)

type Job struct {
	Host     string
	Checksum blockfs.Checksum
}

type Response struct {
	Job   *Job
	Error error
}

func getTask(fs *blockfs.FS, job *Job) error {
	url := job.Host + "/" + string(job.Checksum)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return errors.New("GET " + url + ": " + string(resp.StatusCode))
	}
	b := make([]byte, blockfs.BlockSize)
	length, err := resp.Body.Read(b)
	block := blockfs.NewBlock(b[:length])
	if block.Checksum != job.Checksum {
		return errors.New("invalid checksum")
	}
	err = fs.WriteBlock(block)
	return err
}

func produceResponses(fs *blockfs.FS, jobs <-chan *Job, dst chan<- *Response) {
	for job := range jobs {
		err := getTask(fs, job)
		dst <- &Response{
			Job: job,
			Error: err,
		}
	}
}

func Get(fs *blockfs.FS, addr string, force bool) {
	logger := log.New(os.Stderr, "", log.Ltime)
	hashlist, err := blockfs.NewHashList(os.Stdin)
	checkErr(err)
	if !force {
		hashlist = fs.MissingBlocks(hashlist)
	}

	jobs := make(chan *Job, 4)
	rets := make(chan *Response, 4)
	done := make(chan bool)

	for i := 0; i < 4; i++ {
		go produceResponses(fs, jobs, rets)
	}

	for _, checksum := range hashlist {
		jobs <- &Job{
			Checksum: checksum,
			Host: addr,
		}
	}
	close(jobs)
	go func() {
		failed := false
		for i := 0; i < len(hashlist); i++ {
			res := <-rets
			if res.Error != nil {
				logger.Println("error:", res.Error)
				failed = true
				continue
			}
			logger.Println("ok:", res.Job.Checksum)
		}
		done <- failed
	}()
	ok := <-done
	if ok {
		os.Exit(1)
	}
}

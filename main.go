package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

var jsonfile = "./repo.json"
var binaryfile = "./repo.bin"

func main() {

	c := make(chan Commit)
	go CreateRepo(c)

	start := time.Now()
	writeJSON(c)
	elapsed := time.Since(start)
	log.Printf("Writing JSON took %s", elapsed)

	start = time.Now()
	readJSON()
	elapsed = time.Since(start)
	log.Printf("Reading JSON took %s", elapsed)
	log.Printf("JSON file size: %s", fileSize(jsonfile))

	c = make(chan Commit)
	go CreateRepo(c)

	start = time.Now()
	writeBinary(c)
	elapsed = time.Since(start)
	log.Printf("Writing binary took %s", elapsed)

	start = time.Now()
	readBinary()
	elapsed = time.Since(start)
	log.Printf("Reading binary took %s", elapsed)
	log.Printf("Binary file size: %s", fileSize(binaryfile))
}

func writeJSON(c <-chan Commit) {
	file, _ := os.Create(jsonfile)
	defer file.Close()

	writer := bufio.NewWriter(file)

	for commit := range c {
		c, _ := json.Marshal(commit)
		writer.WriteString(string(c) + "\n")
	}
}

func writeBinary(c <-chan Commit) {
	file, _ := os.Create(binaryfile)
	defer file.Close()

	writer := bufio.NewWriter(file)

	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)

	for commit := range c {
		enc.Encode(commit)
		binary.Write(writer, binary.LittleEndian, buf.Bytes())
	}
}

func readJSON() {
	file, _ := os.Open(jsonfile)
	defer file.Close()

	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')
		if err != nil || err == io.EOF {
			break
		}
		var commit Commit
		json.Unmarshal([]byte(line), &commit)
	}
}

func readBinary() {
	file, _ := os.Open(binaryfile)
	defer file.Close()

	reader := bufio.NewReader(file)
	dec := gob.NewDecoder(reader)

	for {
		var commit Commit
		err := dec.Decode(&commit)
		if err == io.EOF {
			break
		}
		// do something with commit
	}
}

func fileSize(filename string) string {
	f, _ := os.Stat(filename)
	if f.Size() < 1000000 {
		return fmt.Sprintf("%dKB", f.Size()/1024)
	}
	return fmt.Sprintf("%dMB", f.Size()/(1024*1024))
}

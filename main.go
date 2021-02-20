package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"encoding/json"
	"fmt"
	commit "github.com/alimgiray/repo_info_poc/proto"
	"github.com/golang/protobuf/proto"
	"io"
	"log"
	"os"
	"time"
)

const jsonfile = "./repo.json"
const binaryfile = "./repo.bin"
const protofile = "./repo.protobinary"

func main() {

	// JSON
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

	// Binary
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

	// Proto
	c = make(chan Commit)
	go CreateRepo(c)

	start = time.Now()
	writeProto(c)
	elapsed = time.Since(start)
	log.Printf("Writing proto took %s", elapsed)

	start = time.Now()
	readProto()
	elapsed = time.Since(start)
	log.Printf("Reading proto took %s", elapsed)
	log.Printf("Proto file size: %s", fileSize(protofile))
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

func writeProto(c <-chan Commit) {
	file, _ := os.Create(protofile)
	defer file.Close()

	b := make([]byte, 4)
	for commit := range c {
		pb, _ := proto.Marshal(ProtoFromCommit(commit))

		binary.LittleEndian.PutUint32(b, uint32(len(pb)))

		file.Write(b)
		file.Write(pb)
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
		var c Commit
		json.Unmarshal([]byte(line), &c)
	}
}

func readBinary() {
	file, _ := os.Open(binaryfile)
	defer file.Close()

	reader := bufio.NewReader(file)
	dec := gob.NewDecoder(reader)

	for {
		var c Commit
		err := dec.Decode(&c)
		if err == io.EOF {
			break
		}
		// do something with commit
	}
}

func readProto() {
	file, _ := os.Open(protofile)
	defer file.Close()

	size := make([]byte, 4)
	for {
		var c Commit

		_, err := file.Read(size)
		if err == io.EOF {
			break
		}

		message := make([]byte, binary.LittleEndian.Uint32(size))
		_, _ = file.Read(message)

		pb := &commit.ProtoCommit{}
		_ = proto.Unmarshal(message, pb)

		c = CommitFromProto(pb)
		_ = c
	}
}

func fileSize(filename string) string {
	f, _ := os.Stat(filename)
	if f.Size() < 1000000 {
		return fmt.Sprintf("%dKB", f.Size()/1024)
	}
	return fmt.Sprintf("%dMB", f.Size()/(1024*1024))
}

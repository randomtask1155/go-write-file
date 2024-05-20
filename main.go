package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"time"
)

var (
	fsize       = 10   // mb
	bsize       = 3884 // kb
	filePath    = "go-write-file-test.txt"
	interval    time.Duration // second
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ01234567890~!@#$%^&*()_+`-=[];',./{}|:<>?/")
	goWriteFile = true
)

func RandStringRunes(n int) byte {
	return byte(letterRunes[rand.Intn(len(letterRunes))])
}

func generateFile(f string) error {

	fh, err := os.OpenFile(f, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer fh.Close()

	l := (fsize * 1024 * 1024) / bsize

	for i := 0; i < l; i++ {
		var chunk []byte
		for r := 0; r <= bsize; r++ {
			chunk = append(chunk, RandStringRunes(1))
		}
		chunk = append(chunk, '\n')
		_, err := fh.Write(chunk)
		if err != nil {
			return err
		}
	}
	return nil
}

// setIntVar will overwrite int v with environment setting if exists
func setIntVar(e string, v *int) {
	s := os.Getenv(e)
	if s != "" {
		i, err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
		*v = i
	}
}

// setStringVar will overwrite string v with environment setting if exists
func setStringVar(e string, v *string) {
	s := os.Getenv(e)
	if s != "" {
		*v = s
	}
}

func main() {
	interval = 10 * time.Second

	setIntVar("FILE_SIZE", &fsize)
	setIntVar("BLOCK_SIZE", &bsize)
	setStringVar("FILE_PATH", &filePath)

	for {
		if goWriteFile {

			// delete file if exists
			if _, err := os.Stat(filePath); err == nil {
				err := os.Remove(filePath)
				if err != nil {
					fmt.Printf("failed to delete file: %s\n", err)
				}
			}

			err := generateFile(filePath)
			if err != nil {
				fmt.Println(err)
			}

			b, err := exec.Command("grep", "-Pa", `\x00`, filePath).CombinedOutput()
			if err == nil { // grep returns 0 we found null bytes
				fmt.Println("GREP FOUND NULL BYTES")
				goWriteFile = false
			} else {
				if len(b) > 0 {
					fmt.Printf("%s\n", b)
				}
				fmt.Printf("no null bytes found: grep: %s\n", err)
			}
		}
		time.Sleep(interval)

	}

}

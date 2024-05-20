package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"
)

var (
	fsize       = 10   // mb
	bsize       = 3884 // kb
	batchSize   = 2
	filePath    = "go-write-file-test.txt"
	interval    time.Duration // second
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ01234567890~!@#$%^&*()_+`-=[];',./{}|:<>?/")
	goWriteFile = true
	wg          *sync.WaitGroup
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

// scanForNull return true if null found
func scanForNull(f string) bool {

	b, err := exec.Command("grep", "-Pa", `\x00`, f).CombinedOutput()
	if err == nil { // grep returns 0 we found null bytes
		fmt.Printf("GREP FOUND NULL BYTES in file %s\n", f)
		return true
	}

	// no nulls found
	if len(b) > 0 {
		fmt.Printf("%s\n", b)
	}
	fmt.Printf("no null bytes found in file %s: grep: %s\n", f, err)
	return false
}

func runWorkload(f string) {
	defer wg.Done()
	// delete file if exists and does not contain nulls
	if _, err := os.Stat(f); err == nil {
		if scanForNull(f) {
			fmt.Printf("DETECTED NULL BYTES IN FILE %s\n", f)
			goWriteFile = false
			return
		}
		err := os.Remove(f)
		if err != nil {
			fmt.Printf("failed to delete file: %s\n", err)
		}
	}

	err := generateFile(f)
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	interval = 10 * time.Second

	setIntVar("FILE_SIZE", &fsize)
	setIntVar("BLOCK_SIZE", &bsize)
	setIntVar("BATCH_SIZE", &batchSize)
	setStringVar("FILE_PATH", &filePath)

	wg = new(sync.WaitGroup)
	for {
		if goWriteFile {
			for i := 0; i < batchSize; i++ {
				wg.Add(1)
				go runWorkload(fmt.Sprintf("%s-%d", filePath, i))

			}
		}
		wg.Wait()
		time.Sleep(interval)
	}

}

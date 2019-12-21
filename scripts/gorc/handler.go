package gorc

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"
)

var group sync.WaitGroup
var pi chan string = make(chan string, 2)
var exit chan bool = make(chan bool, 1)
var exitSub chan bool = make(chan bool, 1)

func Download(url string) (err error) {
	flag.Parse()
	fl := assign(url)
	if !fl {
		return errors.New("")
	}
	go removeCache()
	log.Println("start download")
	previous := time.Now()
	for key, meta := range Context.fileNames {
		if checkBlockStat(key, meta) {
			continue
		}
		log.Println("file", key, "start", meta.end-meta.start+1)
		group.Add(1)
		go goBT(Context.file.url, key, meta)
	}
	time.Sleep(2 * time.Second)
	goBar(Context.file.length, previous)
	group.Wait()
	//log.Println("start unzip")
	err = createFileOnly(Context.file.filePath)
	if err != nil {
		log.Println(err.Error())
		panic(err)
	}

	for i := len(Context.tempList) - 1; i >= 0; i-- {
		err = appendToFile(Context.file.filePath, string(readFile(Context.tempList[i])))
		if err != nil {
			log.Println(err.Error(), "download request failed,please retry")
			return
		}
		if i == 0 {
			exit <- true
		}
	}

	flag := <-exit
	if flag {
		for _, file := range Context.tempList {
			deleteFile(file)
		}
		log.Println("download completed")
		return
	}
	log.Println("download request failed,please retry")
	return
}
func goBT(url string, address string, b *block) {
	l, err := sendGet(url, address, b.start, b.end)
	if err != nil || l != (b.end-b.start+1) {
		log.Println("下载重试中")
		if b.count > attempt {
			pi <- b.id
			err = nil
		}
		if b.count <= attempt {
			b.count++
			goBT(url, address, b)
		}
	}
	if err == nil {
		group.Done()
	}
}
func removeCache() {
	for {
		select {
		case str := <-pi:
			p := filePath(str)
			deleteFile(p)
			exit <- false
			exitSub <- false
		case <-exitSub:
			break
		}
	}
}
func bar(count, size int) string {
	str := ""
	for i := 0; i < size; i++ {
		if i < count {
			str += "="
		} else {
			str += " "
		}
	}
	return str
}

func goBar(length int64, t time.Time) {
	for {
		var sum int64 = 0
		for key, _ := range Context.fileNames {
			sum += getFileSize(key)
		}
		percent := getPercent(sum, length)
		result, _ := strconv.Atoi(percent)
		str := "working " + percent + "%" + "[" + bar(result, 100) + "] " + " " + fmt.Sprintf("%.f", getCurrentSize(t)) + "s"
		fmt.Printf("\r%s", str)
		time.Sleep(1 * time.Second)
		if sum == length {
			fmt.Println("")
			break
		}
	}
}
func getPercent(a int64, b int64) string {
	result := float64(a) / float64(b) * 100
	return fmt.Sprintf("%.f", result)
}
func getCurrentSize(t time.Time) float64 {
	return time.Now().Sub(t).Seconds()
}

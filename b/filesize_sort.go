package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

var ch = make(chan int, 0)
var monitorStats = &Monitor{}
var f *os.File

func main() {

	dir := ``
	log := ""

	flag.StringVar(&dir, "dir", "", "姓名")
	flag.StringVar(&log, "log", "", "日志文件绝对路径")
	//解析命令行参数
	flag.Parse()

	fmt.Println(dir, log)
	if log == "" {
		panic("log empty")
	}

	var err error
	f, err = os.OpenFile(log, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
	if err != nil {
		panic(err)
	}

	fmt.Println(dir)
	go monitor()
	test1(dir)
}

func test1(dir string) {
	if dir == "" {
		panic("dir empty")
	}
	monitorStats.Dir = dir
	monitorStats.StartTime = time.Now()
	var data = make(map[string]int64) // map[filepath]size
	i := 0

	//	filepath.Walk(`c:\phpStudy`, func(pathStr string, info os.FileInfo, err error) error {
	//	filepath.Walk(`C:\phpStudy\PHPTutorial`, func(pathStr string, info os.FileInfo, err error) error {
	filepath.Walk(dir, func(pathStr string, info os.FileInfo, err error) error {
		pathStr = filepath.ToSlash(pathStr)
		//		fmt.Println(pathStr, info)
		if info == nil {
			os.Stderr.WriteString("info is nil, pathStr=" + pathStr + "\n")
			return nil
		}
		if info.IsDir() {
			pathStr = getPath(pathStr)
			data[pathStr] = info.Size()
		}
		//		i++
		if i > 100 {
			return fmt.Errorf("over")
		}

		return nil
	})

	monitorStats.PathCount = len(data)
	monitorStats.LeaveCalcPathCount = len(data)
	for pathStr, _ := range data {
		data[pathStr] = getPathSize(pathStr)
		monitorStats.HadCalcPathCount++
		monitorStats.LeaveCalcPathCount--
	}

	sortList := sortMapByValue(data)

	for _, v := range sortList {
		f.WriteString(fmt.Sprint(v.Key, " ", formatFileSize(v.Value), "\n"))
	}

}

type Monitor struct {
	Dir                string //初始目录
	PathCount          int    // walk过滤后计算出的要计算文件夹大小的目录总数量
	HadCalcPathCount   int    //已经计算完成的目录总数
	LeaveCalcPathCount int    //剩余待计算的目录总数
	StartTime          time.Time
	UseTimes           string //程序执行时间
}

func monitor() {
	for {
		monitorStats.UseTimes = time.Since(monitorStats.StartTime).String()
		fmt.Printf("monitor: %+v\n", *monitorStats)
		time.Sleep(time.Second * 5)
	}
}

func getPath(pathStr string) string {
	for strings.Count(pathStr, "/") > 3 {
		//		fmt.Println("before", pathStr)
		pathStr = filepath.Dir(pathStr)
		pathStr = filepath.ToSlash(pathStr)
		//		fmt.Println("after", pathStr)
	}
	return pathStr
}

func getPathSize(pathStr string) int64 {
	//	pathStr = filepath.FromSlash(pathStr)
	fmt.Println("getPathSize: ", pathStr)
	//	cmd := exec.Command(`C:/MinGW/msys/bin/du.exe `, "-s", pathStr)

	//powershell -noprofile -command "ls -r e:/i4Tools7 | measure -s Length"
	cmd := exec.Command(`C:\Windows\SysWOW64\WindowsPowerShell\v1.0\powershell.exe `,
		"-noprofile", "-command", "ls -r", pathStr, "| measure -s Length")

	//	cmd.Stdin = strings.NewReader("some input")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	//	fmt.Println("outString: ", out.String())

	if err != nil {
		fmt.Println("err: ", err, " pathStr: ", pathStr)
		//		log.Fatal(pathStr, err)
	}
	//in all caps: "9149\tC:/phpStudy/PHPTutorial/Apache\n"
	//	fmt.Printf("in all caps: %q\n", out.String())
	outStr := out.String()

	reg := regexp.MustCompile(`Sum      : ([0-9]+)`)
	s := reg.FindAllStringSubmatch(outStr, 1)
	if s == nil || len(s) == 0 {
		return 0
	}
	fmt.Println("readline: ", s[0][1])
	outStr = strings.TrimSpace(s[0][1])

	size, err := strconv.ParseInt(outStr, 10, 64)
	if err != nil {
		return 0
		//todo
		panic(err)
	}
	return size
}

type Pair struct {
	Key   string
	Value int64
}

// A slice of Pairs that implements sort.Interface to sort by Value.
type PairList []Pair

func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }

// A function to turn a map into a PairList, then sort and return it.
func sortMapByValue(m map[string]int64) PairList {

	p := make(PairList, len(m))
	i := 0
	for k, v := range m {
		p[i] = Pair{Key: k, Value: v}
		i++
	}

	sort.Sort(sort.Reverse(p))
	return p
}

func formatFileSize(fileSize int64) (size string) {
	if fileSize < 1024 {
		//return strconv.FormatInt(fileSize, 10) + "B"
		return fmt.Sprintf("%.2fB", float64(fileSize)/float64(1))
	} else if fileSize < (1024 * 1024) {
		return fmt.Sprintf("%.2fKB", float64(fileSize)/float64(1024))
	} else if fileSize < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fMB", float64(fileSize)/float64(1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fGB", float64(fileSize)/float64(1024*1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fTB", float64(fileSize)/float64(1024*1024*1024*1024))
	} else { //if fileSize < (1024 * 1024 * 1024 * 1024 * 1024 * 1024)
		return fmt.Sprintf("%.2fEB", float64(fileSize)/float64(1024*1024*1024*1024*1024))
	}
}

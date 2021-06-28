package main

import (
	//"encoding/hex"
	"flag"
	"fmt"
	"github.com/vua/vfmt"
	"math"
	//"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

var color = [4]string{"#4B0082", "#006400", "#FF1493", "#808080"}

type entry struct {
	name    string
	cnt     int
	size    int64
	color   string
	display bool
}

func main() {
	var target string
	flag.StringVar(&target, "h", "/", "du.exe -h target-dir")
	flag.Parse()
	dirs, err := os.ReadDir(target)
	if err != nil {
		panic(err)
	}
	entrys := make([]entry, len(dirs))
	nameLenMax := 0
	for i, dir := range dirs {
		info, _ := dir.Info()
		var size int64
		var name string
		if info.IsDir() {
			size, _ = DirSize(filepath.Join(target, info.Name()))
			name = info.Name() + "/"
		} else {
			size = info.Size()
			name = info.Name()
		}
		entrys[i].size = size
		entrys[i].name = name
	}
	sort.Slice(entrys, func(i, j int) bool {
		aSize, bSize := entrys[i].size, entrys[j].size
		if aSize < bSize {
			return true
		}
		return false
	})

	for len(entrys) > 0 {
		if entrys[0].size == 0 {
			entrys = entrys[1:]
		} else {
			break
		}
	}

	var baseSize int64
	cntSum := 0
	for i, entry := range entrys {
		entrys[i].color = color[i%4]
		size, name := entry.size, entry.name
		if i == 0 {
			baseSize = size
			entrys[i].cnt = 1
		} else {
			entrys[i].cnt = int(math.Log10(float64(size/baseSize))) + 1
		}
		entrys[i].name = name + "(" + strconv.FormatFloat(float64(size)/1024.0/1024.0,'f',4,32) + "MB)"
		if len(entrys[i].name) > nameLenMax {
			nameLenMax = len(entrys[i].name)
		}
		cntSum += entrys[i].cnt
	}

	// 组数上限
	gNum := 150 / nameLenMax
	// 分组
	group := make([][]entry, 0)
	l := len(entrys)
	tmp := 0
	k := 0
	for ; k < l; k++ {
		tmp += entrys[l-1-k].cnt
		if cntSum/tmp < gNum {
			break
		}
	}
	group = append(group, entrys[l-1-k:])
	row := tmp
	prevCntSum := tmp
	cnt := 0
	r := l - 2 - k
	for i := l - 2 - k; i >= -1; i-- {
		if i == -1 || cnt+entrys[i].cnt > prevCntSum+prevCntSum/8 {
			group = append(group, entrys[i+1:r+1])
			r = i
			prevCntSum = cnt
			cnt = 0
			if i == -1 {
				break
			}
		}
		cnt += entrys[i].cnt
	}
	cnts := make([]int, l)
	format := "%-" + strconv.Itoa(nameLenMax) + "s"
	common := fmt.Sprintf(format, "")
	for i := 0; i < row; i++ {
		for k, g := range group {
			if cnts[k] == len(g) {
				fmt.Printf(common)
				continue
			}
			if !g[cnts[k]].display {
				tmp := g[cnts[k]].name
				format := "%-" + strconv.Itoa(nameLenMax-(len(tmp)-len([]rune(tmp)))/2) + "s"
				name := fmt.Sprintf(format, g[cnts[k]].name)
				vfmt.Printf("@[%s::bg%s|bold]", name, g[cnts[k]].color)
				g[cnts[k]].display = true
			} else {
				vfmt.Printf("@[%s::bg%s]", common, g[cnts[k]].color)
			}
			g[cnts[k]].cnt -= 1
			if g[cnts[k]].cnt == 0 {
				cnts[k] += 1
			}
		}
		fmt.Println()
	}
}

//func color() string {
//	rand.Seed(rand.Int63())
//	data := make([]byte, 3)
//	for i := 0; i < 3; i++ {
//		data[i] = byte(rand.Intn(100))
//	}
//	return hex.EncodeToString(data)
//}
func DirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

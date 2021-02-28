package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/oklog/ulid"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

type (
	Shuffler struct {
		ImportPath string
		Help bool
		Count int
		Name string
	}
	MetaData struct {
		FileName string
		Fixed []int
		Column []string
		Data [][]string
	}
)

const helpMessage = `help message`

func main() {
	var shuffler Shuffler
	var metaData MetaData
	var importPath string
	var name string
	flag.StringVar(&importPath, "p", "", "exe mode")
	flag.StringVar(&name, "n", "", "exe mode")
	flag.IntVar(&shuffler.Count, "c", 0, "row count")
	flag.BoolVar(&shuffler.Help, "h", false, "help mode")
	flag.Parse()

	if shuffler.Help {
		fmt.Println(helpMessage)
		return
	}

	if strings.HasSuffix(importPath, ".csv") {
		shuffler.ImportPath = strings.TrimRight(importPath, ".csv")
	} else {
		log.Fatal(errors.New("please select a csv file"))
	}

	if name != "" {
		shuffler.Name = strings.Split(name, ".")[0]
	}

	if err := shuffler.setMeta(&metaData); err != nil {
		log.Fatal(err)
	}
	if err := metaData.Run(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("create: ", metaData.FileName)
	fmt.Print("\n...complete")
	return
}

func (s *Shuffler) setMeta(data *MetaData) error {
	if s.ImportPath == "" {
		return errors.New("file not found")
	}
	// fileOpen Reader作成
	file, err := os.Open(s.ImportPath+".csv")
	if err != nil || file == nil {
		return err
	}
	defer func() {
		if file != nil {
			if err := file.Close(); err != nil {
				log.Fatal(err)
			}
		}
	}()
	reader := csv.NewReader(file)
	lineAll, err := reader.ReadAll()
	if err != nil {
		return err
	}

	var fName string

    // fileName生成
    if s.Name == "" {
		t := time.Now()
		entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
		fName = s.ImportPath+"_"+ulid.MustNew(ulid.Timestamp(t), entropy).String()+".csv"
	} else {
		fName = s.Name+".csv"
	}

	// metaData作成
	if data.Fixed, err = s.conversation(lineAll[0]); err != nil {
		return err
	}
	data.FileName = fName
	data.Column = lineAll[0]
	if s.Count != 0 {
		if s.Count > len(lineAll) {
			data.Data = compensate(s.Count - len(lineAll)+1, lineAll[1:])
		} else {
			data.Data = lineAll[1:s.Count+1]
		}
	} else {
		data.Data = lineAll[1:]
	}

	return nil
}

func compensate(diff int, lines [][]string) [][]string {
	var binder [][]string
	for i:=0; i < diff; i++ {
		binder = append(binder, lines[rand.Intn(len(lines))])
	}
	return append(lines, binder...)
}

func (s *Shuffler) conversation(column []string) ([]int, error) {
	var fixed []int
	var fixedStr []string
	prompt := &survey.MultiSelect{
		Message: "Select any columns that should be fixed, please",
		Options: column,
	}
	if err := survey.AskOne(prompt, &fixedStr); err != nil {
		return nil, err
	}

	for _, str := range fixedStr {
		for index, columnName := range column {
			if str == columnName { fixed = append(fixed, index) }
		}
	}

	return fixed, nil
}

func (m *MetaData) Run() error {
	shuffle(m.Data, m.Fixed, m.Column)

	records := [][]string {m.Column}

	records = append(records, m.Data...)

	f, err := os.Create(m.FileName)
	if err != nil {
		return err
	}
	defer func() {
		if f != nil {
			if err := f.Close(); err != nil {
				log.Fatal(err)
			}
		}
	}()

	// writer作成
	w := csv.NewWriter(f)
	err = w.WriteAll(records)
	if err != nil {
		return err
	}
	w.Flush()

	return nil
}

// ここから重そう
func shuffle(data [][]string, fixed []int, columnCnt []string) {
	n := len(data)
	for i := n - 1; i >= 0; i-- {
		j := rand.Intn(i + 1)
		for index := range columnCnt {
			if contains(fixed, index) { continue }
			data[i][index], data[j][index] = data[j][index], data[i][index]
		}
	}
}

func contains(s []int, e int) bool {
	for _, v := range s {
		if e == v {
			return true
		}
	}
	return false
}
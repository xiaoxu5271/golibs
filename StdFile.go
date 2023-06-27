package golibs

import (
	"os"
	"path/filepath"
	"strings"
)

type StdFile struct {
	myDir     string // 程序运行目录
	parentDir string // 程序运行目录的父目录
	myName    string // 程序名称
}

func (s *StdFile) Init() *StdFile {
	var e error = nil
	run := filepath.Dir(os.Args[0])
	s.myDir, e = filepath.Abs(run)
	if nil != e {
		panic(e)
	}

	spt := strings.Split(os.Args[0], "/")
	s.myName = spt[len(spt)-1]

	s.parentDir = filepath.Dir(s.myDir)
	return s
}

func (s *StdFile) GetMyName() string {
	return s.myName
}

func (s *StdFile) GetMyDir() string {
	return s.myDir
}

func (s *StdFile) GetParentDir() string {
	return s.parentDir
}

// equipment_manage

func (s *StdFile) GetEquCfgFile(name string) string {
	return filepath.Join(s.parentDir, "cfg", name)
}

func (s *StdFile) GetEquTmpFile(name string) string {
	return filepath.Join(s.parentDir, "tmp", name)
}

func (s *StdFile) GetEquipmentID() string {
	prun := strings.Split(s.myDir, "/")
	for i := len(prun) - 1; i > 0; i-- {
		// app id and device id is the father of app
		if prun[i] == "app" && i > 0 {
			return prun[i-1]
		}
	}
	return ""
}

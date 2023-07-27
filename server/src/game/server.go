package game

import (
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"server_logic/csvs"
	"sync"
	"syscall"
)

type Server struct {
	Wait        sync.WaitGroup
	BanWordBase []string //配置生成
	Lock        sync.RWMutex
}

var server *Server

func GetServer() *Server {
	if server == nil {
		server = new(Server)
	}
	return server
}

func (s *Server) Start() {
	csvs.CheckLoadCsv()

	//go GetManageBanWord().Run()

	playerTest := NewTestPlayer()
	go playerTest.Run()

	go s.SignalHandle()

	s.Wait.Wait()
	fmt.Println("服务器关闭")

}

func (s *Server) Stop() {
	GetManageBanWord().Close()
}

func (s *Server) AddGo() {
	s.Wait.Add(1)
}

func (s *Server) GoDone() {
	s.Wait.Done()
}

func (s *Server) IsBanWord(txt string) bool {
	s.Lock.RLock()
	defer s.Lock.RUnlock()
	for _, v := range s.BanWordBase {
		match, _ := regexp.MatchString(v, txt)
		fmt.Println(match, v)
		if match {
			return match
		}
	}
	return false
}

func (s *Server) UpBanWord(banWord []string) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	s.BanWordBase = banWord
}
func (s *Server) SignalHandle() {
	chanSignal := make(chan os.Signal)
	signal.Notify(chanSignal, syscall.SIGINT)
	for {
		select {
		case <-chanSignal:
			s.Stop()
		}
	}
}

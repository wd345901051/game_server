package game

import (
	"fmt"
	"regexp"
	"server_logic/csvs"
	"time"
)

var manageBanWord *ManageBanWord

type ManageBanWord struct {
	BanWordBase  []string //配置生成
	BanWordExtra []string // 更新
	MsgChan      chan int
}

func GetManageBanWord() *ManageBanWord {
	if manageBanWord == nil {
		manageBanWord = new(ManageBanWord)
		manageBanWord.BanWordBase = []string{"外挂", "工具"}
		manageBanWord.BanWordExtra = []string{"原神"}
		manageBanWord.MsgChan = make(chan int)
	}
	return manageBanWord
}

func (mg *ManageBanWord) IsBanWord(txt string) bool {
	for _, v := range mg.BanWordBase {
		match, _ := regexp.MatchString(v, txt)
		fmt.Println(match, v)
		if match {
			return match
		}
	}
	for _, v := range mg.BanWordExtra {
		match, _ := regexp.MatchString(v, txt)
		fmt.Println(match, v)
		if match {
			return match
		}
	}
	return false
}

func (mg *ManageBanWord) Run() {
	GetServer().AddGo()
	mg.BanWordBase = csvs.GetBanWordBase()
	// 基础词库的更新
	ticker := time.NewTicker(time.Second * 1)
	for {
		select {
		case <-ticker.C:
			if time.Now().Unix()%10 == 0 {
				fmt.Println("更新一次词库")
				GetServer().UpBanWord(mg.BanWordBase)
			}
		case _, ok := <-mg.MsgChan:
			if !ok {
				GetServer().GoDone()
				return
			}
		}
	}
}

func (mg *ManageBanWord) Close() {
	close(mg.MsgChan)
}

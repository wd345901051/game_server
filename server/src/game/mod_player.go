package game

import (
	"fmt"
	"server_logic/csvs"
	"time"
)

type ShowRole struct {
	RoleId    int
	RoleLevel int
}

type ModPlayer struct {
	UserId         int
	Icon           int
	Card           int
	Sign           string
	Name           string
	PlayerLevel    int
	PlayerExp      int   //经验
	WorldLevel     int   //世界等级
	WorldLevelNow  int   //当前世界等级
	WorldLevelCool int64 //操作世界等级冷却时间
	Birth          int
	ShowTeam       []*ShowRole //展示整容
	HideShowTeam   int         //是否展示整容 1yes:0no
	ShowCard       []int       //展示名片
	// 看不见的字段
	Prohibit int //封禁状态
	IsGM     int //GM账号标志
}

func (mp *ModPlayer) SetIcon(iconId int, player *Player) {

	if !player.ModIcon.IsHasIcon(iconId) {
		fmt.Println("没有头像", iconId)
		//通知客户端操作非法
		return
	}
	player.ModPlayer.Icon = iconId
	fmt.Println("当前图标：", player.ModPlayer.Icon)
}

func (mp *ModPlayer) SetCard(cardId int, player *Player) {

	if !player.ModCard.IsHasCard(cardId) {
		//通知客户端操作非法
		fmt.Println("该名片不存在", cardId)
		return
	}
	player.ModPlayer.Card = cardId
	fmt.Println("当前名片：", player.ModPlayer.Card)

}

func (mp *ModPlayer) SetName(name string, player *Player) {

	player.ModPlayer.Name = name
	fmt.Println("当前名称：", player.ModPlayer.Name)

}

func (mp *ModPlayer) SetSign(sign string, player *Player) {

	player.ModPlayer.Sign = sign
	fmt.Println("当前签名：", player.ModPlayer.Sign)
}

func (mp *ModPlayer) AddExp(exp int, player *Player) {
	mp.PlayerExp += exp

	for {
		config := csvs.GetNowLevelConfig(mp.PlayerLevel)
		if config == nil {
			break
		}

		if config.PlayerExp == 0 {
			break
		}

		// 是否完成任务 todo
		if config.ChapterId > 0 && !player.ModUniqueTask.IsTaskFinish(config.ChapterId) {
			break
		}
		if mp.PlayerExp >= config.PlayerExp {
			mp.PlayerLevel += 1
			mp.PlayerExp -= config.PlayerExp
		} else {
			break
		}
	}
	fmt.Println("当前等级：", mp.PlayerLevel, "---当前经验", mp.PlayerExp)
}

func (mp *ModPlayer) ReduceWorldLevel(player *Player) {
	if mp.WorldLevel < csvs.REDUCE_WORLD_LEVEL_START {
		fmt.Println("操作失败,---当前世界等级", mp.WorldLevel)
		return
	}
	if mp.WorldLevel-mp.WorldLevelNow >= csvs.REDUCE_WORLD_LEVEL_MAX {
		fmt.Println("操作失败,---当前世界等级", mp.WorldLevel, "-----真实的世界等级", mp.WorldLevelNow)
		return
	}

	if time.Now().Unix() < mp.WorldLevelCool {
		fmt.Println("操作失败,---冷却中")
		return
	}
	mp.WorldLevelNow -= 1
	mp.WorldLevelCool = time.Now().Unix() + csvs.REDUCE_WORLD_LEVEL_COOL_TIME
	fmt.Println("操作成功,---当前世界等级", mp.WorldLevel, "-----真实的世界等级", mp.WorldLevelNow)
}

func (mp *ModPlayer) ReturnWorldLevel(player *Player) {
	if mp.WorldLevelNow == mp.WorldLevel {
		fmt.Println("操作成功,---当前世界等级", mp.WorldLevel, "-----真实的世界等级", mp.WorldLevelNow)
		return
	}
	if time.Now().Unix() < mp.WorldLevelCool {
		fmt.Println("操作失败,---冷却中")
		return
	}
	mp.WorldLevelNow += 1
	mp.WorldLevelCool = 0
	fmt.Println("操作成功,---当前世界等级", mp.WorldLevel, "-----真实的世界等级", mp.WorldLevelNow)
}

func (mp *ModPlayer) SetBirth(birth int, player *Player) {
	if mp.Birth > 0 {
		fmt.Println("生日已设置")
		return
	}
	mouth := birth / 100
	day := birth % 100

	switch mouth {
	case 1, 3, 5, 7, 8, 10, 12:
		if day <= 0 || day > 31 {
			fmt.Println(mouth, "月没有", day, "日！")
			return
		}
	case 4, 6, 9, 11:
		if day <= 0 || day > 30 {
			fmt.Println(mouth, "月没有", day, "日！")
			return
		}
	case 2:
		if day <= 0 || day > 29 {
			fmt.Println(mouth, "月没有", day, "日！")
			return
		}
	default:
		fmt.Println("没有", mouth, "月！")
	}

	mp.Birth = birth
	fmt.Println("设置成功，生日为：", mouth, "月", day, "日")
	if mp.IsBirthDay() {
		fmt.Println("今天是你的生日，生日快乐!!")
	} else {
		fmt.Println("期待你生日到来")
	}
}

func (mp *ModPlayer) IsBirthDay() bool {
	month := int(time.Now().Month())
	day := time.Now().Day()
	if month == mp.Birth/100 && day == mp.Birth%100 {
		return true
	}
	return false
}

func (mp *ModPlayer) SetShowCard(showCard []int, player *Player) {
	if len(showCard) > csvs.SHOW_SIZE {
		return
	}
	cardExist := make(map[int]int)
	newList := make([]int, 0)
	for _, cardId := range showCard {
		_, ok := cardExist[cardId]
		if ok {
			continue
		}
		if !player.ModCard.IsHasCard(cardId) {
			continue
		}
		newList = append(newList, cardId)
		cardExist[cardId] = 1
	}
	mp.ShowCard = newList
	fmt.Println(mp.ShowCard)
}

// 设置展示整容
func (mp *ModPlayer) SetShowTeam(showRole []int, player *Player) {
	if len(showRole) > csvs.SHOW_SIZE {
		return
	}
	roleExist := make(map[int]int)
	newList := make([]*ShowRole, 0)
	for _, roleId := range showRole {
		_, ok := roleExist[roleId]
		if ok {
			continue
		}
		if !player.ModRole.IsHasRole(roleId) {
			continue
		}
		showRole := new(ShowRole)
		showRole.RoleId = roleId
		showRole.RoleLevel = player.ModRole.GetRoleLevel(roleId)
		newList = append(newList, showRole)
		roleExist[roleId] = 1
	}
	mp.ShowTeam = newList
	fmt.Println(mp.ShowTeam)
}

// 设置是否展示整容
func (mp *ModPlayer) SetHideShowTeam(isHide int, player *Player) {
	if isHide != csvs.LOGIC_FALSE && isHide != csvs.LOGIC_TRUE {
		return
	}
	mp.HideShowTeam = isHide
}

// 设置封禁状态
func (mp *ModPlayer) SetProhibit(prohibit int) {
	mp.Prohibit = prohibit
}

// 设置GM账号
func (mp *ModPlayer) SetIsGM(isGM int) {
	mp.IsGM = isGM
}

func (mp *ModPlayer) IsCanEnter() bool {
	return mp.Prohibit < int(time.Now().Unix())
}

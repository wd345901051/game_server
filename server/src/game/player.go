package game

import (
	"fmt"
	"server_logic/csvs"
)

type Player struct {
	ModPlayer     *ModPlayer
	ModIcon       *ModIcon
	ModCard       *ModCard
	ModUniqueTask *ModUniqueTask
	ModRole       *ModRole
	ModBag        *ModBag
	ModWeapon     *ModWeapon
	ModRelics     *ModRelics
	ModCook       *ModCook
	ModHome       *ModHome
	ModPool       *ModPool
	ModMap        *ModMap
}

const (
	TASK_STATE_INIT   = 0
	TASK_STATE_DOING  = 1
	TASK_STATE_FINISH = 2
)

func NewTestPlayer() *Player {
	player := new(Player)
	player.ModPlayer = new(ModPlayer)
	player.ModIcon = new(ModIcon)
	player.ModIcon.IconInfo = make(map[int]*Icon)
	player.ModCard = new(ModCard)
	player.ModCard.CardInfo = make(map[int]*Card)

	player.ModUniqueTask = new(ModUniqueTask)
	player.ModUniqueTask.MyTaskInfo = make(map[int]*TaskInfo)
	// player.ModUniqueTask.Locker = new(sync.RWMutex)
	player.ModBag = new(ModBag)
	player.ModBag.BagInfo = make(map[int]*ItemInfo)

	player.ModRole = new(ModRole)
	player.ModRole.RoleInfo = make(map[int]*RoleInfo)

	player.ModWeapon = new(ModWeapon)
	player.ModWeapon.WeaponInfo = make(map[int]*Weapon)

	player.ModRelics = new(ModRelics)
	player.ModRelics.RelicsInfo = make(map[int]*Relics)

	player.ModCook = new(ModCook)
	player.ModCook.CookInfo = make(map[int]*Cook)

	player.ModHome = new(ModHome)
	player.ModHome.HomeItemInfo = make(map[int]*HomeItem)

	player.ModPool = new(ModPool)
	player.ModPool.UpPoolInfo = new(PoolInfo)

	player.ModMap = new(ModMap)
	player.ModMap.InitData()

	// ************************
	player.ModPlayer.PlayerLevel = 1
	player.ModPlayer.WorldLevel = 1
	player.ModPlayer.WorldLevelNow = 1
	// ************************
	return player
}

// 对外接口

// 设置头像
func (p *Player) RecvSetIcon(iconId int) {
	p.ModPlayer.SetIcon(iconId, p)
}

// 设置名片
func (p *Player) RecvSetCard(cardId int) {
	p.ModPlayer.SetCard(cardId, p)
}

// 设置名称
func (p *Player) RecvSetName(name string) {
	if GetManageBanWord().IsBanWord(name) {
		return
	}
	p.ModPlayer.SetName(name, p)
}

// 设置名片
func (p *Player) RecvSetSign(sign string) {
	if GetManageBanWord().IsBanWord(sign) {
		return
	}
	p.ModPlayer.SetSign(sign, p)
}

// 降低世界等级
func (p *Player) ReduceWorldLevel() {
	p.ModPlayer.ReduceWorldLevel(p)
}

// 返回世界等级
func (p *Player) ReturnWorldLevel() {
	p.ModPlayer.ReturnWorldLevel(p)
}

// 设置生日
func (p *Player) SetBirth(birth int) {
	p.ModPlayer.SetBirth(birth, p)
}

// 设置生日
func (p *Player) SetShowCard(showCard []int) {
	p.ModPlayer.SetShowCard(showCard, p)
}

// 设置展示整容
func (p *Player) SetShowTeam(showRole []int) {
	p.ModPlayer.SetShowTeam(showRole, p)
}

// 设置是否展示整容
func (p *Player) SetHideShowTeam(isHide int) {
	p.ModPlayer.SetHideShowTeam(isHide, p)
}

//func (p *Player) SetEventState(state int) {
//	p.ModMap.SetEventState(state, p)
//}

func (self *Player) Run() {
	fmt.Println("从0开始写原神服务器------测试工具v0.3")
	fmt.Println("作者:B站------golang大海葵")
	fmt.Println("模拟用户创建成功OK------开始测试")
	fmt.Println("↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓")
	for {
		fmt.Println(self.ModPlayer.Name, ",欢迎来到提瓦特大陆,请选择功能：1基础信息2背包3角色(优菈UP池)模拟抽卡4地图(未开放)")
		var modChoose int
		fmt.Scan(&modChoose)
		switch modChoose {
		case 1:
			self.HandleBase()
		case 2:
			self.HandleBag()
		case 3:
			self.HandlePool()
		case 4:
			self.HandleMap()
		}
	}
}

//基础信息
func (self *Player) HandleBase() {
	for {
		fmt.Println("当前处于基础信息界面,请选择操作：0返回1查询信息2设置名字3设置签名4头像5名片6设置生日")
		var action int
		fmt.Scan(&action)
		switch action {
		case 0:
			return
		case 1:
			self.HandleBaseGetInfo()
		case 2:
			self.HandleBagSetName()
		case 3:
			self.HandleBagSetSign()
		case 4:
			self.HandleBagSetIcon()
		case 5:
			self.HandleBagSetCard()
		case 6:
			self.HandleBagSetBirth()
		}
	}
}

func (self *Player) HandleBaseGetInfo() {
	fmt.Println("名字:", self.ModPlayer.Name)
	fmt.Println("等级:", self.ModPlayer.PlayerLevel)
	fmt.Println("大世界等级:", self.ModPlayer.WorldLevelNow)
	if self.ModPlayer.Sign == "" {
		fmt.Println("签名:", "未设置")
	} else {
		fmt.Println("签名:", self.ModPlayer.Sign)
	}

	if self.ModPlayer.Icon == 0 {
		fmt.Println("头像:", "未设置")
	} else {
		fmt.Println("头像:", csvs.GetItemConfig(self.ModPlayer.Icon), self.ModPlayer.Icon)
	}

	if self.ModPlayer.Card == 0 {
		fmt.Println("名片:", "未设置")
	} else {
		fmt.Println("名片:", csvs.GetItemConfig(self.ModPlayer.Card), self.ModPlayer.Card)
	}

	if self.ModPlayer.Birth == 0 {
		fmt.Println("生日:", "未设置")
	} else {
		fmt.Println("生日:", self.ModPlayer.Birth/100, "月", self.ModPlayer.Birth%100, "日")
	}
}

func (self *Player) HandleBagSetName() {
	fmt.Println("请输入名字:")
	var name string
	fmt.Scan(&name)
	self.RecvSetName(name)
}

func (self *Player) HandleBagSetSign() {
	fmt.Println("请输入签名:")
	var sign string
	fmt.Scan(&sign)
	self.RecvSetSign(sign)
}

func (self *Player) HandleBagSetIcon() {
	for {
		fmt.Println("当前处于基础信息--头像界面,请选择操作：0返回1查询头像背包2设置头像")
		var action int
		fmt.Scan(&action)
		switch action {
		case 0:
			return
		case 1:
			self.HandleBagSetIconGetInfo()
		case 2:
			self.HandleBagSetIconSet()
		}
	}
}

func (self *Player) HandleBagSetIconGetInfo() {
	fmt.Println("当前拥有头像如下:")
	for _, v := range self.ModIcon.IconInfo {
		config := csvs.GetItemConfig(v.IconId)
		if config != nil {
			fmt.Println(config.ItemName, ":", config.ItemId)
		}
	}
}

func (self *Player) HandleBagSetIconSet() {
	fmt.Println("请输入头像id:")
	var icon int
	fmt.Scan(&icon)
	self.RecvSetIcon(icon)
}

func (self *Player) HandleBagSetCard() {
	for {
		fmt.Println("当前处于基础信息--名片界面,请选择操作：0返回1查询名片背包2设置名片")
		var action int
		fmt.Scan(&action)
		switch action {
		case 0:
			return
		case 1:
			self.HandleBagSetCardGetInfo()
		case 2:
			self.HandleBagSetCardSet()
		}
	}
}

func (self *Player) HandleBagSetCardGetInfo() {
	fmt.Println("当前拥有名片如下:")
	for _, v := range self.ModCard.CardInfo {
		config := csvs.GetItemConfig(v.CardId)
		if config != nil {
			fmt.Println(config.ItemName, ":", config.ItemId)
		}
	}
}

func (self *Player) HandleBagSetCardSet() {
	fmt.Println("请输入名片id:")
	var card int
	fmt.Scan(&card)
	self.RecvSetCard(card)
}

func (self *Player) HandleBagSetBirth() {
	if self.ModPlayer.Birth > 0 {
		fmt.Println("已设置过生日!")
		return
	}
	fmt.Println("生日只能设置一次，请慎重填写,输入月:")
	var month, day int
	fmt.Scan(&month)
	fmt.Println("请输入日:")
	fmt.Scan(&day)
	self.ModPlayer.SetBirth(month*100+day, self)
}

//背包
func (self *Player) HandleBag() {
	for {
		fmt.Println("当前处于基础信息界面,请选择操作：0返回1增加物品2扣除物品3使用物品")
		var action int
		fmt.Scan(&action)
		switch action {
		case 0:
			return
		case 1:
			self.HandleBagAddItem()
		case 2:
			self.HandleBagRemoveItem()
		case 3:
			self.HandleBagUseItem()
		}
	}
}

//抽卡
func (self *Player) HandlePool() {
	for {
		fmt.Println("当前处于模拟抽卡界面,请选择操作：0返回1十连抽2单抽(可选次数)3五星爆点测试4十连多黄测试5视频原版函数6角色英雄信息7单抽仓检")
		var action int
		fmt.Scan(&action)
		switch action {
		case 0:
			return
		case 1:
			self.ModPool.HandleUpPoolTen(self)
		case 2:
			fmt.Println("请输入抽卡次数,最大值1亿(最大耗时约30秒):")
			var times int
			fmt.Scan(&times)
			self.ModPool.HandleUpPoolSingle(times, self)
		case 3:
			fmt.Println("请输入抽卡次数,最大值1亿(最大耗时约30秒):")
			var times int
			fmt.Scan(&times)
			self.ModPool.HandleUpPoolTimesTest(times)
		case 4:
			fmt.Println("请输入抽卡次数,最大值1亿(最大耗时约30秒):")
			var times int
			fmt.Scan(&times)
			self.ModPool.HandleUpPoolFiveTest(times)
		case 5:
			self.ModPool.DoUpPool(50)
		case 6:
			self.ModRole.HandleSendRoleInfo()
		case 7:
			fmt.Println("请输入抽卡次数,最大值1亿(最大耗时约30秒):")
			var num int
			fmt.Scan(&num)
			self.ModPool.DoUpPoolCheck(num, self)
		}
	}
}

func (self *Player) HandleBagAddItem() {
	itemId := 0
	itemNum := 0
	fmt.Println("物品ID")
	fmt.Scan(&itemId)
	fmt.Println("物品数量")
	fmt.Scan(&itemNum)
	self.ModBag.AddItem(itemId, int64(itemNum), self)
}

func (self *Player) HandleBagRemoveItem() {
	itemId := 0
	itemNum := 0
	fmt.Println("物品ID")
	fmt.Scan(&itemId)
	fmt.Println("物品数量")
	fmt.Scan(&itemNum)
	self.ModBag.RemoveItemToBag(itemId, int64(itemNum), self)
}

func (self *Player) HandleBagUseItem() {
	itemId := 0
	itemNum := 0
	fmt.Println("物品ID")
	fmt.Scan(&itemId)
	fmt.Println("物品数量")
	fmt.Scan(&itemNum)
	self.ModBag.UserItem(itemId, int64(itemNum), self)
}

//地图
func (self *Player) HandleMap() {
	fmt.Println("向着星辰与深渊,欢迎来到冒险家协会！")
	fmt.Println("当前位置:", "蒙德城")
	fmt.Println("地图模块还没写到......")
}

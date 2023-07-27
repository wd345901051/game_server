package game

import (
	"fmt"
	"server_logic/csvs"
	"time"
)

type RoleInfo struct {
	RoleId   int
	GetTimes int
	// 等级经验，圣遗物
	RelicsInfo []int
	WeaponInfo int
}

type ModRole struct {
	RoleInfo  map[int]*RoleInfo
	HpPool    int
	HpCalTime int64
}

func (mr *ModRole) IsHasRole(roleId int) bool {
	return true
}
func (mr *ModRole) GetRoleLevel(roleId int) int {
	return 80
}

func (mr *ModRole) AddItem(roleId int, num int64, player *Player) {
	config := csvs.GetRoleConfig(roleId)
	if config == nil {
		fmt.Println("配置不存在roleId:", roleId)
		return
	}
	for i := 0; i < int(num); i++ {
		if _, ok := mr.RoleInfo[roleId]; !ok {
			mr.RoleInfo[roleId] = &RoleInfo{
				RoleId:   roleId,
				GetTimes: 1,
			}
			player.ModBag.AddItemToBag(roleId, 1)
		} else {
			// 判断实际获得的东西
			// 判断是否转换成其他材料
			mr.RoleInfo[roleId].GetTimes++
			if mr.RoleInfo[roleId].GetTimes <= csvs.ADD_ROLE_NORMAL_TIME_MAX && mr.RoleInfo[roleId].GetTimes >= csvs.ADD_ROLE_NORMAL_TIME_MIN {
				player.ModBag.AddItemToBag(config.Stuff, config.StuffNum)
				player.ModBag.AddItemToBag(config.StuffItem, config.StuffItemNum)
			} else {
				player.ModBag.AddItemToBag(config.MaxStuffItem, config.MaxStuffItemNum)
			}
		}
	}
	itemConfig := csvs.GetItemConfig(roleId)
	if itemConfig != nil {
		fmt.Println("获得角色", itemConfig.ItemName, "获得次数", mr.RoleInfo[roleId].GetTimes)
	}
	player.ModIcon.CheckGetIcon(roleId, player)
	player.ModCard.CheckGetCard(roleId, 10)
}

func (mr *ModRole) HandleSendRoleInfo() {
	fmt.Println("当前拥有角色信息如下")
	for _, v := range mr.RoleInfo {
		v.SendRoleInfo()
	}
}
func (ri *RoleInfo) SendRoleInfo() {
	fmt.Println(csvs.GetItemConfig(ri.RoleId).ItemName, "累计获得次数:", ri.GetTimes)
}

func (ri *ModRole) GetRoleInfoForPoolCheck() (map[int]int, map[int]int) {
	fiveInfo := make(map[int]int)
	fourInfo := make(map[int]int)

	for _, v := range ri.RoleInfo {
		roleConfig := csvs.GetRoleConfig(v.RoleId)
		if roleConfig == nil {
			continue
		}
		if roleConfig.Star == 5 {
			fiveInfo[roleConfig.RoleId] = v.GetTimes
		}
		if roleConfig.Star == 4 {
			fourInfo[roleConfig.RoleId] = v.GetTimes
		}
	}
	return fiveInfo, fourInfo
}
func (mr *ModRole) CalHpPool() {
	if mr.HpCalTime == 0 {
		mr.HpCalTime = time.Now().Unix()
	}
	calTime := time.Now().Unix() - mr.HpCalTime
	mr.HpPool += int(calTime) * 10
	mr.HpCalTime = time.Now().Unix()
	fmt.Println("当前血池回复量:", mr.HpPool)
}

func (mr *ModRole) WearRelics(roleInfo *RoleInfo, relics *Relics, player *Player) {
	relicsConfig := csvs.GetRelicsConfig(relics.RelicsId)
	if relicsConfig == nil {
		return
	}
	mr.CheckRelicsPos(roleInfo, relicsConfig.Pos)
	if relicsConfig.Pos < 0 || relicsConfig.Pos > len(roleInfo.RelicsInfo) {
		return
	}

	oldRelicsKeyId := roleInfo.RelicsInfo[relicsConfig.Pos-1]
	if oldRelicsKeyId > 0 {
		oldRelics := player.ModRelics.RelicsInfo[oldRelicsKeyId]
		if oldRelics != nil {
			oldRelics.RoleId = 0
		}
		roleInfo.RelicsInfo[relicsConfig.Pos-1] = 0
	}
	oldRoleId := relics.RoleId
	if oldRoleId > 0 {
		oldRole := player.ModRole.RoleInfo[oldRoleId]
		if oldRole != nil {
			oldRole.RelicsInfo[relicsConfig.Pos-1] = 0
		}
		relics.RoleId = 0
	}

	roleInfo.RelicsInfo[relicsConfig.Pos-1] = relics.KeyId
	relics.RoleId = roleInfo.RoleId

	if oldRelicsKeyId > 0 && oldRoleId > 0 {
		oldRelics := player.ModRelics.RelicsInfo[oldRelicsKeyId]
		oldRole := player.ModRole.RoleInfo[oldRoleId]
		if oldRelics != nil && oldRole != nil {
			mr.WearRelics(oldRole, oldRelics, player)
		}
	}

	roleInfo.ShowInfo(player)
}

func (mr *ModRole) CheckRelicsPos(roleInfo *RoleInfo, pos int) {
	nowSize := len(roleInfo.RelicsInfo)
	needAdd := pos - nowSize
	for i := 0; i < needAdd; i++ {
		roleInfo.RelicsInfo = append(roleInfo.RelicsInfo, 0)
	}
}

func (ri *RoleInfo) ShowInfo(player *Player) {
	fmt.Println(fmt.Sprintf("当前角色：%s,当前角色ID:%d", csvs.GetItemName(ri.RoleId), ri.RoleId))
	suitMap := make(map[int]int64)

	for _, v := range ri.RelicsInfo {
		relicsNow := player.ModRelics.RelicsInfo[v]
		if relicsNow == nil {
			fmt.Println(fmt.Sprintf("当前部位：%d未穿戴", v))
		}
		fmt.Println(fmt.Sprintf("当前部位：%d,当前装备:%s", v, csvs.GetItemName(relicsNow.RelicsId)))
		relicsNowConfig := csvs.GetRelicsConfig(relicsNow.RelicsId)
		if relicsNowConfig != nil {
			suitMap[relicsNowConfig.Type]++
		}
	}
	suitSkill := make([]int, 0)
	for suit, num := range suitMap {
		for _, config := range csvs.ConfigRelicsSuitMap[suit] {
			if num >= config.Num {
				suitSkill = append(suitSkill, config.SuitSkill)
			}
		}
	}
	for _, v := range suitSkill {
		fmt.Println(fmt.Sprintf("激活套装效果:%d", v))
	}
}

func (mr *ModRole) TakeOffRelics(roleInfo *RoleInfo, relics *Relics, player *Player) {
	relicsConfig := csvs.GetRelicsConfig(relics.RelicsId)
	if relicsConfig == nil {
		return
	}
	mr.CheckRelicsPos(roleInfo, relicsConfig.Pos)
	if relicsConfig.Pos < 0 || relicsConfig.Pos > len(roleInfo.RelicsInfo) {
		return
	}
	if roleInfo.RelicsInfo[relicsConfig.Pos-1] != relics.KeyId {
		fmt.Println(fmt.Sprintf("当前部位：%d未穿戴", relics.KeyId))
		return
	}

	roleInfo.RelicsInfo[relicsConfig.Pos-1] = 0
	relics.RoleId = 0
	roleInfo.ShowInfo(player)
}

func (mr *ModRole) WearWeapon(roleInfo *RoleInfo, weapon *Weapon, player *Player) {
	weaponConfig := csvs.GetWeaponConfig(weapon.WeaponId)
	if weaponConfig == nil {
		return
	}

	//判断武器和角色是否匹配
	roleConfig := csvs.GetRoleConfig(roleInfo.RoleId)
	if roleConfig.Type != weaponConfig.Type {
		fmt.Println("武器和英雄不匹配")
	}

	oldWeaponKey := 0
	if roleInfo.WeaponInfo > 0 {
		oldWeaponKey = roleInfo.WeaponInfo
		roleInfo.WeaponInfo = 0
		oldWeapon := player.ModWeapon.WeaponInfo[oldWeaponKey]
		if oldWeapon != nil {
			oldWeapon.RoleId = 0
		}
	}
	oldRoleId := 0
	if weapon.RoleId > 0 {
		oldRoleId = weapon.RoleId
		weapon.RoleId = 0
		oldRole := player.ModRole.RoleInfo[oldRoleId]
		if oldRole != nil {
			oldRole.WeaponInfo = 0
		}
	}

	roleInfo.WeaponInfo = weapon.KeyId
	weapon.RoleId = roleInfo.RoleId
	if roleInfo.WeaponInfo > 0 && weapon.RoleId > 0 {
		oldWeapon := player.ModWeapon.WeaponInfo[oldWeaponKey]
		oldRole := player.ModRole.RoleInfo[oldRoleId]
		if oldRole.RoleId > 0 && oldWeapon.WeaponId > 0 {
			mr.WearWeapon(oldRole, oldWeapon, player)
		}
	}
}

func (mr *ModRole) TakeOffWeapon(roleInfo *RoleInfo, weapon *Weapon, player *Player) {
	weaponConfig := csvs.GetWeaponConfig(weapon.WeaponId)
	if weaponConfig == nil {
		return
	}

	if roleInfo.WeaponInfo != weapon.KeyId {
		fmt.Println("角色未装备该武器")
		return
	}
	roleInfo.WeaponInfo = 0
	weapon.RoleId = 0
}

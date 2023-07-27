package game

import (
	"fmt"
	"server_logic/csvs"
)

type Weapon struct {
	WeaponId    int
	KeyId       int
	Star        int
	Exp         int
	RoleId      int
	Level       int
	RefineLevel int
	StarLevel   int
}

type ModWeapon struct {
	WeaponInfo map[int]*Weapon
	MasKey     int
}

func (mw *ModWeapon) AddItem(itemId int, num int64) {
	config := csvs.GetWeaponConfig(itemId)
	if config == nil {
		fmt.Println("该武器不存在")
		return
	}
	if len(mw.WeaponInfo)+int(num) > csvs.WEAPON_MAX_COUNT {
		fmt.Println("超过最大值")
		return
	}
	for i := int64(0); i < num; i++ {
		weapon := new(Weapon)
		weapon.WeaponId = itemId
		mw.MasKey++
		weapon.KeyId = mw.MasKey
		mw.WeaponInfo[weapon.KeyId] = weapon
		fmt.Println("获得武器,", csvs.GetItemConfig(itemId).ItemName, "-------武器编号:", weapon.KeyId)
	}
	return
}

func (mw *ModWeapon) WeaponUp(keyId int, player *Player) {
	weapon := mw.WeaponInfo[keyId]
	if weapon == nil {
		return
	}
	weaponConfig := csvs.GetWeaponConfig(weapon.WeaponId)
	if weaponConfig == nil {
		return
	}
	for {
		nextLevelConfig := csvs.GetWeaponLevelConfig(weaponConfig.Star, weapon.Level+1)
		if nextLevelConfig == nil {
			fmt.Println("返还经验:", weapon.Exp)
			weapon.Exp = 0
			break
		}
		if weapon.StarLevel < nextLevelConfig.NeedStarLevel {
			fmt.Println("返还经验:", weapon.Exp)
			weapon.Exp = 0
			break
		}
		if weapon.Exp < nextLevelConfig.NeedExp {
			break
		}
		weapon.Level++
		weapon.Exp -= nextLevelConfig.NeedExp
	}
	weapon.ShowInfo()
}

func (w *Weapon) ShowInfo() {
	fmt.Println(fmt.Sprintf("当前等级%d,当前经验:%d,当前突破等级:%d", w.Level, w.RefineLevel, w.RefineLevel))
}

func (mw *ModWeapon) WeaponUpStar(keyId int, player *Player) {
	weapon := mw.WeaponInfo[keyId]
	if weapon == nil {
		return
	}
	weaponConfig := csvs.GetWeaponConfig(weapon.WeaponId)
	if weaponConfig == nil {
		return
	}
	nextStarConfig := csvs.GetWeaponStarConfig(weaponConfig.Star, weapon.Star+1)
	if nextStarConfig == nil {
		return
	}
	// 验证物品够不够
	if weapon.Level < nextStarConfig.Level {
		fmt.Println("武器等级不足")
		return
	}
	weapon.StarLevel++
	weapon.ShowInfo()
}

func (mw *ModWeapon) WeaponUpRefine(keyId, targetKeyId int, player *Player) {
	if keyId == targetKeyId {
		fmt.Println("错误的材料")
		return
	}
	weapon := mw.WeaponInfo[keyId]
	if weapon == nil {
		return
	}

	targetWeapon := mw.WeaponInfo[targetKeyId]
	if targetWeapon == nil {
		return
	}
	if weapon.WeaponId != targetWeapon.WeaponId {
		fmt.Println("错误的材料")
		return
	}
	if weapon.RefineLevel >= csvs.WEAPON_MAX_REFINE {
		fmt.Println("超过了最大精炼等级")
		return
	}
	weapon.RefineLevel++
	delete(mw.WeaponInfo, targetKeyId)
	weapon.ShowInfo()
}

package game

import (
	"fmt"
	"server_logic/csvs"
)

type ItemInfo struct {
	ItemId  int
	ItemNum int64
}

type ModBag struct {
	BagInfo map[int]*ItemInfo
}

func (mb *ModBag) AddItem(itemId int, num int64, player *Player) {
	itemConfig := csvs.GetItemConfig(itemId)
	if itemConfig == nil {
		fmt.Println(itemId, "物品不存在")
		return
	}
	switch itemConfig.SortType {
	//case csvs.ITEMTYPE_NORMAL:
	//	self.AddItemToBag(itemId, num)
	case csvs.ITEMTYPE_ROLE:
		player.ModRole.AddItem(itemId, num, player)
	case csvs.ITEMTYPE_ICON:
		player.ModIcon.AddItem(itemId, player)
	case csvs.ITEMTYPE_CARD:
		player.ModCard.AddItem(itemId, 1)
	case csvs.ITEMTYPE_WEAPON:
		player.ModWeapon.AddItem(itemId, num)
	case csvs.ITEMTYPE_RELICS:
		player.ModRelics.AddItem(itemId, num)
	case csvs.ITEMTYPE_COOK:
		player.ModCook.AddItem(itemId)
	case csvs.ITEMTYPE_HOME_ITEM:
		player.ModHome.AddItem(itemId, num, player)
	default: // 同普通
		mb.AddItemToBag(itemId, num)
	}

}
func (mb *ModBag) AddItemToBag(itemId int, num int64) {
	_, ok := mb.BagInfo[itemId]
	if ok {
		mb.BagInfo[itemId].ItemNum += num
	} else {
		mb.BagInfo[itemId] = &ItemInfo{
			ItemId:  itemId,
			ItemNum: num,
		}
	}
	config := csvs.GetItemConfig(itemId)
	if config == nil {
		fmt.Println("不存在该物品")
		return
	}
	fmt.Println("获取物品", config.ItemName, "数量", num, "---- 当前数量", mb.BagInfo[itemId].ItemNum)
}

func (mb *ModBag) RemoveItem(itemId int, num int64) {
	itemConfig := csvs.GetItemConfig(itemId)
	if itemConfig == nil {
		fmt.Println(itemId, "物品不存在")
		return
	}
	switch itemConfig.SortType {
	case csvs.ITEMTYPE_NORMAL:
		mb.RemoveItemToBagGM(itemId, num)
	default: // 同普通
		//self.AddItemToBag(itemId, 1)
	}

}

func (mb *ModBag) RemoveItemToBagGM(itemId int, num int64) {
	_, ok := mb.BagInfo[itemId]
	if ok {
		mb.BagInfo[itemId].ItemNum -= num
	} else {
		mb.BagInfo[itemId] = &ItemInfo{
			ItemId:  itemId,
			ItemNum: 0 - num,
		}
	}
	config := csvs.GetItemConfig(itemId)
	if config == nil {
		fmt.Println("不存在该物品")
		return
	}
	fmt.Println("扣除物品", config.ItemName, "数量", num, "---- 当前数量", mb.BagInfo[itemId].ItemNum)
}

func (mb *ModBag) RemoveItemToBag(itemId int, num int64, player *Player) {
	config := csvs.GetItemConfig(itemId)
	if config == nil {
		fmt.Println("不存在该物品")
		return
	}
	if !mb.HasEnoughItem(itemId, num) {
		var nowNum int64
		if v, ok := mb.BagInfo[itemId]; ok {
			nowNum = v.ItemNum
		}
		fmt.Println(config.ItemName, "数量不足", num, "---- 当前数量", nowNum)
		return
	}
	_, ok := mb.BagInfo[itemId]
	if ok {
		mb.BagInfo[itemId].ItemNum -= num
	} else {
		mb.BagInfo[itemId] = &ItemInfo{
			ItemId:  itemId,
			ItemNum: 0 - num,
		}
	}
	if config == nil {
		fmt.Println("不存在该物品")
		return
	}
	fmt.Println("扣除物品", config.ItemName, "数量", num, "---- 当前数量", mb.BagInfo[itemId].ItemNum)
}

func (mb *ModBag) HasEnoughItem(itemId int, num int64) bool {
	v, ok := mb.BagInfo[itemId]
	if !ok {
		return false
	} else {
		if v.ItemNum < num {
			return false
		}
	}
	return true
}

func (mb *ModBag) UserItem(itemId int, num int64, player *Player) {
	itemConfig := csvs.GetItemConfig(itemId)
	if itemConfig == nil {
		fmt.Println(itemId, "物品不存在")
		return
	}

	if !mb.HasEnoughItem(itemId, num) {
		var nowNum int64
		if v, ok := mb.BagInfo[itemId]; ok {
			nowNum = v.ItemNum
		}
		fmt.Println(itemConfig.ItemName, "数量不足", num, "---- 当前数量", nowNum)
		return
	}

	switch itemConfig.SortType {
	//case csvs.ITEMTYPE_NORMAL:
	//	self.AddItemToBag(itemId, num)
	case csvs.ITEMTYPE_COOKBOOK:
		mb.UserCookBook(itemId, num, player)
	case csvs.ITEMTYPE_FOOD:
		// 给英雄加属性
		player.ModCook.AddItem(itemId)
	default: //
		fmt.Println(itemId, "此物品无法使用")
		return
	}

}

func (mb *ModBag) UserCookBook(itemId int, num int64, player *Player) {
	cookBookConfig := csvs.GetCookBookConfig(itemId)
	if cookBookConfig == nil {
		fmt.Println("物品不存在")
		return
	}

	mb.RemoveItem(itemId, num)
	mb.AddItem(cookBookConfig.Reward, num, player)
}

func (mb *ModBag) GetItemNum(itemId int) int64 {
	itemConfig := csvs.GetItemConfig(itemId)
	if itemConfig == nil {
		return 0
	}
	if _, ok := mb.BagInfo[itemId]; !ok {
		return 0
	}
	return mb.BagInfo[itemId].ItemNum
}

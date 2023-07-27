package game

import (
	"fmt"
	"server_logic/csvs"
)

type HomeItem struct {
	HomeItemId  int
	HomeItemNum int64
	KeyId       int
}

type ModHome struct {
	HomeItemInfo map[int]*HomeItem
}

func (mh *ModHome) AddItem(itemId int, num int64, player *Player) {
	_, ok := mh.HomeItemInfo[itemId]
	if ok {
		mh.HomeItemInfo[itemId].HomeItemNum += num
	} else {
		mh.HomeItemInfo[itemId] = &HomeItem{
			HomeItemId:  itemId,
			HomeItemNum: num,
		}
	}
	config := csvs.GetItemConfig(itemId)
	if config == nil {
		fmt.Println("不存在该物品")
		return
	}
	fmt.Println("获取家居物品", config.ItemName, "数量", num, "---- 当前数量", mh.HomeItemInfo[itemId].HomeItemNum)

}

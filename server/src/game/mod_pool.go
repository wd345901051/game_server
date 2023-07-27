package game

import (
	"fmt"
	"math/rand"
	"server_logic/csvs"
)

type PoolInfo struct {
	PoolId        int
	FiveStarTimes int
	FourStarTimes int
	IsMustUp      int
}

type ModPool struct {
	UpPoolInfo *PoolInfo
}

func (mp *ModPool) AddTimes() {
	mp.UpPoolInfo.FiveStarTimes++
	mp.UpPoolInfo.FourStarTimes++
}

func (mp *ModPool) DoUpPool(num int) {
	res := make(map[int]int)
	fourNum := 0
	fiveNum := 0
	for i := 0; i < num; i++ {
		mp.AddTimes()
		dropGroup := csvs.ConfigDropGroupMap[1000]
		if dropGroup == nil {
			return
		}
		if mp.UpPoolInfo.FiveStarTimes > csvs.FIVE_STAR_TIMES_LIMIT || mp.UpPoolInfo.FourStarTimes > csvs.FOUR_STAR_TIMES_LIMIT {
			newDropGroup := new(csvs.DropGroup)
			newDropGroup.DropId = dropGroup.DropId
			newDropGroup.WeightAll = dropGroup.WeightAll
			addFiveWeight := (mp.UpPoolInfo.FiveStarTimes - csvs.FIVE_STAR_TIMES_LIMIT) * csvs.FIVE_STAR_TIMES_LIMIT_EACH_VALUE
			if addFiveWeight < 0 {
				addFiveWeight = 0
			}
			addFourWeight := (mp.UpPoolInfo.FourStarTimes - csvs.FOUR_STAR_TIMES_LIMIT) * csvs.FOUR_STAR_TIMES_LIMIT_EACH_VALUE
			if addFourWeight < 0 {
				addFourWeight = 0
			}
			for _, config := range dropGroup.DropConfigs {
				newConfig := new(csvs.ConfigDrop)
				newConfig.Result = config.Result
				newConfig.IsEnd = config.IsEnd
				newConfig.DropId = config.DropId
				if config.Result == 10001 {
					newConfig.Weight = config.Weight + addFiveWeight
				} else if config.Result == 10002 {
					newConfig.Weight = config.Weight + addFourWeight
				} else if config.Result == 10003 {
					newConfig.Weight = config.Weight - addFiveWeight - addFourWeight
				}
				newDropGroup.DropConfigs = append(newDropGroup.DropConfigs, newConfig)
			}
			dropGroup = newDropGroup
		}
		roleIdConfig := mp.GetRandDropNew(dropGroup)
		if roleIdConfig != nil {
			roleConfig := csvs.GetRoleConfig(roleIdConfig.Result)
			if roleConfig != nil {
				if roleConfig.Star == 5 {
					mp.UpPoolInfo.FiveStarTimes = 0
					fiveNum++
					if mp.UpPoolInfo.IsMustUp == csvs.LOGIC_TRUE {
						dropGroup := csvs.ConfigDropGroupMap[100012]
						if dropGroup != nil {
							roleIdConfig = mp.GetRandDropNew(dropGroup)
							if roleConfig == nil {
								fmt.Println("数据异常")
								return
							}
						}
					}
					if roleIdConfig.DropId == 100012 {
						mp.UpPoolInfo.IsMustUp = csvs.LOGIC_FALSE
					} else {
						mp.UpPoolInfo.IsMustUp = csvs.LOGIC_TRUE
					}
				}
				if roleConfig.Star == 4 {
					mp.UpPoolInfo.FourStarTimes = 0
					fourNum++
				}
			}
			res[roleIdConfig.Result]++
		}
	}
	for v, k := range res {
		fmt.Println(csvs.GetItemConfig(v).ItemName, ":有", k, "个")
	}
	fmt.Println(fourNum, "----", fiveNum)
}

func (mp *ModPool) GetRandDropNew(dropGroup *csvs.DropGroup) *csvs.ConfigDrop {
	ranNum := rand.Intn(dropGroup.WeightAll)
	randNow := 0
	for _, v := range dropGroup.DropConfigs {
		randNow += v.Weight
		if ranNum < randNow {
			if v.IsEnd == csvs.LOGIC_TRUE {
				return v
			}
			dropGroup := csvs.ConfigDropGroupMap[v.Result]
			if dropGroup == nil {
				return nil
			}
			return mp.GetRandDropNew(dropGroup)
		}
	}
	return nil
}

func (self *ModPool) HandleUpPoolTen(player *Player) {
	for i := 0; i < 10; i++ {
		self.AddTimes()
		dropGroup := csvs.ConfigDropGroupMap[1000]
		if dropGroup == nil {
			return
		}

		if self.UpPoolInfo.FiveStarTimes > csvs.FIVE_STAR_TIMES_LIMIT || self.UpPoolInfo.FourStarTimes > csvs.FOUR_STAR_TIMES_LIMIT {
			newDropGroup := new(csvs.DropGroup)
			newDropGroup.DropId = dropGroup.DropId
			newDropGroup.WeightAll = dropGroup.WeightAll
			addFiveWeight := (self.UpPoolInfo.FiveStarTimes - csvs.FIVE_STAR_TIMES_LIMIT) * csvs.FIVE_STAR_TIMES_LIMIT_EACH_VALUE
			if addFiveWeight < 0 {
				addFiveWeight = 0
			}
			addFourWeight := (self.UpPoolInfo.FourStarTimes - csvs.FOUR_STAR_TIMES_LIMIT) * csvs.FOUR_STAR_TIMES_LIMIT_EACH_VALUE
			if addFourWeight < 0 {
				addFourWeight = 0
			}

			for _, config := range dropGroup.DropConfigs {
				newConfig := new(csvs.ConfigDrop)
				newConfig.Result = config.Result
				newConfig.DropId = config.DropId
				newConfig.IsEnd = config.IsEnd
				if config.Result == 10001 {
					newConfig.Weight = config.Weight + addFiveWeight
				} else if config.Result == 10002 {
					newConfig.Weight = config.Weight + addFourWeight
				} else if config.Result == 10003 {
					newConfig.Weight = config.Weight - addFiveWeight - addFourWeight
				}
				newDropGroup.DropConfigs = append(newDropGroup.DropConfigs, newConfig)
			}
			dropGroup = newDropGroup
		}
		roleIdConfig := self.GetRandDropNew(dropGroup)
		if roleIdConfig != nil {
			roleConfig := csvs.GetRoleConfig(roleIdConfig.Result)
			if roleConfig != nil {
				if roleConfig.Star == 5 {
					self.UpPoolInfo.FiveStarTimes = 0
					if self.UpPoolInfo.IsMustUp == csvs.LOGIC_TRUE {
						dropGroup := csvs.ConfigDropGroupMap[100012]
						if dropGroup != nil {
							roleIdConfig = self.GetRandDropNew(dropGroup)
							if roleIdConfig == nil {
								fmt.Println("数据异常")
								return
							}
						}
					}
					if roleIdConfig.DropId == 100012 {
						self.UpPoolInfo.IsMustUp = csvs.LOGIC_FALSE
					} else {
						self.UpPoolInfo.IsMustUp = csvs.LOGIC_TRUE
					}
				} else if roleConfig.Star == 4 {
					self.UpPoolInfo.FourStarTimes = 0
				}
			}
			//fmt.Println(fmt.Sprintf("第%d抽抽中:%s", i+1, csvs.GetItemConfig(roleIdConfig.Result).ItemName))
			player.ModBag.AddItem(roleIdConfig.Result, 1, player)
		}
	}
	if self.UpPoolInfo.IsMustUp == csvs.LOGIC_FALSE {
		fmt.Println(fmt.Sprintf("当前处于小保底区间！"))
	} else {
		fmt.Println(fmt.Sprintf("当前处于大保底区间！"))
	}
	fmt.Println(fmt.Sprintf("当前累计未出5星次数：%d", self.UpPoolInfo.FiveStarTimes))
	fmt.Println(fmt.Sprintf("当前累计未出4星次数：%d", self.UpPoolInfo.FourStarTimes))

}

func (self *ModPool) HandleUpPoolSingle(times int, player *Player) {
	if times <= 0 || times > 100000000 {
		fmt.Println("请输入正确的数值(1~100000000)")
		return
	} else {
		fmt.Println(fmt.Sprintf("累计抽取%d次,结果如下:", times))
	}
	result := make(map[int]int)
	fourNum := 0
	fiveNum := 0
	for i := 0; i < times; i++ {
		self.AddTimes()
		dropGroup := csvs.ConfigDropGroupMap[1000]
		if dropGroup == nil {
			return
		}

		if self.UpPoolInfo.FiveStarTimes > csvs.FIVE_STAR_TIMES_LIMIT || self.UpPoolInfo.FourStarTimes > csvs.FOUR_STAR_TIMES_LIMIT {
			newDropGroup := new(csvs.DropGroup)
			newDropGroup.DropId = dropGroup.DropId
			newDropGroup.WeightAll = dropGroup.WeightAll
			addFiveWeight := (self.UpPoolInfo.FiveStarTimes - csvs.FIVE_STAR_TIMES_LIMIT) * csvs.FIVE_STAR_TIMES_LIMIT_EACH_VALUE
			if addFiveWeight < 0 {
				addFiveWeight = 0
			}
			addFourWeight := (self.UpPoolInfo.FourStarTimes - csvs.FOUR_STAR_TIMES_LIMIT) * csvs.FOUR_STAR_TIMES_LIMIT_EACH_VALUE
			if addFourWeight < 0 {
				addFourWeight = 0
			}

			for _, config := range dropGroup.DropConfigs {
				newConfig := new(csvs.ConfigDrop)
				newConfig.Result = config.Result
				newConfig.DropId = config.DropId
				newConfig.IsEnd = config.IsEnd
				if config.Result == 10001 {
					newConfig.Weight = config.Weight + addFiveWeight
				} else if config.Result == 10002 {
					newConfig.Weight = config.Weight + addFourWeight
				} else if config.Result == 10003 {
					newConfig.Weight = config.Weight - addFiveWeight - addFourWeight
				}
				newDropGroup.DropConfigs = append(newDropGroup.DropConfigs, newConfig)
			}
			dropGroup = newDropGroup
		}
		roleIdConfig := self.GetRandDropNew(dropGroup)
		if roleIdConfig != nil {
			roleConfig := csvs.GetRoleConfig(roleIdConfig.Result)
			if roleConfig != nil {
				if roleConfig.Star == 5 {
					self.UpPoolInfo.FiveStarTimes = 0
					fiveNum++
					if self.UpPoolInfo.IsMustUp == csvs.LOGIC_TRUE {
						dropGroup := csvs.ConfigDropGroupMap[100012]
						if dropGroup != nil {
							roleIdConfig = self.GetRandDropNew(dropGroup)
							if roleIdConfig == nil {
								fmt.Println("数据异常")
								return
							}
						}
					}
					if roleIdConfig.DropId == 100012 {
						self.UpPoolInfo.IsMustUp = csvs.LOGIC_FALSE
					} else {
						self.UpPoolInfo.IsMustUp = csvs.LOGIC_TRUE
					}
				} else if roleConfig.Star == 4 {
					self.UpPoolInfo.FourStarTimes = 0
					fourNum++
				}
			}
			result[roleIdConfig.Result]++
			player.ModBag.AddItem(roleIdConfig.Result, 1, player)
		}
	}

	for k, v := range result {
		fmt.Println(fmt.Sprintf("抽中%s次数：%d", csvs.GetItemConfig(k).ItemName, v))
	}
	fmt.Println(fmt.Sprintf("抽中4星：%d", fourNum))
	fmt.Println(fmt.Sprintf("抽中5星：%d", fiveNum))
}

func (self *ModPool) HandleUpPoolTimesTest(times int) {
	if times <= 0 || times > 100000000 {
		fmt.Println("请输入正确的数值(1~100000000)")
		return
	} else {
		fmt.Println(fmt.Sprintf("累计抽取%d次,结果如下:", times))
	}
	resultEach := make(map[int]int)
	for i := 0; i < times; i++ {
		self.AddTimes()
		dropGroup := csvs.ConfigDropGroupMap[1000]
		if dropGroup == nil {
			return
		}

		if self.UpPoolInfo.FiveStarTimes > csvs.FIVE_STAR_TIMES_LIMIT || self.UpPoolInfo.FourStarTimes > csvs.FOUR_STAR_TIMES_LIMIT {
			newDropGroup := new(csvs.DropGroup)
			newDropGroup.DropId = dropGroup.DropId
			newDropGroup.WeightAll = dropGroup.WeightAll
			addFiveWeight := (self.UpPoolInfo.FiveStarTimes + 1 - csvs.FIVE_STAR_TIMES_LIMIT) * csvs.FIVE_STAR_TIMES_LIMIT_EACH_VALUE
			if addFiveWeight < 0 {
				addFiveWeight = 0
			}
			addFourWeight := (self.UpPoolInfo.FourStarTimes + 1 - csvs.FOUR_STAR_TIMES_LIMIT) * csvs.FOUR_STAR_TIMES_LIMIT_EACH_VALUE
			if addFourWeight < 0 {
				addFourWeight = 0
			}

			for _, config := range dropGroup.DropConfigs {
				newConfig := new(csvs.ConfigDrop)
				newConfig.Result = config.Result
				newConfig.DropId = config.DropId
				newConfig.IsEnd = config.IsEnd
				if config.Result == 10001 {
					newConfig.Weight = config.Weight + addFiveWeight
				} else if config.Result == 10002 {
					newConfig.Weight = config.Weight + addFourWeight
				} else if config.Result == 10003 {
					newConfig.Weight = config.Weight - addFiveWeight - addFourWeight
				}
				newDropGroup.DropConfigs = append(newDropGroup.DropConfigs, newConfig)
			}
			dropGroup = newDropGroup
		}
		roleIdConfig := self.GetRandDropNew(dropGroup)
		if roleIdConfig != nil {
			roleConfig := csvs.GetRoleConfig(roleIdConfig.Result)
			if roleConfig != nil {
				if roleConfig.Star == 5 {
					resultEach[self.UpPoolInfo.FiveStarTimes]++
					self.UpPoolInfo.FiveStarTimes = 0
					if self.UpPoolInfo.IsMustUp == csvs.LOGIC_TRUE {
						dropGroup := csvs.ConfigDropGroupMap[100012]
						if dropGroup != nil {
							roleIdConfig = self.GetRandDropNew(dropGroup)
							if roleIdConfig == nil {
								fmt.Println("数据异常")
								return
							}
						}
					}
					if roleIdConfig.DropId == 100012 {
						self.UpPoolInfo.IsMustUp = csvs.LOGIC_FALSE
					} else {
						self.UpPoolInfo.IsMustUp = csvs.LOGIC_TRUE
					}
				} else if roleConfig.Star == 4 {
					self.UpPoolInfo.FourStarTimes = 0
				}
			}
		}
	}

	for k, v := range resultEach {
		fmt.Println(fmt.Sprintf("第%d抽抽出5星的次数：%d", k, v))
	}
}

func (self *ModPool) HandleUpPoolFiveTest(times int) {
	if times <= 0 || times > 100000000 {
		fmt.Println("请输入正确的数值(1~100000000)")
		return
	} else {
		fmt.Println(fmt.Sprintf("累计抽取%d次,结果如下:", times))
	}
	resultEachTest := make(map[int]int)
	fiveTest := 0
	for i := 0; i < times; i++ {
		self.AddTimes()
		if i%10 == 0 {
			fiveTest = 0
		}
		dropGroup := csvs.ConfigDropGroupMap[1000]
		if dropGroup == nil {
			return
		}

		if self.UpPoolInfo.FiveStarTimes > csvs.FIVE_STAR_TIMES_LIMIT || self.UpPoolInfo.FourStarTimes > csvs.FOUR_STAR_TIMES_LIMIT {
			newDropGroup := new(csvs.DropGroup)
			newDropGroup.DropId = dropGroup.DropId
			newDropGroup.WeightAll = dropGroup.WeightAll
			addFiveWeight := (self.UpPoolInfo.FiveStarTimes - csvs.FIVE_STAR_TIMES_LIMIT) * csvs.FIVE_STAR_TIMES_LIMIT_EACH_VALUE
			if addFiveWeight < 0 {
				addFiveWeight = 0
			}
			addFourWeight := (self.UpPoolInfo.FourStarTimes - csvs.FOUR_STAR_TIMES_LIMIT) * csvs.FOUR_STAR_TIMES_LIMIT_EACH_VALUE
			if addFourWeight < 0 {
				addFourWeight = 0
			}

			for _, config := range dropGroup.DropConfigs {
				newConfig := new(csvs.ConfigDrop)
				newConfig.Result = config.Result
				newConfig.DropId = config.DropId
				newConfig.IsEnd = config.IsEnd
				if config.Result == 10001 {
					newConfig.Weight = config.Weight + addFiveWeight
				} else if config.Result == 10002 {
					newConfig.Weight = config.Weight + addFourWeight
				} else if config.Result == 10003 {
					newConfig.Weight = config.Weight - addFiveWeight - addFourWeight
				}
				newDropGroup.DropConfigs = append(newDropGroup.DropConfigs, newConfig)
			}
			dropGroup = newDropGroup
		}
		roleIdConfig := self.GetRandDropNew(dropGroup)
		if roleIdConfig != nil {
			roleConfig := csvs.GetRoleConfig(roleIdConfig.Result)
			if roleConfig != nil {
				if roleConfig.Star == 5 {
					fiveTest++
					self.UpPoolInfo.FiveStarTimes = 0
					if self.UpPoolInfo.IsMustUp == csvs.LOGIC_TRUE {
						dropGroup := csvs.ConfigDropGroupMap[100012]
						if dropGroup != nil {
							roleIdConfig = self.GetRandDropNew(dropGroup)
							if roleIdConfig == nil {
								fmt.Println("数据异常")
								return
							}
						}
					}
					if roleIdConfig.DropId == 100012 {
						self.UpPoolInfo.IsMustUp = csvs.LOGIC_FALSE
					} else {
						self.UpPoolInfo.IsMustUp = csvs.LOGIC_TRUE
					}
				} else if roleConfig.Star == 4 {
					self.UpPoolInfo.FourStarTimes = 0
				}
			}
		}
		if i%10 == 9 {
			resultEachTest[fiveTest]++
		}
	}

	for k, v := range resultEachTest {
		fmt.Println(fmt.Sprintf("10连%d黄次数：%d", k, v))
	}
}

func (mp *ModPool) DoUpPoolCheck(num int, player *Player) {
	res := make(map[int]int)
	fourNum := 0
	fiveNum := 0
	for i := 0; i < num; i++ {
		mp.AddTimes()
		dropGroup := csvs.ConfigDropGroupMap[1000]
		if dropGroup == nil {
			return
		}
		if mp.UpPoolInfo.FiveStarTimes > csvs.FIVE_STAR_TIMES_LIMIT || mp.UpPoolInfo.FourStarTimes > csvs.FOUR_STAR_TIMES_LIMIT {
			newDropGroup := new(csvs.DropGroup)
			newDropGroup.DropId = dropGroup.DropId
			newDropGroup.WeightAll = dropGroup.WeightAll
			addFiveWeight := (mp.UpPoolInfo.FiveStarTimes - csvs.FIVE_STAR_TIMES_LIMIT) * csvs.FIVE_STAR_TIMES_LIMIT_EACH_VALUE
			if addFiveWeight < 0 {
				addFiveWeight = 0
			}
			addFourWeight := (mp.UpPoolInfo.FourStarTimes - csvs.FOUR_STAR_TIMES_LIMIT) * csvs.FOUR_STAR_TIMES_LIMIT_EACH_VALUE
			if addFourWeight < 0 {
				addFourWeight = 0
			}
			for _, config := range dropGroup.DropConfigs {
				newConfig := new(csvs.ConfigDrop)
				newConfig.Result = config.Result
				newConfig.IsEnd = config.IsEnd
				newConfig.DropId = config.DropId
				if config.Result == 10001 {
					newConfig.Weight = config.Weight + addFiveWeight
				} else if config.Result == 10002 {
					newConfig.Weight = config.Weight + addFourWeight
				} else if config.Result == 10003 {
					newConfig.Weight = config.Weight - addFiveWeight - addFourWeight
				}
				newDropGroup.DropConfigs = append(newDropGroup.DropConfigs, newConfig)
			}
			dropGroup = newDropGroup
		}

		fiveInfo, fourInfo := player.ModRole.GetRoleInfoForPoolCheck()

		roleIdConfig := mp.GetRandDropNewNew1(dropGroup, fiveInfo, fourInfo)
		if roleIdConfig != nil {
			roleConfig := csvs.GetRoleConfig(roleIdConfig.Result)
			if roleConfig != nil {
				if roleConfig.Star == 5 {
					mp.UpPoolInfo.FiveStarTimes = 0
					fiveNum++
					if mp.UpPoolInfo.IsMustUp == csvs.LOGIC_TRUE {
						dropGroup := csvs.ConfigDropGroupMap[100012]
						if dropGroup != nil {
							roleIdConfig = mp.GetRandDropNewNew1(dropGroup, fiveInfo, fourInfo)
							if roleConfig == nil {
								fmt.Println("数据异常")
								return
							}
						}
					}
					if roleIdConfig.DropId == 100012 {
						mp.UpPoolInfo.IsMustUp = csvs.LOGIC_FALSE
					} else {
						mp.UpPoolInfo.IsMustUp = csvs.LOGIC_TRUE
					}
				}
				if roleConfig.Star == 4 {
					mp.UpPoolInfo.FourStarTimes = 0
					fourNum++
				}
			}
			res[roleIdConfig.Result]++
			player.ModBag.AddItem(roleIdConfig.Result, 1, player)
		}
	}
	for v, k := range res {
		fmt.Println(csvs.GetItemConfig(v).ItemName, ":有", k, "个")
	}
	fmt.Println(fourNum, "----", fiveNum)
}

func (mp *ModPool) GetRandDropNewNew(dropGroup *csvs.DropGroup, fiveInfo, fourInfo map[int]int) *csvs.ConfigDrop {
	for _, v := range dropGroup.DropConfigs {
		if _, ok := fiveInfo[v.Result]; ok {
			index := 0
			maxGetTime := 0
			for k, config := range dropGroup.DropConfigs {
				_, nowOK := fiveInfo[config.Result]
				if !nowOK {
					continue
				}
				if maxGetTime < fiveInfo[config.Result] {
					maxGetTime = fiveInfo[config.Result]
					index = k
				}
			}
			return dropGroup.DropConfigs[index]
		}

		if _, ok := fourInfo[v.Result]; ok {
			index := 0
			maxGetTime := 0
			for k, config := range dropGroup.DropConfigs {
				_, nowOK := fourInfo[config.Result]
				if !nowOK {
					continue
				}
				if maxGetTime < fourInfo[config.Result] {
					maxGetTime = fourInfo[config.Result]
					index = k
				}
			}
			return dropGroup.DropConfigs[index]
		}
	}

	ranNum := rand.Intn(dropGroup.WeightAll)
	randNow := 0
	for _, v := range dropGroup.DropConfigs {
		randNow += v.Weight
		if ranNum < randNow {
			if v.IsEnd == csvs.LOGIC_TRUE {
				return v
			}
			dropGroup := csvs.ConfigDropGroupMap[v.Result]
			if dropGroup == nil {
				return nil
			}
			return mp.GetRandDropNewNew(dropGroup, fiveInfo, fourInfo)
		}
	}
	return nil
}

func (mp *ModPool) GetRandDropNewNew1(dropGroup *csvs.DropGroup, fiveInfo, fourInfo map[int]int) *csvs.ConfigDrop {
	for _, v := range dropGroup.DropConfigs {
		if _, ok := fiveInfo[v.Result]; ok {
			index := 0
			minGetTime := 0
			for k, config := range dropGroup.DropConfigs {
				_, nowOK := fiveInfo[config.Result]
				if !nowOK {
					index = k
					break
				}
				if minGetTime == 0 || minGetTime > fiveInfo[config.Result] {
					minGetTime = fiveInfo[config.Result]
					index = k
				}
			}
			return dropGroup.DropConfigs[index]
		}

		if _, ok := fourInfo[v.Result]; ok {
			index := 0
			minGetTime := 0
			for k, config := range dropGroup.DropConfigs {
				_, nowOK := fourInfo[config.Result]
				if !nowOK {
					index = k
					break
				}
				if minGetTime == 0 || minGetTime > fourInfo[config.Result] {
					minGetTime = fourInfo[config.Result]
					index = k
				}
			}
			return dropGroup.DropConfigs[index]
		}
	}

	ranNum := rand.Intn(dropGroup.WeightAll)
	randNow := 0
	for _, v := range dropGroup.DropConfigs {
		randNow += v.Weight
		if ranNum < randNow {
			if v.IsEnd == csvs.LOGIC_TRUE {
				return v
			}
			dropGroup := csvs.ConfigDropGroupMap[v.Result]
			if dropGroup == nil {
				return nil
			}
			return mp.GetRandDropNewNew1(dropGroup, fiveInfo, fourInfo)
		}
	}
	return nil
}

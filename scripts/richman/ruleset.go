package richman

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

var (
	LuckRules map[LUCK_CARD_TYPE_ENUM]func(room *GameRoom, user *User) (err error)
	NewsRules map[NEWS_CARD_TYPE_ENUM]func(room *GameRoom, user *User) (err error)
)

//初始化规则处理函数
func InitRuleset() {
	LuckRules = map[LUCK_CARD_TYPE_ENUM]func(room *GameRoom, user *User) (err error){
		LUCK_CARD_TYPE__NO1:  LuckCardNO1,
		LUCK_CARD_TYPE__NO2:  LuckCardNO2,
		LUCK_CARD_TYPE__NO3:  LuckCardNO3,
		LUCK_CARD_TYPE__NO4:  LuckCardNO4,
		LUCK_CARD_TYPE__NO5:  LuckCardNO5,
		LUCK_CARD_TYPE__NO6:  LuckCardNO6,
		LUCK_CARD_TYPE__NO7:  LuckCardNO7,
		LUCK_CARD_TYPE__NO8:  LuckCardNO8,
		LUCK_CARD_TYPE__NO9:  LuckCardNO9,
		LUCK_CARD_TYPE__NO10: LuckCardNO10,
		LUCK_CARD_TYPE__NO11: LuckCardNO11,
		LUCK_CARD_TYPE__NO12: LuckCardNO12,
		LUCK_CARD_TYPE__NO13: LuckCardNO13,
	}
	NewsRules = map[NEWS_CARD_TYPE_ENUM]func(room *GameRoom, user *User) (err error){
		NEWS_CARD_TYPE__NO1:  NewsCardNO1,
		NEWS_CARD_TYPE__NO2:  NewsCardNO2,
		NEWS_CARD_TYPE__NO3:  NewsCardNO3,
		NEWS_CARD_TYPE__NO4:  NewsCardNO4,
		NEWS_CARD_TYPE__NO5:  NewsCardNO5,
		NEWS_CARD_TYPE__NO6:  NewsCardNO6,
		NEWS_CARD_TYPE__NO7:  NewsCardNO7,
		NEWS_CARD_TYPE__NO8:  NewsCardNO8,
		NEWS_CARD_TYPE__NO9:  NewsCardNO9,
		NEWS_CARD_TYPE__NO10: NewsCardNO10,
		NEWS_CARD_TYPE__NO11: NewsCardNO11,
		NEWS_CARD_TYPE__NO12: NewsCardNO12,
	}
}

// 运气卡处理函数
// 遗失钱包，你失去1000元，位于你后方的第一位玩家获得1000元
func LuckCardNO1(room *GameRoom, user *User) (err error) {

	for u := range room.Money {
		if u == user {
			room.Money[user] -= 1000

			// 下一个玩家 + 1000
			next := room.Player.Next().Value.(*User)
			room.Money[next] += 1000

			room.Broadcast <- fmt.Sprintf("运气： @%s 遗失钱包，丢失1000元，@%s 捡到，意外获得1000元", user.Name, next.Name)
		}
	}

	return nil
}

// 黑历史被查，立即移动到监狱，并停留三回合
func LuckCardNO2(room *GameRoom, user *User) (err error) {
	for {
		// 如果在监狱中，停留加三回合
		if room.GameMap.Map[user].Value.(MapElement).Flag == GAME_FLAG__PRISION {
			room.Prision[user] += 3
			break
		}
		// 移动到监狱
		room.GameMap.Map[user] = room.GameMap.Map[user].Next()
	}

	room.Broadcast <- fmt.Sprintf("运气： @%s 黑历史被查，被送进监狱 3 天", user.Name)

	return nil
}

// 社会主义春风吹过，你可以立即免费升级一块
func LuckCardNO3(room *GameRoom, user *User) (err error) {
	for i, land := range room.GameMap.PlayerElement[user] {
		if land.Level >= int(LAND_LEVEL__MAX) {
			continue
		}

		room.GameMap.PlayerElement[user][i].Level++
		room.GameMap.PlayerElement[user][i].RentFee += land.Fee
		room.GameMap.PlayerElement[user][i].Fee += int64(float64(land.Fee) * 0.5)

		room.Broadcast <- fmt.Sprintf("运气： 改革春风吹满地，%s 获得一次升级地产的机会，%s 升级到 %d", user.Name, land.Name, room.GameMap.PlayerElement[user][i].Level)
		break
	}

	return nil
}

// 双十一期间，疯狂消费，支付1000元
func LuckCardNO4(room *GameRoom, user *User) (err error) {
	room.Money[user] -= 1000

	room.Broadcast <- fmt.Sprintf("运气： @%s 双十一疯狂消费，支付1000元", user.Name)

	return nil
}

// 前往九寨沟旅游，支付500元
func LuckCardNO5(room *GameRoom, user *User) (err error) {
	room.Money[user] -= 500

	room.Broadcast <- fmt.Sprintf("运气： @%s 前往九寨沟旅游，支付500元", user.Name)

	return nil
}

// 潜入银行内部系统，从每位玩家手中收取1000元
func LuckCardNO6(room *GameRoom, user *User) (err error) {
	var cost int64 = 0
	for u := range room.Money {
		if u != user {
			cost += 500
			room.Money[u] -= 500
		}
	}
	room.Money[user] += cost

	room.Broadcast <- fmt.Sprintf("运气： @%s 潜入银行内部系统，从每位玩家手中收取500元", user.Name)

	return nil
}

// 在香港乘坐豪华游轮，支付1000元，并立即移动到起点领取奖励
func LuckCardNO7(room *GameRoom, user *User) (err error) {
	if room.Money[user] >= 1000 {
		room.Money[user] -= 1000

		for {
			if room.GameMap.Map[user].Value.(MapElement).Flag == GAME_FLAG__START_POINT {
				room.Money[user] += SEND_MONY
				break
			}
			// 移动到起点
			room.GameMap.Map[user] = room.GameMap.Map[user].Next()
		}

		room.Broadcast <- fmt.Sprintf("运气： @%s 在香港乘坐豪华游轮，支付1000元，达到起点并领取 %d 元", user.Name, SEND_MONY)
	} else {
		room.Broadcast <- fmt.Sprintf("运气： @%s 因资金不足，在香港乘坐豪华游轮时被赶下船", user.Name)
	}

	return nil
}

// 立即移动到你的左边玩家的位置，并按该结果结算
func LuckCardNO8(room *GameRoom, user *User) (err error) {
	// 玩家往左移动，直到找到最近的一个玩家
	var uu *User
	var location MapElement
MOVE:
	for i := 0; i < room.GameMap.Map[user].Len(); i++ {
		for u, rl := range room.GameMap.Map {
			if u != user {
				if room.GameMap.Map[user].Value.(MapElement).IsEqual(rl.Value.(MapElement)) {
					uu = u
					location = room.GameMap.Map[user].Value.(MapElement)
					break MOVE
				}
			}
		}

		room.GameMap.Map[user] = room.GameMap.Map[user].Prev()
	}

	room.do(user, location)

	room.Broadcast <- fmt.Sprintf("运气： @%s 移动到了左边玩家 @%s 的位置", user.Name, uu.Name)

	return nil
}

// 立即移动到你的右边玩家的位置，并按该结果结算
func LuckCardNO9(room *GameRoom, user *User) (err error) {
	// 玩家往左移动，直到找到最近的一个玩家
	var uu *User
	var location MapElement
MOVE:
	for i := 0; i < room.GameMap.Map[user].Len(); i++ {
		for u, rl := range room.GameMap.Map {
			if u != user {
				if room.GameMap.Map[user].Value.(MapElement).IsEqual(rl.Value.(MapElement)) {
					uu = u
					location = room.GameMap.Map[user].Value.(MapElement)
					break MOVE
				}
			}
		}

		room.GameMap.Map[user] = room.GameMap.Map[user].Next()
	}

	room.do(user, location)

	room.Broadcast <- fmt.Sprintf("运气： @%s 移动到了右边玩家 @%s 的位置", user.Name, uu.Name)

	return nil
}

// 彩票中奖，获得5000元
func LuckCardNO10(room *GameRoom, user *User) (err error) {
	room.Money[user] += 5000

	room.Broadcast <- fmt.Sprintf("运气： @%s 彩票中奖，获得5000元", user.Name)

	return nil
}

// 获得额外遗产，获得3000元
func LuckCardNO11(room *GameRoom, user *User) (err error) {
	room.Money[user] += 1000

	room.Broadcast <- fmt.Sprintf("运气： @%s 获得额外财产，获得3000元", user.Name)

	return nil
}

// 购买最新款私人坐骑，支付500元，并立即额外进行一回合的行动
func LuckCardNO12(room *GameRoom, user *User) (err error) {
	room.Money[user] -= 500

	if room.Direction == 0 {
		room.Player = room.Player.Prev()
	} else {
		room.Player = room.Player.Next()
	}

	room.Broadcast <- fmt.Sprintf("运气： @%s 购买最新款私人坐骑，支付500元，额外进行一回合行动", user.Name)

	return nil
}

// 被车祸撞伤入院 5 回合
func LuckCardNO13(room *GameRoom, user *User) (err error) {

	for {
		// 如果在医院中，停留加五回合
		if room.GameMap.Map[user].Value.(MapElement).Flag == GAME_FLAG__HOSPITAL {
			room.Hospital[user] += 5
			break
		}
		// 移动到医院
		room.GameMap.Map[user] = room.GameMap.Map[user].Next()
	}

	room.Broadcast <- fmt.Sprintf("运气： @%s 被车撞伤，入院 5 天", user.Name)

	return nil
}

/**********************************************************************************************/

// 新闻卡处理函数
// 社会发放福利，每位玩家获得2000元
func NewsCardNO1(room *GameRoom, user *User) (err error) {
	room.Player.Do(func(u interface{}) {
		room.Money[u.(*User)] += 2000
	})

	room.Broadcast <- "新闻：社会发放福利，每位玩家获得 1000 元"

	return nil
}

// 住院病人提前出院
func NewsCardNO2(room *GameRoom, user *User) (err error) {
	for u := range room.Hospital {
		room.Hospital[u] = 0
	}

	room.Broadcast <- "新闻：住院病人提前出院"

	return nil
}

// 地震，所有玩家地产降低一级
func NewsCardNO3(room *GameRoom, user *User) (err error) {

	for u, lands := range room.GameMap.PlayerElement {
		for idx, land := range lands {
			if land.Level > int(LAND_LEVEL__START1) {
				room.GameMap.PlayerElement[u][idx].Level--
			}
		}
	}

	room.Broadcast <- "新闻：突发地震，所有玩家地产降一级"

	return nil
}

// 政府公告，所有人交税 5%
func NewsCardNO4(room *GameRoom, user *User) (err error) {

	reply := "新闻：政府公告，所有人交税 5%\n"
	for u, money := range room.Money {
		tax := int64(float64(money) * 0.05)
		room.Money[u] -= tax

		reply += fmt.Sprintf("@%s 缴税 %d 元\n", u.Name, tax)
	}

	reply = strings.TrimSuffix(reply, "\n")

	room.Broadcast <- reply

	return nil
}

// 随机一处房屋闹鬼 地价下跌30%
func NewsCardNO5(room *GameRoom, user *User) (err error) {

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	idxs := r.Perm(len(room.Map))

	var land MapElement
	for j := 0; j < len(room.Map); j++ {
		if room.Map[idxs[j]].Flag == GAME_FLAG__LAND {
			land = room.Map[idxs[j]]

			room.Map[idxs[j]].Fee -= int64(float64(land.Fee) * 0.3)

			break
		}
	}

	for u, lands := range room.GameMap.PlayerElement {
		for i, l := range lands {
			if land.IsEqual(l) {
				room.GameMap.PlayerElement[u][i].Fee -= int64(float64(l.Fee) * 0.3)
				room.GameMap.PlayerElement[u][i].RentFee -= int64(float64(l.RentFee) * 0.3)
			}
		}
	}

	room.Broadcast <- fmt.Sprintf("新闻：%s 处房屋闹鬼，地价下跌30%%", land.Name)

	return nil
}

// 随机地皮地价调涨 30%
func NewsCardNO6(room *GameRoom, user *User) (err error) {

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	idxs := r.Perm(len(room.Map))

	var land MapElement
	for j := 0; j < len(room.Map); j++ {
		if room.Map[idxs[j]].Flag == GAME_FLAG__LAND {
			land := room.Map[idxs[j]]

			room.Map[idxs[j]].Fee += int64(float64(land.Fee) * 0.3)

			break
		}
	}

	for u, lands := range room.GameMap.PlayerElement {
		for i, l := range lands {
			if land.IsEqual(l) {
				room.GameMap.PlayerElement[u][i].Fee += int64(float64(l.Fee) * 0.3)
				room.GameMap.PlayerElement[u][i].RentFee += int64(float64(l.RentFee) * 0.3)
			}
		}
	}

	room.Broadcast <- fmt.Sprintf("新闻：政府公告，%s 处地皮地价调涨 30%%", land.Name)

	return nil
}

// 政府公开补助土地少者1500元
func NewsCardNO7(room *GameRoom, user *User) (err error) {
	// 最少土地
	var min int = 999
	for _, lands := range room.GameMap.PlayerElement {
		tempMin := len(lands)
		if tempMin < min {
			min = tempMin
		}
	}

	// 找到土地最少的玩家
	var users = make([]*User, 0)
	room.Player.Do(func(i interface{}) {
		if len(room.GameMap.PlayerElement[i.(*User)]) == min {
			users = append(users, i.(*User))
			room.Money[i.(*User)] += 1500
		}
	})

	reply := "新闻：政府公开补助土地最少者"

	for _, u := range users {
		reply += " @" + u.Name
	}

	reply += " 获得 1500 元"

	room.Broadcast <- reply

	return nil
}

// 无名慈善家资助，每位玩家可以立即免费赎回一块地
func NewsCardNO8(room *GameRoom, user *User) (err error) {

	for u, lands := range room.GameMap.PlayerElement {
		for i, land := range lands {
			if land.Enable == 0 {
				room.GameMap.PlayerElement[u][i].Enable = 1
				break
			}
		}
	}

	room.Broadcast <- "新闻：无名慈善家资助，每位玩家免费赎回一块地"

	return nil
}

// 百年一遇特大暴雨，所有玩家原地停留一回合
func NewsCardNO9(room *GameRoom, user *User) (err error) {

	room.Player.Do(func(i interface{}) {
		room.Stay[i.(*User)] += 1
	})

	room.Broadcast <- "新闻：百年一遇特大暴雨，所有玩家原地停留一回合"

	return nil
}

// 所有玩家缴纳个人所得税，每块地产 +500 元，
func NewsCardNO10(room *GameRoom, user *User) (err error) {
	for u, lands := range room.GameMap.PlayerElement {
		for range lands {
			room.Money[u] -= 500
		}
	}

	room.Broadcast <- "新闻：所有玩家缴纳个人所得税，每块地产500元"

	return nil
}

// 政府公开表扬土地最多的人获得1000
func NewsCardNO11(room *GameRoom, user *User) (err error) {

	// 最多土地
	var max int
	for _, lands := range room.GameMap.PlayerElement {
		tempMax := len(lands)
		if tempMax > max {
			max = tempMax
		}
	}

	// 找到土地最多的玩家
	var users = make([]*User, 0)
	room.Player.Do(func(i interface{}) {
		if len(room.GameMap.PlayerElement[i.(*User)]) == max {
			users = append(users, i.(*User))
			room.Money[i.(*User)] += 1000
		}
	})

	reply := "新闻：政府公开表扬土地最多的人，"

	for _, u := range users {
		reply += " @" + u.Name
	}

	reply += " 获得 1000 元"

	room.Broadcast <- reply

	return nil
}

// 交通严打，所有行人交款1500
func NewsCardNO12(room *GameRoom, user *User) (err error) {
	var users []*User
	for u, num := range room.Prision {
		if num <= 0 {
			users = append(users, u)
		}
	}

	for _, u := range users {
		room.Money[u] -= 1500
	}

	room.Broadcast <- "新闻：交通严打，所有行人交款1500"

	return nil
}

package richman

import (
	"bytes"
	"container/ring"
	"fmt"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"sort"
	"strconv"
	"sync"
	"text/template"
	"time"
)

type GameRoom struct {
	// 房间唯一标识
	ID string

	// 玩家状态
	User map[*User]GAME_USER_ENUM

	// 游戏玩家及操作顺序(环形)
	Player *ring.Ring

	// 房主
	HomeOwner *User

	// 运气牌池
	LuckCards map[LUCK_CARD_TYPE_ENUM]bool

	// 新闻卡池
	NewsCards map[NEWS_CARD_TYPE_ENUM]bool

	// 房间最大游戏人数
	MaxGameNumber int

	// 房间游戏人数
	GameNumber int

	// 房间使用地图
	Map []MapElement

	// 地图位置
	GameMap GameMap

	// 钱
	Money map[*User]int64

	// 开始通知
	StartChan chan bool

	// 结束通知
	StopChan chan bool

	//房间状态，用于判断房间是否可用，游戏中的不可用，销毁的不可用，创建的时候可用
	RoomStatus GAMEROOM_STATUS_ENUM
	//房间状态锁
	Mu sync.Mutex

	// 监狱
	Prision map[*User]int

	// 医院
	Hospital map[*User]int

	// 原地停留
	Stay map[*User]int

	// 旅馆
	Hotel map[*User]int

	// 投资
	Invest map[*User]map[GAME_FLAG_ENUM]int

	// 移动方向 0 正向 1 反向
	Direction int

	// 广播消息
	Broadcast chan string

	// 骰子是否可用
	DiceStatus bool

	History []*history
}

// 历史记录
type history struct {
	Dice map[*User]int    // 骰子点数
	Logs map[*User]string // 玩家记录
}

func NewRoom(num int) *GameRoom {
	uid, _ := uuid.NewV4()
	room := &GameRoom{
		ID:            uid.String(),
		User:          make(map[*User]GAME_USER_ENUM),
		Money:         make(map[*User]int64),
		LuckCards:     make(map[LUCK_CARD_TYPE_ENUM]bool, int(LUCK_CARD_TYPE__MAX)),
		NewsCards:     make(map[NEWS_CARD_TYPE_ENUM]bool, int(NEWS_CARD_TYPE__MAX)),
		StartChan:     make(chan bool),
		StopChan:      make(chan bool),
		Prision:       make(map[*User]int),
		Hospital:      make(map[*User]int),
		Stay:          make(map[*User]int),
		Hotel:         make(map[*User]int),
		MaxGameNumber: num,
		GameNumber:    0,
		DiceStatus:    false,
		History:       make([]*history, 0),
		Broadcast:     make(chan string),
	}

	gameMap := GameMap{
		PlayerElement: make(map[*User][]MapElement),
		Map:           make(map[*User]*ring.Ring),
	}

	room.GameMap = gameMap

	return room
}

// 初始化运气卡
func (room *GameRoom) InitLuckCardMap() {
	for i := int(LUCK_CARD_TYPE__MIN) + 1; i < int(LUCK_CARD_TYPE__MAX); i++ {
		room.LuckCards[LUCK_CARD_TYPE_ENUM(i)] = true
	}
}

// 初始化新闻卡
func (room *GameRoom) InitNewsCardMap() {
	for i := int(NEWS_CARD_TYPE__MIN) + 1; i < int(NEWS_CARD_TYPE__MAX); i++ {
		room.NewsCards[NEWS_CARD_TYPE_ENUM(i)] = true
	}
}

// 运气卡
func (room *GameRoom) LuckCard() (cardNo int) {
	cardNo = RandNumber() % int(LUCK_CARD_TYPE__MAX-1)
	if room.LuckCards[LUCK_CARD_TYPE_ENUM(cardNo)] == true {
		room.LuckCards[LUCK_CARD_TYPE_ENUM(cardNo)] = false
		return cardNo
	} else {
		var flag = false
		for idx, d := range room.LuckCards {
			if d == true {
				flag = true
				return int(idx)
			}
		}
		if flag == false {
			room.InitLuckCardMap()
			return room.LuckCard()
		}
	}
	return cardNo
}

// 新闻卡
func (room *GameRoom) NewsCard() (cardNo int) {
	cardNo = RandNumber() % int(NEWS_CARD_TYPE__MAX-1)
	if room.NewsCards[NEWS_CARD_TYPE_ENUM(cardNo)] == true {
		room.NewsCards[NEWS_CARD_TYPE_ENUM(cardNo)] = false
		return cardNo
	} else {
		var flag = false
		for idx, d := range room.NewsCards {
			if d == true {
				flag = true
				return int(idx)
			}
		}
		if flag == false {
			room.InitNewsCardMap()
			return room.NewsCard()
		}
	}
	return cardNo
}

// 获取房间状态
func (room *GameRoom) GetRoomStatus() (roomStatus GAMEROOM_STATUS_ENUM) {
	room.Mu.Lock()
	defer room.Mu.Unlock()
	return room.RoomStatus
}

// 设置房间状态
func (room *GameRoom) SetRoomStatus(roomStatus GAMEROOM_STATUS_ENUM) {
	room.Mu.Lock()
	defer room.Mu.Unlock()
	room.RoomStatus = roomStatus
}

// 设置房间地图
func (room *GameRoom) SetRoomMap(playMap []MapElement) {
	room.Map = playMap

}

// 获取当前行动用户
func (room *GameRoom) CurrentPlayer() *User {
	room.Mu.Lock()
	defer room.Mu.Unlock()

	return room.Player.Value.(*User)
}

// 加入房间
func (room *GameRoom) Join(user *User, status GAME_USER_ENUM) {

	// 检查房间是否已经关闭
	if room.RoomStatus == GAMEROOM_STATUS__DISABLE {

		return
	}

	// 检查玩家是否加入其它房间
	if user.Room != nil {
		return
	}

	if status == GAME_USER__PLAYER {
		room.GameNumber++

		// 房间最大人数不为0则为限制人数房间
		if room.MaxGameNumber > 0 {
			if room.GameNumber == room.MaxGameNumber {
				// 如果房间游戏人数和房间容量相同，设置房间可用
				room.RoomStatus = GAMEROOM_STATUS__ENABLE

			} else if room.GameNumber > room.MaxGameNumber {
				// 如果房间游戏人数超过房间容量，游戏房间禁止加入
				room.GameNumber--
				return
			}
		}
	}

	// 玩家和旁观者加入房间
	user.Room = room
	room.User[user] = status

}

// 暂离
func (room *GameRoom) AFK(user *User) {

	if status, ok := room.User[user]; ok && status == GAME_USER__PLAYER {

		room.User[user] = GAME_USER__AFK
	}

	if room.RoomStatus == GAMEROOM_STATUS__GAMESTART {
		// 如果是当前行动用户，
		room.listen(user)
	}
}

// 退出房间
func (room *GameRoom) Quit(user *User) {
	room.Mu.Lock()
	defer room.Mu.Unlock()

	if _, ok := room.User[user]; ok {

		if room.User[user] == GAME_USER__PLAYER || room.User[user] == GAME_USER__AFK {
			room.GameNumber -= 1
		}

		// 释放玩家房间信息
		user.Room = nil

		// 设置玩家为退出状态
		room.User[user] = GAME_USER__EXIT

		// 如果游戏已经开始，清空玩家资产，清除玩家，判断游戏是否结束
		if room.RoomStatus == GAMEROOM_STATUS__GAMESTART {

			if _, ok := room.GameMap.PlayerElement[user]; ok {
				delete(room.GameMap.PlayerElement, user)
			}

			for i := 0; i < room.Player.Len(); i++ {
				if room.Player.Value.(*User) == user {
					room.Player = room.Player.Prev()
					room.Player.Unlink(1)
				}
				room.Player = room.Player.Next()
			}

			room.CheckDone()

		}

	}

	if room.GameNumber <= 0 {
		room.StopChan <- true
	}

}

// 关闭房间
func (room *GameRoom) Close() {

	fmt.Println("游戏结束")

	if room.RoomStatus == GAMEROOM_STATUS__DISABLE {
		return
	}

	room.RoomStatus = GAMEROOM_STATUS__DISABLE
	if len(room.User) > 0 {
		for user := range room.User {
			user.Room = nil
		}
	}
}

// 开始指令
func (room *GameRoom) Start() {
	room.Mu.Lock()
	defer room.Mu.Unlock()

	// 检查房间是否可用
	if room.RoomStatus != GAMEROOM_STATUS__ENABLE {
		room.Broadcast <- "未创建游戏"
		return
	}

	if room.Map == nil {
		room.Broadcast <- "未设置地图"
		return
	}

	// 设置游戏玩家
	room.Player = ring.New(room.GameNumber)

	// 读取地图元素
	rMap := ring.New(len(room.Map))
	for _, m := range room.Map {

		rMap.Value = m
		rMap = rMap.Next()
	}

	// 初始化玩家地产
	room.GameMap.PlayerElement = make(map[*User][]MapElement, room.GameNumber)

	for user, state := range room.User {
		// 如果用户状态为玩家或者托管，则初始化为游戏状态
		if state == GAME_USER__PLAYER || state == GAME_USER__AFK {
			// 初始化游戏玩家
			room.Player.Value = user
			room.Player = room.Player.Next()

			// 初始化玩家金币数量
			room.Money[user] = INITIAL_MONEY

			// 初始化玩家停留回合数
			room.Prision[user] = 0
			room.Hospital[user] = 0
			room.Stay[user] = 0

			// 正向移动
			room.Direction = 0

			// 初始化地图
			room.GameMap.Map[user] = rMap

			room.GameMap.PlayerElement[user] = nil
		}
	}

	// 初始化运气牌池和新闻卡池
	room.InitLuckCardMap()
	room.InitNewsCardMap()

	// 设置房间状态为游戏开始状态
	room.RoomStatus = GAMEROOM_STATUS__GAMESTART

	// 骰子设置为可用状态
	room.DiceStatus = true

	// 发送开始通知
	room.StartChan <- true

	// 监听当前用户
	room.listen(room.Player.Value.(*User))
}

// 监听用户
func (room *GameRoom) listen(user *User) {
	if user != room.Player.Value.(*User) {
		return
	}

	// 检查玩家状态
	if room.User[user] == GAME_USER__AFK {
		room.entrust(user)
		return
	}

	room.Broadcast <- "请玩家 @" + user.Name + " 掷骰子"
}

// 托管操作
func (room *GameRoom) entrust(user *User) {

	if user != room.Player.Value.(*User) {
		return
	}

	dice := ShakeDice()

	room.Broadcast <- "@" + user.Name + " 处于托管状态，自动掷骰子，点数 " + strconv.Itoa(dice)

	// 延迟
	time.Sleep(time.Second)

	room.Dice(dice, user)

}

// 检查是否停留
func (room *GameRoom) isStay(user *User) {
	time.Sleep(time.Second)
	flag := false
	// 检查是否在监狱中
	if num, ok := room.Prision[user]; ok && num > 0 {
		room.Prision[user]--

		room.Broadcast <- "@" + user.Name + " 在监狱中，禁止行动"

		flag = true

	}

	// 检查是否在医院
	if num, ok := room.Hospital[user]; ok && num > 0 {
		room.Hospital[user]--

		room.Broadcast <- "@" + user.Name + " 在医院中，好好休息"

		flag = true
	}

	// 检查是否在旅馆
	if num, ok := room.Hotel[user]; ok && num > 0 {
		room.Hotel[user]--

		room.Broadcast <- "@" + user.Name + " 在旅馆中，停留一回合"

		flag = true
	}

	// 检查是否停留
	if num, ok := room.Stay[user]; ok && num > 0 {

		room.Broadcast <- "@" + user.Name + " 停留一回合"

		// 执行当前地图事件
		location := room.GameMap.Map[user].Value.(MapElement)
		room.do(user, location)

		// 停留事件减一次
		room.Stay[user]--

		flag = true
	}

	if flag {
		if room.Direction == 0 {
			room.Player = room.Player.Next()
		} else {
			room.Player = room.Player.Prev()
		}

		room.isStay(room.Player.Value.(*User))

	}

}

// 玩家掷骰子
func (room *GameRoom) Dice(dice int, user *User) {

	// 房间未开始不做任何操作
	if room.RoomStatus != GAMEROOM_STATUS__GAMESTART {
		fmt.Println("房间关闭")
		return
	}

	if !room.DiceStatus {
		fmt.Println("骰子不可用")
		return
	}

	// 非当前游戏玩家不做任何操作
	if room.Player.Value != nil && room.Player.Value.(*User) != user {
		fmt.Println(user.Name, "非当前玩家")
		return
	}

	// 骰子设置为不可用状态
	room.DiceStatus = false

	// 检查可赎回地产
	var lands []MapElement
	for _, land := range room.GameMap.PlayerElement[user] {
		halfFee := int64(room.Money[user] / 2)

		if halfFee >= land.Fee {
			break
		}

		if land.Enable == 0 {
			lands = append(lands, land)
		}
	}
	room.LandRedeem(user, lands)

	// 移动方向
	if room.Direction == 1 {
		dice = -dice
	}

	// 玩家移动
	room.userMove(dice, user)

	// 将环形队列指向下一个玩家
	if room.Direction == 0 {
		room.Player = room.Player.Next()
	} else {
		room.Player = room.Player.Prev()
	}

	// 队列移动后检查玩家是否可以行动
	room.isStay(room.Player.Value.(*User))

	// 检查游戏是否结束
	room.CheckDone()

	// 等待
	time.Sleep(time.Second)

	// 将骰子设置为可用状态
	room.DiceStatus = true

	// 监听下一个用户
	room.listen(room.Player.Value.(*User))
}

// 用户移动到指定位置
func (room *GameRoom) userMove(step int, user *User) {
	// 非当前游戏用户不做任何操作
	if user != room.Player.Value.(*User) {
		return
	}

	// 等待
	time.Sleep(time.Second)

	// 用户移动 setp 步
	room.GameMap.Map[user] = room.GameMap.Map[user].Move(step)

	// 获取用户位置
	location := room.GameMap.Map[user].Value.(MapElement)

	logrus.WithFields(logrus.Fields{
		"操作用户": user.Name,
		"当前玩家": room.Player.Value.(*User).Name,
		"位置":   location.Name,
	}).Infof("移动 %d 步", step)

	// + 检查是否经过起点，经过起点发钱
	room.do(user, location)

}

// 掷完骰子后，检查需要做的动作
func (room *GameRoom) do(user *User, location MapElement) {
	switch location.Flag {
	case GAME_FLAG__LAND:
		// 位置为土地，其他人的土地支付佣金，自己的地升级土地，空地购买土地
		flag := false
		for u, lands := range room.GameMap.PlayerElement {
			for _, land := range lands {
				// 判断土地是否可用
				if location.IsEqual(land) {
					flag = true
					if u == user {
						logrus.WithFields(logrus.Fields{
							"玩家": user.Name,
							"位置": land.Name,
							"等级": land.Level,
						}).Info("升级地产")

						// 自己的土地，升级
						room.UpdateLand(user, land)
						return
					} else {
						logrus.WithFields(logrus.Fields{
							"地产所有者": u.Name,
							"玩家":    user.Name,
							"位置":    land.Name,
							"费用":    land.RentFee,
						}).Info("支付租金")

						// 别人的土地，支付租金
						room.Payment(user, u, land)
						return
					}
				}
			}
		}
		// 购买空地
		if !flag {
			logrus.WithFields(logrus.Fields{
				"玩家": user.Name,
				"位置": location.Name,
			}).Info("购买空地")

			room.BuyLand(user, location)
		}
	case GAME_FLAG__LUCK:
		// 位置为运气卡
		cardNo := LUCK_CARD_TYPE_ENUM(room.LuckCard())
		LuckRules[cardNo](room, user)

	case GAME_FLAG__NEWS:
		// 位置为新闻卡
		cardNo := NEWS_CARD_TYPE_ENUM(room.NewsCard())

		NewsRules[cardNo](room, user)

	case GAME_FLAG__SECURITIES_CENTER:
		// 位置为证券中心,获得你拥有地产数量*500元的奖励
		landNum := len(room.GameMap.PlayerElement[user])

		room.Money[user] += int64(500 * landNum)

		room.Broadcast <- fmt.Sprintf("@%s 到达证券中心，获得 %d 元奖励", user.Name, int64(500*landNum))

	case GAME_FLAG__PRISION:
		// 位置为监狱
		room.Prision[user] += 1
		room.Broadcast <- "@" + user.Name + " 被捕入狱，停留一回合"
	case GAME_FLAG__HOTEL:
		// 位置为饭店，付款1000元并停留一回合
		room.Money[user] += 1000
		room.Hotel[user] = 1

		room.Broadcast <- fmt.Sprintf("@%s 走到了饭店，付款 1000 元休息一回合", user.Name)
	case GAME_FLAG__TAX_CENTER:
		// 位置为税务
		// 每块地交税 300
		landNum := len(room.GameMap.PlayerElement[user])
		room.Money[user] -= int64(300 * landNum)
		room.Broadcast <- fmt.Sprintf("@%s 走到税务中心，交税 %d", user.Name, 300*landNum)
	case GAME_FLAG__HOSPITAL:
		// 医院
		room.Hospital[user] += 1
		room.Broadcast <- "@" + user.Name + " 进医院看病，停留一回合"

	case GAME_FLAG__POWER_STATION:
		// 电力公司，交电费，基数300，每处房产 +300 元
		landNum := len(room.GameMap.PlayerElement[user])

		room.Money[user] -= int64(300*landNum + 300)

		room.Broadcast <- fmt.Sprintf("@%s 到达电力公司，需缴电费 %d", user.Name, int64(300*landNum+300))

	case GAME_FLAG__CONSTRUCTION_COMPANY:
		// 地产公司，免费升级一块地产
		if len(room.GameMap.PlayerElement[user]) <= 0 {
			room.Broadcast <- fmt.Sprintf("@%s 到达地产公司，但是没有地产可升级", user.Name)
			return
		}
		for i, land := range room.GameMap.PlayerElement[user] {
			if land.Level >= int(LAND_LEVEL__MAX) {
				continue
			}

			room.GameMap.PlayerElement[user][i].Level++
			room.GameMap.PlayerElement[user][i].RentFee += int64(float64(land.Fee) * 0.8)
			room.GameMap.PlayerElement[user][i].Fee += int64(float64(land.Fee) * 0.5)

			room.Broadcast <- fmt.Sprintf("@%s 到达地产公司，免费升级地产 %s", user.Name, land.Name)

			break
		}

	case GAME_FLAG__CONTINENTAL_TRANSPORTION:
		// 大陆运输 获得2000元
		room.Money[user] += 2000

		room.Broadcast <- "@" + user.Name + " 到大陆运输运输货物，赚到2000元"

	case GAME_FLAG__AIR_TRANSPORTION:
		// 航空运输 获得10000元
		room.Money[user] += 10000

		room.Broadcast <- "@" + user.Name + " 到航空运输运输货物，赚到10000元"

	case GAME_FLAG__OCEAN_TRANSPORTION:
		// 大洋运输 获得5000元
		room.Money[user] += 5000

		room.Broadcast <- "@" + user.Name + " 到大洋运输运输货物，赚到5000元"

	case GAME_FLAG__TV_STATION:
		// 电视台，播报两次新闻
		room.Broadcast <- "@" + user.Name + " 到达电视台，播放两条新闻"

		NewsRules[NEWS_CARD_TYPE_ENUM(room.NewsCard())](room, user)

		time.Sleep(500 * time.Millisecond)

		NewsRules[NEWS_CARD_TYPE_ENUM(room.NewsCard())](room, user)

	case GAME_FLAG__SEWAGE_TREATMENT:
		// 污水处理厂，处理废水花费800
		room.Money[user] -= 2000

		room.Broadcast <- "@" + user.Name + " 到污水处理厂，处理废水花费2000"

	case GAME_FLAG__START_POINT:
		room.Broadcast <- "@" + user.Name + " 到达起点，奖励 " + strconv.Itoa(SEND_MONY)
	}
}

// 获取升级地产费用
func GetLandUpdateFee(land MapElement) int64 {
	return land.Fee + int64(float64(land.Level)*0.2+float64(land.Level))*land.Fee
}

// 升级地产
func (room *GameRoom) UpdateLand(user *User, land MapElement) {
	// 抵押房产或者房产不可用
	if land.Enable == 0 {
		return
	}
	// 确认是否升级，最高等级不升级
	if land.Level == int(LAND_LEVEL__MAX) {
		room.Broadcast <- fmt.Sprintf("地产所有人 @%s，地产已经是最高等级，不可升级", user.Name)
		return
	}
	// 升级地产
	// 升级时所需费用
	fee := GetLandUpdateFee(land)

	// 钱够，扣除钱
	if room.Money[user] > fee {
		// 钱足够，扣钱
		room.Money[user] -= fee
		for idx, l := range room.GameMap.PlayerElement[user] {
			if l.IsEqual(land) {
				room.GameMap.PlayerElement[user][idx].Level++
				room.GameMap.PlayerElement[user][idx].RentFee += int64(float64(land.Fee) * 0.8)
				room.GameMap.PlayerElement[user][idx].Fee += int64(float64(l.Fee) * 0.5)
			}
		}

		reply := fmt.Sprintf("到达 『%s』，地产所有人 @%s 支付 %d 费用升级地产，当前地产等级为 %d", land.Name, user.Name, fee, land.Level+1)
		room.Broadcast <- reply
	} else {
		// 钱不够，不升级地产
		reply := fmt.Sprintf("到达 『%s』，地产所有人 @%s，需支付 %d 费用升级地产，您当前资金为 %d，升级资金不足", land.Name, user.Name, fee, room.Money[user])
		room.Broadcast <- reply
	}
}

// 支付佣金
func (room *GameRoom) Payment(user, owner *User, land MapElement) {
	// 抵押房产或者房产不可用不用支付
	if land.Enable == 0 {
		return
	}

	// 检查是否在监狱中
	if num, ok := room.Prision[owner]; ok && num > 0 {

		room.Broadcast <- fmt.Sprintf("到达 『%s』，@%s 在监狱中，不需要支付租金", land.Name, owner.Name)

		return
	}

	// 检查是否在医院
	if num, ok := room.Hospital[owner]; ok && num > 0 {

		room.Broadcast <- fmt.Sprintf("到达 『%s』，@%s 在医院中，不需要支付租金", land.Name, owner.Name)

		return
	}

	// 路过别人的地，需要支付租金
	if room.Money[user] > land.RentFee {
		// 资金足够，支付租金
		room.Money[user] -= land.RentFee
		room.Money[owner] += land.RentFee

		reply := fmt.Sprintf("到达 『%s』，地产属于 @%s，需要支付租金 %d", land.Name, owner.Name, land.RentFee)
		room.Broadcast <- reply
	} else {
		// 资金不足
		umap := room.GameMap.PlayerElement[user]
		if len(umap) <= 0 {
			// 没有房产，直接破产
			room.Broke(user)
			return
		}

		room.Broadcast <- fmt.Sprintf("到达 『%s』，当前资金不足，需要抵押房产以支付租金", land.Name)

		// 获取需要抵押的房产
		sort.Slice(umap, func(i, j int) bool {
			return umap[i].Fee < umap[j].Fee
		})

		sort.Slice(umap, func(i, j int) bool {
			return umap[i].Level < umap[j].Level
		})

		var cost = room.Money[user]
		var lands []MapElement
		for _, land := range umap {
			// 房产处于抵押状态
			if land.Enable == 0 {
				continue
			}
			cost += land.Fee
			lands = append(lands, land)

			if cost >= land.RentFee {
				break
			}
		}

		// 抵押地产
		room.LandImpawn(user, lands)

		// 抵押后费用足够，支付费用
		if room.Money[user] > land.RentFee {
			// 资金足够，支付租金
			room.Money[user] -= land.RentFee
			room.Money[owner] += land.RentFee

			reply := fmt.Sprintf("支付租金 %d，资产剩余 %d", land.RentFee, room.Money[user])
			room.Broadcast <- reply
		} else {
			// 费用不足，破产
			room.Broke(user)
			return
		}
	}
}

// 购买地产
func (room *GameRoom) BuyLand(user *User, location MapElement) {
	// 土地不在玩家地产中
	// 空地
	if room.Money[user] > location.Fee {
		// 钱足够，扣钱
		room.Money[user] -= location.Fee

		location.Level = 1
		location.RentFee = int64(float64(location.Fee) * 0.8)
		location.Enable = 1
		room.GameMap.PlayerElement[user] = append(room.GameMap.PlayerElement[user], location)

		room.Broadcast <- fmt.Sprintf("到达 『%s』，玩家 @%s 购买地产，花费 %d", location.Name, user.Name, location.Fee)
	} else {
		// 钱不够，不做任何操作
		room.Broadcast <- fmt.Sprintf("到达 『%s』，资金不够了，买不起地", location.Name)
	}
}

// 用户地产抵押
func (room *GameRoom) LandImpawn(user *User, mapList []MapElement) {

	for _, ml := range mapList {

		// 房产处于抵押状态
		if ml.Enable == 0 {
			continue
		}

		for index, land := range room.GameMap.PlayerElement[user] {
			if land.IsEqual(ml) {
				room.Money[user] += land.Fee
				room.GameMap.PlayerElement[user][index].Enable = 0

				//reply := fmt.Sprintf("@%s 的『%s』 处房产被抵押，获得资金 %d", user.Name, land.Name, land.Fee)
				//room.Broadcast <- reply
			}
		}
	}
}

// 用户地产赎回
func (room *GameRoom) LandRedeem(user *User, mapList []MapElement) {
	for _, ml := range mapList {
		// 房产未被抵押
		if ml.Enable == 1 {
			continue
		}

		for idx, land := range room.GameMap.PlayerElement[user] {

			if land.IsEqual(ml) {
				// 判断地产是否是同一个地产，如果是同一个地产，则把地产赎回，并根据地产计算费用
				// 支付费用
				if room.Money[user] >= room.GameMap.PlayerElement[user][idx].Fee {
					room.Money[user] = room.Money[user] - room.GameMap.PlayerElement[user][idx].Fee
					//把地产变为可用
					room.GameMap.PlayerElement[user][idx].Enable = 1

					//reply := fmt.Sprintf("@%s 的『%s』 处房产被赎回，花费资金 %d", user.Name, land.Name, land.Fee)
					//room.Broadcast <- reply
				}
			}
		}

	}

	return
}

// 用户破产
func (room *GameRoom) Broke(user *User) {
	fmt.Println(user.Name, "破产")
	// 玩家退出游戏
	if room.Player.Value.(*User) == user {
		room.Player = room.Player.Prev()
		room.Player.Unlink(1)
		room.Player = room.Player.Next()
	}

	// 游戏人数减一
	room.GameNumber -= 1

	// 玩家设置为破产状态
	room.User[user] = GAME_USER__BROKE

	// 释放玩家所在房间
	user.Room = nil

	// 清空玩家地产
	if _, ok := room.GameMap.PlayerElement[user]; ok {
		delete(room.GameMap.PlayerElement, user)
	}

	room.Broadcast <- "玩家 @" + user.Name + " 破产"

}

// 判断游戏输赢
func (room *GameRoom) CheckDone() {

	// 判断玩家状态
	if room.Player.Len() <= 1 {
		user := room.Player.Value.(*User)

		reply := fmt.Sprintf("游戏结束\n恭喜玩家 @%s 获胜，最终资产为 %d", user.Name, room.Money[user])
		room.Broadcast <- reply

		room.Close()
	}
}

// 生成房间状态图片
func (room *GameRoom) genRoomStateImg() (path string, err error) {

	type data struct {
		State    GAME_USER_ENUM
		Money    int64              // 玩家资金
		Lands    []MapElement       // 玩家房产
		Location map[MapElement]int // 如果在停留状态，标明停留次数
	}

	var datas = make(map[*User]data, len(room.User))
	for user, state := range room.User {

		var loc = make(map[MapElement]int, 1)
		land := room.GameMap.Map[user].Value.(MapElement)

		if room.Prision[user] > 0 {
			loc[land] = room.Prision[user]
		} else if room.Hospital[user] > 0 {
			loc[land] = room.Hospital[user]
		} else if room.Hotel[user] > 0 {
			loc[land] = room.Hotel[user]
		} else if room.Stay[user] > 0 {
			loc[land] = room.Stay[user]
		} else {
			loc[land] = 0
		}

		d := data{
			State:    state,
			Money:    room.Money[user],
			Lands:    room.GameMap.PlayerElement[user],
			Location: loc,
		}

		datas[user] = d
	}

	buf := new(bytes.Buffer)

	tmpl, err := template.New("room_status").Parse(RoomStateHtml)
	if err != nil {
		return "", err
	}

	tmpl.Execute(buf, &datas)

	opt := ImageOptions{
		Input:      "-",
		Html:       buf.String(),
		Format:     "svg",
		Output:     "1.svg",
		Width:      960,
		Quality:    80,
		BinaryPath: "wkhtmltoimage.exe",
	}

	_, err = GenerateImage(&opt)
	if err != nil {
		return "", err
	}

	return "1.svg", nil
}

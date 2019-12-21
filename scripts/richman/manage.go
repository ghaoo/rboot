package richman

import (
	"fmt"
	"github.com/ghaoo/rboot"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

var mapNo = 1

var GM *Game

func Go() {
	// 创建游戏
	GM = New()

}

// 创建游戏
func CreateGameRoom(in rboot.Message, bot *rboot.Robot) (msg []rboot.Message) {
	if in.Mate["GroupMsg"].(bool) {
		if !in.Mate["AtMe"].(bool) {
			return nil
		}
	}

	room, ok := GM.GameRooms[in.From.ID]
	// 如果房间存在并且房间不处于关闭状态，则不作任何操作
	if ok && room.RoomStatus != GAMEROOM_STATUS__DISABLE {
		return []rboot.Message{{Content: "游戏中..."}}
	}

	// 获取用户唯一标识
	userID := in.Sender.ID
	// 根据用户ID获取用户信息
	user, ok := GM.User[userID]

	if ok && user.Room != nil {
		return []rboot.Message{{Content: "正在游戏中，不允许创建游戏"}}
	} else {
		user = &User{
			ID:   in.Sender.ID,
			Name: in.Sender.Name,
		}
	}

	// 创建有限人数游戏时需要提供人数 num，仅限群组
	num := 0
	var err error

	if len(bot.Match) >= 2 {

		if !in.Mate["GroupMsg"].(bool) {
			return []rboot.Message{{Content: "多人游戏请在群组创建。"}}
		}

		num, err = strconv.Atoi(bot.Match[1])
		if num < 2 {
			return []rboot.Message{{Content: "人数小于2，创建游戏失败。\n单人游戏可私聊创建。"}}
		}
		if err != nil {
			logrus.Error(err)
			return
		}
	}

	room = NewRoom(num)

	// 使用第一张地图
	room.SetRoomMap(GM.Map[mapNo])

	room.HomeOwner = user
	// 将用户注册到游戏房间
	room.Join(user, GAME_USER__PLAYER)

	// 房间序号+1并注册到游戏
	//roomNum := len(GM.GameRooms) + 1
	// 使用群ID或者个人ID作为房间编号
	GM.GameRooms[in.From.ID] = room
	GM.User[in.Sender.ID] = user

	if !in.Mate["GroupMsg"].(bool) {
		// 将ai加入游戏
		ai := &User{
			ID:   in.To.ID,
			Name: in.To.Name,
		}
		room.Join(ai, GAME_USER__PLAYER)
		// 将AI设置为托管状态
		room.AFK(ai)

		msg = []rboot.Message{
			{Content: fmt.Sprintf("游戏创建成功！")},
		}
	} else {
		msg = []rboot.Message{
			{Content: fmt.Sprintf("游戏创建成功！\n其他参加游戏的同学请在群里@我并回复 “加入游戏”")},
		}
	}

	// 房间设置为可用状态
	room.RoomStatus = GAMEROOM_STATUS__ENABLE

	go func() {
		// 等待玩家进入，超时取消
		timer := time.NewTimer(3 * time.Minute)
	Loop:
		for {
			select {

			case broadcast := <-room.Broadcast:
				// 监听广播消息
				bot.SendText(broadcast, in.From)

			case <-timer.C:

				room.Close()

				bot.SendText("超时！游戏房间取消", in.From)

				break Loop

			case <-room.StartChan:

				if !timer.Stop() {
					<-timer.C
				}

				bot.SendText("请玩家做好准备，游戏开始", in.From)

			case <-room.StopChan:

				room.Close()

				bot.SendText("游戏房间关闭", in.From)

				break Loop
			}
		}
	}()

	return
}

func StopGame(in rboot.Message) []rboot.Message {
	room, _ := GM.GameRooms[in.From.ID]

	room.StopChan <- true

	return nil
}

// 监听游戏开始后用户操作
/*func listenUser(room *GameRoom, user *User) {

	// 开始监听用户操作，等待30秒
	timer := time.NewTimer(30 * time.Second)
	for {
		select {
		// 若没有任何动作则将用户设置为托管状态
		case <-timer.C:
			room.AFK(user)
			// 通知
			room.Broadcast <- "玩家 " + user.Name + " 30分钟未操作，系统托管"

			dice := ShakeDice()
			notify := Notify{
				Type:    GAME_NOTIFY__TYPE__SHAKE_DICE,
				Player:  user,
				Message: NotifyDice{DiceNum: dice},
			}
			room.Notify <- notify
		}
	}
}*/

// 加入游戏
func JoinRoom(in rboot.Message, bot *rboot.Robot) []rboot.Message {

	if !in.Mate["GroupMsg"].(bool) && !in.Mate["AtMe"].(bool) {
		return nil
	}

	room, ok := GM.GameRooms[in.From.ID]
	// 如果房间不存在或者游戏已经结束，提示创建游戏
	if !ok || room.RoomStatus == GAMEROOM_STATUS__DISABLE {
		return []rboot.Message{{Content: "游戏不存在，请先创建游戏"}}
	}

	// 如果房间游戏已经开始，检查是否是离开状态，如果是
	if room.RoomStatus == GAMEROOM_STATUS__GAMESTART {
		return nil
	}

	// 获取用户唯一标识
	userID := in.Sender.ID
	// 根据用户ID获取用户信息
	user, ok := GM.User[userID]

	// 判断用户是否在游戏中
	if ok && user.Room != nil {
		return []rboot.Message{{Content: "@" + in.Sender.Name + " 正在游戏中，不允许加入其他游戏"}}
	} else {
		user = &User{
			ID:   in.Sender.ID,
			Name: in.Sender.Name,
		}
	}

	GM.User[userID] = user

	room.Join(user, GAME_USER__PLAYER)

	return []rboot.Message{{Content: "@" + in.Sender.Name + " 加入游戏成功，请等待游戏开始"}}
}

// 暂离游戏
func AFK(in rboot.Message, bot *rboot.Robot) []rboot.Message {
	if in.Mate["GroupMsg"].(bool) {
		if !in.Mate["AtMe"].(bool) {
			return nil
		}
	}

	room, ok := GM.GameRooms[in.From.ID]
	// 如果房间不存在或者游戏已经结束，提示创建游戏
	if !ok || room.RoomStatus == GAMEROOM_STATUS__DISABLE {
		return []rboot.Message{{Content: "游戏不存在，请先创建游戏"}}
	}

	// 获取用户唯一标识
	userID := in.Sender.ID
	// 根据用户ID获取用户信息
	user, ok := GM.User[userID]

	if !ok || user.Room == nil {
		return []rboot.Message{{Content: "@" + in.Sender.Name + " 未参加游戏"}}
	}

	if room.User[user] != GAME_USER__PLAYER {
		return []rboot.Message{{Content: "@" + in.Sender.Name + " 已经处于托管状态"}}
	}

	go room.AFK(user)

	return []rboot.Message{{Content: "@" + in.Sender.Name + " 暂时离开，游戏将由系统托管"}}
}

// 退出游戏房间
func QuitGame(in rboot.Message, bot *rboot.Robot) []rboot.Message {
	if in.Mate["GroupMsg"].(bool) {
		if !in.Mate["AtMe"].(bool) {
			return nil
		}
	}

	room, ok := GM.GameRooms[in.From.ID]
	// 如果房间不存在或者游戏已经结束，提示创建游戏
	if !ok || room.RoomStatus == GAMEROOM_STATUS__DISABLE {
		return []rboot.Message{{Content: "游戏不存在，请先创建游戏"}}
	}

	// 获取用户唯一标识
	userID := in.Sender.ID
	// 根据用户ID获取用户信息
	user, ok := GM.User[userID]

	if !ok || user.Room == nil {
		return []rboot.Message{{Content: "@" + in.Sender.Name + " 未参加游戏"}}
	}

	room.Quit(user)

	reply := fmt.Sprintf("@%s 退出游戏", in.Sender.Name)

	return []rboot.Message{{Content: reply}}
}

// 开始游戏
func StartGame(in rboot.Message, bot *rboot.Robot) []rboot.Message {
	if in.Mate["GroupMsg"].(bool) {
		if !in.Mate["AtMe"].(bool) {
			return nil
		}
	}

	room, ok := GM.GameRooms[in.From.ID]
	// 如果房间不存在或者游戏已经结束，提示创建游戏
	if !ok || room.RoomStatus == GAMEROOM_STATUS__DISABLE {
		return []rboot.Message{{Content: "游戏不存在，请先创建游戏"}}
	}

	// 获取用户唯一标识
	userID := in.Sender.ID
	// 根据用户ID获取用户信息
	user, ok := GM.User[userID]

	if !ok || user.Room == nil {
		return []rboot.Message{{Content: "@" + in.Sender.Name + " 未参加游戏"}}
	}

	if room.GameNumber < 2 {
		return []rboot.Message{{Content: "游戏人数不足！游戏人数最少为2人..."}}
	}

	if room.MaxGameNumber > 0 && room.MaxGameNumber > room.GameNumber {
		reply := fmt.Sprintf("游戏人数不足！还差 %d 人才能开始游戏...", room.MaxGameNumber-room.GameNumber)
		return []rboot.Message{{Content: reply}}
	}

	if room.RoomStatus != GAMEROOM_STATUS__ENABLE {
		return []rboot.Message{{Content: "开始游戏失败，房间状态为不可用状态"}}
	}

	if room.HomeOwner != user {
		return []rboot.Message{{Content: "@" + room.HomeOwner.Name + " 有人催你开始游戏"}}
	}

	// 游戏开始
	room.Start()

	// 监听当前用户
	// player := room.CurrentPlayer()

	return nil
}

// 掷骰子
func Dice(in rboot.Message, bot *rboot.Robot) []rboot.Message {

	room, ok := GM.GameRooms[in.From.ID]
	// 如果房间不存在或者游戏已经结束，不做任何操作
	if !ok || room.RoomStatus == GAMEROOM_STATUS__DISABLE {
		return nil
	}

	if !room.DiceStatus {
		return nil
	}

	// 获取用户唯一标识
	userID := in.Sender.ID
	// 根据用户ID获取用户信息
	user, ok := GM.User[userID]
	// 非游戏玩家不做任何操作
	if !ok || user.Room == nil {
		return nil
	}

	// 非当前游戏玩家不做任何操作
	if room.Player.Value != nil && room.Player.Value.(*User) != user {
		return nil
	}

	// 获取骰子点数
	if len(bot.Match) < 2 {
		return nil
	}

	var dice int

	dice, err := strconv.Atoi(bot.Match[1])
	if err != nil {
		logrus.Errorf("骰子解析失败：%v", err)
		return nil
	}

	dice -= 3

	room.Dice(dice, user)

	return nil
}

// 查看用户信息
func Look(in rboot.Message, bot *rboot.Robot) []rboot.Message {
	room, ok := GM.GameRooms[in.From.ID]
	// 如果房间不存在或者游戏已经结束，不做任何操作
	if !ok || room.RoomStatus == GAMEROOM_STATUS__DISABLE {
		fmt.Println("房间关闭中")
		return nil
	}

	// 获取用户唯一标识
	userID := in.Sender.ID
	// 根据用户ID获取用户信息
	user, ok := GM.User[userID]
	// 非游戏玩家不做任何操作
	if !ok || user.Room == nil {
		fmt.Println("非游戏玩家")
		return nil
	}
	reply := "用户: " + user.Name + "\n资产: " + strconv.Itoa(int(room.Money[user])) + "\n地产: "
	for _, land := range room.GameMap.PlayerElement[user] {
		reply += land.Name + " "
	}

	return []rboot.Message{{Content: reply}}
}

// 查看房间玩家状态
func RoomInfo(in rboot.Message, bot *rboot.Robot) []rboot.Message {

	room, ok := GM.GameRooms[in.From.ID]
	// 如果房间不存在或者游戏已经结束，不做任何操作
	if !ok || room.RoomStatus != GAMEROOM_STATUS__GAMESTART {
		fmt.Println("房间未开始")
		return nil
	}

	reply := ""

	for user, lands := range room.GameMap.PlayerElement {
		reply += "用户: " + user.Name + "\n位置: " + room.GameMap.Map[user].Value.(MapElement).Name + "\n资产: " + strconv.Itoa(int(room.Money[user])) + "\n地产: "
		for _, land := range lands {
			reply += fmt.Sprintf("%s[%d] ", land.Name, land.Level)
		}

		reply += "\n\n"
	}

	reply = strings.TrimSuffix(reply, "\n\n")

	return []rboot.Message{{Content: reply}}

}

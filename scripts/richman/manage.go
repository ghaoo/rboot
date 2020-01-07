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
func CreateGameRoom(in rboot.Message, bot *rboot.Robot) (msg *rboot.Message) {
	if group, _ := strconv.ParseBool(in.Header.Get("GroupMsg")); group {
		if atme, _ := strconv.ParseBool(in.Header.Get("AtMe")); !atme {
			return nil
		}
	}

	room, ok := GM.GameRooms[in.From]
	// 如果房间存在并且房间不处于关闭状态，则不作任何操作
	if ok && room.RoomStatus != GAMEROOM_STATUS__DISABLE {
		return rboot.NewMessage("游戏中...")
	}

	// 获取用户唯一标识
	userID := in.Sender
	// 根据用户ID获取用户信息
	user, ok := GM.User[userID]

	if ok && user.Room != nil {
		return rboot.NewMessage("正在游戏中，不允许创建游戏")
	} else {
		user = &User{
			ID:   in.Sender,
			Name: bot.GetUserName(in.Sender),
		}
	}

	// 创建有限人数游戏时需要提供人数 num，仅限群组
	num := 0
	var err error

	args := bot.Args

	if len(args) >= 2 {

		if group, _ := strconv.ParseBool(in.Header.Get("GroupMsg")); !group {
			return rboot.NewMessage("多人游戏请在群组创建。")
		}

		num, err = strconv.Atoi(args[1])
		if num < 2 {
			return rboot.NewMessage("人数小于2，创建游戏失败。\n单人游戏可私聊创建。")
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
	GM.GameRooms[in.From] = room
	GM.User[in.Sender] = user

	if group, _ := strconv.ParseBool(in.Header.Get("GroupMsg")); !group {
		// 将ai加入游戏
		ai := &User{
			ID:   in.To,
			Name: bot.GetUserName(in.To),
		}
		room.Join(ai, GAME_USER__PLAYER)
		// 将AI设置为托管状态
		room.AFK(ai)

		msg = rboot.NewMessage("游戏创建成功！")
	} else {
		msg = rboot.NewMessage("游戏创建成功！\n其他参加游戏的同学请在群里@我并回复 “加入游戏”")
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

func StopGame(in rboot.Message) *rboot.Message {
	room, _ := GM.GameRooms[in.From]

	room.StopChan <- true

	return nil
}

// 加入游戏
func JoinRoom(in rboot.Message, bot *rboot.Robot) *rboot.Message {

	if group, _ := strconv.ParseBool(in.Header.Get("GroupMsg")); !group {
		return nil
	}

	if atme, _ := strconv.ParseBool(in.Header.Get("AtMe")); !atme {
		return nil
	}

	room, ok := GM.GameRooms[in.From]
	// 如果房间不存在或者游戏已经结束，提示创建游戏
	if !ok || room.RoomStatus == GAMEROOM_STATUS__DISABLE {
		return rboot.NewMessage("游戏不存在，请先创建游戏")
	}

	// 如果房间游戏已经开始，检查是否是离开状态，如果是
	if room.RoomStatus == GAMEROOM_STATUS__GAMESTART {
		return nil
	}

	// 获取用户唯一标识
	userID := in.Sender
	// 根据用户ID获取用户信息
	user, ok := GM.User[userID]

	// 判断用户是否在游戏中
	if ok && user.Room != nil {
		return rboot.NewMessage("@" + bot.GetUserName(in.Sender) + " 正在游戏中，不允许加入其他游戏")
	} else {
		user = &User{
			ID:   in.Sender,
			Name: bot.GetUserName(in.Sender),
		}
	}

	GM.User[userID] = user

	room.Join(user, GAME_USER__PLAYER)

	return rboot.NewMessage("@" + bot.GetUserName(in.Sender) + " 加入游戏成功，请等待游戏开始")
}

// 暂离游戏
func AFK(in rboot.Message, bot *rboot.Robot) *rboot.Message {

	if group, _ := strconv.ParseBool(in.Header.Get("GroupMsg")); group {
		if atme, _ := strconv.ParseBool(in.Header.Get("AtMe")); !atme {
			return nil
		}
	}

	room, ok := GM.GameRooms[in.From]
	// 如果房间不存在或者游戏已经结束，提示创建游戏
	if !ok || room.RoomStatus == GAMEROOM_STATUS__DISABLE {
		return rboot.NewMessage("游戏不存在，请先创建游戏")
	}

	// 获取用户唯一标识
	userID := in.Sender
	// 根据用户ID获取用户信息
	user, ok := GM.User[userID]

	sender := bot.GetUserName(in.Sender)

	if !ok || user.Room == nil {
		return rboot.NewMessage("@" + sender + " 未参加游戏")
	}

	if room.User[user] != GAME_USER__PLAYER {
		return rboot.NewMessage("@" + sender + " 已经处于托管状态")
	}

	go room.AFK(user)

	return rboot.NewMessage("@" + sender + " 暂时离开，游戏将由系统托管")
}

// 退出游戏房间
func QuitGame(in rboot.Message, bot *rboot.Robot) *rboot.Message {

	if group, _ := strconv.ParseBool(in.Header.Get("GroupMsg")); group {
		if atme, _ := strconv.ParseBool(in.Header.Get("AtMe")); !atme {
			return nil
		}
	}

	room, ok := GM.GameRooms[in.From]
	// 如果房间不存在或者游戏已经结束，提示创建游戏
	if !ok || room.RoomStatus == GAMEROOM_STATUS__DISABLE {
		return rboot.NewMessage("游戏不存在，请先创建游戏")
	}

	// 获取用户唯一标识
	userID := in.Sender
	// 根据用户ID获取用户信息
	user, ok := GM.User[userID]

	sender := bot.GetUserName(in.Sender)

	if !ok || user.Room == nil {
		return rboot.NewMessage("@" + sender + " 未参加游戏")
	}

	room.Quit(user)

	reply := fmt.Sprintf("@%s 退出游戏", sender)

	return rboot.NewMessage(reply)
}

// 开始游戏
func StartGame(in rboot.Message, bot *rboot.Robot) *rboot.Message {

	if group, _ := strconv.ParseBool(in.Header.Get("GroupMsg")); group {
		if atme, _ := strconv.ParseBool(in.Header.Get("AtMe")); !atme {
			return nil
		}
	}

	room, ok := GM.GameRooms[in.From]
	// 如果房间不存在或者游戏已经结束，提示创建游戏
	if !ok || room.RoomStatus == GAMEROOM_STATUS__DISABLE {
		return rboot.NewMessage("游戏不存在，请先创建游戏")
	}

	// 获取用户唯一标识
	userID := in.Sender
	// 根据用户ID获取用户信息
	user, ok := GM.User[userID]

	if !ok || user.Room == nil {
		return rboot.NewMessage("@" + bot.GetUserName(in.Sender) + " 未参加游戏")
	}

	if room.GameNumber < 2 {
		return rboot.NewMessage("游戏人数不足！游戏人数最少为2人...")
	}

	if room.MaxGameNumber > 0 && room.MaxGameNumber > room.GameNumber {
		reply := fmt.Sprintf("游戏人数不足！还差 %d 人才能开始游戏...", room.MaxGameNumber-room.GameNumber)
		return rboot.NewMessage(reply)
	}

	if room.RoomStatus != GAMEROOM_STATUS__ENABLE {
		return rboot.NewMessage("开始游戏失败，房间状态为不可用状态")
	}

	if room.HomeOwner != user {
		return rboot.NewMessage("@" + room.HomeOwner.Name + " 有人催你开始游戏")
	}

	// 游戏开始
	room.Start()

	return nil
}

// 掷骰子
func Dice(in rboot.Message, bot *rboot.Robot) *rboot.Message {

	room, ok := GM.GameRooms[in.From]
	// 如果房间不存在或者游戏已经结束，不做任何操作
	if !ok || room.RoomStatus == GAMEROOM_STATUS__DISABLE {
		return nil
	}

	if !room.DiceStatus {
		return nil
	}

	// 获取用户唯一标识
	userID := in.Sender
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

	var dice int

	switch bot.Ruleset {
	case `shake`:
		args := bot.Args

		fmt.Println(len(args))

		// 获取骰子点数
		if len(args) < 2 {
			return nil
		}

		fmt.Println(args)

		dice, err := strconv.Atoi(args[1])
		if err != nil {
			logrus.Errorf("骰子解析失败：%v", err)
			return nil
		}

		dice -= 3

	case `roll`:
		dice = ShakeDice()

		bot.SendText(fmt.Sprintf("点数: %d", dice), in.From)
	}

	room.Dice(dice, user)

	return nil
}

// 查看用户信息
func Look(in rboot.Message, bot *rboot.Robot) *rboot.Message {
	room, ok := GM.GameRooms[in.From]
	// 如果房间不存在或者游戏已经结束，不做任何操作
	if !ok || room.RoomStatus == GAMEROOM_STATUS__DISABLE {
		fmt.Println("房间关闭中")
		return nil
	}

	// 获取用户唯一标识
	userID := in.Sender
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

	return rboot.NewMessage(reply)
}

// 查看房间玩家状态
func RoomInfo(in rboot.Message, bot *rboot.Robot) *rboot.Message {

	room, ok := GM.GameRooms[in.From]
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

	return rboot.NewMessage(reply)

}

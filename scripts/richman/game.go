package richman

import (
	"container/ring"
	"encoding/json"
	"io/ioutil"
	"os"
)

type GameMap struct {
	PlayerElement map[*User][]MapElement // 玩家所拥有的地产
	Map           map[*User]*ring.Ring   // 玩家地图元素
}

type Game struct {
	// 用户
	User map[string]*User

	// 房间列表
	GameRooms map[string]*GameRoom

	// 地图
	Map map[int][]MapElement
}

func New() *Game {
	game := &Game{
		User:      make(map[string]*User),
		GameRooms: make(map[string]*GameRoom),
		Map:       make(map[int][]MapElement),
	}

	game.Map = game.ReadMaps()

	// 初始化幸运卡和新闻卡
	InitRuleset()

	return game
}

// 读取地图
func (g *Game) ReadMaps() map[int][]MapElement {

	var maps map[int][]MapElement

	file := os.Getenv(`MAP_FILE`)

	if file == "" {
		file = "maps/map.json"
	}

	b, err := ioutil.ReadFile(file)

	if err != nil {
		panic("地图读取失败:" + err.Error())
		return nil
	}

	err = json.Unmarshal(b, &maps)

	if err != nil {
		panic("地图解析失败: " + err.Error())
		return nil
	}

	return maps
}

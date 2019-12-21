package richman

type User struct {
	ID      string
	Name    string    // 玩家名称
	Session string    // 游戏会话
	Room    *GameRoom // 玩家所在房间
}

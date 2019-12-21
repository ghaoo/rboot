package richman

const (
	// 初始化游戏，每个用户的钱
	INITIAL_MONEY = 100000
	// 每次过起点，送给用户的钱
	SEND_MONY = 10000
)

type GAME_USER_ENUM int

const (
	_ GAME_USER_ENUM = iota
	// 玩家
	GAME_USER__PLAYER
	// 暂时退出
	GAME_USER__AFK
	// 破产用户
	GAME_USER__BROKE
	// 退出用户(不再发通知)
	GAME_USER__EXIT
	// 旁观者
	GAME_USER__WATCHER
)

// 地图元素角色
type GAME_FLAG_ENUM int

const (
	_ GAME_FLAG_ENUM = iota
	// 起点 1
	GAME_FLAG__START_POINT
	// 土地 2
	GAME_FLAG__LAND
	// 运气 3
	GAME_FLAG__LUCK
	// 新闻 4
	GAME_FLAG__NEWS
	// 用户 5
	GAME_FLAG__USER
	// 证券中心 6
	GAME_FLAG__SECURITIES_CENTER
	// 监狱 7
	GAME_FLAG__PRISION
	// 饭店 8
	GAME_FLAG__HOTEL
	// 税务 9
	GAME_FLAG__TAX_CENTER
	// 医院 10
	GAME_FLAG__HOSPITAL
)

const (
	// 投资
	GAME_FLAG__INVESTMENT_START = 100 + iota

	// 电力公司 101
	GAME_FLAG__POWER_STATION
	// 建筑公司 102
	GAME_FLAG__CONSTRUCTION_COMPANY
	// 大陆运输 103
	GAME_FLAG__CONTINENTAL_TRANSPORTION
	// 电视台 104
	GAME_FLAG__TV_STATION
	// 航空运输 105
	GAME_FLAG__AIR_TRANSPORTION
	// 污水处理 106
	GAME_FLAG__SEWAGE_TREATMENT
	// 大洋运输 107
	GAME_FLAG__OCEAN_TRANSPORTION
	// 结束
	GAME_FLAG__INVESTMENT_END
)

//土地星级
type LAND_LEVELS_ENUM int

const (
	// 最小星级
	LAND_LEVEL__MIN LAND_LEVELS_ENUM = iota
	// 一星级
	LAND_LEVEL__START1
	// 二星级
	LAND_LEVEL__START2
	// 三星级
	LAND_LEVEL__START3
	// 四星级
	LAND_LEVEL__START4
	// 五星级
	LAND_LEVEL__START5
	//最大星级
	LAND_LEVEL__MAX = LAND_LEVEL__START5
)

// 运气卡
type LUCK_CARD_TYPE_ENUM int

const (
	LUCK_CARD_TYPE__MIN LUCK_CARD_TYPE_ENUM = iota
	LUCK_CARD_TYPE__NO1
	LUCK_CARD_TYPE__NO2
	LUCK_CARD_TYPE__NO3
	LUCK_CARD_TYPE__NO4
	LUCK_CARD_TYPE__NO5
	LUCK_CARD_TYPE__NO6
	LUCK_CARD_TYPE__NO7
	LUCK_CARD_TYPE__NO8
	LUCK_CARD_TYPE__NO9
	LUCK_CARD_TYPE__NO10
	LUCK_CARD_TYPE__NO11
	LUCK_CARD_TYPE__NO12
	LUCK_CARD_TYPE__NO13
	LUCK_CARD_TYPE__MAX
)

// 新闻卡
type NEWS_CARD_TYPE_ENUM int

const (
	NEWS_CARD_TYPE__MIN NEWS_CARD_TYPE_ENUM = iota
	NEWS_CARD_TYPE__NO1
	NEWS_CARD_TYPE__NO2
	NEWS_CARD_TYPE__NO3
	NEWS_CARD_TYPE__NO4
	NEWS_CARD_TYPE__NO5
	NEWS_CARD_TYPE__NO6
	NEWS_CARD_TYPE__NO7
	NEWS_CARD_TYPE__NO8
	NEWS_CARD_TYPE__NO9
	NEWS_CARD_TYPE__NO10
	NEWS_CARD_TYPE__NO11
	NEWS_CARD_TYPE__NO12
	NEWS_CARD_TYPE__MAX
)

type GAMEROOM_STATUS_ENUM int

const (
	_ GAMEROOM_STATUS_ENUM = iota
	// 房间可用(可开始)
	GAMEROOM_STATUS__ENABLE
	// 房间游戏开始
	GAMEROOM_STATUS__GAMESTART
	//
	// 房间不可用
	GAMEROOM_STATUS__DISABLE
)

const RoomStateHtml = `
<!DOCTYPE html>
<html lang="zh-CN">
<head>
	<meta charset="UTF-8">
	<meta http-equiv="content-type" content="text/html">
    <link href="https://cdn.bootcss.com/twitter-bootstrap/3.4.1/css/bootstrap.min.css" rel="stylesheet">
</head>
<body>
<table class="table table-bordered">
    <tbody>
	<tr>
		<th width="10px"></th>
		{{range $user, $data := .}}
        <th>{{$user.Name}}
		
		{{if eq $data.State 1}}<span class="label label-success">游戏中</span>
		{{else if eq $data.State 2}}<span class="label label-primary">托管中</span>
		{{else if eq $data.State 3}}<span class="label label-danger">破产</span>
		{{end}}
		</th>
        {{end}}
    </tr>

    <tr>
        <td>位置</td>
		{{range $user, $data := .}}
			{{range $loc, $n := $data.Location}}
        	<td>{{$loc.Name}} {{if gt $n 0}}<small class="text-danger">停留{{$n}}回合</small>{{end}}</td>
        	{{end}}
        {{end}}
    </tr>

	<tr>
        <td>资金</td>
        {{range $user, $data := .}}
        <td>{{$data.Money}}</td>
        {{end}}
    </tr>

    <tr>
        <td>地产</td>
        {{range $user, $data := .}}
        <td>
            {{range $land := $data.Lands}}
            {{$land.Name}}
            {{end}}
        </td>
        {{end}}
    </tr>


    </tbody></table>
</body>
</html>
`

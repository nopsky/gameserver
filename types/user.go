package types

type UserInfo struct {
	Uid      uint64 // 用户id
	Username string // 用户名
	Role_Id  uint8  // 角色ID
	Cash     int32
	Diamond  int32
	State    uint8 // 玩家状态，0:空闲，1:游戏中，2:准备中 3:游戏中并掉线了
}

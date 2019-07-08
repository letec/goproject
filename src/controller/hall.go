package controller

import (
	"fmt"
	"goproject/src/common"
	"goproject/src/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

// Table for user
type Table struct {
	ID            string `bson:"ID"`
	PlayerA       string `bson:"PlayerA"`
	PlayerB       string `bson:"PlayerB"`
	PlayerC       string `bson:"PlayerC"`
	PlayerAStatus string `bson:"PlayerAStatus"`
	PlayerBStatus string `bson:"PlayerBStatus"`
	PlayerCStatus string `bson:"PlayerCStatus"`
	Status        string `bson:"Status"`
}

// User info set
type User struct {
	UserID   string `bson:"UserID"`
	UserName string `bson:"UserName"`
	Avatar   string `bson:"Avatar"`
	Win      string `bson:"Win"`
	Lose     string `bson:"Lose"`
	Score    string `bson:"Score"`
	Status   string `bson:"Status"`
}

// GetSeatList 获取房间列表
func GetSeatList(c *gin.Context) {
	info, _ := c.Get("MAP")
	params := info.(map[string]string)
	if !common.CheckParamsExist([]string{"gameCode"}, params) {
		c.JSON(http.StatusOK, gin.H{"result": false, "msg": "参数缺失"})
		return
	}
	gameCode := params["gameCode"]
	if !common.InSlice(gameCode, AllGame) {
		c.JSON(http.StatusOK, gin.H{"result": false, "msg": "游戏代号错误"})
		return
	}
	con := model.Mongo.Copy()
	defer con.Close()

	tables := []Table{}
	con.DB("online").C("HALL_" + gameCode + "_Tables").Find(nil).All(&tables)

	users := []User{}
	con.DB("online").C("HALL_" + gameCode + "_Users").Find(nil).All(&users)
	userList := []User{}
	for i := 0; i < len(users); i++ {
		temp, _ := model.GetUserInfoInHall([]string{users[i].UserID}, gameCode)
		if len(temp) > 0 {
			users[i].UserName = temp[0]["UserName"].(string)
			users[i].Avatar = temp[0]["Avatar"].(string)
			users[i].Score = temp[0]["Score"].(string)
			users[i].Win = temp[0][gameCode+"Win"].(string)
			users[i].Lose = temp[0][gameCode+"Lose"].(string)
			userList = append(userList, users[i])
		}
	}
	c.JSON(http.StatusOK, gin.H{"result": true, "msg": "", "data": gin.H{"tables": tables, "users": userList}})
}

// GetIntoHall 进入房间
func GetIntoHall(c *gin.Context) {
	userid, _ := c.Get("userid")
	info, _ := c.Get("MAP")
	params := info.(map[string]string)
	gameCode := params["gameCode"]
	if !common.InSlice(gameCode, AllGame) {
		c.JSON(http.StatusOK, gin.H{"result": false, "msg": "游戏代号错误"})
		return
	}
	collection := fmt.Sprintf("HALL_%v_Users", gameCode)

	conn := model.Mongo.Copy()
	defer conn.Close()

	nums, _ := conn.DB("online").C(collection).Find(bson.M{"UserID": userid.(string)}).Count()
	if nums > 0 {
		c.JSON(http.StatusOK, gin.H{"result": false, "msg": "您已经在房间中", "data": "EXIST"})
		return
	}

	inHall := model.RedisGet("USER_IN_HALL_" + userid.(string))
	if inHall != "" {
		conn.DB("online").C(inHall).Remove(bson.M{"UserID": userid.(string)})
	}

	err := conn.DB("online").C(collection).Insert(User{UserID: userid.(string), Win: "0", Lose: "0", Score: "0", Status: "0"})
	if err == nil {
		model.RedisSet("USER_IN_HALL_"+userid.(string), collection, "99999")
		c.JSON(http.StatusOK, gin.H{"result": true, "msg": "成功进入", "data": ""})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": false, "msg": "进入房间失败!", "data": "FAIL"})
}

// SeatDown player seat
func SeatDown(c *gin.Context) {
	info, _ := c.Get("MAP")
	params := info.(map[string]string)
	if !common.CheckParamsExist([]string{"gameCode", "tableID", "seat"}, params) {
		c.JSON(http.StatusOK, gin.H{"result": false, "msg": "参数缺失"})
		return
	}
	gameCode := params["gameCode"]
	if !common.InSlice(gameCode, AllGame) {
		c.JSON(http.StatusOK, gin.H{"result": false, "msg": "游戏代号错误"})
		return
	}
	userid, _ := c.Get("userid")
	tableID := params["tableID"]
	seat := params["seat"]

	conn := model.Mongo.Copy()
	defer conn.Close()

	onSeat := model.RedisGet("USER_ON_SEAT_" + userid.(string))
	if onSeat != "" {
		c.JSON(http.StatusOK, gin.H{"result": false, "msg": "您已经在一个位置坐下!"})
		return
	}

	collection := fmt.Sprintf("HALL_%v_Tables", gameCode)
	err := conn.DB("online").C(collection).Update(bson.M{"ID": tableID, seat: ""}, bson.M{"$set": bson.M{seat: userid}})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"result": false, "msg": "这个座位已经有人!"})
		return
	}

	model.RedisSet("USER_ON_SEAT_"+userid.(string), seat, "99999")

	c.JSON(http.StatusOK, gin.H{"result": true, "msg": "您已经进入房间!"})
	return
}

// StandUP 站起来
func StandUP(c *gin.Context) {
	info, _ := c.Get("MAP")
	params := info.(map[string]string)
	if !common.CheckParamsExist([]string{"gameCode", "tableID", "seat"}, params) {
		c.JSON(http.StatusOK, gin.H{"result": false, "msg": "参数缺失"})
		return
	}
	gameCode := params["gameCode"]
	if !common.InSlice(gameCode, AllGame) {
		c.JSON(http.StatusOK, gin.H{"result": false, "msg": "游戏代号错误"})
		return
	}
	userid, _ := c.Get("userid")
	tableID := params["tableID"]
	seat := params["seat"]

	conn := model.Mongo.Copy()
	defer conn.Close()

	collection := fmt.Sprintf("HALL_%v_Tables", gameCode)
	conn.DB("online").C(collection).Update(bson.M{"ID": tableID, seat: userid.(string)}, bson.M{"$set": bson.M{seat: "", seat + "Status": "0"}})

	model.RedisDel("USER_ON_SEAT_" + userid.(string))
}

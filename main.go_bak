package main

import (
	"gin_project/logging"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)


var db *gorm.DB


//用户 -》 restful http接口 -> 序列化 -》 入库 》 反序列化 》  restful http接口 》 用户
type (
	//定义原始的数据库字段
	UserInfoModel struct {
		//对数据库进行初始化
		gorm.Model
		Name string `json:"name"`
		Sex string `json:"sex"`
		Phone int `json:"phone"`
		City string `json:"city"`
	}
	//处理返回的字段
	transformedUserInfo struct {
		ID uint `json:"id"`
		Name string `json:"name"`
		Sex string `json:"sex"`
		City string `json:"city"`
		Phone int `json:"phone"`
	}

)

type userinfo struct {
	Name string `json:"name"`
	Sex string `json:"sex"`
	Phone string `json:"phone"`
	City string `json:"city"`
}

//初始化配置
func init() {
	var err error
	dsn := "root:root@tcp(192.168.3.123:3306)/lufflysex?charset=utf8&parseTime=True&loc=Local"
	db,err = gorm.Open(mysql.Open(dsn),&gorm.Config{})
	//db,err := gorm.Open(mysql.Open(dsn))

	if err != nil{
		panic("failed to connect database")
	}
	//在数据库里初始化表
	db.AutoMigrate(&UserInfoModel{})

}

//创建用户POST
func createUser(c *gin.Context)  {
	var i userinfo
	err := c.BindJSON(&i)
	if err != nil{
		log.Println(err)
		return
	}
	phone , _ := strconv.Atoi(i.Phone)
	u := UserInfoModel{
		Name:  i.Name,
		Sex:   i.Sex,
		Phone: phone,
		City: i.City,
	}
	db.Create(&u)
	c.JSON(http.StatusCreated,i)
}


//查询用户信息GET
func fetchAllUsers(c *gin.Context)  {
	var users []UserInfoModel
	db.Find(&users)
	if len(users) <= 0 {
		c.JSON(http.StatusNotFound,gin.H{"status":http.StatusNotFound,"message":"No user found"})
		return
	}
	c.JSON(http.StatusOK,&users)
}

//根据用户id找到单独的用户信息GET
func fetchSingleUsers(c *gin.Context)  {
	var users UserInfoModel
	userID := c.Param("id")
	db.First(&users,userID)
	if users.ID == 0 {
		c.JSON(http.StatusNotFound,gin.H{"status":http.StatusNotFound,"message":"No user found"})
		return
	}
	_user := transformedUserInfo{ID: users.ID,Name: users.Name,Sex: users.Sex,Phone: users.Phone,City: users.City}
	c.JSON(http.StatusOK,gin.H{"status":http.StatusOK,"data":_user})
}

//更新用户信息PUT
func updateUser(c *gin.Context)  {
	var users UserInfoModel
	var putInfo userinfo
	userID := c.Param("id")
	db.First(&users,userID)
	if users.ID == 0 {
		c.JSON(http.StatusNotFound,gin.H{"status":http.StatusNotFound,"message":"No user found"})
		return
	}
	c.BindJSON(&putInfo)
	phone,_ := strconv.Atoi(putInfo.Phone)
	db.Model(&users).Update("city",putInfo.City)
	db.Model(&users).Update("phone",phone)
	db.Model(&users).Update("name",putInfo.Name)
	db.Model(&users).Update("sex",putInfo.Sex)
	c.JSON(http.StatusOK,gin.H{"status":http.StatusOK,"message":"User updated successfully"})
}

//以下为分页操作
type Pagination struct {
	Limit int `form:"limit,omitempty;query:limit"`
	Page int `form:"page,omitempty;query:page"`
	Total int64 `form:"total"`
	Results interface{} `form:"results"`
}

func (p *Pagination)GetOffset() int  {
	return (p.GetPage() - 1) * p.Getlimit()
}

func (p *Pagination)Getlimit() int {
	if p.Limit == 0 {
		p.Limit = 10
	}
	return p.Limit
}

func (p *Pagination) GetPage() int  {
	if p.Page == 0{
		p.Page = 1
	}
	return p.Page
}

func paginate(value interface{},pagination *Pagination,db *gorm.DB) func(*gorm.DB) *gorm.DB {
	var totalRows int64
	db.Model(value).Count(&totalRows)
	pagination.Total = totalRows
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.Getlimit())
	}
}

//分页查询GET
func fetchAllUserByPaging(c *gin.Context)  {
	var users []UserInfoModel
	var pagination Pagination
	user := c.MustGet(gin.AuthUserKey).(string)

	logging.DefaultLogger().Debug(user)
	if err := c.ShouldBindQuery(&pagination); err != nil{
		c.JSON(http.StatusNotFound,gin.H{"status":http.StatusNotFound,"message":"no user found"})
		return
	}
	db.Scopes(paginate(users,&pagination,db)).Find(&users)
	pagination.Results = users
	c.JSON(http.StatusOK,&pagination)
}


func main()  {
	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	router.Use(cors.Default())
	v1 := router.Group("/api/v1/user",gin.BasicAuth(gin.Accounts{
		"abc":"123",
	}))
	{
		v1.POST("/",createUser)
		v1.GET("/",fetchAllUsers)
		v1.GET("/:id",fetchSingleUsers)
		v1.PUT("/:id",updateUser)
		v1.GET("paging/",fetchAllUserByPaging)
	}

	router.Run()


}
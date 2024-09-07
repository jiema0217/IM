package routers

func InitUserInfo() {
	r := Routes.Group("/user_info")
	{
		r.POST("/add", nil)
	}
}

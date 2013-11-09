package routers

type UserRouter struct {
	baseRouter
}

func (this *UserRouter) Home() {
	this.TplNames = "user/home.html"
}

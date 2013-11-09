package routers

type GeneralRouter struct {
	baseRouter
}

/*func (this *GeneralRouter) NestPrepare() {

}*/

func (this *GeneralRouter) PublicHome() {
	this.TplNames = "public/home.html"

	this.SetPageTitleFromKey("pagetitles.public_home")
}

func (this *GeneralRouter) UserHome() {
	if this.CheckLoginRedirect() {
		return
	}

	this.Data["IsHome"] = true
	this.TplNames = "user/home.html"
	this.SetPageTitleFromKey("pagetitles.user_home")
}

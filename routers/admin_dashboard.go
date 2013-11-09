package routers

type AdminDashboardRouter struct {
	BaseAdminRouter
}

func (this *AdminDashboardRouter) Get() {
	this.Data["consoleAdmin"] = true
	this.TplNames = "admin/dashboard.html"
	this.SetPageTitleFromKey("pagetitles.admin_dashboard")
}

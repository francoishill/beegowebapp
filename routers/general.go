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

func (this *GeneralRouter) Acknowledgements() {
	this.TplNames = "public/acknowledgements.html"

	this.SetPageTitleFromKey("pagetitles.public_acknowledgements")

	type acknowledgement struct {
		Name         string
		Url          string
		ExtraDetails string //Optional
	}

	basedOnBeego := this.Locale.Tr("based_on") + " Beego"

	listOfAcknowledgements := []acknowledgement{}
	listOfAcknowledgements = append(listOfAcknowledgements, acknowledgement{Name: "Beego", Url: "https://github.com/astaxie/beego"})
	listOfAcknowledgements = append(listOfAcknowledgements, acknowledgement{Name: "Wetalk", Url: "https://github.com/beego/wetalk", ExtraDetails: basedOnBeego})
	listOfAcknowledgements = append(listOfAcknowledgements, acknowledgement{Name: "Angular JS", Url: "http://angularjs.org/"})
	listOfAcknowledgements = append(listOfAcknowledgements, acknowledgement{Name: "Bootstrap", Url: "http://getbootstrap.com"})
	listOfAcknowledgements = append(listOfAcknowledgements, acknowledgement{Name: "Jquery", Url: "http://jquery.com/"})
	listOfAcknowledgements = append(listOfAcknowledgements, acknowledgement{Name: "Sass", Url: "http://sass-lang.com/"})

	this.Data["AcknowledgementsList"] = listOfAcknowledgements
}

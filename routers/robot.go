package routers

import (
	"bytes"
	"html/template"
	"strings"

	"./../utils"
	"github.com/astaxie/beego"
)

var robotTxt string

const robotTpl = `{{$disallow := .Disallow}}{{range .Uas}}User-Agent: {{.}}
Disallow: {{$disallow}}

{{end}}User-Agent: *
Disallow: /
`

// RobotRouter implemented global settings for all other routers.
type RobotRouter struct {
	beego.Controller
}

// Get implemented Prepare method for RobotRouter.
func (this *RobotRouter) Get() {
	if len(robotTxt) == 0 {
		// Generate "robot.txt".
		t := template.New("robotTpl")
		t.Parse(robotTpl)
		uas := strings.Split(utils.Cfg_General.MustValue("robot", "uas"), "|")

		data := make(map[string]interface{})
		data["Uas"] = uas
		data["Disallow"] = utils.Cfg_General.MustValue("robot", "disallow")

		buf := new(bytes.Buffer)
		t.Execute(buf, data)
		robotTxt = buf.String()
	}

	this.Ctx.WriteString(robotTxt)
	return
}

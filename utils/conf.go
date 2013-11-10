// Package utils implemented some useful functions.
package utils

import (
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/Unknwon/goconfig"
	"github.com/howeyc/fsnotify"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/orm"
	"github.com/beego/compress"
	"github.com/beego/i18n"

	"github.com/francoishill/runsass"

	//"strconv"
	"time"

	"./../mailer"
)

const (
	APP_VER = "0.0.1.0000"
)

var (
	AppName             string
	AppDescription      string
	AppKeywords         string
	AppVer              string
	AppHost             string
	AppUrl              string
	AppLogo             string
	AvatarURL           string
	SecretKey           string
	IsProMode           bool
	MailUser            string
	MailFrom            string
	ActivationCodeLives int
	ResetPwdCodeLives   int
	//LoginRememberDays   int
	DateFormat     string
	DateTimeFormat string
	//RealtimeRenderMD bool
	CompressSettings *compress.Settings
	Langs            []string
)

var (
	Cfg_General         *goconfig.ConfigFile
	Cfg_MachineSpecific *goconfig.ConfigFile
	Cache               cache.Cache
)

var (
	AppConfPath                 = "conf/app_general.ini"
	AppConfPath_MachineSpecific = "conf/app_machine_specific.ini"
	CompressConfPath            = "conf/compress.json"
)

// LoadConfig loads configuration file.
func LoadConfig() *goconfig.ConfigFile {
	var err error

	if fh, _ := os.OpenFile(AppConfPath, os.O_RDONLY|os.O_CREATE, 0600); fh != nil {
		fh.Close()
	}
	if fh, _ := os.OpenFile(AppConfPath_MachineSpecific, os.O_RDONLY|os.O_CREATE, 0600); fh != nil {
		fh.Close()
	}

	// Load configuration, set app version and log level.
	Cfg_General, err = goconfig.LoadConfigFile(AppConfPath)
	Cfg_General.BlockMode = false
	if err != nil {
		panic("Failed to load configuration (general) file: " + err.Error())
	}

	Cfg_MachineSpecific, err = goconfig.LoadConfigFile(AppConfPath_MachineSpecific)
	Cfg_MachineSpecific.BlockMode = false
	if err != nil {
		panic("Failed to load configuration (machine specific) file: " + err.Error())
	}

	// Trim 4th part.
	AppVer = strings.Join(strings.Split(APP_VER, ".")[:3], ".")

	configWatcher()
	reloadConfig()
	reloadConfig_MachineSpecific()

	IsProMode = beego.RunMode == "pro"

	runSassCommand() //Must be after IsProMode

	if IsProMode {
		beego.SetLevel(beego.LevelInfo)
	} else {
		beego.SetLevel(beego.LevelDebug)
	}

	// cache system
	Cache, err = cache.NewCache("memory", `{"interval":360}`)

	// session settings
	beego.SessionOn = true
	beego.SessionProvider = Cfg_General.MustValue("app", "session_provider")
	beego.SessionSavePath = Cfg_General.MustValue("app", "session_path")
	beego.SessionName = Cfg_General.MustValue("app", "session_name")

	beego.EnableXSRF = true
	// xsrf token expire time
	beego.XSRFExpire = 86400 * 365

	driverName := Cfg_General.MustValue("orm", "driver_name")
	dataSource := Cfg_General.MustValue("orm", "data_source")
	maxIdle := Cfg_General.MustInt("orm", "max_idle_conn")
	maxOpen := Cfg_General.MustInt("orm", "max_open_conn")

	orm.DefaultTimeLoc = time.UTC

	// set default database
	orm.RegisterDataBase("default", driverName, dataSource, maxIdle, maxOpen)
	orm.RunCommand()

	/*Rather call via commandline, for help call: ./main orm syncdb -h
	Refer to: http://beego.me/docs/Models_Cmd
	err = orm.RunSyncdb("default", false, false)
	if err != nil {
		beego.Error(err)
	}*/

	settingLocales()
	settingCompress()

	return Cfg_General
}

func reloadConfig() {
	AppName = Cfg_General.MustValue("app", "app_name")
	beego.AppName = AppName

	AppLogo = Cfg_General.MustValue("app", "app_logo")
	AppDescription = Cfg_General.MustValue("app", "description")
	AppKeywords = Cfg_General.MustValue("app", "keywords")
	AvatarURL = Cfg_General.MustValue("app", "avatar_url")
	DateFormat = Cfg_General.MustValue("app", "date_format")
	DateTimeFormat = Cfg_General.MustValue("app", "datetime_format")

	MailUser = Cfg_General.MustValue("app", "mail_user")
	MailFrom = Cfg_General.MustValue("app", "mail_from")

	SecretKey = Cfg_General.MustValue("app", "secret_key")
	ActivationCodeLives = Cfg_General.MustInt("app", "activation_code_live_days")
	ResetPwdCodeLives = Cfg_General.MustInt("app", "resetpwd_code_live_days")
	//LoginRememberDays = Cfg_General.MustInt("app", "login_remember_days")
	//RealtimeRenderMD = Cfg_General.MustBool("app", "realtime_render_markdown")
}

func reloadConfig_MachineSpecific() {
	beego.RunMode = Cfg_MachineSpecific.MustValue("beego", "run_mode")
	beego.HttpPort = Cfg_MachineSpecific.MustInt("beego", "http_port_"+beego.RunMode)

	AppHost = Cfg_MachineSpecific.MustValue("app", "app_host_"+beego.RunMode)
	AppUrl = Cfg_MachineSpecific.MustValue("app", "app_url_"+beego.RunMode)

	beego.EnableGzip = Cfg_MachineSpecific.MustBool("app", "app_enable_gzip")

	// set mailer connect args
	mailer.MailHost = Cfg_MachineSpecific.MustValue("mailer", "host")
	mailer.AuthUser = Cfg_MachineSpecific.MustValue("mailer", "user")
	mailer.AuthPass = Cfg_MachineSpecific.MustValue("mailer", "pass")

	orm.Debug = Cfg_MachineSpecific.MustBool("orm", "debug_log")
}

func settingLocales() {
	// autoload locales with locale_LANG.ini files
	dirs, _ := ioutil.ReadDir("conf")
	for _, info := range dirs {
		if !info.IsDir() {
			name := info.Name()
			if filepath.HasPrefix(name, "locale_") {
				if filepath.Ext(name) == ".ini" {
					lang := name[7 : len(name)-4]
					if len(lang) > 0 {
						if err := i18n.SetMessage(lang, "conf/"+name); err != nil {
							panic("Fail to set message file: " + err.Error())
						}
						continue
					}
				}
				beego.Error("locale ", name, " not loaded")
			}
		}
	}

	Langs = i18n.ListLangs()
}

func settingCompress() {
	setting, err := compress.LoadJsonConf(CompressConfPath, IsProMode, AppUrl)
	if err != nil {
		beego.Error(err)
		return
	}

	setting.RunCommand()

	if IsProMode {
		setting.RunCompress(true, false, true)
	}

	beego.AddFuncMap("compress_js", setting.Js.CompressJs)
	beego.AddFuncMap("compress_css", setting.Css.CompressCss)
}

func configWatcher() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic("Failed start app watcher: " + err.Error())
	}

	go func() {
		for {
			select {
			case event := <-watcher.Event:
				switch filepath.Ext(event.Name) {
				case ".ini":
					beego.Info(event)

					if err := Cfg_General.Reload(); err != nil {
						beego.Error("Conf general reload: ", err)
					}

					if err := Cfg_MachineSpecific.Reload(); err != nil {
						beego.Error("Conf machine specific reload: ", err)
					}

					if err := i18n.ReloadLangs(); err != nil {
						beego.Error("Conf Reload: ", err)
					}

					reloadConfig()
					reloadConfig_MachineSpecific()
					beego.Info("Config Reloaded")

				case ".json":
					if event.Name == CompressConfPath {
						settingCompress()
						beego.Info("Beego Compress Reloaded")
					}
				case ".scss", ".sass":
					runSassCommand()
				}
			}
		}
	}()

	if err := watcher.WatchFlags("conf", fsnotify.FSN_MODIFY); err != nil {
		beego.Error(err)
	}
}

func runSassCommand() {
	sett := runsass.Settings{
		SourceDir:      `C:\Francois\Other\_myapp_clone_beegowebapp\static_source\scss`,
		DestinationDir: `C:\Francois\Other\_myapp_clone_beegowebapp\static_source\css`,
	}
	//sett.RunCommand()

	runCmd := runsass.RunSassAll{Cache: false, Verbose: true}
	if IsProMode {
		runCmd.Style = "compressed"
	} else {
		runCmd.Style = "nested"
	}
	runCmd.Run(&sett)
}

func IsMatchHost(uri string) bool {
	if len(uri) == 0 {
		return false
	}

	u, err := url.ParseRequestURI(uri)
	if err != nil {
		return false
	}

	if u.Host != AppHost {
		return false
	}

	return true
}

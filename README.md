Important to CLONE first
========================
After installing this code, do not edit any of it. Use the cloning tool provided here:
[https://github.com/francoishill/clonebeegowebapp](https://github.com/francoishill/clonebeegowebapp).

# What is this?

An open source project to create a skeleton of a [Beego](https://github.com/astaxie/beego) website.

The aim is to have all the basics of a website, like User registration, login, password reset, etc. It will also include clientside libraries like AngularJS and Bootstrap.


This package was initially forked from [wetalk](https://github.com/beego/wetalk) and stripped down to give a basic web application.


# Included in this package

The stripped down version includes the following features:

0. **Authentication** (register, activate email, reset password, login, XSRF, etc)
0. **Admin** side to manage the registered users
0. **Migration** (makes use [Beego](https://github.com/astaxie/beego)/orm syncdb) to do migration and create log files
0. **Localization** using [i18n](https://github.com/beego/i18n)
0. **Two config** files, one with general settings (like app name) and the other with machine specific settings (like runmode, http_port, etc)
0. **[Compression](https://github.com/beego/compress)** of static files (like css, js, etc)
0. **MVC** pattern is used which is part of Beego splits our code up into Routers (controllers), Models and Views
0. We make use of a Master view (in views/master.html) to define the default layout of the site, ``{{.LayoutContent}}`` defines where the other templates are rendered
0. Making use of [Beego](https://github.com/astaxie/beego) also gives us the ORM capability


We have also included the following libraries into the design (they are laid out in the conf/compress.json file):

0. [AngularJS](http://angularjs.org/), created by Google and is the missing link between HTML and Javascript
0. [Bootstrap](http://getbootstrap.com/2.3.2/), created by Twitter and is a front-end framework
0. [jQuery](http://jquery.com/), do we really need to say what jQuery is?
0. [SASS](http://sass-lang.com/), a CSS-precompiler


# Installation

0. #### Here are the dependencies to be installed.

```bash
go get github.com/astaxie/beego
go get github.com/Unknwon/goconfig
go get github.com/beego/i18n
go get github.com/howeyc/fsnotify
go get github.com/beego/compress
go get github.com/go-sql-driver/mysql
go get github.com/francoishill/beegowebapp
```

0. #### Now install the cloning tool

    ```bash
    cd $GOPATH/src/github.com/francoishill/beegowebapp
    ./install_cloner.sh
    ```

0. #### Additional (non-default) packages that can be used by this application/beego:

    0. memcache: https://github.com/youtube/vitess
    0. redis: github.com/garyburd/redigo/redis
    0. x2jXML: github.com/clbanning/x2j
    0. goyaml2: github.com/wendal/goyaml2
    0. postgres: github.com/lib/pq
    0. sqlite3: github.com/mattn/go-sqlite3
    0. websockets: github.com/garyburd/go-websocket/websocket


# Quick Start

0. #### Generate your first web app (based on the [Beego](https://github.com/astaxie/beego) framework)

    ```bash
    cd /temp
    clonebeegowebapp blank "myfirstapp"
    ```
    
0. #### Important changes required in file *"app\_general.ini"* found in the **"conf"** folder.

    ```ini
    [app]
    app_name = My Web Application

    session_path = sess_db_username:sess_db_password@/sess_db_name?charset=utf8
    session_name = myapp_session
    mail_user = Myapp Mailer
    mail_from = admin@myapp.com
    secret_key = 1234abc890def567
    
    [orm]
    data_source = main_db_username:main_db_password@/main_db_name?charset=utf8
    ```
    
0. #### Have a look at the file **"app\_machine\_specific.ini"** too.
    
    
0. #### If you configured mysql to be the session provider, create its required table with the following SQL:

    ```sql
    CREATE TABLE `session` (
        `session_key` char(64) NOT NULL,
        `session_data` blob,
        `session_expiry` int(11) unsigned NOT NULL,
        PRIMARY KEY (`session_key`)
    ) ENGINE=MyISAM DEFAULT CHARSET=utf8;
    ```
    
0. #### Run the migration to generate the required tables.

    Run the script `db/migrate_windows.bat` (or `db/migrate_linux.sh` if on linux) to create all the tables for models registerd with `orm.RegisterModel`. By default this will be like the *User* model.
    
    This migration status will be logged into a timestamped file inside *db/migrations/windows* or *db/migrations/linux*.
    
    **Double check** this log to ensure everything ran smoothly.
    
    
## Running on another machine

If you want to run the web app on another machine you will only need to copy the executable and three folders (in development mode it is four folders).

In production mode: 

    conf (the complete folder)
    static (the complete folder)
    views (the complete folder)
    ??.exe (your applications' executable, on windows it will have the extension .exe)
    
In **development** mode is the same as production but you will also need to copy the *static_source* folder.


# Contribution

Please feel free to contribute or to contact me.

## License

[The MIT License (MIT)](http://opensource.org/licenses/MIT).


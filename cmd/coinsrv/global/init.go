package global

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"neo-dev/configure"
	"neo-dev/utils/errorutil"
	"net/http"
)

var (
	WebApp *echo.Echo
	DB     *sqlx.DB
)

type CustomContext struct {
	echo.Context
	DB *sqlx.DB
}

func NewDB() (*sqlx.DB, error) {
	if DB != nil {
		return DB, nil
	}
	var db *sqlx.DB
	db, err := configure.InitMysql()
	DB = db
	return DB, err
}

func NewWebApp() (*echo.Echo, error) {
	if WebApp != nil {
		return WebApp, nil
	}
	WebApp := echo.New()
	WebApp.HideBanner = true
	WebApp.HTTPErrorHandler = httpErrorHandler
	WebApp.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[${time_rfc3339}] method=${method} uri=${uri} status=${status} remote_ip=${remote_ip} \n",
	}))
	WebApp.Logger.(*log.Logger).SetHeader(`[${time_rfc3339_nano}] ${short_file}@{line}`)
	if configure.Debug {
		WebApp.Logger.SetLevel(log.DEBUG)
	} else {
		WebApp.Logger.SetLevel(log.INFO)
	}

	db, err := NewDB()
	if err != nil {
		WebApp.Logger.Fatal("连接数据库失败: ", err.Error())
		return WebApp, err
	}
	WebApp.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &CustomContext{c, db}
			return h(cc)
		}
	})
	return WebApp, nil
}

func httpErrorHandler(err error, c echo.Context) {
	var (
		code = http.StatusInternalServerError
		msg  interface{}
	)
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		msg = he.Message
		if he.Internal != nil {
			msg = fmt.Sprintf("%v, %v", err, he.Internal)
		}
	} else if WebApp.Debug {
		msg = err.Error()
	} else {
		msg = http.StatusText(code)
	}
	if !c.Response().Committed {
		if c.Request().Method == echo.HEAD {
			err = c.NoContent(code)
		} else {
			err = c.JSON(code, map[string]string{"resultCode": errorutil.System_Error, "type": "unknow", "msg": msg.(string)})
		}
	}
}

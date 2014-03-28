package controllers

import (
	"github.com/acsellers/multitemplate/revel/app/controllers"
	"github.com/revel/revel"
)

func init() {
	revel.OnAppStart(func() {
		InitDB()
		multitemplate.Init()
		multitemplate.DefaultLayout[multitemplate.HTML] = "layouts/app.html"
	})
	revel.InterceptMethod((*GorpController).Begin, revel.BEFORE)
	revel.InterceptMethod(Application.AddUser, revel.BEFORE)
	revel.InterceptMethod(Hotels.checkUser, revel.BEFORE)
	revel.InterceptMethod((*GorpController).Commit, revel.AFTER)
	revel.InterceptMethod((*GorpController).Rollback, revel.FINALLY)
}

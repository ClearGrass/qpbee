// Copyright 2013 bee authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package apiapp

import (
	"fmt"
	"os"
	path "path/filepath"
	"strings"

	"github.com/ClearGrass/qpbee/cmd/commands"
	"github.com/ClearGrass/qpbee/cmd/commands/version"
	"github.com/ClearGrass/qpbee/generate"
	"github.com/ClearGrass/qpbee/logger"
	"github.com/ClearGrass/qpbee/utils"
)

var CmdApiapp = &commands.Command{
	// CustomFlags: true,
	UsageLine: "api [appname]",
	Short:     "Creates a Beego API application",
	Long: `
  The command 'api' creates a Beego API application.

  {{"Example:"|bold}}
      $ bee api [appname] [-tables=""] [-driver=mysql] [-conn=root:@tcp(127.0.0.1:3306)/test]

  If 'conn' argument is empty, the command will generate an example API application. Otherwise the command
  will connect to your database and generate models based on the existing tables.

  The command 'api' creates a folder named [appname] with the following structure:

	    ├── main.go
	    ├── {{"conf"|foldername}}
	    │     └── app.conf
	    ├── {{"controllers"|foldername}}
	    │     └── object.go
	    │     └── user.go
        │     └── common_controller.go
	    ├── {{"routers"|foldername}}
	    │     └── router.go
	    ├── {{"tests"|foldername}}
	    │     └── default_test.go
	    ├── {{"models"|foldername}}
	    │       └── object.go
	    │       └── user.go
		├── {{"dao"|foldername}}
		│      
		├── {{"utils"|foldername}}
        │       └── contants.go
		├── {{"filter"|foldername}}
        │   
        └── {{"service"|foldername}}
`,
	PreRun: func(cmd *commands.Command, args []string) { version.ShowShortVersionBanner() },
	Run:    createAPI,
}
var apiconf = `appname = {{.Appname}}
httpport = 8080
runmode = dev
autorender = false
copyrequestbody = true
EnableDocs = true
`
var apiapolloconf = `{
    "appId":"{{.Appname}}",
    "cluster":"DEV",
    "namespaceNames":["application"],
    "ip":"54.222.185.168:8081"
}

`
var apiMaingo = `package main

import (
	_ "{{.Appname}}/routers"
    "{{.Appname}}/utils"
	"github.com/astaxie/beego"
    //"github.com/philchia/agollo"
	//"github.com/ClearGrass/code/util"
    "github.com/astaxie/beego/logs"

)
func init() {
	
	//log
	//beego.BeeLogger.DelLogger("console")
	logs.SetLogger(logs.AdapterFile,"{\"filename\":\"/opt/logs/{{.Path}}/server.log\",\"level\":7,\"daily\":true,\"maxdays\":30}")
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(4)
	logs.Async()
    
    //agollo.StartWithConfFile(util.GetApolloConfigFile(beego.AppConfig.String("appname")))

	//init constants
	utils.InitConstants()

}

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}
`

var apiMainconngo = `package main

import (
	_ "{{.Appname}}/routers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	{{.DriverPkg}}
)

func init() {
	orm.RegisterDataBase("default", "{{.DriverName}}", "{{.conn}}")
}

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}

`

var apirouter = `// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"{{.Appname}}/controllers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"

)

func init() {
	beego.GlobalControllerRouter["{{.Appname}}/controllers:CheckController"] = append(beego.GlobalControllerRouter["{{.Appname}}/controllers:CheckController"],
		beego.ControllerComments{
			Method:           "Get",
			Router:           "/",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Params:           nil})
	checkNs := beego.NewNamespace("/check",
		beego.NSInclude(
			&controllers.CheckController{},
		))
	beego.AddNamespace(checkNs)

	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/object",
			beego.NSInclude(
				&controllers.ObjectController{},
			),
		),
		beego.NSNamespace("/user",
			beego.NSInclude(
				&controllers.UserController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
`

var APIModels = `package models

import (
	"errors"
	"strconv"
	"time"
)

var (
	Objects map[string]*Object
)

type Object struct {
	ObjectId   string
	Score      int64
	PlayerName string
}

func init() {
	Objects = make(map[string]*Object)
	Objects["hjkhsbnmn123"] = &Object{"hjkhsbnmn123", 100, "astaxie"}
	Objects["mjjkxsxsaa23"] = &Object{"mjjkxsxsaa23", 101, "someone"}
}

func AddOne(object Object) (ObjectId string) {
	object.ObjectId = "astaxie" + strconv.FormatInt(time.Now().UnixNano(), 10)
	Objects[object.ObjectId] = &object
	return object.ObjectId
}

func GetOne(ObjectId string) (object *Object, err error) {
	if v, ok := Objects[ObjectId]; ok {
		return v, nil
	}
	return nil, errors.New("ObjectId Not Exist")
}

func GetAll() map[string]*Object {
	return Objects
}

func Update(ObjectId string, Score int64) (err error) {
	if v, ok := Objects[ObjectId]; ok {
		v.Score = Score
		return nil
	}
	return errors.New("ObjectId Not Exist")
}

func Delete(ObjectId string) {
	delete(Objects, ObjectId)
}

`

var APIModels2 = `package models

import (
	"errors"
	"strconv"
	"time"
)

var (
	UserList map[string]*User
)

func init() {
	UserList = make(map[string]*User)
	u := User{"user_11111", "astaxie", "11111", Profile{"male", 20, "Singapore", "astaxie@gmail.com"}}
	UserList["user_11111"] = &u
}

type User struct {
	Id       string
	Username string
	Password string
	Profile  Profile
}

type Profile struct {
	Gender  string
	Age     int
	Address string
	Email   string
}

func AddUser(u User) string {
	u.Id = "user_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	UserList[u.Id] = &u
	return u.Id
}

func GetUser(uid string) (u *User, err error) {
	if u, ok := UserList[uid]; ok {
		return u, nil
	}
	return nil, errors.New("User not exists")
}

func GetAllUsers() map[string]*User {
	return UserList
}

func UpdateUser(uid string, uu *User) (a *User, err error) {
	if u, ok := UserList[uid]; ok {
		if uu.Username != "" {
			u.Username = uu.Username
		}
		if uu.Password != "" {
			u.Password = uu.Password
		}
		if uu.Profile.Age != 0 {
			u.Profile.Age = uu.Profile.Age
		}
		if uu.Profile.Address != "" {
			u.Profile.Address = uu.Profile.Address
		}
		if uu.Profile.Gender != "" {
			u.Profile.Gender = uu.Profile.Gender
		}
		if uu.Profile.Email != "" {
			u.Profile.Email = uu.Profile.Email
		}
		return u, nil
	}
	return nil, errors.New("User Not Exist")
}

func Login(username, password string) bool {
	for _, u := range UserList {
		if u.Username == username && u.Password == password {
			return true
		}
	}
	return false
}

func DeleteUser(uid string) {
	delete(UserList, uid)
}
`

var apiControllers = `package controllers

import (
	"{{.Appname}}/models"
	"encoding/json"
)

// Operations about object
type ObjectController struct {
	BaseController
}

// @Title Create
// @Description create object
// @Param	body		body 	models.Object	true		"The object content"
// @Success 200 {string} models.Object.Id
// @Failure 403 body is empty
// @router / [post]
func (o *ObjectController) Post() {
	var ob models.Object
	json.Unmarshal(o.Ctx.Input.RequestBody, &ob)
	objectid := models.AddOne(ob)
	data := map[string]string{"ObjectId": objectid}
	o.resSuccessJson(data)
}

// @Title Get
// @Description find object by objectid
// @Param	objectId		path 	string	true		"the objectid you want to get"
// @Success 200 {object} models.Object
// @Failure 403 :objectId is empty
// @router /:objectId [get]
func (o *ObjectController) Get() {
	objectId := o.Ctx.Input.Param(":objectId")
	if objectId != "" {
		ob, err := models.GetOne(objectId)
		if err != nil {
			o.resServerErrorJson()
		} else {
			o.resSuccessJson(ob)
		}
	}
}

// @Title GetAll
// @Description get all objects
// @Success 200 {object} models.Object
// @Failure 403 :objectId is empty
// @router / [get]
func (o *ObjectController) GetAll() {
	obs := models.GetAll()
	o.resSuccessJson(obs)
}

// @Title Update
// @Description update the object
// @Param	objectId		path 	string	true		"The objectid you want to update"
// @Param	body		body 	models.Object	true		"The body"
// @Success 200 {object} models.Object
// @Failure 403 :objectId is empty
// @router /:objectId [put]
func (o *ObjectController) Put() {
	objectId := o.Ctx.Input.Param(":objectId")
	var ob models.Object
	json.Unmarshal(o.Ctx.Input.RequestBody, &ob)

	err := models.Update(objectId, ob.Score)
	if err != nil {
		o.resServerErrorJson()
	} else {
		o.resSuccessJson("update success!")
	}
}

// @Title Delete
// @Description delete the object
// @Param	objectId		path 	string	true		"The objectId you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 objectId is empty
// @router /:objectId [delete]
func (o *ObjectController) Delete() {
	objectId := o.Ctx.Input.Param(":objectId")
	models.Delete(objectId)
	o.resSuccessJson("delete success!")
}

`
var apiControllers2 = `package controllers

import (
	"{{.Appname}}/models"
	"encoding/json"

)

// Operations about Users
type UserController struct {
	BaseController
}

// @Title CreateUser
// @Description create users
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {int} models.User.Id
// @Failure 403 body is empty
// @router / [post]
func (u *UserController) Post() {
	var user models.User
	json.Unmarshal(u.Ctx.Input.RequestBody, &user)
	uid := models.AddUser(user)
	result := map[string]string{"uid": uid}
	u.resSuccessJson(result)
}

// @Title GetAll
// @Description get all Users
// @Success 200 {object} models.User
// @router / [get]
func (u *UserController) GetAll() {
	users := models.GetAllUsers()
	u.resSuccessJson(users)
}

// @Title Get
// @Description get user by uid
// @Param	uid		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.User
// @Failure 403 :uid is empty
// @router /:uid [get]
func (u *UserController) Get() {
	uid := u.GetString(":uid")
	if uid != "" {
		user, err := models.GetUser(uid)
		if err != nil {
			u.resServerErrorJson()
		} else {
			u.resSuccessJson(user)
		}
	}
}


// @Title Update
// @Description update the user
// @Param	uid		path 	string	true		"The uid you want to update"
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {object} models.User
// @Failure 403 :uid is not int
// @router /:uid [put]
func (u *UserController) Put() {
	uid := u.GetString(":uid")
	if uid != "" {
		var user models.User
		json.Unmarshal(u.Ctx.Input.RequestBody, &user)
		uu, err := models.UpdateUser(uid, &user)
		if err != nil {
			u.resServerErrorJson()
		} else {
			u.resSuccessJson(uu)
		}
	}
}

// @Title Delete
// @Description delete the user
// @Param	uid		path 	string	true		"The uid you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 uid is empty
// @router /:uid [delete]
func (u *UserController) Delete() {
	uid := u.GetString(":uid")
	models.DeleteUser(uid)
	u.resSuccessJson("delete success!")
}

// @Title Login
// @Description Logs user into the system
// @Param	username		query 	string	true		"The username for login"
// @Param	password		query 	string	true		"The password for login"
// @Success 200 {string} login success
// @Failure 403 user not exist
// @router /login [get]
func (u *UserController) Login() {
	username := u.GetString("username")
	password := u.GetString("password")
	if models.Login(username, password) {
		u.resSuccessJson("login success")
	} else {
		u.resSuccessJson("user not exist")
	}
}

// @Title logout
// @Description Logs out current logged in user session
// @Success 200 {string} logout success
// @router /logout [get]
func (u *UserController) Logout() {
	u.resSuccessJson("logout success")
}

`

var apiControllerCommon = `package controllers

import (
	"github.com/ClearGrass/api-common-code/resultcode"
	"github.com/ClearGrass/api-common-code/util"
	"github.com/astaxie/beego"
)


var Signature = "signature"

type BaseController struct {
	beego.Controller
}

func (self *BaseController) resSuccessJson(data interface{}) {
	self.resJson(resultcode.ApiSuccess, data)
}

func (self *BaseController) resParamErrorJson() {
	self.resJson(resultcode.ApiParamError, nil)
}

func (self *BaseController) resIllegalRequestJson() {
	self.resJson(resultcode.ApiIllegalRequest, nil)
}

func (self *BaseController) resServerErrorJson() {
	self.resJson(resultcode.ApiServerError, nil)
}

func (self *BaseController) resErrorResultJson(code resultcode.ResultCode) {
	self.resJson(code, nil)
}

func (self *BaseController) resJson(resultCode resultcode.ResultCode, data interface{}) {
	result := util.GenrateJsonResult(resultCode, data)
	self.Data["json"] = result
	self.ServeJSON()
}

`
var apiControllerCheck = `package controllers

type CheckController struct {
	BaseController
}

func (cc *CheckController)Get() {
	result := make(map[string]string)
	result["status"] = "ok"
	cc.resSuccessJson(result)
}
`

var apimlogs = `package mlogs

import (
	"github.com/astaxie/beego/logs"
	"encoding/json"
	"strings"
	"fmt"
)

func formatLog(f interface{}, v ...interface{}) string {
	var msg string
	switch f.(type) {
	case string:
		msg = f.(string)
		if len(v) == 0 {
			return msg
		}
		if strings.Contains(msg, "%") && !strings.Contains(msg, "%%") {
			//format string
		} else {
			//do not contain format char
			msg += strings.Repeat(" %v", len(v))
		}
	default:
		msg = fmt.Sprint(f)
		if len(v) == 0 {
			return msg
		}
		msg += strings.Repeat(" %v", len(v))
	}
	return fmt.Sprintf(msg, v...)
}

func parsh(msg string) map[string]string {
	result := make(map[string]string)
	index := strings.Index(msg, " ")
	if index == -1 {
		logs.Notice(msg)
		return nil
	}
	result["method"] = msg[0:index]

	if index+1 > len(msg) {
		logs.Notice(msg)
		return nil
	}
	msg = msg[index+1:]
	index = strings.Index(msg, ",[")
	if index == -1 {
		result["desc"] = msg
		return result
	}
	result["desc"] = msg[0:index]

	if index+2 > len(msg) {
		logs.Notice(msg)
		return nil
	}
	msg = msg[index+1:]
	for {
		if strings.Index(msg, "[") != 0 {
			logs.Notice(msg)
			return nil
		}
		msg = msg[1:]
		index = strings.Index(msg, "|")
		if index == -1 {
			logs.Notice(msg)
			return nil
		}
		tag := msg[0:index]
		if index+1 > len(msg) {
			logs.Notice(msg)
			return nil
		}
		msg = msg[index+1:]
		index = strings.Index(msg, "]")
		if index == -1 {
			logs.Notice(msg)
			return nil
		}
		val := msg[0:index]
		result[tag] = val
		if index+1 >= len(msg) {
			break
		}
		msg = msg[index+1:]
	}

	return result
}

func Debug(f interface{}, v ...interface{}) {
	msg := formatLog(f, v...)
	str, err := json.Marshal(parsh(msg))
	if err == nil {
		logs.Debug(string(str))
	}
}

func Info(f interface{}, v ...interface{}) {
	msg := formatLog(f, v...)
	str, err := json.Marshal(parsh(msg))
	if err == nil {
		logs.Info(string(str))
	}
}

func Warning(f interface{}, v ...interface{}) {
	msg := formatLog(f, v...)
	str, err := json.Marshal(parsh(msg))
	if err == nil {
		logs.Warning(string(str))
	}
}

func Alert(f interface{}, v ...interface{}) {
	msg := formatLog(f, v...)
	str, err := json.Marshal(parsh(msg))
	if err == nil {
		logs.Alert(string(str))
	}
}

func Emergency(f interface{}, v ...interface{}) {
	msg := formatLog(f, v...)
	str, err := json.Marshal(parsh(msg))
	if err == nil {
		logs.Emergency(string(str))
	}
}

func Critical(f interface{}, v ...interface{}) {
	msg := formatLog(f, v...)
	str, err := json.Marshal(parsh(msg))
	if err == nil {
		logs.Critical(string(str))
	}
}

func Notice(f interface{}, v ...interface{}) {
	msg := formatLog(f, v...)
	str, err := json.Marshal(parsh(msg))
	if err == nil {
		logs.Notice(string(str))
	}
}

func Error(f interface{}, v ...interface{}) {
	msg := formatLog(f, v...)
	str, err := json.Marshal(parsh(msg))
	if err == nil {
		logs.Error(string(str))
	}
}
`

var apiTests = `package test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"runtime"
	"path/filepath"
	_ "{{.Appname}}/routers"

	"github.com/astaxie/beego"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	_, file, _, _ := runtime.Caller(1)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".." + string(filepath.Separator))))
	beego.TestBeegoInit(apppath)
}

// TestGet is a sample to run an endpoint test
func TestGet(t *testing.T) {
	r, _ := http.NewRequest("GET", "/v1/object", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	beego.Trace("testing", "TestGet", "Code[%d]\n%s", w.Code, w.Body.String())

	Convey("Subject: Test Station Endpoint\n", t, func() {
	        Convey("Status Code Should Be 200", func() {
	                So(w.Code, ShouldEqual, 200)
	        })
	        Convey("The Result Should Not Be Empty", func() {
	                So(w.Body.Len(), ShouldBeGreaterThan, 0)
	        })
	})
}

`

var apiConstants = `package utils

import (
	//"github.com/philchia/agollo"
)

var LocationSecretKey string

func InitConstants() {
	//LocationSecretKey = agollo.GetStringValue("location.secret.key", "")
}

`


func init() {
	CmdApiapp.Flag.Var(&generate.Tables, "tables", "List of table names separated by a comma.")
	CmdApiapp.Flag.Var(&generate.SQLDriver, "driver", "Database driver. Either mysql, postgres or sqlite.")
	CmdApiapp.Flag.Var(&generate.SQLConn, "conn", "Connection string used by the driver to connect to a database instance.")
	commands.AvailableCommands = append(commands.AvailableCommands, CmdApiapp)
}

func createAPI(cmd *commands.Command, args []string) int {
	output := cmd.Out()

	if len(args) < 1 {
		beeLogger.Log.Fatal("Argument [appname] is missing")
	}

	if len(args) > 1 {
		err := cmd.Flag.Parse(args[1:])
		if err != nil {
			beeLogger.Log.Error(err.Error())
		}
	}

	appPath, packPath, err := utils.CheckEnv(args[0])
	if err != nil {
		beeLogger.Log.Fatalf("%s", err)
	}
	if generate.SQLDriver == "" {
		generate.SQLDriver = "mysql"
	}

	beeLogger.Log.Info("Creating API...")

	os.MkdirAll(appPath, 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", appPath, "\x1b[0m")

	os.Mkdir(path.Join(appPath, "conf"), 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(appPath, "conf"), "\x1b[0m")
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(appPath, "conf", "app.conf"), "\x1b[0m")
	utils.WriteToFile(path.Join(appPath, "conf", "app.conf"),
		strings.Replace(apiconf, "{{.Appname}}", path.Base(args[0]), -1))
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(appPath, "conf", "apollo.properties"), "\x1b[0m")
	utils.WriteToFile(path.Join(appPath, "conf", "apollo.properties"),
		strings.Replace(apiapolloconf, "{{.Appname}}", path.Base(args[0]), -1))

	os.Mkdir(path.Join(appPath, "controllers"), 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(appPath, "controllers"), "\x1b[0m")

	os.MkdirAll(path.Join(appPath, "tests/conf"), 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(appPath, "tests/conf"), "\x1b[0m")
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(appPath, "tests/conf", "app.conf"), "\x1b[0m")
	utils.WriteToFile(path.Join(appPath, "tests/conf", "app.conf"),
		strings.Replace(apiconf, "{{.Appname}}", path.Base(args[0]), -1))

	os.MkdirAll(path.Join(appPath, "mlogs"), 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(appPath, "mlogs", "mlogs.go"), "\x1b[0m")
	utils.WriteToFile(path.Join(appPath, "mlogs", "mlogs.go"),
		strings.Replace(apimlogs, "{{.Appname}}", path.Base(args[0]), -1))

	os.Mkdir(path.Join(appPath,"dao"), 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(appPath, "dao"), "\x1b[0m")
	os.Mkdir(path.Join(appPath,"service"), 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(appPath, "service"), "\x1b[0m")
	os.Mkdir(path.Join(appPath,"filter"), 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(appPath, "filter"), "\x1b[0m")
	os.Mkdir(path.Join(appPath,"utils"), 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(appPath, "utils"), "\x1b[0m")

	if generate.SQLConn != "" {
		fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(appPath, "main.go"), "\x1b[0m")
		mainGoContent := strings.Replace(apiMainconngo, "{{.Appname}}", packPath, -1)
		mainGoContent = strings.Replace(mainGoContent, "{{.DriverName}}", string(generate.SQLDriver), -1)
		if generate.SQLDriver == "mysql" {
			mainGoContent = strings.Replace(mainGoContent, "{{.DriverPkg}}", `_ "github.com/go-sql-driver/mysql"`, -1)
		} else if generate.SQLDriver == "postgres" {
			mainGoContent = strings.Replace(mainGoContent, "{{.DriverPkg}}", `_ "github.com/lib/pq"`, -1)
		}
		utils.WriteToFile(path.Join(appPath, "main.go"),
			strings.Replace(
				mainGoContent,
				"{{.conn}}",
				generate.SQLConn.String(),
				-1,
			),
		)
		beeLogger.Log.Infof("Using '%s' as 'driver'", generate.SQLDriver)
		beeLogger.Log.Infof("Using '%s' as 'conn'", generate.SQLConn)
		beeLogger.Log.Infof("Using '%s' as 'tables'", generate.Tables)
		generate.GenerateAppcode(string(generate.SQLDriver), string(generate.SQLConn), "3", string(generate.Tables), appPath)
	} else {
		os.Mkdir(path.Join(appPath, "models"), 0755)
		fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(appPath, "models"), "\x1b[0m")
		os.Mkdir(path.Join(appPath, "routers"), 0755)
		fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(appPath, "routers")+string(path.Separator), "\x1b[0m")

		fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(appPath, "controllers", "object.go"), "\x1b[0m")
		utils.WriteToFile(path.Join(appPath, "controllers", "object.go"),
			strings.Replace(apiControllers, "{{.Appname}}", packPath, -1))

		fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(appPath, "controllers", "user.go"), "\x1b[0m")
		utils.WriteToFile(path.Join(appPath, "controllers", "user.go"),
			strings.Replace(apiControllers2, "{{.Appname}}", packPath, -1))

		fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(appPath, "controllers", "common_controller.go"), "\x1b[0m")
		utils.WriteToFile(path.Join(appPath, "controllers", "common_controller.go"),
			strings.Replace(apiControllerCommon, "{{.Appname}}", packPath, -1))

		fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(appPath, "controllers", "check_controller.go"), "\x1b[0m")
		utils.WriteToFile(path.Join(appPath, "controllers", "check_controller.go"),
			strings.Replace(apiControllerCheck, "{{.Appname}}", packPath, -1))

		fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(appPath, "tests", "default_test.go"), "\x1b[0m")
		utils.WriteToFile(path.Join(appPath, "tests", "default_test.go"),
			strings.Replace(apiTests, "{{.Appname}}", packPath, -1))

		fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(appPath, "routers", "router.go"), "\x1b[0m")
		utils.WriteToFile(path.Join(appPath, "routers", "router.go"),
			strings.Replace(apirouter, "{{.Appname}}", packPath, -1))

		fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(appPath, "models", "object.go"), "\x1b[0m")
		utils.WriteToFile(path.Join(appPath, "models", "object.go"), APIModels)

		fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(appPath, "models", "user.go"), "\x1b[0m")
		utils.WriteToFile(path.Join(appPath, "models", "user.go"), APIModels2)

		fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(appPath, "utils", "constants.go"), "\x1b[0m")
		utils.WriteToFile(path.Join(appPath, "utils" ,"constants.go"),
			strings.Replace(apiConstants, "{{.Appname}}", packPath, -1))

		fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(appPath, "main.go"), "\x1b[0m")
		utils.WriteToFile(path.Join(appPath, "main.go"),
			strings.Replace(strings.Replace(apiMaingo, "{{.Appname}}", packPath, -1), "{{.Path}}", path.Base(packPath), -1))


	}
	beeLogger.Log.Success("New API successfully created!")
	return 0
}

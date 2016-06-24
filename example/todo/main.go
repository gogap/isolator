package main

import (
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"github.com/go-xorm/core"
	"github.com/gogap/isolator"
	"github.com/gogap/isolator/example/todo/logic"
	"github.com/gogap/isolator/example/todo/repo"
	"github.com/gogap/isolator/extension/xorm"

	_ "github.com/go-sql-driver/mysql"
	dbxorm "github.com/go-xorm/xorm"
)

var (
	taskManager = logic.TaskManager{}
	userManager = logic.UserManager{}
)

var (
	xormEngines = make(xorm.XORMEngines)
)

func main() {

	var err error

	var todoXORMEngine *dbxorm.Engine

	if todoXORMEngine, err = dbxorm.NewEngine("mysql", "root:@/todo?charset=utf8"); err != nil {
		fmt.Println(err)
		return
	}

	todoXORMEngine.SetColumnMapper(core.LintGonicMapper)

	xormEngines["todo"] = todoXORMEngine

	objectBuilder := isolator.NewClassicObjectBuilder()

	if err = objectBuilder.RegisterObjects(
		repo.NewUserRepo(),
		repo.NewTaskRepo(),
	); err != nil {
		fmt.Println(err)
		return
	}

	isor := isolator.NewIsolator(
		isolator.IsolatorObjectBuilder(objectBuilder),
	)

	isor.ObjectsSessionOptions(
		[]isolator.Object{
			(*repo.UserRepo)(nil),
			(*repo.TaskRepo)(nil),
		},
		xormEngines.NewXORMSession("todo", true),
		isolator.SessionOnSuccess(xorm.OnXORMSessionSuccess),
		isolator.SessionOnError(xorm.OnXORMSessionError),
	)

	userID := randomdata.SillyName()
	userName := randomdata.FullName(randomdata.RandomGender)

	err = isor.Invoke(userManager.AddUser, userID, userName).End(
		func(e error) {
			if e != nil {
				fmt.Println(e)
			}
			return
		},
	)

	if err != nil {
		return
	}

	isor.Invoke(taskManager.AddTask, userID, randomdata.Street(), randomdata.Street()).End(
		func(taskID string, err error) {
			if err != nil {
				fmt.Println(err)
				return
			}
		},
	)
}

package main

import (
	"errors"
	"fmt"

	"github.com/gogap/isolator"
)

var (
	userManager = UserManagement{}
)

type UserInfo struct {
	Username string
	Age      int
	Sex      int
	PhoneNum string
}

type UserRepo struct {
	session *isolator.Session
}

func (p *UserRepo) Derive(session *isolator.Session) (obj isolator.Object, err error) {
	obj = &UserRepo{session: session}
	return
}

func (p *UserRepo) GetUserInfo(userName string) UserInfo {
	return UserInfo{userName, 10000, 1, "13400000000"}
}

type UserManagement struct {
}

func (p *UserManagement) IsKing(userName string) (interface{}, error) {

	if len(userName) == 0 {
		return nil, errors.New("username is empty")
	}

	return func(userRepo *UserRepo) (bool, error) {
		ui := userRepo.GetUserInfo(userName)
		return ui.Age >= 10000, nil
	}, nil
}

func main() {

	objectBuilder := isolator.NewClassicObjectBuilder()

	if err := objectBuilder.RegisterObjects(new(UserRepo)); err != nil {
		fmt.Println(err)
		return
	}

	isor := isolator.NewIsolator(
		isolator.IsolatorObjectBuilder(objectBuilder),
	)

	onErrorRollbackFn := func(session *isolator.Session, err error) {
		fmt.Printf("## Session on error, id: %s\n## error: %s\n", session.ID, err)
		fmt.Println("## Rolling back ...")
	}

	isor.ObjectSessionOptions(
		(*UserRepo)(nil),
		isolator.SessionOnError(onErrorRollbackFn),
	)

	var isKing bool
	var err error

	fmt.Println("================== Normal ==================")
	fmt.Println("* Calling UserManager.IsKing")
	isor.Invoke(userManager.IsKing, "zeal").End(
		func(v bool, e error) {
			isKing = v
			err = e
		})

	if err != nil {
		fmt.Println("== Error Received:", err)
		return
	}

	fmt.Printf("-> king? %v \n", isKing)
	fmt.Println("================== OnError ==================")
	fmt.Println("* Calling UserManager.IsKing")
	isor.Invoke(userManager.IsKing, "").End(
		func(v bool, e error) {
			isKing = v
			err = e
		})

	if err != nil {
		fmt.Println("== Error Received:", err)
		return
	}
}

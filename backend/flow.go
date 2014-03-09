package backend

import (
	"fmt"
	"log"

	"github.com/cocaine/cocaine-flow/common"
)

const null_arg = 0

type AppUplodaInfo struct {
	Path    string `codec:"path"`
	App     string `codec:"app"`
	Version string `codec:"-"`
}

func (a *AppUplodaInfo) Fullname() string {
	return fmt.Sprintf("%s@%s", a.App, a.Version)
}

func (a *AppUplodaInfo) RoutingGroup() string {
	return a.App
}

type AppUploadTask struct {
	Username string `codec:"user"`
	Docker   string `codec:"docker"`
	Registry string `codec:"registry"`
	AppUplodaInfo
}

func isStringInSlice(name string, slice []string) (b bool) {
	b = false
	for _, item := range slice {
		if item == name {
			b = true
			return
		}
	}
	return
}

type cocainebackend struct {
	app appWrapper
	common.Context
}

/*
	ProfileController implementation
*/

func (b *cocainebackend) ProfileList() (list []string, err error) {
	err = b.app.Call("profile-list", null_arg, &list)
	return
}

func (b *cocainebackend) ProfileRead(name string) (pf map[string]interface{}, err error) {
	err = b.app.Call("profile-read", name, &pf)
	return
}

/*
	HostController impl
*/

func (b *cocainebackend) HostAdd(name string) (err error) {
	err = b.app.Call("host-add", name)
	return
}

func (b *cocainebackend) HostList() (hosts []string, err error) {
	err = b.app.Call("host-list", null_arg, &hosts)
	return
}

func (b *cocainebackend) HostRemove(name string) (err error) {
	err = b.app.Call("host-remove", name)
	return
}

/*
	RunlistController impl
*/

func (b *cocainebackend) RunlistList() (list []string, err error) {
	err = b.app.Call("runlist-list", null_arg, &list)
	return
}

func (b *cocainebackend) RunlistRead(name string) (runlist map[string]string, err error) {
	err = b.app.Call("runlist-read", name, &runlist)
	return
}

/*
	GroupController impl
*/

func (b *cocainebackend) GroupList() (list []string, err error) {
	err = b.app.Call("group-list", null_arg, &list)
	return
}

func (b *cocainebackend) GroupCreate(name string) (err error) {
	err = b.app.Call("group-create", name)
	return
}

func (b *cocainebackend) GroupRemove(name string) (err error) {
	err = b.app.Call("group-remove", name)
	return
}

func (b *cocainebackend) GroupRead(name string) (gr map[string]interface{}, err error) {
	err = b.app.Call("group-read", name, &gr)
	return
}

func (b *cocainebackend) GroupPushApp(name string, app string, weight int) (err error) {
	task := map[string]interface{}{
		"name":   name,
		"app":    app,
		"weight": weight,
	}
	err = b.app.Call("group-pushapp", task)
	return
}

func (b *cocainebackend) GroupPopApp(name string, app string) (err error) {
	task := map[string]interface{}{
		"name": name,
		"app":  app,
	}
	err = b.app.Call("group-popapp", task)
	return
}

func (b *cocainebackend) GroupRefresh(name ...string) (err error) {
	var groupname string
	if len(name) == 1 {
		groupname = name[0]
	}
	err = b.app.Call("group-refresh", groupname)
	return
}

/*
	Crashlog impl
*/

func (b *cocainebackend) CrashlogList(name string) (crashlogs []string, err error) {
	err = b.app.Call("crashlog-list", name, &crashlogs)
	return
}

func (b *cocainebackend) CrashlogView(name string, timestamp int) (crashlog string, err error) {
	task := map[string]interface{}{
		"name":      name,
		"timestamp": timestamp,
	}
	err = b.app.Call("crashlog-view", task, &crashlog)
	return
}

/*
	ApplicationController impl
*/

func (b *cocainebackend) ApplicationList(username string) (apps []string, err error) {
	err = b.app.Call("user-app-list", username, &apps)
	return
}

func (b *cocainebackend) ApplicationUpload(username string, info AppUplodaInfo) (<-chan string, <-chan error, error) {
	task := AppUploadTask{
		Username:      username,
		Docker:        b.Context.DockerEndpoint(),
		Registry:      b.Context.RegistryEndpoint(),
		AppUplodaInfo: info,
	}
	stream, err := b.app.StreamCall("user-upload", task)
	if err != nil {
		return nil, nil, err
	}

	uploadError := make(chan error, 1)
	ans := make(chan string, 10)

	onSuccess := func() {

		routingGroupName := info.RoutingGroup()

		rgs, err := b.GroupList()
		if err != nil {
			log.Println(err)
			return
		}

		if !isStringInSlice(routingGroupName, rgs) {
			b.GroupCreate(routingGroupName)
		}

		err = b.GroupPushApp(
			routingGroupName,
			info.Fullname(), 0)
		if err != nil {
			log.Println(err)
		}

	}

	go func() {
		//close response stream
		defer close(ans)
		for {
			var logdata string
			select {
			case res, ok := <-stream:
				if !ok {
					uploadError <- nil
					close(uploadError)
					// update routing groups with new application
					onSuccess()
					return
				}

				if res.Err() != nil {
					uploadError <- res.Err()
					close(uploadError)
					return
				}

				extracterr := res.Extract(&logdata)
				if extracterr != nil {
					/*
						Should I log this situation???
					*/
					continue
				}

				if len(logdata) > 0 {
					ans <- logdata
				}
			}
		}
	}()
	return ans, uploadError, nil
}

/*
	Buildlogs impl
*/

func (b *cocainebackend) BuildLogList(username string) (buildlogs []string, err error) {
	err = b.app.Call("user-buildlog-list", username, &buildlogs)
	return
}

func (b *cocainebackend) BuildLogRead(id string) (buildlog string, err error) {
	err = b.app.Call("user-buildlog-read", id, &buildlog)
	return
}

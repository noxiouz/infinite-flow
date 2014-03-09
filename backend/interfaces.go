package backend

/*
	Content interfaces
*/

type Cocaine interface {
	CrashlogController
	GroupController
	HostController
	ProfileController
	RunlistController
	ApplicationController
	BuildLogController
}

type ProfileController interface {
	ProfileList() ([]string, error)
	ProfileRead(name string) (map[string]interface{}, error)
}

type HostController interface {
	HostAdd(name string) error
	HostRemove(name string) error
	HostList() ([]string, error)
}

type RunlistController interface {
	RunlistRead(name string) (map[string]string, error)
	RunlistList() ([]string, error)
}

type GroupController interface {
	GroupList() ([]string, error)
	GroupRead(name string) (map[string]interface{}, error)
	GroupCreate(name string) error
	GroupRemove(name string) error

	GroupPushApp(name string, app string, weight int) error
	GroupPopApp(name string, app string) error
	GroupRefresh(name ...string) error
}

type CrashlogController interface {
	CrashlogList(name string) ([]string, error)
	CrashlogView(name string, timestamp int) (string, error)
}

type ApplicationController interface {
	ApplicationList() ([]string, error)
	ApplicationUpload(info AppUplodaInfo) (<-chan string, <-chan error, error)
}

type BuildLogController interface {
	BuildLogList() ([]string, error)
	BuildLogRead(id string) (string, error)
}

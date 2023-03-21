package model

var (
	DefaultFrontdoor = CreateFrontdoor("localhost", 10021)
	DefaultBackdoor  = CreateBackdoor("localhost", 10022)
)

type Frontdoor struct {
	Endpoint
}

type Backdoor struct {
	Endpoint
}

type Application struct {
	Endpoint
}

func CreateFrontdoor(host string, port int) Frontdoor {
	return Frontdoor{
		Endpoint: Endpoint{
			Host:       host,
			ListenPort: port,
		},
	}
}

func CreateBackdoor(host string, port int) Backdoor {
	return Backdoor{
		Endpoint: Endpoint{
			Host:       host,
			ListenPort: port,
		},
	}
}

func CreateApplication(host string, port int) Application {
	return Application{
		Endpoint: Endpoint{
			Host:       host,
			ListenPort: port,
		},
	}
}

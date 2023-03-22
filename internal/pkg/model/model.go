package model

var DefaultBackdoor = CreateBackdoor("localhost", 10022)

type Frontdoor struct {
	Endpoint
}

type Backdoor struct {
	Endpoint
}

type Application struct {
	Endpoint
}

func CreateFrontdoor(protocol, host string, port int) Frontdoor {
	return Frontdoor{
		Endpoint: Endpoint{
			Name:       "frontdoor",
			Protocol:   protocol,
			Host:       host,
			ListenPort: port,
		},
	}
}

func CreateBackdoor(host string, port int) Backdoor {
	return Backdoor{
		Endpoint: Endpoint{
			Name:       "backdoor",
			Protocol:   "tcp",
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

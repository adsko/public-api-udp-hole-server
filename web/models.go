package web

type register struct {
	Name string
	Addr []string
	Port string
	Password string
}


type connectBody struct {
	Server string
	Addr []string
	Port string
}

type connectResponse struct {
	Addr []string
	Port string
	ClientID uint64
	ClientSecret string
	ServerSecret string
}

type openConnection struct {
	Addr     []string
	Port     string
	ClientID uint64
	ClientSecret string
	ServerSecret string
}

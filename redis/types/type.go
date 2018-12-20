package types

//Server - containing server i/o and connecting ServerConnHandler and ServCmndsHandler
type Server struct {
	HandFlds []string
	Rslt     chan string
}

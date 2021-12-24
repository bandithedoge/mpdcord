package main

// define default config
type Config struct {
	ID                         int
	Address, Network, Password string
	Remaining                  bool
}

var DefaultConfig = Config{
	ID:        922175995828654100,
	Address:   "localhost:6600",
	Network:   "tcp",
	Password:  "",
	Remaining: false,
}

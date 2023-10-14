package my_web

type Middleware func(next HandleFunc) HandleFunc

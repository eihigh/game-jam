package main

type level struct {
	y        float64
	fireVel  float64
	interval int
}

var levels = []*level{
	0: {
		y:        0,
		fireVel:  0.3 / 60,
		interval: 100,
	},
	1: {
		y:        10,
		fireVel:  1.0 / 60,
		interval: 80,
	},
	2: {
		y:        50,
		fireVel:  1.1 / 60,
		interval: 60,
	},
	3: {
		y:        80,
		fireVel:  1.4 / 60,
		interval: 45,
	},
	4: {
		y:        150,
		fireVel:  1.5 / 60,
		interval: 30,
	},
}

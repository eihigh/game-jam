package main

type level struct {
	y           float64
	fireVel     float64
	intervalMin int
	intervalMax int
}

var levels = []*level{
	0: {
		y:           0,
		fireVel:     0.3 / 60,
		intervalMin: 100,
		intervalMax: 120,
	},
	1: {
		y:           10,
		fireVel:     1.2 / 60,
		intervalMin: 60,
		intervalMax: 120,
	},
	2: {
		y:           50,
		fireVel:     1.3 / 60,
		intervalMin: 30,
		intervalMax: 60,
	},
	3: {
		y:           80,
		fireVel:     1.4 / 60,
		intervalMin: 15,
		intervalMax: 45,
	},
}

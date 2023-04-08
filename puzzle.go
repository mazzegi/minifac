package minifac

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/mazzegi/minifac/err"
	"github.com/mazzegi/minifac/grid"
)

func LoadPuzzleFromFile(file string) (*Puzzle, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("open %q: %w", file, err)
	}
	defer f.Close()
	return LoadPuzzle(f)
}

func LoadPuzzle(r io.Reader) (*Puzzle, error) {
	var pzl Puzzle
	err := json.NewDecoder(r).Decode(&pzl)
	if err != nil {
		return nil, fmt.Errorf("json.decode: %w", err)
	}
	return &pzl, nil
}

type Puzzle struct {
	Size      grid.Size `json:"size"`
	Producers []struct {
		Resource Resource      `json:"resource"`
		Position grid.Position `json:"position"`
	} `json:"producers"`
	Assemblers []struct {
		Resource Resource      `json:"resource"`
		Position grid.Position `json:"position"`
	} `json:"assemblers"`
	Trashbins []struct {
		Position grid.Position `json:"position"`
	} `json:"trashbins"`
	Finalizers []struct {
		Resource Resource      `json:"resource"`
		Position grid.Position `json:"position"`
	} `json:"finalizers"`
	Conveyors []struct {
		Direction grid.Direction `json:"direction"`
		Position  grid.Position  `json:"position"`
	} `json:"conveyors"`
	Obstacles []struct {
		Type     ObstacleType  `json:"type"`
		Position grid.Position `json:"position"`
	} `json:"obstacles"`
}

func (p *Puzzle) Universe(cfg *Config) (*Universe, error) {
	// add 2 in each direction for bounding walls
	size := grid.S(p.Size.DX+2, p.Size.DY+2)
	u := NewUniverse(size)

	errgrp := err.NewGroup()
	for x := 0; x < size.DX; x++ {
		errgrp.Handle(u.AddObject(NewObstacle("wall", ObstacleWall), grid.P(x, 0)))
		errgrp.Handle(u.AddObject(NewObstacle("wall", ObstacleWall), grid.P(x, size.DY-1)))
	}
	for y := 1; y < size.DY-1; y++ {
		errgrp.Handle(u.AddObject(NewObstacle("wall", ObstacleWall), grid.P(0, y)))
		errgrp.Handle(u.AddObject(NewObstacle("wall", ObstacleWall), grid.P(size.DX-1, y)))
	}

	for _, p := range p.Producers {
		cp := cfg.Producers[p.Resource]
		errgrp.Handle(u.AddObject(
			NewIncarnationProducer("prod", p.Resource, cp.Rate, cp.StockCapacity),
			p.Position,
		))
	}
	for _, a := range p.Assemblers {
		ca := cfg.Assemblers[a.Resource]
		rec := Receipt{
			Input:          ca.Input,
			Output:         a.Resource,
			ProductionTime: ca.ProductionTime,
		}
		errgrp.Handle(u.AddObject(
			NewAssembler("ass", rec, ca.InputStockCapacity, ca.OutputStockCapacity),
			a.Position,
		))
	}
	for _, b := range p.Trashbins {
		errgrp.Handle(u.AddObject(NewTrashbin("trash"), b.Position))
	}
	for _, f := range p.Finalizers {
		errgrp.Handle(u.AddObject(NewFinalizer("fin", f.Resource), f.Position))
	}
	for _, c := range p.Conveyors {
		errgrp.Handle(u.AddObject(
			NewConveyor("conv", c.Direction, 1),
			c.Position,
		))
	}
	for _, o := range p.Obstacles {
		errgrp.Handle(u.AddObject(NewObstacle("obs", o.Type), o.Position))
	}

	err := errgrp.Error()
	if err != nil {
		return nil, err
	}
	return u, nil
}

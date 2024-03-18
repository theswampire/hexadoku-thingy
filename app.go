package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"math"
)

// App struct
type App struct {
	ctx    context.Context
	sudoku *Sudoku
}

// NewApp creates a new App application struct
func NewApp() *App {
	sudoku := NewSudoku(16)
	sudoku.Values = [][]int{
		{5, 10, -1, 13, -1, 15, -1, 6, -1, -1, 2, -1, -1, -1, -1, 11},
		{14, -1, -1, 9, 12, 0, 5, -1, -1, -1, 3, -1, 8, 15, -1, -1},
		{1, 6, 8, 3, -1, -1, -1, -1, -1, 7, -1, -1, -1, 4, -1, -1},
		{-1, -1, 2, 12, 3, 8, 4, -1, 1, -1, 9, -1, 14, -1, 10, -1},
		{-1, -1, -1, 5, 10, -1, 7, 1, -1, -1, -1, -1, -1, -1, -1, -1},
		{-1, 12, 7, -1, -1, -1, 0, 14, -1, 3, 5, -1, 2, -1, 4, 15},
		{13, 1, 9, -1, 4, -1, -1, -1, -1, -1, -1, 11, 7, 12, -1, -1},
		{-1, -1, 6, 11, -1, -1, -1, -1, 10, -1, -1, -1, -1, -1, 0, -1},
		{10, 11, 13, 0, 7, 3, -1, -1, 8, 12, 14, -1, 5, -1, -1, 2},
		{-1, -1, -1, -1, -1, -1, 12, 15, 9, 0, -1, -1, 11, -1, -1, -1},
		{4, 7, 5, 14, -1, 2, -1, -1, -1, 10, -1, -1, -1, 0, -1, -1},
		{-1, -1, -1, -1, -1, -1, -1, -1, -1, 6, -1, 5, 10, 13, -1, 8},
		{-1, -1, -1, -1, -1, 5, -1, -1, -1, 4, -1, 1, 15, -1, 8, 7},
		{-1, -1, -1, -1, 9, 7, -1, 12, 3, 5, -1, -1, -1, 2, -1, -1},
		{-1, -1, -1, -1, -1, 4, 15, -1, 2, -1, -1, 14, 9, -1, 5, 13},
		{-1, -1, -1, -1, -1, -1, 14, -1, -1, -1, -1, 13, -1, -1, -1, 0},
	}
	return &App{sudoku: sudoku}
}

// startup is called when the app starts. The context is saved,
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	runtime.EventsOn(a.ctx, "request_possibles", func(optionalData ...interface{}) {
		a.findPossiblesAll()
	})
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) NewSudoku(size int) {
	a.sudoku = NewSudoku(size)
}

func (a *App) GetSudoku() *Sudoku {
	return a.sudoku
}

func (a *App) InitCell(x, y, value int) error {
	err := a.sudoku.InitCell(x, y, value)

	if !a.sudoku.CheckConstraints() {
		runtime.EventsEmit(a.ctx, fmt.Sprintf("invalid_field:%d-%d:%d", a.sudoku.Size, x, y))
	}

	a.findPossiblesAll()

	return err
}

func (a *App) findPossiblesAll() {
	s := a.sudoku

	for x := 0; x < s.Size; x++ {
		for y := 0; y < s.Size; y++ {
			go a.findPossibles(x, y)
		}
	}
}

func (a *App) findPossibles(i, j int) {
	possibles := a.sudoku.findPossible(i, j)

	runtime.EventsEmit(a.ctx, fmt.Sprintf("possibility_update:%d-%d", i, j), possibles)
}

func (a *App) LockCells() *[][]bool {
	return a.sudoku.LockCells()
}

func (a *App) UnlockCells() *[][]bool {
	return a.sudoku.UnlockCells()
}

type Sudoku struct {
	Size      int `json:"size"`
	blockSize int
	Values    [][]int  `json:"values"`
	Locked    [][]bool `json:"locked"`
}

func NewSudoku(size int) *Sudoku {
	sudoku := Sudoku{Size: size, Values: make([][]int, size), Locked: make([][]bool, size)}
	for i := 0; i < size; i++ {
		sudoku.Values[i] = make([]int, size)
		sudoku.Locked[i] = make([]bool, size)

		for j := 0; j < size; j++ {
			sudoku.Values[i][j] = -1
		}
	}
	sudoku.blockSize = int(math.Sqrt(float64(size)))
	return &sudoku
}

func (s *Sudoku) validInitConstraints(x, y, value int) bool {
	return x < s.Size && y < s.Size && value < s.Size && value >= -1
}

func (s *Sudoku) InitCell(x, y, value int) error {
	if !s.validInitConstraints(x, y, value) {
		return errors.New("sudoku size constraints violated")
	}

	s.Values[x][y] = value

	return nil
}

func (s *Sudoku) LockCells() *[][]bool {
	for i := 0; i < s.Size; i++ {
		for j := 0; j < s.Size; j++ {
			s.Locked[i][j] = s.Values[i][j] != -1
		}

	}

	return &s.Locked
}

func (s *Sudoku) UnlockCells() *[][]bool {
	for i := 0; i < s.Size; i++ {
		for j := 0; j < s.Size; j++ {
			s.Locked[i][j] = false
		}

	}

	return &s.Locked
}

func (s *Sudoku) CheckConstraints() bool {

	for i := 0; i < s.blockSize; i++ {
		for j := 0; j < s.blockSize; j++ {
			if !s.checkBlock(i, j) {
				return false
			}
		}
	}

	for i := 0; i < s.Size; i++ {
		if !s.checkRow(i) || !s.checkColumn(i) {
			return false
		}
	}

	return true
}
func (s *Sudoku) checkBlock(i, j int) bool {
	values := make([]int, s.Size)

	for x := 0; x < s.blockSize; x++ {
		for y := 0; y < s.blockSize; y++ {
			v := s.Values[i*s.blockSize+x][j*s.blockSize+y]
			if v != -1 {
				values[v] += 1
			}
		}
	}

	for _, element := range values {
		if element > 1 {
			return false
		}
	}

	return true
}

func (s *Sudoku) checkRow(r int) bool {
	values := make([]int, s.Size)

	for i := 0; i < s.Size; i++ {
		v := s.Values[i][r]

		if v != -1 {
			values[v] += 1
		}
	}

	for _, element := range values {
		if element > 1 {
			return false
		}
	}
	return true
}
func (s *Sudoku) checkColumn(c int) bool {
	values := make([]int, s.Size)

	for i := 0; i < s.Size; i++ {
		v := s.Values[c][i]

		if v != -1 {
			values[v] += 1
		}
	}

	for _, element := range values {
		if element > 1 {
			return false
		}
	}
	return true
}

func (s *Sudoku) findPossible(i, j int) *[]int {
	values := make([]bool, s.Size)
	for index := range values {
		values[index] = true
	}

	blockI := i - i%s.blockSize
	blockJ := j - j%s.blockSize

	/*
		fmt.Printf("%d blockI: %d\n", i, blockI)
		fmt.Printf("%d blockJ: %d\n", j, blockJ)
	*/

	for x := 0; x < s.blockSize; x++ {
		for y := 0; y < s.blockSize; y++ {
			v := s.Values[blockI+x][blockJ+y]
			if v != -1 {
				values[v] = false
			}
		}
	}

	for r := 0; r < s.Size; r++ {
		v := s.Values[r][j]
		if v != -1 {
			values[v] = false
		}
	}

	for c := 0; c < s.Size; c++ {
		v := s.Values[i][c]
		if v != -1 {
			values[v] = false
		}
	}

	count := 0
	for _, element := range values {
		if element {
			count++
		}
	}
	possible := make([]int, count)
	counter := 0
	for index, element := range values {
		if element {
			possible[counter] = index
			counter++
		}
	}
	/*
		fmt.Printf("found for %v \n", values)
		fmt.Printf("found for %v \n", possible)
	*/

	return &possible
}

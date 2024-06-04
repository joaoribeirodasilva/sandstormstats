package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"time"
)

type Server struct {
	ID            string
	LastLine      uint
	LastTimestamp time.Time
	Games         []Game
}

type Game struct {
	ID           string
	Start        time.Time
	End          time.Time
	GameScenario string
	GameClass    string
	Players      map[string]Player
}

type Player struct {
	ID    uint64
	Name  string
	Role  string
	Team  uint
	Stats PlayerStats
}

type ObjectiveData struct {
	DateTime  time.Time
	Objective uint
}

type PlayerObjectives struct {
	Captures  ObjectiveData
	Destroies ObjectiveData
}

type PlayerStats struct {
	Enter      time.Time
	Exit       time.Time
	Kills      []Kill
	Deaths     []Kill
	Objectives []PlayerObjectives
}

type Kill struct {
	DateTime time.Time
	Role     string
	Shared   bool
}

func LogRead(path string) {
	f, err := os.Open(path)
	if err != nil {
		fmt.Printf("error opening file: %v\n", err)
		os.Exit(1)
	}
	regexpCapture, _ := regexp.Compile(`LogGameplayEvents: Display: Objective \d was .* by .*`)
	regexpDestroy, _ := regexp.Compile(`LogGameplayEvents: Display: Objective \d owned .* by .*`)
	regexpScenario, _ := regexp.Compile(`LogLoad: Warning: Scenario '`)
	regexpClass, _ := regexp.Compile(`LogLoad: Game class is '.*`)
	regexpStartRound, _ := regexp.Compile(`LogGameMode: Display: State: PreRound -> RoundActive`)
	regexpEndRound, _ := regexp.Compile(`LogGameplayEvents: Display: Round \d* Over:`)
	regexpKill, _ := regexp.Compile(`LogGameplayEvents: Display: .* killed`)

	r := bufio.NewReader(f)
	//s, e := Readln(r)
	for {

		s, e := Readln(r)

		match := regexpKill.MatchString(s)
		if match {
			fmt.Println(s)
			continue
		}

		match = regexpCapture.MatchString(s)
		if match {
			fmt.Println(s)
			continue
		}

		match = regexpDestroy.MatchString(s)
		if match {
			fmt.Println(s)
			continue
		}

		match = regexpScenario.MatchString(s)
		if match {
			fmt.Println(s)
			continue
		}

		match = regexpClass.MatchString(s)
		if match {
			fmt.Println(s)
			continue
		}

		match = regexpStartRound.MatchString(s)
		if match {
			fmt.Println(s)
			continue
		}

		match = regexpEndRound.MatchString(s)
		if match {
			fmt.Println(s)
			continue
		}

		if e != nil {
			return
		}
	}
}

func Readln(r *bufio.Reader) (string, error) {
	var (
		isPrefix bool  = true
		err      error = nil
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln), err
}

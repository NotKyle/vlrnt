package structs

import (
	"time"
)

type Role struct {
	ID   int    `sql:"primary_key"`
	Name string `sql:"type:text"`

	IsAgent      bool `sql:"type:boolean"`
	IsController bool `sql:"type:boolean"`
	IsDuelist    bool `sql:"type:boolean"`
	IsInitiator  bool `sql:"type:boolean"`
	IsSentinel   bool `sql:"type:boolean"`
}

type Player struct {
	ID int `sql:"primary_key"`

	Name     string `sql:"type:text"`
	NameReal string `sql:"type:text"`

	Team Team `sql:"-"`
	Role Role `sql:"-"`

	IsLeader bool `sql:"type:boolean"`
	IsCoach  bool `sql:"type:boolean"`
}

type Team struct {
	ID      int      `sql:"primary_key"`
	Name    string   `sql:"type:text"`
	Players []Player `sql:"-"`
}

type Map struct {
	ID   int    `sql:"primary_key"`
	Name string `sql:"type:text"`
}

type Match struct {
	ID    int    `sql:"int"`
	VLRID string `sql:"primary_key"`
	URL   string `sql:"type:text"`

	Team1 Team `sql:"int"`
	Team2 Team `sql:"int"`

	StartTime time.Time `sql:"type:datetime"`
	EndTime   time.Time `sql:"type:datetime"`

	FinalScore int
	Duration   int

	Region Region `sql:"int"`

	WinningTeam Team `sql:"int"`
	LosingTeam  Team `sql:"int"`

	MapPick Team `sql:"int"`
}

type Round struct {
	ID          int `sql:"primary_key"`
	RoundNumber int `sql:"type:integer"`

	Match Match `sql:"-"`
	Map   Map   `sql:"-"`

	RoundStartTime time.Time `sql:"type:datetime"`
	RoundEndTime   time.Time `sql:"type:datetime"`
	RoundDuration  time.Time `sql:"type:time"`

	Kills []Player `sql:"-"`

	Winner Team `sql:"-"`
	Loser  Team `sql:"-"`

	IsPistolRound   bool `sql:"type:boolean"`
	IsEcoRound      bool `sql:"type:boolean"`
	IsForceBuyRound bool `sql:"type:boolean"`
	IsGunRound      bool `sql:"type:boolean"`

	IsBombPlanted bool `sql:"type:boolean"`
	IsBombDefused bool `sql:"type:boolean"`

	IsClutchRound bool `sql:"type:boolean"`

	IsAceRound bool   `sql:"type:boolean"`
	PlayerAced Player `sql:"-"`

	IsOTRound bool `sql:"type:boolean"`
	OTNumber  int  `sql:"type:integer"`

	Attacker Team `sql:"-"`
	Defender Team `sql:"-"`
}

type Region struct {
	ID    int    `sql:"primary_key"`
	Name  string `sql:"type:text"`
	Teams []Team `sql:"-"`
}

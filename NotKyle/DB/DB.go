package db

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/sunary/sqlize"
	_ "modernc.org/sqlite"
	"notkyle.org/vlrnt/structs"
)

// Open establishes a database connection
func Open() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "file:db.sqlite?cache=shared&mode=rwc")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	return db, nil
}

// ApplyMigrations reads and executes all SQL files in the specified folder
func ApplyMigrations(db *sql.DB, folder string) error {
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		return fmt.Errorf("failed to read migration folder: %w", err)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") {
			path := filepath.Join(folder, file.Name())
			sqlContent, err := ioutil.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to read migration file %s: %w", file.Name(), err)
			}

			_, err = db.Exec(string(sqlContent))
			if err != nil {
				return fmt.Errorf("failed to execute migration %s: %w", file.Name(), err)
			}

			log.Printf("Applied migration: %s", file.Name())
		}
	}
	return nil
}

// CreateTables generates migrations and applies them to the database
func ensureFolderExists(folder string) error {
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		if err := os.MkdirAll(folder, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create migration folder: %w", err)
		}
	}
	return nil
}

func CreateTables(db *sql.DB) error {
	migrationFolder := "./migrations"

	// Ensure migration folder exists
	if err := ensureFolderExists(migrationFolder); err != nil {
		return err
	}

	newMigration := sqlize.NewSqlize(
		sqlize.WithSqlTag("sql"),
		sqlize.WithMigrationFolder(migrationFolder),
	)

	// Ensure foreign keys are enabled
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	// Auto increment primary key ID
	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS migrations (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT)"); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Generate migrations from structs
	if err := newMigration.FromObjects(
		structs.Role{},
		structs.Player{},
		structs.Team{},
		structs.Map{},
		structs.Match{},
		structs.Round{},
	); err != nil {
		return fmt.Errorf("failed to generate migrations: %w", err)
	}

	// Write migration files
	if err := newMigration.WriteFiles("initial migration"); err != nil {
		return fmt.Errorf("failed to write migration files: %w", err)
	}

	// Apply migrations
	if err := ApplyMigrations(db, migrationFolder); err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}

// CreateDB initializes the database and sets up tables
func CreateDB() {
	db, err := Open()
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	if err := CreateTables(db); err != nil {
		log.Fatalf("Error creating tables: %v", err)
	}
	log.Println("Database initialized successfully")
}

func AddMatch(db *sql.DB, match structs.Match) error {
	// fmt.Println("[1] Adding match to database", match.URL)
	// fmt.Println("[2] Match URL:", match.URL)

	// If ID exists, use next available ID
	var id int = 0

	idRow := db.QueryRow("SELECT MAX(id) FROM match")
	idRow.Scan(&id)

	// Increment the id
	id = id + 1

	_, err := db.Exec(
		"INSERT INTO match (id, url, team1, team2, start_time, end_time, final_score, duration, region, winning_team, losing_team, map_pick) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		id, match.URL, match.Team1.Name, match.Team2.Name, match.StartTime, match.EndTime, 0, 0, 0, "", "", "",
	)

	// fmt.Println("[2] Added match to database", res)
	// fmt.Println("[3] Error:", err)

	if err != nil {
		return fmt.Errorf("failed to add match: %w", err)
	}

	// log.Printf("Added match: %s vs %s at %s", match.Team1, match.Team2, match.URL)
	return nil
}

func GetMatch(db *sql.DB, url string) (structs.Match, error) {
	var match structs.Match
	row := db.QueryRow("SELECT id, url, team1, team2, start_time, end_time FROM match WHERE url = ?", url)
	err := row.Scan(&match.ID, &match.URL, &match.Team1, &match.Team2, &match.StartTime, &match.EndTime)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No match found for URL: %s", url)
			return structs.Match{}, nil // Return an empty Match and no error
		}
		return structs.Match{}, fmt.Errorf("failed to get match: %w", err)
	}
	return match, nil
}

func GetMatchByID(db *sql.DB, ID int) (structs.Match, error) {
	var match structs.Match

	row := db.QueryRow("SELECT id, url, team1, team2, start_time, end_time FROM match WHERE id = ?", ID)
	err := row.Scan(&match.ID, &match.URL, &match.Team1, &match.Team2, &match.StartTime, &match.EndTime)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No match found for ID: %d", ID)
			return structs.Match{}, nil // Return an empty Match and no error
		}
		return structs.Match{}, fmt.Errorf("failed to get match: %w", err)
	}

	return match, nil
}

func GetMatches(db *sql.DB) ([]structs.Match, error) {
	rows, err := db.Query("SELECT id, url, team1, team2, start_time, end_time FROM match")

	// Convert rows to structs
	var matches []structs.Match

	for rows.Next() {
		var match structs.Match
		err := rows.Scan(&match.ID, &match.URL, &match.Team1.Name, &match.Team2.Name, &match.StartTime, &match.EndTime)
		if err != nil {
			return nil, fmt.Errorf("failed to get matches: %w", err)
		}
		matches = append(matches, match)
	}

	if err != nil {
		return nil, fmt.Errorf("Failed to get matches: %w", err)
	}

	defer rows.Close()

	return matches, nil
}

func GetTeam(db *sql.DB, name string) (structs.Team, error) {
	var team structs.Team
	row := db.QueryRow("SELECT * FROM teams WHERE name = ?", name)
	err := row.Scan(&team.Name)
	if err != nil {
		return structs.Team{}, fmt.Errorf("failed to get team: %w", err)
	}
	return team, nil
}

func AddTeam(db *sql.DB, team structs.Team) error {
	_, err := db.Exec(
		"INSERT INTO teams (name) VALUES (?)",
		team.Name,
	)
	if err != nil {
		return fmt.Errorf("failed to add team: %w", err)
	}
	return nil
}

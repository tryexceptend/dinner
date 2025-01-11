package storagesqlite

import (
	"database/sql"
	"dinner/internal/domain/models"
	"fmt"
	"log/slog"
	"time"
)

type Storage struct {
	log *slog.Logger
	db  *sql.DB
}

// Конструктор Storage
func New(log *slog.Logger, storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	fmt.Print()
	// Указываем путь до файла БД
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{
		log: log,
		db:  db,
	}, nil
}

func (s *Storage) GetFoods() ([]models.Food, error) {
	const op = "storagesqlite.GetFoods"

	rows, err := s.db.Query("SELECT name, category from foods")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	foods := []models.Food{}

	for rows.Next() {
		var food models.Food
		_ = rows.Scan(&food.Name, &food.Category)
		foods = append(foods, food)
	}

	return foods, nil
}

func (s *Storage) SaveRequest(userId int64) error {
	const op = "storagesqlite.SaveRequest"

	stmt, err := s.db.Prepare("INSERT INTO history(userId, dt) VALUES(?, ?)")
	if err != nil {
		s.log.Error("sql prepare", slog.Any("error", err))
		return fmt.Errorf("%s: %w", op, err)
	}
	_, errExec := stmt.Exec(userId, time.Now())
	if errExec != nil {
		s.log.Error("sql exec", slog.Any("error", errExec))
		return fmt.Errorf("%s: %w", op, errExec)
	}
	return nil
}

func (s *Storage) IsLimit(userId int64) (bool, error) {
	const op = "storagesqlite.IsLimit"

	stmt, err := s.db.Prepare("SELECT count(userId) FROM history WHERE userId==? AND dt>=?")
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRow(userId, time.Now().Add((-1)*time.Duration(24)*time.Hour))

	if row == nil {
		return true, nil
	}

	cnt := 0
	row.Scan(&cnt)

	return cnt < 10, nil
}

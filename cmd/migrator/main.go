// Точка входа для запуска мигратора
package main

import (
	"errors"
	"flag"
	"fmt"

	// Библиотека для миграций
	"github.com/golang-migrate/migrate/v4"
	// Драйвер для выполнения миграций SQLite 3
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	// Драйвер для получения миграций из файлов
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var storagePath, migrationsPath, migrationsTable string

	// Получаем необходимые значения из флагов запуска

	// Путь до файла БД
	// Его достаточно, т.к. мы используем SQLite, другие креды не нужны
	flag.StringVar(&storagePath, "storage-path", "", "path to storage")
	// Путь до папки с миграциями
	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	// Таблица, в которой будет храниться информация о миграциях.
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "name of migrations table")
	// Выполняем парсинг флагов
	flag.Parse()

	// Валидация параметров
	if storagePath == "" {
		// Паника, если путь к БД не указан
		panic("storage-path is required")
	}
	if migrationsPath == "" {
		// Паника, если путь в файлам миграции не указан
		panic("migrations-path is required")
	}

	// Создаем объект мигратора, передав креды нашей БД
	m, err := migrate.New("file://"+migrationsPath,
		fmt.Sprintf("sqlite3://%s?x-migrations-table=%s", storagePath, migrationsTable),
	)
	if err != nil {
		panic(err)
	}

	// Выполняем миграции до последней версии
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")
			return
		}
		panic(err)
	}
}

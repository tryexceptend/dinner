CREATE TABLE categories (
	id INTEGER NOT NULL PRIMARY KEY,
	name TEXT
);

CREATE TABLE foods (
	id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	category INTEGER NOT NULL,
    CONSTRAINT foods_categories_FK FOREIGN KEY (category) REFERENCES categories(id) ON DELETE RESTRICT ON UPDATE RESTRICT
);

INSERT INTO categories
(id, name)
VALUES
(1,'Суп'),
(2,'Салат'),
(3,'Мясо'),
(4,'Гарнир')
;



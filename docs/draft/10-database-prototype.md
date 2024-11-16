# Проектируем базу данных


Cпроектируем базу данных, учитывая все типы привычек, их конфигурацию, систему баллов и историю выполнения.

Мы будем использовать SQLite, так как это лёгкая и встроенная база данных, которая идеально подходит для небольших проектов.

## Основные сущности:

1.	Пользователи (users) — хранит информацию о пользователях бота.
2.	Привычки (habits) — хранит конфигурацию каждой привычки.
3.	История выполнения (records) — хранит записи выполнения привычек, включая время выполнения и начисленные баллы.
4.	Напоминания (reminders) — хранит настройки напоминаний для привычек.
5.	Прогресс (progress) — используется для накопительных привычек (cumulative).

## Схема базы данных

### Таблица users

Таблица содержит информацию о пользователях, которые взаимодействуют с ботом.

Описание данных:

- id — уникальный идентификатор пользователя в базе данных.
- telegram_id — идентификатор пользователя в Telegram, используется для связи с ботом.
- username — имя пользователя в Telegram (если доступно).
- created_at — дата и время регистрации пользователя в системе.

Когда используются данные:

- При первом запуске бот добавляет нового пользователя в таблицу.
- Используется для идентификации пользователя при запросах к другим таблицам.

```sql
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY,
    telegram_id INTEGER UNIQUE,
    username TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```


### Таблица habits

Таблица содержит все привычки пользователя и их конфигурацию.

Описание данных:

- id — уникальный идентификатор привычки.
- user_id — ссылка на пользователя, которому принадлежит привычка.
- name — название привычки.
- type — тип привычки (simple, counter, time, duration, cumulative, periodic, checklist, random).
- target — целевое значение (например, количество отжиманий или шагов).
- target_time — целевое время (например, 6:00 для привычек типа time).
- max_time — максимальное время, после которого баллы не начисляются (для привычек типа time).
- unit — единица измерения (например, “ml”, “steps”).
- points — максимальное количество баллов за выполнение привычки.
- points_mode — режим начисления баллов (fixed, proportional, time_based).
- created_at — дата создания привычки.

Когда используются данные:

- Когда пользователь добавляет новую привычку, она записывается в эту таблицу.
- При выполнении привычки бот считывает настройки из таблицы для расчёта баллов и проверки выполнения.

```sql
CREATE TABLE IF NOT EXISTS habits (
    id INTEGER PRIMARY KEY,
    user_id INTEGER,
    name TEXT,
    type TEXT,
    target INTEGER DEFAULT 0,
    target_time TEXT DEFAULT NULL,
    max_time TEXT DEFAULT NULL,
    unit TEXT DEFAULT NULL,
    points INTEGER DEFAULT 0,
    points_mode TEXT DEFAULT 'fixed',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```


### Таблица records

Таблица хранит историю выполнения привычек пользователем.

Описание данных:

- id — уникальный идентификатор записи.
- habit_id — ссылка на привычку, выполнение которой фиксируется.
- user_id — ссылка на пользователя, который выполнил привычку.
- value — значение выполнения (например, количество отжиманий, шагов или длительность медитации).
- timestamp — время выполнения привычки.
- points — количество начисленных баллов за выполнение.

Когда используются данные:

- Когда пользователь отмечает выполнение привычки через команду /done, создаётся новая запись в таблице.
- Данные используются для расчёта статистики и отображения прогресса пользователю.

```sql
CREATE TABLE IF NOT EXISTS records (
    id INTEGER PRIMARY KEY,
    habit_id INTEGER,
    user_id INTEGER,
    value INTEGER DEFAULT 0,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    points INTEGER DEFAULT 0,
    FOREIGN KEY (habit_id) REFERENCES habits(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```


### Таблица reminders

Таблица хранит настройки напоминаний для каждой привычки пользователя.

Описание данных:

- id — уникальный идентификатор напоминания.
- habit_id — ссылка на привычку, для которой установлено напоминание.
- user_id — ссылка на пользователя, которому принадлежит напоминание.
- time — время, в которое бот должен отправить напоминание (формат HH:MM).
- days — дни недели, когда должно отправляться напоминание (например, “mon,tue,wed”).

Когда используются данные:

- Когда пользователь настраивает напоминание для привычки, оно записывается в эту таблицу.
- Бот регулярно проверяет время и отправляет напоминания пользователю в указанные дни.

```sql
CREATE TABLE IF NOT EXISTS reminders (
    id INTEGER PRIMARY KEY,
    habit_id INTEGER,
    user_id INTEGER,
    time TEXT,
    days TEXT,
    FOREIGN KEY (habit_id) REFERENCES habits(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```

### Таблица progress

Таблица используется для хранения накопительного прогресса по привычкам типа cumulative.

Описание данных:

- id — уникальный идентификатор прогресса.
- habit_id — ссылка на привычку, для которой отслеживается прогресс.
- user_id — ссылка на пользователя, который выполняет привычку.
- accumulated_value — накопленное значение (например, количество шагов за неделю).
- target — целевое значение (например, 70000 шагов).
- last_updated — время последнего обновления прогресса.

Когда используются данные:

- Когда пользователь добавляет значение для накопительной привычки (например, количество шагов), таблица обновляется.
- Используется для отслеживания прогресса и расчёта баллов при достижении цели.

```sql
CREATE TABLE IF NOT EXISTS progress (
    id INTEGER PRIMARY KEY,
    habit_id INTEGER,
    user_id INTEGER,
    accumulated_value INTEGER DEFAULT 0,
    target INTEGER,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (habit_id) REFERENCES habits(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```


## Пример запросов

1. Добавление новой привычки:
```sql
INSERT INTO habits (user_id, name, type, target, points, points_mode)
VALUES (1, 'Отжимания', 'counter', 30, 20, 'proportional');
```

2. Отметка выполнения привычки:
```sql
INSERT INTO records (habit_id, user_id, value, points)
VALUES (1, 1, 20, 13);
```

3. Получение статистики по привычкам:
```sql
SELECT h.name, SUM(r.points) as total_points
FROM habits h
JOIN records r ON h.id = r.habit_id
WHERE h.user_id = 1
GROUP BY h.name;
```


---

Prev:: [Детальные сценарии всех типов привычек](9-scenario-done-types.md)

Next:: []()

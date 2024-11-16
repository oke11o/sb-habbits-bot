# Repositories

## 1. Интерфейс HabitRepo

Для работы с привычками.

```go
type HabitRepo interface {
    CreateHabit(ctx context.Context, habit model.Habit) (model.Habit, error)
    UpdateHabit(ctx context.Context, habit model.Habit) error
    UpsertHabit(ctx context.Context, habit model.Habit) (model.Habit, error)
    
    GetHabitByID(ctx context.Context, habitID int64) (model.Habit, error)
    GetHabitsByUserID(ctx context.Context, userID int64) ([]model.Habit, error)
    GetHabitByName(ctx context.Context, userID int64, name string) (model.Habit, error)
    
    DeleteHabitByID(ctx context.Context, habitID int64) error
    DeleteHabitByName(ctx context.Context, userID int64, name string) error
}
```


## 2. Интерфейс ReminderRepo

Для работы с напоминаниями.

```go
type ReminderRepo interface {
    CreateReminder(ctx context.Context, reminder model.Reminder) (model.Reminder, error)
    UpdateReminder(ctx context.Context, reminder model.Reminder) error
    DeleteReminder(ctx context.Context, reminderID int64) error
    
    GetRemindersByHabitID(ctx context.Context, habitID int64) ([]model.Reminder, error)
    GetRemindersByUserID(ctx context.Context, -userID int64) ([]model.Reminder, error)
}
```

## 3. Интерфейс ProgressRepo

Для работы с накопительным прогрессом.
```go
type ProgressRepo interface {
    CreateProgress(progress *model.Progress) error
    UpdateProgress(progress *model.Progress) error
    GetProgressByHabitID(habitID int64) (*model.Progress, error)
    GetProgressByUserID(userID int64) ([]*model.Progress, error)
    DeleteProgress(habitID int64) error
}
```

## 4. Интерфейс RecordRepo

Для работы с историей выполнения привычек.
```go
type RecordRepo interface {
    CreateRecord(record *model.Record) error
    GetRecordsByHabitID(habitID int64, limit int) ([]*model.Record, error)
    GetRecordsByUserID(userID int64, limit int) ([]*model.Record, error)
    GetLatestRecordByHabitID(habitID int64) (*model.Record, error)
    DeleteRecord(recordID int64) error
}
```




---

Prev:: [Проектируем базу данных](11-repositories.md)

Next:: []()


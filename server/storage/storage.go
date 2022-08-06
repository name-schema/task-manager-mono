package storage

import (
	"context"
	"os"
	"time"

	"github.com/keyslapperdev/task-manager-mono/server/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var _SQLITE_DB_NAME = os.Getenv("SQLITE_DB_NAME")

type DataMgr interface {
	CreateTask(context.Context, models.Task) uint
	GetTasks(context.Context) []models.Task
	GetTaskByID(context.Context, uint) models.Task
	UpdateTask(context.Context, models.Task) models.Task
	CloseTask(context.Context, models.Task) models.Task
	DeleteTask(context.Context, models.Task)
}

type DBStorage struct {
	*gorm.DB
}

func NewDBStorer(inMemory bool) DataMgr {
	dsn := _SQLITE_DB_NAME

	if inMemory {
		dsn = ":memory:"
	}

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("error opening db: " + err.Error())
	}

	db.AutoMigrate(models.StatusTypes)
	db.AutoMigrate(models.Comment{})
	db.AutoMigrate(models.Task{})

	return &DBStorage{db.Debug()}
}

func (dbs DBStorage) CreateTask(ctx context.Context, task models.Task) uint {
	dbs.WithContext(ctx).Create(&task)
	return task.ID
}

func (dbs DBStorage) GetTasks(ctx context.Context) []models.Task {
	var tasks []models.Task
	dbs.WithContext(ctx).Preload("Comments").Find(&tasks)
	return tasks
}

func (dbs DBStorage) GetTaskByID(ctx context.Context, taskID uint) models.Task {
	task := models.Task{ID: taskID}
	dbs.WithContext(ctx).Preload("Comments").Find(&task)
	return task
}

//  Note, when manually updating (rolling json by hand), if you don't put the createdAt
//   time the value will n00k itself. I'm okay with this behavior because when the
//   frontend is built out, the entire object will be sent back. If something changes
//   with that plan between now and then, I'll address it then.
func (dbs DBStorage) UpdateTask(ctx context.Context, task models.Task) models.Task {
	dbs.WithContext(ctx).Save(&task)
	return task
}

func (dbs DBStorage) CloseTask(ctx context.Context, task models.Task) models.Task {
	task.StatusID = models.StatusClosed.ID
	task.ClosedAt = time.Now()

	dbs.WithContext(ctx).Debug().Save(&task)
	return task
}

func (dbs DBStorage) DeleteTask(ctx context.Context, task models.Task) {
	dbs.WithContext(ctx).Debug().Delete(&task)
}

package schedule

import (
	"encoding/json"
	"fmt"
	scheduleModel "intehub/internal/app/models/schedule"
	appService "intehub/internal/app/service/app"
	pushService "intehub/internal/app/service/push"
	workflowService "intehub/internal/app/service/workflow"
	"time"

	"github.com/robfig/cron/v3"
)

type service struct {
	model           scheduleModel.Model
	appService      appService.Service
	pushService     pushService.Service
	workflowService workflowService.Service
	cron            *cron.Cron
	cronEntries     map[uint]cron.EntryID
}

func New(model scheduleModel.Model, appSvc appService.Service, pushSvc pushService.Service, workflowSvc workflowService.Service) Service {
	return &service{
		model:           model,
		appService:      appSvc,
		pushService:     pushSvc,
		workflowService: workflowSvc,
		cron:            cron.New(cron.WithSeconds()),
		cronEntries:     make(map[uint]cron.EntryID),
	}
}

func (s *service) GetTasks(userID uint) ([]*scheduleModel.ScheduleTask, error) {
	return s.model.GetTasks(userID)
}

func (s *service) GetTask(id uint) (*scheduleModel.ScheduleTask, error) {
	return s.model.GetTask(id)
}

func (s *service) CreateTask(task *scheduleModel.ScheduleTask) error {
	if err := s.model.CreateTask(task); err != nil {
		return err
	}
	if task.Enabled {
		if err := s.addTaskToCron(task); err != nil {
			return err
		}
	}
	return nil
}

func (s *service) UpdateTask(task *scheduleModel.ScheduleTask) error {
	if err := s.model.UpdateTask(task); err != nil {
		return err
	}
	s.removeTaskFromCron(task.ID)
	if task.Enabled {
		if err := s.addTaskToCron(task); err != nil {
			return err
		}
	}
	return nil
}

func (s *service) DeleteTask(id uint) error {
	s.removeTaskFromCron(id)
	return s.model.DeleteTask(id)
}

func (s *service) ToggleTask(id uint) error {
	task, err := s.model.GetTask(id)
	if err != nil {
		return err
	}
	if err := s.model.ToggleTask(id); err != nil {
		return err
	}
	if task.Enabled {
		s.removeTaskFromCron(id)
	} else {
		task.Enabled = true
		if err := s.addTaskToCron(task); err != nil {
			return err
		}
	}
	return nil
}

func (s *service) GetLogs(taskID *uint, page, pageSize int) ([]*scheduleModel.ScheduleLog, int64, error) {
	return s.model.GetLogs(taskID, page, pageSize)
}

func (s *service) ExecuteTask(task *scheduleModel.ScheduleTask) error {
	startTime := time.Now()
	log := &scheduleModel.ScheduleLog{
		TaskID:     task.ID,
		TaskName:   task.Name,
		Status:     "success",
		StartedAt:  startTime,
		FinishedAt: startTime,
	}

	var err error
	switch task.TaskType {
	case "app":
		if task.AppID == nil {
			err = fmt.Errorf("应用ID不能为空")
		} else {
			var input map[string]interface{}
			if task.FieldData != "" {
				if jsonErr := json.Unmarshal([]byte(task.FieldData), &input); jsonErr != nil {
					err = fmt.Errorf("解析输入参数失败: %v", jsonErr)
				}
			}
			if err == nil {
				_, err = s.appService.RunWithInput(*task.AppID, input)
			}
		}
	case "workflow":
		if task.WorkflowID == nil {
			err = fmt.Errorf("工作流ID不能为空")
		} else {
			var input map[string]interface{}
			if task.FieldData != "" {
				if jsonErr := json.Unmarshal([]byte(task.FieldData), &input); jsonErr != nil {
					err = fmt.Errorf("解析输入参数失败: %v", jsonErr)
				}
			}
			if err == nil {
				_, err = s.workflowService.Run(*task.WorkflowID, input)
			}
		}
	case "push":
		if task.ConfigID == nil {
			err = fmt.Errorf("推送配置ID不能为空")
		} else {
			var data map[string]interface{}
			if task.FieldData != "" {
				if jsonErr := json.Unmarshal([]byte(task.FieldData), &data); jsonErr != nil {
					err = fmt.Errorf("解析推送数据失败: %v", jsonErr)
				}
			}
			if err == nil {
				err = s.pushService.Send(*task.ConfigID, data)
			}
		}
	default:
		err = fmt.Errorf("不支持的任务类型: %s", task.TaskType)
	}

	finishTime := time.Now()
	log.FinishedAt = finishTime
	log.Duration = finishTime.Sub(startTime).Milliseconds()

	if err != nil {
		log.Status = "error"
		log.Message = err.Error()
	} else {
		log.Message = "执行成功"
	}

	s.model.CreateLog(log)

	now := time.Now()
	task.LastRunAt = &now
	s.model.UpdateTask(task)

	return err
}

func (s *service) StartScheduler() error {
	tasks, err := s.model.GetTasks(0)
	if err != nil {
		return err
	}

	for _, task := range tasks {
		if task.Enabled {
			if err := s.addTaskToCron(task); err != nil {
				return fmt.Errorf("启动任务 %s 失败: %v", task.Name, err)
			}
		}
	}

	s.cron.Start()
	return nil
}

func (s *service) StopScheduler() {
	s.cron.Stop()
}

func (s *service) addTaskToCron(task *scheduleModel.ScheduleTask) error {
	if entryID, exists := s.cronEntries[task.ID]; exists {
		s.cron.Remove(entryID)
	}

	entryID, err := s.cron.AddFunc(task.CronExpr, func() {
		s.ExecuteTask(task)
	})

	if err != nil {
		return fmt.Errorf("添加定时任务失败: %v", err)
	}

	s.cronEntries[task.ID] = entryID
	return nil
}

func (s *service) removeTaskFromCron(taskID uint) {
	if entryID, exists := s.cronEntries[taskID]; exists {
		s.cron.Remove(entryID)
		delete(s.cronEntries, taskID)
	}
}

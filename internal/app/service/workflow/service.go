package workflow

import (
	"encoding/json"
	"errors"
	"fmt"
	workflowModel "intehub/internal/app/models/workflow"
	appService "intehub/internal/app/service/app"
	"time"
)

type service struct {
	workflowModel workflowModel.Model
	appService    appService.Service
}

func New(workflowModel workflowModel.Model, appService appService.Service) Service {
	return &service{
		workflowModel: workflowModel,
		appService:    appService,
	}
}

func (s *service) Create(workflow *workflowModel.Workflow) error {
	return s.workflowModel.Create(workflow)
}

func (s *service) GetByID(id uint) (*workflowModel.Workflow, error) {
	return s.workflowModel.GetByID(id)
}

func (s *service) List(userID uint) ([]*workflowModel.Workflow, error) {
	return s.workflowModel.List(userID)
}

func (s *service) Update(workflow *workflowModel.Workflow) error {
	return s.workflowModel.Update(workflow)
}

func (s *service) Delete(id uint) error {
	return s.workflowModel.Delete(id)
}

func (s *service) Run(id uint, input map[string]interface{}) (*workflowModel.WorkflowLog, error) {
	// 获取工作流定义
	workflow, err := s.workflowModel.GetByID(id)
	if err != nil {
		return nil, errors.New("工作流不存在")
	}

	if !workflow.Enabled {
		return nil, errors.New("工作流已禁用")
	}

	// 解析节点配置
	var nodes []workflowModel.WorkflowNode
	if err := json.Unmarshal([]byte(workflow.Nodes), &nodes); err != nil {
		return nil, fmt.Errorf("解析工作流配置失败: %v", err)
	}

	// 创建执行日志
	startTime := time.Now()
	inputJSON, _ := json.Marshal(input)
	log := &workflowModel.WorkflowLog{
		WorkflowID: workflow.ID,
		Name:       workflow.Name,
		Status:     "running",
		Input:      string(inputJSON),
		StartedAt:  startTime,
	}
	if err := s.workflowModel.CreateLog(log); err != nil {
		return nil, err
	}

	// 执行工作流
	nodeLogs, output, execErr := s.executeWorkflow(nodes, input)

	// 更新日志
	finishTime := time.Now()
	log.FinishedAt = finishTime
	log.Duration = finishTime.Sub(startTime).Milliseconds()

	nodeLogsJSON, _ := json.Marshal(nodeLogs)
	log.NodeLogs = string(nodeLogsJSON)

	if execErr != nil {
		log.Status = "error"
		log.Error = execErr.Error()
	} else {
		log.Status = "success"
		outputJSON, _ := json.Marshal(output)
		log.Output = string(outputJSON)
	}

	s.workflowModel.UpdateLog(log)

	return log, execErr
}

func (s *service) executeWorkflow(nodes []workflowModel.WorkflowNode, input map[string]interface{}) ([]workflowModel.NodeLog, map[string]interface{}, error) {
	if len(nodes) == 0 {
		return nil, nil, errors.New("应用流没有配置节点")
	}

	// 构建节点映射
	nodeMap := make(map[string]*workflowModel.WorkflowNode)
	var startNode *workflowModel.WorkflowNode

	for i := range nodes {
		nodeMap[nodes[i].ID] = &nodes[i]
		// 查找开始节点
		if nodes[i].Type == "start" {
			startNode = &nodes[i]
		}
	}

	// 如果没有开始节点，使用第一个节点
	if startNode == nil {
		startNode = &nodes[0]
	}

	// 执行节点
	var nodeLogs []workflowModel.NodeLog
	currentData := input
	currentNode := startNode

	for currentNode != nil {
		// 跳过开始和结束节点的实际执行
		if currentNode.Type == "start" {
			nodeLog := &workflowModel.NodeLog{
				NodeID:     currentNode.ID,
				Status:     "success",
				Input:      input,
				Output:     input,
				StartedAt:  time.Now(),
				FinishedAt: time.Now(),
				Duration:   0,
			}
			nodeLogs = append(nodeLogs, *nodeLog)
		} else if currentNode.Type == "end" {
			nodeLog := &workflowModel.NodeLog{
				NodeID:     currentNode.ID,
				Status:     "success",
				Input:      currentData,
				Output:     currentData,
				StartedAt:  time.Now(),
				FinishedAt: time.Now(),
				Duration:   0,
			}
			nodeLogs = append(nodeLogs, *nodeLog)
		} else {
			// 执行应用节点
			nodeLog, output, err := s.executeNode(currentNode, currentData)
			nodeLogs = append(nodeLogs, *nodeLog)

			if err != nil {
				return nodeLogs, nil, fmt.Errorf("节点 %s 执行失败: %v", currentNode.ID, err)
			}

			// 更新当前数据为节点输出
			currentData = output
		}

		// 获取下一个节点
		if len(currentNode.Next) > 0 {
			nextNodeID := currentNode.Next[0]
			currentNode = nodeMap[nextNodeID]
		} else {
			currentNode = nil
		}
	}

	return nodeLogs, currentData, nil
}

func (s *service) executeNode(node *workflowModel.WorkflowNode, input map[string]interface{}) (*workflowModel.NodeLog, map[string]interface{}, error) {
	startTime := time.Now()

	nodeLog := &workflowModel.NodeLog{
		NodeID:    node.ID,
		AppID:     node.AppID,
		Status:    "running",
		Input:     input,
		StartedAt: startTime,
	}

	// 根据节点类型执行
	var output map[string]interface{}
	var err error

	switch node.Type {
	case "start":
		// 开始节点直接传递数据
		nodeLog.AppName = "开始"
		output = input
	case "end":
		// 结束节点直接传递数据
		nodeLog.AppName = "结束"
		output = input
	case "app":
		output, err = s.executeAppNode(node, input, nodeLog)
	case "transform":
		output, err = s.executeTransformNode(node, input)
	default:
		err = fmt.Errorf("不支持的节点类型: %s", node.Type)
	}

	finishTime := time.Now()
	nodeLog.FinishedAt = finishTime
	nodeLog.Duration = finishTime.Sub(startTime).Milliseconds()

	if err != nil {
		nodeLog.Status = "error"
		nodeLog.Error = err.Error()
		return nodeLog, nil, err
	}

	nodeLog.Status = "success"
	nodeLog.Output = output

	return nodeLog, output, nil
}

func (s *service) executeAppNode(node *workflowModel.WorkflowNode, input map[string]interface{}, nodeLog *workflowModel.NodeLog) (map[string]interface{}, error) {
	// 获取应用信息
	app, _, err := s.appService.GetByIDWithConfigs(node.AppID)
	if err != nil {
		return nil, fmt.Errorf("应用不存在: %v", err)
	}

	nodeLog.AppName = app.Name

	// 合并节点配置和输入数据
	appInput := make(map[string]interface{})
	for k, v := range node.Config {
		appInput[k] = v
	}
	for k, v := range input {
		appInput[k] = v
	}

	// 执行应用
	appLog, err := s.appService.RunWithInput(node.AppID, appInput)
	if err != nil {
		return nil, err
	}

	if appLog.Status != "success" {
		return nil, errors.New(appLog.Error)
	}

	// 解析应用输出
	var output map[string]interface{}
	if err := json.Unmarshal([]byte(appLog.Output), &output); err != nil {
		return nil, fmt.Errorf("解析应用输出失败: %v", err)
	}

	return output, nil
}

func (s *service) executeTransformNode(node *workflowModel.WorkflowNode, input map[string]interface{}) (map[string]interface{}, error) {
	// 数据转换节点：根据配置提取或转换数据
	output := make(map[string]interface{})

	// 从配置中获取映射规则
	if mappings, ok := node.Config["mappings"].(map[string]interface{}); ok {
		for targetKey, sourceKey := range mappings {
			if sourceKeyStr, ok := sourceKey.(string); ok {
				if value, exists := input[sourceKeyStr]; exists {
					output[targetKey] = value
				}
			}
		}
	}

	// 如果没有配置映射，直接传递所有数据
	if len(output) == 0 {
		output = input
	}

	return output, nil
}

func (s *service) GetLogs(workflowID *uint, page, pageSize int) ([]*workflowModel.WorkflowLog, int64, error) {
	return s.workflowModel.GetLogs(workflowID, page, pageSize)
}

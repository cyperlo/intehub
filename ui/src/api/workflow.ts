import request from './request'

export interface WorkflowNode {
  id: string
  type: string
  app_id: number
  config: Record<string, any>
  next: string[]
  position?: { x: number; y: number }
  label?: string
}

export interface Workflow {
  id?: number
  name: string
  description: string
  nodes: string
  enabled: boolean
  user_id?: number
  created_at?: string
  updated_at?: string
}

export interface WorkflowLog {
  id: number
  workflow_id: number
  name: string
  status: string
  input: string
  output: string
  error: string
  node_logs: string
  duration: number
  started_at: string
  finished_at: string
  created_at: string
}

export const getWorkflows = () => {
  return request.get<any, Workflow[]>('/workflows')
}

export const getWorkflow = (id: number) => {
  return request.get<any, Workflow>(`/workflows/${id}`)
}

export const createWorkflow = (data: Workflow) => {
  return request.post<any, Workflow>('/workflows', data)
}

export const updateWorkflow = (id: number, data: Workflow) => {
  return request.put<any, Workflow>(`/workflows/${id}`, data)
}

export const deleteWorkflow = (id: number) => {
  return request.delete(`/workflows/${id}`)
}

export const runWorkflow = (id: number, input?: Record<string, any>) => {
  return request.post(`/workflows/${id}/run`, input || {})
}

export const getWorkflowLogs = (params: any) => {
  return request.get<any, { list: WorkflowLog[]; total: number; page: number }>('/workflows/logs', { params })
}

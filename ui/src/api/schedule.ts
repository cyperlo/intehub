import request from './request'

export interface ScheduleTask {
  id?: number
  name: string
  description: string
  cron_expr: string
  task_type: string
  config_id?: number
  app_id?: number
  workflow_id?: number
  field_data?: string
  enabled: boolean
  last_run_at?: string
  next_run_at?: string
  user_id?: number
  created_at?: string
  updated_at?: string
}

export interface ScheduleLog {
  id: number
  task_id: number
  task_name: string
  status: string
  message: string
  duration: number
  started_at: string
  finished_at: string
  created_at: string
}

export const getScheduleTasks = () => {
  return request.get<any, ScheduleTask[]>('/schedule/tasks')
}

export const getScheduleTask = (id: number) => {
  return request.get<any, ScheduleTask>(`/schedule/tasks/${id}`)
}

export const createScheduleTask = (data: ScheduleTask) => {
  return request.post<any, ScheduleTask>('/schedule/tasks', data)
}

export const updateScheduleTask = (id: number, data: ScheduleTask) => {
  return request.put<any, ScheduleTask>(`/schedule/tasks/${id}`, data)
}

export const deleteScheduleTask = (id: number) => {
  return request.delete(`/schedule/tasks/${id}`)
}

export const toggleScheduleTask = (id: number) => {
  return request.post(`/schedule/tasks/${id}/toggle`)
}

export const getScheduleLogs = (params: any) => {
  return request.get<any, { list: ScheduleLog[]; total: number; page: number }>('/schedule/logs', { params })
}

import request from './request'

export interface App {
  id?: number
  name: string
  description: string
  code: string
  language: string
  enabled: boolean
  user_id?: number
  created_at?: string
  updated_at?: string
}

export interface AppConfig {
  key: string
  value: string
  type: string
  encrypted: boolean
}

export interface AppLog {
  id: number
  app_id: number
  app_name: string
  status: string
  output: string
  error: string
  duration: number
  started_at: string
  finished_at: string
  created_at: string
}

export const getApps = () => {
  return request.get<any, App[]>('/apps')
}

export const getApp = (id: number) => {
  return request.get<any, App>(`/apps/${id}`)
}

export const createApp = (data: App & { configs?: AppConfig[] }) => {
  return request.post<any, App>('/apps', data)
}

export const updateApp = (id: number, data: App & { configs?: AppConfig[] }) => {
  return request.put<any, App>(`/apps/${id}`, data)
}

export const deleteApp = (id: number) => {
  return request.delete(`/apps/${id}`)
}

export const runApp = (id: number, input?: Record<string, any>) => {
  return request.post(`/apps/${id}/run`, input || {})
}

export const getAppLogs = (params: any) => {
  return request.get<any, { list: AppLog[]; total: number; page: number }>('/apps/logs', { params })
}

export interface PublishToStoreRequest {
  display_name: string
  description: string
  icon?: string
  category?: string
  version: string
  author?: string
  tags?: string
  config_schema?: string
}

export const publishToStore = (id: number, data: PublishToStoreRequest) => {
  return request.post(`/apps/${id}/publish`, data)
}

import request from './request'

export interface AppTemplate {
  id?: number
  name: string
  display_name: string
  description: string
  icon: string
  code: string
  language: string
  category: string
  version: string
  author: string
  tags: string
  config_schema: string
  enabled: boolean
  created_at?: string
  updated_at?: string
}

export interface AppConfig {
  id?: number
  app_id: number
  key: string
  value: string
  type: string
  encrypted: boolean
  created_at?: string
  updated_at?: string
}

export interface InstallAppRequest {
  template_id: number
  name: string
  description: string
  user_id: number
  configs?: Record<string, string>
}

// 模板管理
export const getTemplates = (category?: string) => {
  return request.get<any, AppTemplate[]>('/appstore/templates', {
    params: { category }
  })
}

export const getTemplate = (id: number) => {
  return request.get<any, { data: AppTemplate }>(`/appstore/templates/${id}`)
}

export const createTemplate = (data: Partial<AppTemplate>) => {
  return request.post<any, { data: AppTemplate }>('/appstore/templates', data)
}

export const updateTemplate = (id: number, data: Partial<AppTemplate>) => {
  return request.put<any, { data: AppTemplate }>(`/appstore/templates/${id}`, data)
}

export const deleteTemplate = (id: number) => {
  return request.delete(`/appstore/templates/${id}`)
}

// 配置管理
export const getAppConfigs = (appId: number) => {
  return request.get<any, { data: AppConfig[] }>(`/appstore/configs/${appId}`)
}

export const createConfig = (data: Partial<AppConfig>) => {
  return request.post<any, { data: AppConfig }>('/appstore/configs', data)
}

export const updateConfig = (id: number, data: Partial<AppConfig>) => {
  return request.put<any, { data: AppConfig }>(`/appstore/configs/${id}`, data)
}

export const deleteConfig = (id: number) => {
  return request.delete(`/appstore/configs/${id}`)
}

// 从模板安装应用
export const installFromTemplate = (data: InstallAppRequest) => {
  return request.post('/appstore/install', data)
}

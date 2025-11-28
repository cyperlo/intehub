import request from './request'
import type { PushConfig, PushHistory } from '../types'

export const getPushConfigs = () => {
  return request.get<any, PushConfig[]>('/push/configs')
}

export const getPushConfig = (id: number) => {
  return request.get<any, PushConfig>(`/push/configs/${id}`)
}

export const createPushConfig = (data: PushConfig) => {
  return request.post<any, PushConfig>('/push/configs', data)
}

export const updatePushConfig = (id: number, data: PushConfig) => {
  return request.put<any, PushConfig>(`/push/configs/${id}`, data)
}

export const deletePushConfig = (id: number) => {
  return request.delete(`/push/configs/${id}`)
}

export const sendPush = (data: { config_id: number; data?: Record<string, string> }) => {
  return request.post('/push/send', data)
}

export const getPushHistory = (params?: { page?: number; page_size?: number; config_id?: number }) => {
  return request.get<any, { total: number; page: number; page_size: number; data: PushHistory[] }>('/push/history', { params })
}

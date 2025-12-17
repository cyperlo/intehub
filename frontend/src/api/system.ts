import request from './request'

export interface SystemLog {
  id: number
  level: string
  module: string
  action: string
  content: string
  user_id: number
  ip: string
  created_at: string
}

export interface Menu {
  id?: number
  name: string
  path: string
  icon: string
  parent_id: number
  sort: number
  visible: boolean
  roles: string
  created_at?: string
  updated_at?: string
}

export interface User {
  id?: number
  username: string
  nickname: string
  role: string
  password?: string
  created_at?: string
  updated_at?: string
}

export const getSystemLogs = (params: any) => {
  return request.get<any, { list: SystemLog[]; total: number; page: number }>('/system/logs', { params })
}

export const createSystemLog = (data: Partial<SystemLog>) => {
  return request.post('/system/logs', data)
}

export const getMenus = () => {
  return request.get<any, Menu[]>('/system/menus')
}

export const createMenu = (data: Menu) => {
  return request.post<any, Menu>('/system/menus', data)
}

export const updateMenu = (id: number, data: Menu) => {
  return request.put<any, Menu>(`/system/menus/${id}`, data)
}

export const deleteMenu = (id: number) => {
  return request.delete(`/system/menus/${id}`)
}

export const getUsers = () => {
  return request.get<any, User[]>('/system/users')
}

export const createUser = (data: User) => {
  return request.post<any, User>('/system/users', data)
}

export const updateUser = (id: number, data: Partial<User>) => {
  return request.put<any, User>(`/system/users/${id}`, data)
}

export const deleteUser = (id: number) => {
  return request.delete(`/system/users/${id}`)
}

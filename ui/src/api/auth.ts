import request from './request'
import type { LoginForm, LoginResponse, User } from '../types'

export const login = (data: LoginForm) => {
  return request.post<any, LoginResponse>('/auth/login', data)
}

export const logout = () => {
  return request.post('/auth/logout')
}

export const getCurrentUser = () => {
  return request.get<any, User>('/user/current')
}

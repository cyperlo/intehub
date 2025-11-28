export interface User {
  id: number
  username: string
  nickname: string
  role: string
  created_at: string
  updated_at: string
}

export interface PushConfig {
  id?: number
  name: string
  description: string
  url: string
  method: string
  headers: string
  content_type: string
  template: string
  enabled: boolean
  user_id?: number
  created_at?: string
  updated_at?: string
}

export interface PushHistory {
  id: number
  config_id: number
  config_name: string
  url: string
  method: string
  content: string
  status_code: number
  response: string
  success: boolean
  error: string
  duration: number
  user_id: number
  created_at: string
}

export interface LoginForm {
  username: string
  password: string
}

export interface LoginResponse {
  token: string
  user_info: User
}

export type FieldType = 'text' | 'textarea' | 'number' | 'date' | 'select' | 'url' | 'email'

export interface FieldSchema {
  id?: number
  name: string
  key: string
  type: FieldType
  description: string
  required: boolean
  default_value: string
  options: string
  placeholder: string
  validation: string
  user_id?: number
  created_at?: string
  updated_at?: string
}

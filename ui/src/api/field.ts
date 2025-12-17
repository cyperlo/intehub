import request from './request'
import type { FieldSchema } from '../types'

export const getFieldSchemas = () => {
  return request.get<any, FieldSchema[]>('/fields')
}

export const getFieldSchema = (id: number) => {
  return request.get<any, FieldSchema>(`/fields/${id}`)
}

export const createFieldSchema = (data: FieldSchema) => {
  return request.post<any, FieldSchema>('/fields', data)
}

export const updateFieldSchema = (id: number, data: FieldSchema) => {
  return request.put<any, FieldSchema>(`/fields/${id}`, data)
}

export const deleteFieldSchema = (id: number) => {
  return request.delete(`/fields/${id}`)
}

export const getConfigFields = (configId: number) => {
  return request.get<any, FieldSchema[]>(`/push/configs/${configId}/fields`)
}

export const updateConfigFields = (configId: number, fieldIds: number[]) => {
  return request.put(`/push/configs/${configId}/fields`, { field_ids: fieldIds })
}

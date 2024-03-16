import { $authHost, $host, $convertHost } from "."
import { AxiosError } from "axios"

export const convertData = async (file: Blob) => {
  const {data} = await $convertHost.post('/upload', {file}, {headers: {
    'Content-Type': 'multipart/form-data',
  }})
  return data
}

export const getGLBFile = async (fileName: string) => {
  const {data} = await $convertHost.get(`/files/${fileName}`)
  return data
}

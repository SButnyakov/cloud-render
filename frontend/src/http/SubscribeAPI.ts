import { $authHost, $host } from "."
import { AxiosError } from "axios"

export const subscribe = async () => {
  const {data} = await $authHost.post('/subscribe', {})
  return data
}

import { $authHost, $host } from "."
import { User } from "../store/UserStore"

export const getUser = async (): Promise<User> => {
  const {data} = await $authHost.get('/user')
  return data
}

export const editUser = async (email: string, login: string, password: string) => {
  const {data} = await $host.put('/user/edit', {login, email, password}, {headers: {Authorization: `Bearer ${localStorage.getItem('token')}`}})
  return data
}

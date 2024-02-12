import { $host } from "."
import { AxiosError } from "axios"

type RegistrationData = {
  email: string,
  login: string,
  password: string
}

export const registration = async (data: RegistrationData) => {
  const {email, login, password} = data

  await $host.post('signup', {email, login, password})
    .catch(err => {
      throw err
    })
}

export const auth = async (login_or_email: string, password: string) => {
  await $host.post('signin', {login_or_email, password})
      .then(({data}) => {
        localStorage.setItem('token', data.access_token)
        localStorage.setItem('refresh_token', data.refresh_token)
      })
      .catch((error: AxiosError) => {
        throw error
      })
}

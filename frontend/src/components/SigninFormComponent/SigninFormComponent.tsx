import React, { FormEvent, useState } from "react"
import { auth } from "../../http/AuthAPI"
import { observer } from "mobx-react-lite"

import { useStore } from "../../hooks/useStore"
import { useNavigate } from "react-router-dom"

import styles from './SigninFormComponent.module.css'
import { AxiosError } from "axios"
import { SigninResponseCodes } from "../../http/httpTypes"
import { getUser } from "../../http/UserApi"

export const SigninForm = observer(() => {
  const [login, setLogin] = useState('')
  const [password, setPassword] = useState('')

  const [errorMessage, setErrorMessage] = useState('')

  const store = useStore()

  const route = useNavigate()

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault()
    signIn()
  }

  const signIn = async () => {
    await auth(login, password)
      .then(async () => {
        const user = await getUser()

        store.userStore.setUser({email: user.email, login: user.login})
        store.userStore.setIsAuth(true)

        route('/upload')
      })
      .catch((error: AxiosError) => {
        const {response} = error

        if (response?.status === SigninResponseCodes.INTERNAL_SERVER_ERROR) {
          setErrorMessage('Failed to verify login information')
        }

        if (response?.status === SigninResponseCodes.INVALID_CREDENTIALS) {
          setErrorMessage('Wrong login or password')

          setLogin('')
          setPassword('')
        }
      })
  }

  return(
    <div >
      <form className={styles.formContainer} onSubmit={handleSubmit}>
        <div>
          <input 
            name="login" 
            type="text" 
            value={login} 
            placeholder="Login"
            onChange={e => {setLogin(e.target.value)}}
          />
        </div>

        <div>
          <input 
            name="password" 
            type="password" 
            value={password}
            placeholder="Password"
            onChange={e => {setPassword(e.target.value)}}
          />
        </div>
        <div className={styles.buttonContainer}>
          <button disabled={!login || !password}>Log In</button>
          <button onClick={() => {route('/signup')}}>Register</button>
        </div>
      </form>
      <div className={styles.errorBlockMessage}>
        {errorMessage}
      </div>
    </div>
  )
})

export default SigninForm
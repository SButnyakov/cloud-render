import React, { FormEvent, useMemo, useState } from "react"
import { registration } from "../../http/AuthAPI"
import { observer } from "mobx-react-lite"
import { useNavigate } from "react-router-dom"
import { AxiosError } from "axios"

import styles from './SignupFormComponent.module.css'
import { useForm } from "react-hook-form"

type FormData = {
  login: string;
  email: string;
  password: string;
  passwordRepeat: string;
};

/* TODO: Доделать валидацию*/

export const SignupForm = observer(() => {
  const {
    register,
    handleSubmit,
    watch,
    formState: { errors },
  } = useForm<FormData>();

  const [serverError, setServerError] = useState('')
  const [extraErrors, setExtraErrors] = useState('')

  const route = useNavigate()

  const errorList = useMemo(() => {
    const errorMessages = []
    const errorsObject = {...errors}

    for (let [_, value] of Object.entries(errorsObject)) {
      errorMessages.push(value)
    }

    return errorMessages
  }, [errors])

  const onSubmit = async (data: FormData) => {
    if (data.password !== data.passwordRepeat) {
      setExtraErrors('Passwords must be equal!')
    }
    else {
      setServerError('')
      try {
        await registration(data)
        .then(() => {
          route('/signin')
        })
      }
      catch (e: any) {
        setServerError(e.response.data.error)
      }
    }
  }

  return(
    <div>
      <form className={styles.formContainer} onSubmit={handleSubmit(onSubmit)}>
        <div>
          <input
            type="text" 
            id="email"
            placeholder="Email"
            {...register('email', {
              required: 'Input email',
              pattern: {
                value: /^[\w-]+(\.[\w-]+)*@([\w-]+\.)+[a-zA-Z]{2,7}$/,
                message:
                  'Wrong email format!',
              },
            })}
          />
        </div>

        <div>
          <input
            type="text" 
            id="login"
            placeholder="Login"
            {...register('login', {
              required: 'Input Login',
              pattern: {
                value: /^[a-zA-Z0-9]{4,15}$/,
                message:
                  'The login must contain Latin letters and have a maximum length of 15 characters',
              },
            })}
          />
        </div>

        <div>
          <input
            type="password" 
            id="password"
            placeholder="Password"
            {...register('password')}
          />
        </div>

        <div>
          <input
            type="password" 
            id="passwordRepeat"
            placeholder="Confirm password"
            {...register('passwordRepeat')}
          />
        </div>

        <div className={styles.errorBlockMessage}>
          {errorList.map(el => {
           return <p>{el.message}</p> 
          })}

          {serverError && <p>{serverError}</p>}

          {extraErrors ?? null}
        </div>

        <div className={styles.buttonContainer}>
          <button type="submit">Sign Up</button>
          <button onClick={() => {route('/signin')}}>Sign In</button>
        </div>
      </form>
      
    </div>
  )
})

export default SignupForm
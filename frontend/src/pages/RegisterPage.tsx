import React from "react"
import SignupForm from "../components/SignupFormComponent/SignupFormComponent"

import styles from './styles/AuthPage.module.css'



export const RegisterPage = () => {
  return(
    <div className={styles.page}>
      <SignupForm/>
    </div>
  )
}

export default RegisterPage
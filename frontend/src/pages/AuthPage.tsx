import SigninForm from "../components/SigninFormComponent/SigninFormComponent"
import styles from './styles/AuthPage.module.css'

const AuthPage = () => {

  return(
    <div className={styles.page}>
      <SigninForm/>
    </div>
  )
}

export default AuthPage
import { makeAutoObservable } from "mobx"

export type User = {
  login: string,
  email: string,
  expirationDate?: string
}

export class UserStore {
  private _isAuth: boolean
  private _user: User

  constructor() {
    
    if (JSON.parse(localStorage.getItem('user') as string) as User) {
      this._isAuth = true
    }
    else {
      this._isAuth = false
    }
    
    this._user = JSON.parse(localStorage.getItem('user') as string) as User

    makeAutoObservable(this)
  }

  public setExparationDate(date: string) {
    this._user.expirationDate = date

    localStorage.setItem('user', JSON.stringify(this._user))
  }

  public setIsAuth(isAuth: boolean) {
    this._isAuth = isAuth
  }

  public setUser(user: User) {
    this._user = user
    this._user.expirationDate = ''

    localStorage.setItem('user', JSON.stringify(this._user))
  }

  public setEmail(email: string) {
    this._user.email = email

    localStorage.setItem('user', JSON.stringify(this._user))
  }

  public editUser(user: User) {
    this._user = user

    localStorage.setItem('user', JSON.stringify(this._user))
  }

  get isAuth() {
    return this._isAuth
  }

  get user() {
    return this._user
  }
}

export default UserStore
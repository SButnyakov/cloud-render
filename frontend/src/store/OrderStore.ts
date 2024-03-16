import { makeAutoObservable } from "mobx"

export type DownloadLink = {
  String: string;
  Valid: boolean;
}

export type Order = {
  id: number,
  filename: string,
  date: string,
  status: string
  downloadLink: DownloadLink,
}

export class OrderStore {

  private _orders: Order[] = []

  constructor() {
    makeAutoObservable(this)
  }
}

export default OrderStore
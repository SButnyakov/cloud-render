import { $authHost } from "."

/* TODO: Пока в body ручки можно отправлять только сам файл FormData. Позже нужно будет переделать */ 
/* SOLVED: Нужно засовывать в FormData поля format и resolution */
export const sendUploadedFile = async (format: string, resolution: string, uploadfile: Blob) => {
  try {
    const {data} = await $authHost.post('send', {uploadfile, format, resolution}, {headers: {
      'Content-Type': 'multipart/form-data',
    }})
  }
  catch (e) {
    console.error(e)
  }
}

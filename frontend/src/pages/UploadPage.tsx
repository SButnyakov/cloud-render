import React, { useRef, useState } from "react";

import { sendUploadedFile } from "../http/RenderBackendAPI";
import { useStore } from "../hooks/useStore";
import { observer } from "mobx-react-lite";
import { useNavigate } from "react-router-dom";

import styles from './styles/UploadPage.module.css'
import { getOrders } from "../http/OrdersAPI";
import { Order } from "../store/OrderStore";
import ThreeDScene from "../components/ThreeDSceneComponent/ThreeDScene";
import { convertData, getGLBFile } from "../http/ConvertAPI";
import Loader from "../components/loader/Loader";

const UploadPage = observer(() => {
  const [isFileUploaded, setIsFileUploaded] = useState(false)
  const [formatSettings, setFormatSettings] = useState('png')
  const [resolutionSettings, setResolutionSettings] = useState('1920x1080')
  const [errorMesage, setErrorMessage] = useState('')

  const [isOptionBlockVisible, setIsOptionBlockVisible] = useState(true)

  const [isLoadingFile, setIsLoadingFile] = useState(false)

  const [convertedFileName, setConvertedFileName] = useState<string | ArrayBuffer | null>()

  const route = useNavigate()

  const [file, setFile] = useState<string | Blob>('')
  const [fileName, setFileName] = useState('')
  
  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleGetFile = () => {
    setIsFileUploaded(true)
    setIsLoadingFile(false)
  }
  console.log(formatSettings)
  const handleFileUploaded = (event: any) => {
    

    console.log(event.target.files)
    if (event.target.files[0]) {
      setIsLoadingFile(true)
      setFile(event.target.files[0])
      setFileName(event.target.files[0].name)
      
      convertData(event.target.files[0] as Blob)
        .then(async res => {
          setConvertedFileName(res.convertedFileName)
          
        })
    }
  }

  const findCirrentOrderByid = (orders: Order[]) => {
    let maxId = 0

    for (let order of orders) {
      const currentId = order.id

      if (currentId > maxId) {
        maxId = currentId;
      }
    }
    return maxId
  }

  const handleSendFile = async () => {
    setErrorMessage('')

    try {
      await sendUploadedFile(formatSettings, resolutionSettings, file as Blob)
      console.log(formatSettings)
      const orders = await getOrders()

      const newOrder = findCirrentOrderByid(orders)

      route(`/order/${newOrder}`)
    }
    catch (e) {
      console.error(e)
    }
  }

  const handleDropFile = () => {
    if (fileInputRef.current) {
      fileInputRef.current.value = '';
    }

    setIsFileUploaded(false)
    setConvertedFileName('')
  }

  const simulateFileInputClick = () => {
    if (fileInputRef.current !== null) {
      fileInputRef.current.click();
    }
  };


  return(
    <div className={styles.uploadPage} style={{justifyContent: isFileUploaded ? 'flex-start' : 'center'}}>

      {isLoadingFile && <Loader />}

      {convertedFileName &&
        <ThreeDScene glbFile={convertedFileName} handleFileUploaded={handleGetFile}/>
      }

      {(!isLoadingFile && !convertedFileName) &&
        <div className={styles.uploadField} onClick={simulateFileInputClick}>
          UPLOAD
          <input type="file" accept=".blend" onChange={handleFileUploaded} ref={fileInputRef}/>
        </div>
      }

      {isFileUploaded && (
        <div className={styles.preferencesBlock}>
          <div>
            <div className={styles.optionGear} onClick={() => setIsOptionBlockVisible(!isOptionBlockVisible)}>
              ⚙️
            </div>

            <div className={styles.optionBlock} style={{visibility: isOptionBlockVisible ? 'visible' : 'hidden'}}>
              
              <p>Format:</p>
              
              <select
                placeholder="format"
                value={formatSettings}
                onChange={(e) => setFormatSettings(e.target.value)}
              >
                <option value="png">png</option>
                <option value="jpeg">jpeg</option>
              </select>

              <p>Resolution:</p>
              <select
                placeholder="resolution" 
                value={resolutionSettings} 
                onChange={(e) => setResolutionSettings(e.target.value)}
              >
                <option value="1920x1080">1920x1080</option>
                <option value="1280x720">1280x720</option>
              </select>
            </div>
          </div>

          <div className={styles.controlsBlock}>
            <button onClick={() => handleSendFile()}>
                Send File
            </button>

            <button onClick={() => handleDropFile()}>
              DROP
            </button>
          </div>
        </div>
      )}
    </div>
  )
})

export default UploadPage

const express = require('express');
const multer = require('multer');
const { exec } = require('child_process');
const fs = require('fs/promises');

const path = require('path');

const cors = require('cors');

const storageConfig = multer.diskStorage({
  destination: (req, file, cb) =>{
      cb(null, "uploads");
  },
  filename: (req, file, cb) =>{
      cb(null, file.originalname);
  }
});

const upload = multer({ storage: storageConfig }); // сохранять файлы в папку 'uploads'

const app = express();

app.use(cors());

app.use(express.static('/glb_files'));

app.post('/upload', upload.single('file'), (req, res) => {
  console.log(req.file);
  // Запускать bash скрипт с передачей пути к файлу
  runBashScript(req.file.path, req.file.originalname, (resultPath) => {
    // возвращаем пользователю результат
    res.json({ convertedFileName: resultPath });
  });
});

app.get('/files/:fileName', (req, res) => {
  const fileName = req.params.fileName;

  if (fileName.includes('..')) {
      return res.status(400).send("Некорректный запрос");
  }

  const filePath = path.join(__dirname, 'glb_files', fileName);

  res.sendFile(filePath, (err) => {
      if (err) {
          console.log(err);
          res.status(404).send("Файл не найден");
      } else {
          console.log("Файл был успешно отправлен.");
          fs.unlink(filePath);
      }
  });

  
});

app.listen(5500, () => console.log('App is listening on port 5500'));

function runBashScript(filePath, originalFileName, callback) {
  const scriptPath = 'utils/scripts/2gltf2-master/2gltf2.sh';
  exec(`${scriptPath} ${filePath}`, async (error, stdout, stderr) => {
    if (error) {
      console.error(exec `error: ${error}`);
      return;
    }
    const fileName = originalFileName.split('.')[0];
    await fs.unlink(filePath);

    callback(`${fileName}.glb`);
  });
}

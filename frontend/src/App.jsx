import { useState } from 'react'
import './App.css'

function App() {
  const [isDragOver, setIsDragOver] = useState(false)
  const [selectedFile, setSelectedFile] = useState(null)
  const [processingType, setProcessingType] = useState('summarization')
  const [summarizationPrompt, setSummarizationPrompt] = useState('Подведи итог встречи')
  const [participantsCount, setParticipantsCount] = useState('0');


  const handleDragOver = (e) => {
    e.preventDefault()
    setIsDragOver(true)
  }

  const handleDragLeave = (e) => {
    e.preventDefault()
    setIsDragOver(false)
  }

  const handleDrop = (e) => {
    e.preventDefault()
    setIsDragOver(false)
    const files = e.dataTransfer.files
    if (files.length > 0) {
      setSelectedFile(files[0])
    }
  }

  const handleFileSelect = (e) => {
    const files = e.target.files
    if (files.length > 0) {
      setSelectedFile(files[0])
    }
  }

  // собрать FormData (например в handleStart)
  const handleStart = async () => {
    if (!selectedFile) {
      alert('Выберите файл');
      return;
    }

    const fd = new FormData();
    fd.append('file', selectedFile);                       // файл записи
    fd.append('processingType', processingType);           // 'summarization' | 'transcription'
    fd.append('summarizationPrompt', summarizationPrompt); // текст промта
    fd.append('participantsCount', participantsCount);     // число участников

    try {
      const res = await fetch('/api/process', {
        method: 'POST',
        body: fd,
        // НЕ устанавливайте Content-Type вручную — браузер сам выставит нужную границу (boundary)
        headers: {
          // если нужен токен:
          // 'Authorization': `Bearer ${token}`,
        },
      });

      if (!res.ok) throw new Error(`Ошибка: ${res.status}`);
      const data = await res.json();
      console.log('Ответ сервера', data);
    } catch (err) {
      console.error(err);
    }
  }


  return (
    <div className="app">
      <div className="container">
        {/* Header */}
        <header className="header">
          <h1 className="title">Саммаризация встреч</h1>
          <button className="history-button">История</button>
        </header>

        {/* File Upload Section */}
        <section className="upload-section">
          <h2 className="upload-title">Загрузка записи</h2>
          <div 
            className={`upload-area ${isDragOver ? 'drag-over' : ''} ${selectedFile ? 'has-file' : ''}`}
            onDragOver={handleDragOver}
            onDragLeave={handleDragLeave}
            onDrop={handleDrop}
          >
            <div className="upload-content">
              {selectedFile ? (
                <div className="file-selected">
                  <div className="file-icon">📄</div>
                  <div className="file-name">{selectedFile.name}</div>
                  <button 
                    className="remove-file"
                    onClick={() => setSelectedFile(null)}
                  >
                    ×
                  </button>
                </div>
              ) : (
                <>
                  <div className="upload-icon">📥</div>
                  <p className="upload-text">Перетащите файл или нажмите для выбора</p>
                  <p className="file-formats">Поддерживаемые форматы: MP3, MP4</p>
                  <input 
                    type="file" 
                    id="file-input"
                    className="file-input"
                    accept=".mp3,.mp4"
                    onChange={handleFileSelect}
                  />
                  <label htmlFor="file-input" className="browse-button">
                    Выбрать файл
                  </label>
                </>
              )}
            </div>
          </div>
        </section>

        {/* Processing Settings */}
        <section className="settings-section">
          <h2 className="settings-title">Настройки обработки</h2>
          <div className="settings-options">
            <h3 className="process-type">Тип обработки</h3>

            <label className="option">
              <input 
                type="radio" 
                name="processingType"
                checked={processingType === 'transcription'}
                onChange={() => setProcessingType('transcription')}
              />
              <span className="checkmark radio" aria-hidden></span>
              Транскрибация
            </label>

            <label className="option">
              <input 
                type="radio" 
                name="processingType"
                checked={processingType === 'summarization'}
                onChange={() => setProcessingType('summarization')}
              />
              <span className="checkmark radio" aria-hidden></span>
              Саммаризация
            </label>
          </div>

          {/* Новая настройка: Промт для суммаризации */}
          {processingType === 'summarization' && (
            <div className="summarization-prompt">
              <label htmlFor="summarization-prompt-input" className="prompt-label">
                Промт для саммаризации
              </label>
              <input
                id="summarization-prompt-input"
                className="prompt-input"
                type="text"
                value={summarizationPrompt}
                onChange={(e) => setSummarizationPrompt(e.target.value)}
                placeholder="Подведи итог встречи"
                aria-label="Промт для суммаризации"
              />
            </div>
          )}
          <div className="participants-row">
            <label htmlFor="participants-input" className="participants-label">Количество участников</label>
            <input
              id="participants-input"
              className="participants-input"
              type="text"
              placeholder="0"
              value={participantsCount}         
              onChange={(e) => setParticipantsCount(e.target.value)}
              aria-label="Количество участников"
            />
          </div>
          <div className="start-action-row">
            <button className="start-button" onClick={handleStart}>
              <span className="start-icon" aria-hidden>
                ⬆️
              </span>
              Начать обработку
            </button>
          </div>
        </section>
      </div>
    </div>
  )
}

export default App

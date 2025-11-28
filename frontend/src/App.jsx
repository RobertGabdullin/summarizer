import { useState } from 'react'
import './App.css'

function App() {
  const [isDragOver, setIsDragOver] = useState(false)
  const [selectedFile, setSelectedFile] = useState(null)
  const [processingType, setProcessingType] = useState('summarization')
  const [summarizationPrompt, setSummarizationPrompt] = useState('Подведи итог встречи')
  const [participantsCount, setParticipantsCount] = useState('0')

  const [isUploadOpen, setIsUploadOpen] = useState(false)
  const [uploadPercent, setUploadPercent] = useState(0)
  const [isUploadComplete, setIsUploadComplete] = useState(false)
  const [uploadError, setUploadError] = useState(null)
  const xhrRef = useRef(null)

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
  const handleStart = () => {
    if (!selectedFile) {
      alert('Выберите файл');
      return;
    }

    setIsUploadOpen(true)
    setIsUploadComplete(false)
    setUploadPercent(0)
    setUploadError(null)

    const fd = new FormData();
    fd.append('file', selectedFile);
    fd.append('processingType', processingType);
    fd.append('summarizationPrompt', summarizationPrompt);
    fd.append('participantsCount', participantsCount);

    const xhr = new XMLHttpRequest();
    xhrRef.current = xhr;

    xhr.open('POST', '/api/process', true);

    xhr.upload.onprogress = (e) => {
      if (e.lengthComputable) {
        const pct = Math.round((e.loaded / e.total) * 100);
        setUploadPercent(pct);
      }
    };

    xhr.onload = () => {
      if (xhr.status >= 200 && xhr.status < 300) {
        setUploadPercent(100);
        setIsUploadComplete(true);
        try {
          const resp = JSON.parse(xhr.responseText);
          console.log('Ответ сервера:', resp);
        } catch (err) {
          console.log('Ответ сервера (не JSON)');
        }
      } else {
        setUploadError(`Ошибка сервера: ${xhr.status}`);
      }
      xhrRef.current = null;
    };

    xhr.onerror = () => {
      setUploadError('Сетевая ошибка при загрузке');
      xhrRef.current = null;
    };

    xhr.onabort = () => {
      setUploadError('Загрузка отменена');
      xhrRef.current = null;
      setIsUploadOpen(false);
    };

    xhr.send(fd);
  }

  const handleCancelUpload = () => {
    if (xhrRef.current) {
      xhrRef.current.abort();
    } else {
      // если ничего не загружается, просто закрыть окно
      setIsUploadOpen(false);
    }
  }

  const handleCloseToHome = () => {
    // закрывает окно и сбрасывает состояние
    setIsUploadOpen(false);
    setIsUploadComplete(false);
    setUploadPercent(0);
    setUploadError(null);
    setSelectedFile(null);
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
      {isUploadOpen && (
        <div>
          <div
            style={{
              width: '560px',
              maxWidth: '90vw',
              background: '#ffffff',
              borderRadius: '16px',
              padding: '28px',
              boxShadow: '0 8px 30px rgba(2,6,23,0.16)',
              textAlign: 'center',
            }}
          >
            <h3 style={{ margin: 0, fontSize: '18px', color: '#111827' }}>Загрузка записи...</h3>
            <p style={{ marginTop: '8px', marginBottom: '20px', color: '#6b7280', fontSize: '13px' }}>
              Пожалуйста, не закрывайте это окно
            </p>

            {/* Прогресс-бар */}
            <div style={{
              height: '12px',
              background: '#eef2ff',
              borderRadius: '999px',
              overflow: 'hidden',
              marginBottom: '18px'
            }}>
              <div style={{
                width: `${uploadPercent}%`,
                height: '100%',
                transition: 'width 180ms linear',
                background: '#3b82f6'
              }} />
            </div>

            {uploadError && <div style={{ color: '#ef4444', marginBottom: '12px' }}>{uploadError}</div>}

            <div style={{ display: 'flex', justifyContent: 'center', gap: '12px' }}>
              {!isUploadComplete ? (
                <button
                  className="history-button"
                  onClick={handleCancelUpload}
                  style={{ background: '#ef4444' }}
                >
                  Отменить
                </button>
              ) : (
                <button
                  className="history-button"
                  onClick={handleCloseToHome}
                >
                  На главную
                </button>
              )}
            </div>
          </div>
        </div>
      )}
    </div>
  )
}

export default App

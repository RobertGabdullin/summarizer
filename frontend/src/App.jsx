import { useState } from 'react';

function App() {
  const [isDragOver, setIsDragOver] = useState(false);
  const [selectedFile, setSelectedFile] = useState(null);
  const [processingType, setProcessingType] = useState('summarization');
  const [summarizationPrompt, setSummarizationPrompt] = useState('Подведи итог встречи');
  const [participantsCount, setParticipantsCount] = useState('0');

  const handleDragOver = (e) => {
    e.preventDefault();
    setIsDragOver(true);
  };

  const handleDragLeave = (e) => {
    e.preventDefault();
    setIsDragOver(false);
  };

  const handleDrop = (e) => {
    e.preventDefault();
    setIsDragOver(false);
    const files = e.dataTransfer.files;
    if (files.length > 0) {
      setSelectedFile(files[0]);
    }
  };

  const handleFileSelect = (e) => {
    const files = e.target.files;
    if (files.length > 0) {
      setSelectedFile(files[0]);
    }
  };

  const handleStart = async () => {
    if (!selectedFile) {
      alert('Выберите файл');
      return;
    }

    const fd = new FormData();
    fd.append('file', selectedFile);
    fd.append('processingType', processingType);
    fd.append('summarizationPrompt', summarizationPrompt);
    fd.append('participantsCount', participantsCount);

    try {
      const res = await fetch('/api/process', {
        method: 'POST',
        body: fd,
        headers: {},
      });

      if (!res.ok) throw new Error(`Ошибка: ${res.status}`);
      const data = await res.json();
      console.log('Ответ сервера', data);
    } catch (err) {
      console.error(err);
    }
  };

  return (
    <div className="min-h-screen bg-gray-100 flex items-center justify-center p-5 font-sans">
      <div className="bg-white w-full max-w-4xl min-h-[80vh] rounded-2xl shadow-lg p-8 mx-auto">
        {/* Header */}
        <header className="flex justify-between items-center mb-10 pb-5 border-b border-gray-200">
          <h1 className="text-2xl font-semibold text-gray-800">Саммаризация встреч</h1>
          <button className="bg-blue-600 text-white px-5 py-2.5 rounded-lg text-sm font-medium hover:bg-blue-700 transition-colors">
            История
          </button>
        </header>

        {/* File Upload Section */}
        <section className="mb-10">
          <h2 className="text-lg font-semibold text-gray-800 mb-4">Загрузка записи</h2>
          <div
            className={`
              border-2 border-dashed rounded-lg p-10 text-center transition-all duration-300
              ${isDragOver
                ? 'border-blue-600 bg-blue-50 scale-[1.02]'
                : selectedFile
                ? 'border-green-500 bg-green-50'
                : 'border-gray-300 bg-gray-50 hover:border-blue-500 hover:bg-blue-50'
              }
              cursor-pointer
            `}
            onDragOver={handleDragOver}
            onDragLeave={handleDragLeave}
            onDrop={handleDrop}
          >
            <div className="flex flex-col items-center gap-4">
              {selectedFile ? (
                <div className="flex items-center gap-3 p-4 bg-white rounded-lg shadow-md">
                  <div className="text-2xl">📄</div>
                  <div className="text-sm font-medium text-gray-700">{selectedFile.name}</div>
                  <button
                    className="bg-red-500 text-white w-6 h-6 rounded-full flex items-center justify-center text-lg hover:bg-red-600"
                    onClick={() => setSelectedFile(null)}
                  >
                    ×
                  </button>
                </div>
              ) : (
                <>
                  <div className="text-5xl mb-2">📥</div>
                  <p className="text-base text-gray-600">Перетащите файл или нажмите для выбора</p>
                  <p className="text-sm text-gray-500 mt-1">Поддерживаемые форматы: MP3, MP4</p>
                  <input
                    type="file"
                    id="file-input"
                    className="hidden"
                    accept=".mp3,.mp4"
                    onChange={handleFileSelect}
                  />
                  <label
                    htmlFor="file-input"
                    className="bg-blue-600 text-white px-5 py-2.5 rounded-lg text-sm font-medium cursor-pointer hover:bg-blue-700 transition-colors mt-2"
                  >
                    Выбрать файл
                  </label>
                </>
              )}
            </div>
          </div>
        </section>

        {/* Processing Settings */}
        <section className="mt-10">
          <h2 className="text-lg font-semibold text-gray-800 mb-6">Настройки обработки</h2>
          
          <div className="flex flex-col gap-3 mb-6">
            <h3 className="text-base font-semibold text-gray-800 mb-1">Тип обработки</h3>
            
            <label className="flex items-center gap-3 cursor-pointer relative">
              <input
                type="radio"
                name="processingType"
                className="absolute opacity-0"
                checked={processingType === 'transcription'}
                onChange={() => setProcessingType('transcription')}
              />
              <span className={`
                w-5 h-5 border-2 rounded-full flex items-center justify-center transition-all
                ${processingType === 'transcription'
                  ? 'border-blue-600 bg-blue-600'
                  : 'border-gray-300 hover:border-blue-500'
                }
              `}>
                {processingType === 'transcription' && (
                  <span className="w-2 h-2 rounded-full bg-white"></span>
                )}
              </span>
              <span className="text-base text-gray-700">Транскрибация</span>
            </label>

            <label className="flex items-center gap-3 cursor-pointer relative">
              <input
                type="radio"
                name="processingType"
                className="absolute opacity-0"
                checked={processingType === 'summarization'}
                onChange={() => setProcessingType('summarization')}
              />
              <span className={`
                w-5 h-5 border-2 rounded-full flex items-center justify-center transition-all
                ${processingType === 'summarization'
                  ? 'border-blue-600 bg-blue-600'
                  : 'border-gray-300 hover:border-blue-500'
                }
              `}>
                {processingType === 'summarization' && (
                  <span className="w-2 h-2 rounded-full bg-white"></span>
                )}
              </span>
              <span className="text-base text-gray-700">Саммаризация</span>
            </label>
          </div>

          {processingType === 'summarization' && (
            <div className="mb-6">
              <label htmlFor="summarization-prompt-input" className="block text-sm font-semibold text-gray-700 mb-2">
                Промт для саммаризации
              </label>
              <input
                id="summarization-prompt-input"
                className="w-full p-3 border border-gray-300 rounded-lg bg-gray-50 text-gray-900 text-sm focus:border-blue-500 focus:ring-4 focus:ring-blue-100 outline-none transition-all"
                type="text"
                value={summarizationPrompt}
                onChange={(e) => setSummarizationPrompt(e.target.value)}
                placeholder="Подведи итог встречи"
                aria-label="Промт для суммаризации"
              />
            </div>
          )}

          <div className="mb-8">
            <label htmlFor="participants-input" className="block text-sm font-semibold text-gray-700 mb-2">
              Количество участников
            </label>
            <input
              id="participants-input"
              className="w-40 p-3 border border-gray-300 rounded-lg bg-gray-50 text-gray-900 text-sm focus:border-blue-500 focus:ring-4 focus:ring-blue-100 outline-none transition-all"
              type="text"
              placeholder="0"
              value={participantsCount}
              onChange={(e) => setParticipantsCount(e.target.value)}
              aria-label="Количество участников"
            />
          </div>

          <div className="w-full mt-8">
            <button
              className="flex items-center justify-center gap-2 w-full bg-blue-600 text-white py-3 px-4 rounded-lg text-sm font-medium hover:bg-blue-700 transition-colors"
              onClick={handleStart}
            >
              <span className="text-lg">⬆️</span>
              Начать обработку
            </button>
          </div>
        </section>
      </div>
    </div>
  );
}

export default App;
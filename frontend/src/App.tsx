import { useEffect, useState } from 'react'
import { apiClient } from './api/client'

function App() {
  const [message, setMessage] = useState('Loading...')

  useEffect(() => {
    apiClient
      .GET('/api/hello')
      .then(({ data, error }) => {
        if (error) {
          throw new Error('failed to load message')
        }

        setMessage(data.message)
      })
      .catch(() => setMessage('Failed to load message'))
  }, [])

  return <h1>{message}</h1>
}

export default App

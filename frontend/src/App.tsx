import { useEffect, useState } from 'react'
import { apiClient } from './api/client'

function App() {
  const [message, setMessage] = useState('Loading...')
  const [userId, setUserId] = useState<string | null>(null)

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

    apiClient
      .GET('/api/me')
      .then(({ data, error }) => {
        if (error) {
          throw new Error('failed to load session')
        }

        setUserId(data.user_id)
      })
      .catch(() => setUserId(null))
  }, [])

  return (
    <>
      <h1>{message}</h1>
      {userId && <p>Session user: {userId}</p>}
    </>
  )
}

export default App

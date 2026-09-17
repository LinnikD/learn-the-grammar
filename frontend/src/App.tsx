import { useEffect, useState } from 'react'

type HelloResponse = {
  message: string
}

function App() {
  const [message, setMessage] = useState('Loading...')

  useEffect(() => {
    fetch('/api/hello')
      .then((response) => {
        if (!response.ok) {
          throw new Error(`HTTP ${response.status}`)
        }

        return response.json() as Promise<HelloResponse>
      })
      .then((data) => setMessage(data.message))
      .catch(() => setMessage('Failed to load message'))
  }, [])

  return <h1>{message}</h1>
}

export default App

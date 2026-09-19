import { useEffect, useState } from 'react'
import { ApiError, displayError } from './api/errors'
import { loadGreeting, loadSession } from './api/session'

type LoadState<T> =
  | { status: 'loading' }
  | { status: 'success'; data: T }
  | { status: 'error'; error: ApiError }

function ErrorMessage({ error }: { error: ApiError }) {
  return (
    <div role="alert">
      <p>{error.message}</p>
      {error.requestId && (
        <details>
          <summary>Error details</summary>
          <p>Request ID: {error.requestId}</p>
        </details>
      )}
    </div>
  )
}

function App() {
  const [message, setMessage] = useState<LoadState<string>>({
    status: 'loading',
  })
  const [userId, setUserId] = useState<LoadState<string>>({ status: 'loading' })

  useEffect(() => {
    let active = true
    let sessionLoaded = false

    loadSession()
      .then((session) => {
        sessionLoaded = true
        if (active) setUserId({ status: 'success', data: session.user_id })

        return loadGreeting()
      })
      .then((message) => {
        if (active) setMessage({ status: 'success', data: message })
      })
      .catch((cause: unknown) => {
        const error = displayError(cause)
        if (!error || !active) return

        setMessage({ status: 'error', error })
        if (!sessionLoaded) setUserId({ status: 'error', error })
      })

    return () => {
      active = false
    }
  }, [])

  return (
    <>
      {message.status === 'loading' && <h1>Loading...</h1>}
      {message.status === 'success' && <h1>{message.data}</h1>}
      {message.status === 'error' && <ErrorMessage error={message.error} />}
      {userId.status === 'loading' && <p>Loading session...</p>}
      {userId.status === 'success' && <p>Session user: {userId.data}</p>}
      {userId.status === 'error' && <ErrorMessage error={userId.error} />}
    </>
  )
}

export default App

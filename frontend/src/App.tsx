import { useEffect, useState } from 'react'
import { apiClient } from './api/client'
import { ApiError, displayError } from './api/errors'

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
    const controller = new AbortController()
    const options = { signal: controller.signal }
    apiClient
      .GET('/api/hello', options)
      .then(({ data }) => {
        if (!data || typeof data.message !== 'string') {
          throw new ApiError(
            'Something went wrong. Please try again.',
            'unexpected_response',
          )
        }
        if (!controller.signal.aborted)
          setMessage({ status: 'success', data: data.message })
      })
      .catch((cause: unknown) => {
        const error = displayError(cause)
        if (error && !controller.signal.aborted)
          setMessage({ status: 'error', error })
      })

    apiClient
      .GET('/api/me', options)
      .then(({ data }) => {
        if (!data || typeof data.user_id !== 'string') {
          throw new ApiError(
            'Something went wrong. Please try again.',
            'unexpected_response',
          )
        }
        if (!controller.signal.aborted)
          setUserId({ status: 'success', data: data.user_id })
      })
      .catch((cause: unknown) => {
        const error = displayError(cause)
        if (error && !controller.signal.aborted)
          setUserId({ status: 'error', error })
      })
    return () => controller.abort()
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

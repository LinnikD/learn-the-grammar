import { apiClient } from './client'

type Session = { user_id: string }

let loadingSession: Promise<Session> | undefined
let loadingGreeting: Promise<string> | undefined

// loadSession serializes the first session request. This also makes React
// Strict Mode's development-only effect replay reuse the same request.
export function loadSession(): Promise<Session> {
  if (!loadingSession) {
    loadingSession = apiClient
      .GET('/api/me')
      .then(({ data, error }) => {
        if (error || !data) throw new Error('failed to load session')
        return data
      })
      .catch((error: unknown) => {
        loadingSession = undefined
        throw error
      })
  }

  return loadingSession
}

// loadGreeting shares the startup request across React Strict Mode's effect
// replay and follows loadSession in App.
export function loadGreeting(): Promise<string> {
  if (!loadingGreeting) {
    loadingGreeting = apiClient
      .GET('/api/hello')
      .then(({ data, error }) => {
        if (error || !data) throw new Error('failed to load message')
        return data.message
      })
      .catch((error: unknown) => {
        loadingGreeting = undefined
        throw error
      })
  }

  return loadingGreeting
}

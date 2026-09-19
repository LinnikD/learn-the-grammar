import type { components } from './schema.gen'

type ErrorResponse = components['schemas']['ErrorResponse']

export class ApiError extends Error {
  readonly status?: number
  readonly code: string
  readonly requestId?: string

  constructor(
    message: string,
    code: string,
    requestId?: string,
    status?: number,
  ) {
    super(message)
    this.name = 'ApiError'
    this.status = status
    this.code = code
    this.requestId = requestId
  }
}

function isErrorResponse(value: unknown): value is ErrorResponse {
  if (typeof value !== 'object' || value === null) return false
  return (
    'code' in value &&
    typeof value.code === 'string' &&
    'message' in value &&
    typeof value.message === 'string' &&
    'request_id' in value &&
    typeof value.request_id === 'string'
  )
}

export async function responseError(response: Response): Promise<ApiError> {
  let body: unknown
  try {
    body = await response.clone().json()
  } catch {
    // A proxy can return HTML or an empty response instead of our contract.
  }
  if (isErrorResponse(body)) {
    return new ApiError(
      body.message,
      body.code,
      body.request_id,
      response.status,
    )
  }
  return new ApiError(
    'Something went wrong. Please try again.',
    'unexpected_response',
    response.headers.get('X-Request-ID') ?? undefined,
    response.status,
  )
}

export function displayError(error: unknown): ApiError | null {
  if (error instanceof Error && error.name === 'AbortError') return null
  if (error instanceof ApiError) return error
  if (error instanceof TypeError) {
    return new ApiError(
      'Unable to connect to the server. Please try again.',
      'network_error',
    )
  }
  return new ApiError(
    'Something went wrong. Please try again.',
    'unexpected_response',
  )
}

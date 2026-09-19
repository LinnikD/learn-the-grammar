import createClient from 'openapi-fetch'
import type { paths } from './schema.gen'
import { responseError } from './errors'

export const apiClient = createClient<paths>({ baseUrl: '/' })

apiClient.use({
  async onResponse({ response }) {
    if (!response.ok) throw await responseError(response)
  },
})

import { test, expect } from '@playwright/test'

// Isolate error rendering from backend availability and session creation.
test.beforeEach(async ({ page }) => {
  await page.route('**/api/hello', (route) =>
    route.fulfill({
      json: { message: 'Learn The Grammar!' },
    }),
  )
})

test('shows API error and request details for a failed session', async ({
  page,
}) => {
  await page.route('**/api/me', (route) =>
    route.fulfill({
      status: 500,
      json: {
        code: 'internal_error',
        message: 'Something went wrong. Please try again.',
        request_id: 'test-request-id',
      },
    }),
  )
  await page.goto('/')
  await expect(page.getByRole('alert')).toContainText('Something went wrong.')
  await page.getByText('Error details').click()
  await expect(page.getByText('Request ID: test-request-id')).toBeVisible()
  await expect(page.getByText('Loading session...')).toHaveCount(0)
})

test('shows a connection error when the request fails', async ({ page }) => {
  await page.route('**/api/me', (route) => route.abort('failed'))
  await page.goto('/')
  await expect(page.getByRole('alert')).toContainText(
    'Unable to connect to the server.',
  )
})

test('handles an HTML proxy error without exposing its body', async ({
  page,
}) => {
  await page.route('**/api/me', (route) =>
    route.fulfill({
      status: 502,
      contentType: 'text/html',
      body: '<h1>private proxy detail</h1>',
    }),
  )
  await page.goto('/')
  await expect(page.getByRole('alert')).toContainText('Something went wrong.')
  await expect(page.getByText('private proxy detail')).toHaveCount(0)
})

test('handles malformed successful responses', async ({ page }) => {
  await page.route('**/api/me', (route) => route.fulfill({ json: {} }))
  await page.goto('/')
  await expect(page.getByRole('alert')).toContainText('Something went wrong.')
})

test('shows greeting failures while keeping the successful session', async ({
  page,
}) => {
  await page.route('**/api/hello', (route) =>
    route.fulfill({
      status: 500,
      json: {
        code: 'internal_error',
        message: 'Something went wrong. Please try again.',
        request_id: 'greeting-error-id',
      },
    }),
  )
  await page.route('**/api/me', (route) =>
    route.fulfill({ json: { user_id: 'test-user' } }),
  )
  await page.goto('/')
  await expect(page.getByRole('alert')).toContainText('Something went wrong.')
  await expect(page.getByText('Session user: test-user')).toBeVisible()
  await expect(page.getByRole('heading', { name: 'Loading...' })).toHaveCount(0)
})

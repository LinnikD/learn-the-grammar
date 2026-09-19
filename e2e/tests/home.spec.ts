import { test, expect } from '@playwright/test'

test('home page displays Learn The Grammar!', async ({ page }) => {
  await page.goto('/')

  await expect(page.getByText('Learn The Grammar!')).toBeVisible()
})

test('home page displays an auto-created session user id', async ({ page }) => {
  await page.goto('/')

  const uuidPattern =
    /Session user: [0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/i
  const user = page.getByText(uuidPattern)
  await expect(user).toBeVisible()
  const initialUserID = await user.textContent()

  const cookies = await page.context().cookies()
  const sessionCookie = cookies.find((c) => c.name === 'ltg_session')

  expect(sessionCookie).toBeDefined()
  expect(sessionCookie?.httpOnly).toBe(true)
	expect(sessionCookie?.secure).toBe(true)
  expect(sessionCookie?.sameSite).toBe('Lax')

  await page.reload()
  await expect(page.getByText(initialUserID ?? '')).toBeVisible()
})

test('initializes the session before requesting the greeting', async ({
  page,
}) => {
  let sessionRequests = 0
  let helloRequests = 0
  let releaseSession: (() => void) | undefined

  await page.route('**/api/me', async (route) => {
    sessionRequests++
    await new Promise<void>((resolve) => {
      releaseSession = resolve
    })
    await route.fulfill({
      json: { user_id: '00000000-0000-4000-8000-000000000001' },
    })
  })
  await page.route('**/api/hello', async (route) => {
    helloRequests++
    await route.fulfill({ json: { message: 'Learn The Grammar!' } })
  })

  const navigation = page.goto('/')
  await expect.poll(() => releaseSession !== undefined).toBe(true)
  expect(helloRequests).toBe(0)

  releaseSession?.()
  await navigation
  await expect(page.getByText('Learn The Grammar!')).toBeVisible()
  expect(sessionRequests).toBe(1)
  expect(helloRequests).toBe(1)
})

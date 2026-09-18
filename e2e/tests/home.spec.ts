import { test, expect } from '@playwright/test'

test('home page displays Learn The Grammar!', async ({ page }) => {
  await page.goto('/')

  await expect(page.getByText('Learn The Grammar!')).toBeVisible()
})

test('home page displays an auto-created session user id', async ({ page }) => {
  await page.goto('/')

  const uuidPattern =
    /Session user: [0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/i
  await expect(page.getByText(uuidPattern)).toBeVisible()

  const cookies = await page.context().cookies()
  const sessionCookie = cookies.find((c) => c.name === 'ltg_session')

  expect(sessionCookie).toBeDefined()
  expect(sessionCookie?.httpOnly).toBe(true)
  expect(sessionCookie?.secure).toBe(true)
  expect(sessionCookie?.sameSite).toBe('Lax')
})
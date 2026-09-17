import { test, expect } from '@playwright/test'

test('home page displays Hello World', async ({ page }) => {
  await page.goto('/')

  await expect(page.getByText('Hello World')).toBeVisible()
})
import { test, expect } from '@playwright/test'

test('home page displays Learn The Grammar!', async ({ page }) => {
  await page.goto('/')

  await expect(page.getByText('Learn The Grammar!')).toBeVisible()
})
import { expect, type APIRequestContext } from '@playwright/test';

const MAILDEV_BASE = 'http://maildev:1080/maildev';

interface MaildevEmail {
  id: string;
  to: Array<{ address: string }>;
  subject: string;
  html: string;
  text: string;
}

/**
 * Poll maildev for verification email and extract the verification link.
 * MailDev API: GET /email returns all emails.
 */
export async function getVerificationLink(
  request: APIRequestContext,
  recipientEmail: string,
  timeoutMs = 30_000,
): Promise<string> {
  const start = Date.now();

  while (Date.now() - start < timeoutMs) {
    const response = await request.get(`${MAILDEV_BASE}/email`);
    expect(response.ok()).toBeTruthy();

    const emails: MaildevEmail[] = await response.json();
    const verifyEmail = emails.find((e) =>
      e.to.some((t) => t.address === recipientEmail),
    );

    if (verifyEmail) {
      // Extract verification link from email HTML or text
      const linkMatch = verifyEmail.html?.match(/href="([^"]*verify-email[^"]*)"/);
      if (linkMatch) {
        return linkMatch[1];
      }

      // Try text body
      const textMatch = verifyEmail.text?.match(/(https?:\/\/[^\s]*verify-email[^\s]*)/);
      if (textMatch) {
        return textMatch[1];
      }

      // If no verify-email link found, try any link with userId and code params
      const genericMatch = verifyEmail.html?.match(/href="([^"]*userId=[^"]*code=[^"]*)"/);
      if (genericMatch) {
        return genericMatch[1];
      }
    }

    await new Promise((r) => setTimeout(r, 2_000));
  }

  throw new Error(`No verification email found for ${recipientEmail} within ${timeoutMs}ms`);
}

/**
 * Delete all emails from maildev inbox.
 */
export async function clearInbox(request: APIRequestContext): Promise<void> {
  await request.delete(`${MAILDEV_BASE}/email/all`);
}

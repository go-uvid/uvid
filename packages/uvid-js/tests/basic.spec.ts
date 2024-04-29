import { randomUUID } from "node:crypto";
import { test, expect } from "@playwright/test";
import { type EventDTO } from "../lib/types/span";

const apiHost = "http://localhost:3000";
const SpanEndPoint = {
  session: `${apiHost}/span/session`,
  error: `${apiHost}/span/error`,
  http: `${apiHost}/span/http`,
  event: `${apiHost}/span/event`,
  performance: `${apiHost}/span/performance`,
  pageView: `${apiHost}/span/pageview`,
};

const testPageURL = `http://localhost:4000/tests/test`;

test.beforeEach(async ({ page }) => {
  await page.route(`${apiHost}/**/*`, async (route) => {
    const postData = route.request().postData();
    expect(postData).toBeTruthy();
    if (!postData) return;

    if (route.request().url().includes("/span/session")) {
      const uuid = randomUUID();
      await route.fulfill({ body: uuid, status: 204 });
    } else {
      await route.fulfill({
        status: 204,
      });
    }
  });
  // WARN route before goto page, otherwise we can't catch session request
  await Promise.all([
    page.goto(testPageURL),
    page.waitForRequest(SpanEndPoint.session),
    page.waitForRequest(SpanEndPoint.pageView),
  ]);
});

test("basic", async ({ page }) => {
  await Promise.all([
    page.locator("#fid").click(),
    page.waitForRequest(SpanEndPoint.performance),
  ]);

  await Promise.all([
    page.evaluate(async () => {
      const error = new Error("test error");
      await window.uvid.error(error);
      return error.stack;
    }),
    page.waitForRequest(SpanEndPoint.error),
  ]);

  const testEvent: EventDTO = {
    action: "event-action",
    value: "event-value",
  };
  await Promise.all([
    page.evaluate(async (testEvent_) => {
      await window.uvid.event(testEvent_.action, testEvent_.value);
    }, testEvent),
    page.waitForRequest(SpanEndPoint.event),
  ]);

  // Test if session exist after reload
  await page.reload();
  await Promise.all([
    page.locator("#register").click(),
    page.waitForRequest(SpanEndPoint.event),
  ]);
});

import { expect, test } from "vitest";

test("Node.js runtime is available", () => {
    expect(process.version).toMatch(/^v\d+\.\d+\.\d+/);
});

import { defineConfig } from "vitepress";

// https://vitepress.dev/reference/site-config
export default defineConfig({
  title: "UVID",
  description: "Observable Platform for Frontend Websites",
  themeConfig: {
    // https://vitepress.dev/reference/default-theme-config
    nav: [
      { text: "Home", link: "/" },
      { text: "Documentation", link: "/documentation" },
    ],

    sidebar: [],

    socialLinks: [{ icon: "github", link: "https://github.com/go-uvid/uvid" }],
    logo: "logo.svg",
  },
  head: [
    [
      "script",
      { type: "module" },
      `
      import { init } from "https://www.unpkg.com/uvid-js?module";
      
      window.uvid = init({
        host: "https://uvid.deno.dev",
        sessionMeta: {
          from: "uvid-site",
        },
      });
      `,
    ],
  ],
  ignoreDeadLinks: "localhostLinks",
});

// Colour-scheme bootstrap, using GoFastr's convention: the scheme lives on
// <html data-color-scheme="light|dark"> and the visitor's explicit choice in
// localStorage["gofastr.colorScheme"] (light | dark | auto). A choice saved
// under the old "lotus-theme" key is carried over once.
(() => {
  const KEY = "gofastr.colorScheme";
  const LEGACY_KEY = "lotus-theme";
  const root = document.documentElement;
  let stored = "";

  try {
    stored = localStorage.getItem(KEY) || "";
    const legacy = localStorage.getItem(LEGACY_KEY);
    if (!stored && (legacy === "light" || legacy === "dark")) {
      stored = legacy;
      localStorage.setItem(KEY, legacy);
    }
    if (legacy !== null) localStorage.removeItem(LEGACY_KEY);
  } catch (_) {}

  let scheme = stored === "light" || stored === "dark" ? stored : "";
  if (!scheme) {
    scheme = window.matchMedia?.("(prefers-color-scheme: dark)").matches ? "dark" : "light";
  }

  root.setAttribute("data-color-scheme", scheme);
  root.style.colorScheme = scheme;
  const isMenuPage = root.dataset.page === "menu";
  const themeColor = isMenuPage
    ? (scheme === "dark" ? "#0c2a1f" : "#173c2d")
    : (scheme === "dark" ? "#0a1711" : "#f5f0e5");
  document.querySelector('meta[name="theme-color"]')?.setAttribute("content", themeColor);
})();

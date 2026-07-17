(() => {
  const root = document.documentElement;
  let theme = "";

  try {
    const savedTheme = localStorage.getItem("lotus-theme");
    if (savedTheme === "light" || savedTheme === "dark") theme = savedTheme;
  } catch (_) {}

  if (!theme) {
    theme = window.matchMedia?.("(prefers-color-scheme: dark)").matches ? "dark" : "light";
  }

  root.dataset.theme = theme;
  root.style.colorScheme = theme;
  const isMenuPage = root.dataset.page === "menu";
  const themeColor = isMenuPage
    ? (theme === "dark" ? "#0c2a1f" : "#173c2d")
    : (theme === "dark" ? "#0a1711" : "#f5f0e5");
  document.querySelector('meta[name="theme-color"]')?.setAttribute("content", themeColor);
})();

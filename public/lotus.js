// Lotus Vegan's one script, loaded after GoFastr's colour-scheme bootstrap.
(() => {
  const root = document.documentElement;

  // Carry a colour scheme saved by the old site ("lotus-theme") over to
  // GoFastr's key, once. A visitor who already chose under the new key
  // keeps that choice.
  try {
    const legacy = localStorage.getItem("lotus-theme");
    if (legacy !== null) {
      localStorage.removeItem("lotus-theme");
      if (!localStorage.getItem("gofastr.colorScheme") && (legacy === "light" || legacy === "dark")) {
        if (window.__gofastr_colorScheme) window.__gofastr_colorScheme.set(legacy);
        else localStorage.setItem("gofastr.colorScheme", legacy);
      }
    }
  } catch (_) {}

  // Keep the browser chrome colour in step with the page.
  const meta = document.querySelector('meta[name="theme-color"]');
  const sync = () => {
    meta?.setAttribute("content", root.getAttribute("data-color-scheme") === "dark" ? "#0a1711" : "#f5f0e5");
  };
  sync();
  new MutationObserver(sync).observe(root, { attributes: true, attributeFilter: ["data-color-scheme"] });
})();

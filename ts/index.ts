function setupEvents() {
  document.querySelectorAll("a").forEach((el) => {
    if (!el.href.startsWith(window.location.origin)) return;
    if (el.href == window.location.origin + "/random") return;
    el.addEventListener("click", (e) => {
      e.preventDefault();
      // change page
      fetch(el.href)
        .then((resp) => resp.text())
        .then((html) => {
          const doc = new DOMParser().parseFromString(html, "text/html");
          window.history.pushState({}, "", el.href);
          document.title = doc.title;
          document.body = doc.body;
          setupEvents();
          scrollTo({
            top: 0,
            left: 0,
            behavior: "auto",
          });
        });
    });
  });
}

window.addEventListener("popstate", (_) => {
  window.location.reload();
  scrollTo({
    top: 0,
    left: 0,
    behavior: "auto",
  });
});
setupEvents();
